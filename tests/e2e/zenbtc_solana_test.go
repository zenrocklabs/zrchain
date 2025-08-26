package e2e

import (
	"errors"
	"os"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	eventstore "github.com/Zenrock-Foundation/zrchain/v6/contracts/sol-event-store/go-sdk"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc/generated/zenbtc_spl_token"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	zentype "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

var _ = Describe("ZenBTC Solana flow:", func() {
	var env *TestEnv
	var shouldSkip bool
	var requestID uint64
	var solanaAccount string
	var bitcoinAddress string
	var burnTxHash solana.Signature
	var initialBTCBalance float64
	var initialRedemptions int
	var randomBTCAddress string

	BeforeEach(func() {
		env = setupTestEnv(GinkgoT())
		if shouldSkip {
			Skip("Previous test failed, skipping the rest")
		}
	})

	It("creates a solana account", func() {
		r, err := env.Docker.Exec("solana", []string{
			"solana-keygen", "new",
			"--no-passphrase",
			"--force",
			"-o", "/root/.config/solana/id.json",
		})
		Expect(err).ToNot(HaveOccurred())
		solanaAccount, err = extractSolanaPubkey(r)
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Solana account created: %s\n", solanaAccount)
	})

	It("funds the solana account", func() {
		r, err := env.Docker.Exec("solana", []string{
			"solana", "airdrop",
			"100",
			solanaAccount,
		})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Solana account funded: %v\n", r)
	})

	It("creates a new bitcoin key request", func() {
		hash, err := env.Tx.NewZenBTCKeyRequest(
			env.Ctx,
			"workspace1mphgzyhncnzyggfxmv4nmh",
			"keyring1k6vc6vhp6e6l3rxalue9v4ux",
			"bitcoin",
			solanaAccount,
			types.WalletType_WALLET_TYPE_SOLANA,
			"solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD",
			"",
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
		}, "150s", "5s").Should(Equal(zentype.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED))
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
		}, "150s", "5s").Should(Equal(zentype.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED))
		GinkgoWriter.Printf("Mint transaction moved to minted \n")
	})

	It("creates a burn", func() {
		var burnTx zentype.BurnEvent
		var client = rpc.New("http://localhost:8899")
		// Read private key from solana container
		out, err := env.Docker.Exec("solana", []string{
			"cat", "/root/.config/solana/id.json",
		})
		Expect(err).ToNot(HaveOccurred())
		// Create random BTC address
		randomBTCAddress, err = randomBTCRegnetAddress()
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Generated destination BTC address: %s\n", randomBTCAddress)

		signerWallet, err := LoadSolanaPrivateKeyFromJSON(out)
		Expect(err).ToNot(HaveOccurred())
		programId, err := solana.PublicKeyFromBase58(SOLANA_ZENBTC_PROGRAM_ID)
		Expect(err).ToNot(HaveOccurred())
		signer, err := solana.PublicKeyFromBase58(solanaAccount)
		Expect(err).ToNot(HaveOccurred())
		feeWallet, err := solana.PublicKeyFromBase58(SOLANA_FEE_WALLET)
		Expect(err).ToNot(HaveOccurred())
		mintAddress, err := solana.PublicKeyFromBase58(SOLANA_ZENBTC_TOKEN_ADDRESS)
		Expect(err).ToNot(HaveOccurred())
		multisigAddress, err := solana.PublicKeyFromBase58(SOLANA_ZENBTC_MULTISIG)
		Expect(err).ToNot(HaveOccurred())

		latest, err := client.GetLatestBlockhash(env.Ctx, rpc.CommitmentProcessed)
		Expect(err).ToNot(HaveOccurred())
		// Derive EventStore PDAs similar to keeper logic (unwrap path).
		var eventStoreProgram, eventStoreGlobalConfig, zenbtcUnwrapShardPDA solana.PublicKey
		// Allow optional injection via env var (e.g., SOLANA_EVENTSTORE_PROGRAM_ID); empty means placeholders remain zero.
		if esProgramID := os.Getenv("SOLANA_EVENTSTORE_PROGRAM_ID"); esProgramID != "" {
			if parsed, err := solana.PublicKeyFromBase58(esProgramID); err == nil {
				eventStoreProgram = parsed
				gc, gcErr := eventstore.DeriveGlobalConfigPDA(eventStoreProgram)
				if gcErr == nil {
					eventStoreGlobalConfig = gc
				}
				// For tests we don't know the burn event ID beforehand; if desired, set SOLANA_UNWRAP_EVENT_ID to force shard derivation.
				if unwrapIDStr := os.Getenv("SOLANA_UNWRAP_EVENT_ID"); unwrapIDStr != "" {
					if unwrapID, convErr := strconv.ParseUint(unwrapIDStr, 10, 64); convErr == nil && unwrapID > 0 {
						if shard, _, shardErr := eventstore.DeriveZenbtcUnwrapShardPDA(eventStoreProgram, unwrapID); shardErr == nil {
							zenbtcUnwrapShardPDA = shard
						}
					}
				}
			}
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				solzenbtc.Unwrap(
					programId,
					zenbtc_spl_token.UnwrapArgs{
						Value:    uint64(1000000),
						DestAddr: []byte(randomBTCAddress),
					},
					signer,
					mintAddress,
					multisigAddress,
					feeWallet,
					eventStoreProgram,      // eventStore program (zero if not set)
					eventStoreGlobalConfig, // global config PDA (zero if not derivable)
					programId,              // calling program (zenbtc)
					zenbtcUnwrapShardPDA,   // unwrap shard PDA (optional)
				),
			},
			latest.Value.Blockhash,
			solana.TransactionPayer(signer),
		)
		Expect(err).ToNot(HaveOccurred())

		// Sign transaction
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if key.Equals(signerWallet.PublicKey()) {
					return &signerWallet.PrivateKey
				}
				return nil
			},
		)
		Expect(err).ToNot(HaveOccurred())
		burnTxHash, err = client.SendTransactionWithOpts(env.Ctx, tx, rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentProcessed,
		})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Signature of burn: %s\n", burnTxHash.String())
		// fast-forward anvil here to speed up the confirmations needed for EL
		_, err = env.Docker.Exec("anvil", []string{
			"cast", "rpc", "anvil_mine", "50",
		})
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool {
			resp, err := env.Query.BurnEvents(env.Ctx, 0, burnTxHash.String(), 0, "solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD")
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
		Eventually(func() (zentype.BurnStatus, error) {
			resp, err := env.Query.BurnEvents(env.Ctx, 0, burnTxHash.String(), 0, "solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD")
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

	// Mine 100 blocks so the balance moves to spendable
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
