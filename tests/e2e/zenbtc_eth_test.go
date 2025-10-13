package e2e

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/Zenrock-Foundation/zrchain/v6/zenbtc/bindings"
	zentype "github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"
)

var _ = Describe("ZenBTC ETH flow:", func() {
	var env *TestEnv
	var shouldSkip bool
	var requestID uint64
	var bitcoinAddress string
	var ethereumWallet *ecdsa.PrivateKey
	var ethereumAddress common.Address
	var initialBTCBalance float64
	var initialRedemptions int
	var randomBTCAddress string
	var burnTxHash string

	BeforeEach(func() {
		env = setupTestEnv(GinkgoT())
		if shouldSkip {
			Skip("Previous test failed, skipping the rest")
		}
	})

	// Create an ethereum account in anvil
	It("creates a ethereum wallet", func() {
		r, err := env.Docker.Exec("anvil", []string{
			"cast", "wallet", "new",
		})
		Expect(err).ToNot(HaveOccurred())
		ethereumWallet, err = ethereumWalletFromOutput(r)
		Expect(err).ToNot(HaveOccurred())
		ethereumAddress = crypto.PubkeyToAddress(ethereumWallet.PublicKey)
		GinkgoWriter.Printf("Ethereum pub address: %s\n", ethereumAddress.String())
	})

	It("Fund the wallet with 10 ETH", func() {
		_, err := env.Docker.Exec("anvil", []string{
			"cast", "rpc", "anvil_setBalance", ethereumAddress.String(), "0x8AC7230489E80000",
		})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Funded %v with 1 ETH\n", ethereumAddress.String())
	})

	It("creates a new bitcoin key request", func() {
		hash, err := env.Tx.NewZenBTCKeyRequest(
			env.Ctx,
			"workspace1mphgzyhncnzyggfxmv4nmh",
			"keyring1k6vc6vhp6e6l3rxalue9v4ux",
			"bitcoin",
			ethereumAddress.String(),
			types.WalletType_WALLET_TYPE_EVM,
			"eip155:560048",
			"",
			dcttypes.Asset_ASSET_ZENBTC,
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(hash).ToNot(BeEmpty())

		r, err := env.Tx.GetTx(env.Ctx, hash)
		Expect(err).ToNot(HaveOccurred())
		Expect(r.TxResponse).ToNot(BeNil())
		Expect(r.TxResponse.RawLog).To(BeEmpty())

		var requestIDStr string
		for _, event := range r.TxResponse.Events {
			if event.Type == "new_key_request" {
				for _, attr := range event.Attributes {
					if attr.Key == "request_id" {
						requestIDStr = attr.Value
						break
					}
				}
			}
		}
		requestID, err = strconv.ParseUint(requestIDStr, 10, 64)
		Expect(err).ToNot(HaveOccurred())
		Expect(requestID).ToNot(BeNil())
		GinkgoWriter.Printf("Bitcoin Key Request created: %d\n", requestID)
	})

	It("fetches the bitcoin key request within 5 seconds", func() {
		Eventually(func() (uint64, error) {
			req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
			if err != nil {
				return 0, err
			}
			return req.Id, nil
		}, "5s", "1s").Should(Equal(requestID))
		GinkgoWriter.Printf("Bitcoin Key Request fetched: %d\n", requestID)
	})

	It("gets fulfilled within 15 seconds", func() {
		Eventually(func() (string, error) {
			req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
			if err != nil {
				return "", err
			}
			return req.Status, nil
		}, "15s", "1s").Should(Equal(types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED.String()))
		GinkgoWriter.Printf("Bitcoin Key Request fulfilled: %d\n", requestID)
	})

	It("fetches the bitcoin REGNET key address", func() {
		req, err := env.Query.GetKey(env.Ctx, requestID)
		Expect(err).ToNot(HaveOccurred())
		for _, w := range req.Wallets {
			if w.Type == "WALLET_TYPE_BTC_REGNET" {
				bitcoinAddress = w.Address
			}
		}
		Expect(bitcoinAddress).ToNot(BeEmpty())
		GinkgoWriter.Printf("Bitcoin REGNET address: %s\n", bitcoinAddress)
	})

	It("deposits on Bitcoin", func() {
		r, err := env.Docker.Exec("bitcoin", []string{"/app/send.sh", "1.0", bitcoinAddress})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("response docker cmd: %s\n", r)
	})

	It("creates a mint transaction on zrchain", func() {
		var lastCount int
		var err error
		// Initial count of pending mint transactions
		initialResp, err := env.Query.PendingMintTransactions(env.Ctx, 0)
		Expect(err).ToNot(HaveOccurred())
		lastCount = len(initialResp.PendingMintTransactions)

		var newTx zentype.PendingMintTransaction

		Eventually(func() (int, error) {
			resp, err := env.Query.PendingMintTransactions(env.Ctx, 0)
			if err != nil {
				return 0, err
			}

			if len(resp.PendingMintTransactions) > lastCount {
				newTx = *resp.PendingMintTransactions[len(resp.PendingMintTransactions)-1]
			}

			return len(resp.PendingMintTransactions), nil
		}, "30s", "2s").Should(BeNumerically(">", lastCount))

		Expect(newTx.Status).To(Equal(zentype.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED))
		GinkgoWriter.Printf("Mint transaction created with ID %d\n", newTx.Id)
	})

	It("mint gets staked", func() {
		Eventually(func() (zentype.MintTransactionStatus, error) {
			resp, err := env.Query.PendingMintTransactions(env.Ctx, 0)
			if err != nil {
				return 0, err
			}
			if len(resp.PendingMintTransactions) == 0 {
				return 0, nil // Return 0 status if no transactions
			}
			lastTx := *resp.PendingMintTransactions[len(resp.PendingMintTransactions)-1]

			return lastTx.Status, nil
		}, "90s", "5s").Should(Equal(zentype.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED))
		GinkgoWriter.Printf("Mint transaction moved to staked \n")
	})

	It("mint gets minted", func() {
		Eventually(func() (zentype.MintTransactionStatus, error) {
			resp, err := env.Query.PendingMintTransactions(env.Ctx, 0)
			if err != nil {
				return 0, err
			}
			if len(resp.PendingMintTransactions) == 0 {
				return 0, nil // Return 0 status if no transactions
			}
			lastTx := *resp.PendingMintTransactions[len(resp.PendingMintTransactions)-1]

			return lastTx.Status, nil
		}, "90s", "5s").Should(Equal(zentype.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED))
		GinkgoWriter.Printf("Mint transaction moved to minted \n")
	})

	It("creates a burn", func() {
		var burnTx zentype.BurnEvent
		client, err := ethclient.Dial("http://localhost:8545")
		Expect(err).ToNot(HaveOccurred())
		// Create random BTC address
		randomBTCAddress, err = randomBTCRegnetAddress()
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Generated destination BTC address: %s\n", randomBTCAddress)
		// Create transaction signer
		auth, err := bind.NewKeyedTransactorWithChainID(ethereumWallet, big.NewInt(ETHEREUM_CHAIN_ID))
		Expect(err).ToNot(HaveOccurred())
		auth.GasLimit = uint64(300_000)
		auth.Context = env.Ctx

		// Load contract
		contractAddr := common.HexToAddress(ETHEREUM_ZENBTC_CONTRACT)
		zenbtc, err := bindings.NewZenBTC(contractAddr, client)
		Expect(err).ToNot(HaveOccurred())

		// Call unwrap()
		amount := big.NewInt(1e6) // 1 token (6 decimals)
		dest := []byte(randomBTCAddress)

		chainID := "eip155:" + strconv.FormatUint(ETHEREUM_CHAIN_ID, 10)
		tx, err := zenbtc.Unwrap(auth, amount, dest)
		Expect(err).ToNot(HaveOccurred())
		receipt, err := bind.WaitMined(env.Ctx, client, tx)
		Expect(err).ToNot(HaveOccurred())
		burnTxHash = receipt.TxHash.String()
		GinkgoWriter.Printf("Broadcasted TX - hash: %v block: %d\n", burnTxHash, receipt.BlockNumber)
		// fast-forward anvil here to speed up the confirmations needed for EL
		_, err = env.Docker.Exec("anvil", []string{
			"cast", "rpc", "anvil_mine", "50",
		})
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool {
			resp, err := env.Query.BurnEvents(env.Ctx, 0, burnTxHash, 0, chainID)
			Expect(err).ToNot(HaveOccurred())

			if len(resp.BurnEvents) > 0 {
				burnTx = *resp.BurnEvents[0]
				return true
			}

			return false
		}, "120s", "5s").Should(BeTrue())
		GinkgoWriter.Printf("Burn tx created: %v\n", burnTx)
		respRedemptions, err := env.Query.Redemptions(env.Ctx, 0, zentype.RedemptionStatus(-1))
		Expect(err).ToNot(HaveOccurred())
		initialRedemptions = len(respRedemptions.Redemptions)
	})

	It("burn status moves to UNSTAKING", func() {
		chainID := "eip155:" + strconv.FormatUint(ETHEREUM_CHAIN_ID, 10)
		Eventually(func() (zentype.BurnStatus, error) {
			resp, err := env.Query.BurnEvents(env.Ctx, 0, burnTxHash, 0, chainID)
			if err != nil {
				return 0, err
			}
			if len(resp.BurnEvents) == 0 {
				return 0, errors.New("0 burn events returned")
			}
			return resp.BurnEvents[0].Status, nil
		}, "120s", "5s").Should(Equal(zentype.BurnStatus_BURN_STATUS_UNSTAKING))
	})

	// This is needed for the unstake on EL
	It("fast-forward anvil 50 blocks", func() {
		_, err := env.Docker.Exec("anvil", []string{
			"cast", "rpc", "anvil_mine", "50",
		})
		Expect(err).ToNot(HaveOccurred())
	})

	// Redemption gets created
	It("redemption is created", func() {
		var redemption zentype.Redemption
		Eventually(func() bool {
			resp, err := env.Query.Redemptions(env.Ctx, 0, zentype.RedemptionStatus(-1))
			Expect(err).ToNot(HaveOccurred())
			if len(resp.Redemptions) <= int(initialRedemptions) {
				return false
			}
			for _, event := range resp.Redemptions {
				if string(event.Data.DestinationAddress) == randomBTCAddress {
					redemption = event
					return true
				}
			}

			return false
		}, "180s", "10s").Should(BeTrue())
		GinkgoWriter.Printf("Redemption created: %v\n", redemption)
	})

	// Get balance before redemption completes
	It("gets initial BTC balance on destination address", func() {
		out, err := env.Docker.Exec(
			"bitcoin",
			[]string{"/app/balance.sh", randomBTCAddress},
		)
		Expect(err).ToNot(HaveOccurred())
		initialBTCBalance, err = extractBTCBalance(out)
		Expect(err).ToNot(HaveOccurred())
	})

	// Redemption gets completed
	It("redemption is completed", func() {
		var redemption zentype.Redemption
		Eventually(func() bool {
			resp, err := env.Query.Redemptions(env.Ctx, 0, zentype.RedemptionStatus(-1))
			Expect(err).ToNot(HaveOccurred())
			if len(resp.Redemptions) <= int(initialRedemptions) {
				return false
			}
			for _, event := range resp.Redemptions {
				if string(event.Data.DestinationAddress) == randomBTCAddress {
					if event.Status == zentype.RedemptionStatus_COMPLETED {
						redemption = event
						return true
					}
				}
			}

			return false
		}, "180s", "10s").Should(BeTrue())
		GinkgoWriter.Printf("Redemption completed: %v\n", redemption)
	})

	// Mine 100 blocks so the balance moves spendable
	It("mine bitcoin blocks", func() {
		time.Sleep(10 * time.Second)
		_, err := env.Docker.Exec(
			"bitcoin",
			[]string{"/app/mine.sh", "100", "bcrt1qa8fk4zz3j8revg0frgvtt6dkrnw9snhan49cx8"},
		)
		Expect(err).ToNot(HaveOccurred())
	})

	// Get balance after redemption completes
	It("balance on destination address increased", func() {
		out, err := env.Docker.Exec(
			"bitcoin",
			[]string{"/app/balance.sh", randomBTCAddress},
		)
		Expect(err).ToNot(HaveOccurred())
		balance, err := extractBTCBalance(out)
		Expect(err).ToNot(HaveOccurred())
		Expect(balance).To(BeNumerically(">", initialBTCBalance))
		GinkgoWriter.Printf("Balance changed from %f to %f\n", initialBTCBalance, balance)
	})

	AfterEach(func() {
		if CurrentSpecReport().Failed() {
			shouldSkip = true
		}
	})
})
