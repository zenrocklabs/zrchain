package e2e

import (
	"strconv"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	zenextypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Zenex BTC flow:", func() {
	var env *TestEnv
	var shouldSkip bool
	var btcRequestID uint64
	var ecdsaRequestID uint64
	var btcPoolKeyRequestId uint64
	var rockAddress string
	var bitcoinAddress string
	var bitcoinPoolKeyAddress string
	var workspaceAddress string
	var swapID uint64
	var amountRockIn uint64
	BeforeEach(func() {
		env = setupTestEnv(GinkgoT())
		if shouldSkip {
			Skip("Previous test failed, skipping the rest")
		}
	})

	// Temporary to log the params (which are empty for some reason)
	It("gets the zenex params", func() {
		params, err := env.Query.ZenexQueryClient.Params(env.Ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(params.Params).ToNot(BeNil())
		GinkgoWriter.Printf("Zenex params: %v\n", params.Params)
	})

	// Isolated workspace for this test only
	It("creates a new workspace", func() {
		hash, err := env.Tx.NewWorkspace(env.Ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(hash).ToNot(BeEmpty())
		r, err := env.Tx.GetTx(env.Ctx, hash)
		Expect(err).ToNot(HaveOccurred())
		Expect(r.TxResponse).ToNot(BeNil())
		Expect(r.TxResponse.RawLog).To(BeEmpty())
		for _, event := range r.TxResponse.Events {
			if event.Type == "new_workspace" {
				for _, attr := range event.Attributes {
					if attr.Key == "workspace_addr" {
						workspaceAddress = attr.Value
						break
					}
				}
			}
		}
	})

	// Require a new ecdsa key to deposit ROCK
	It("creates a new ecdsa key request", func() {
		hash, err := env.Tx.NewKeyRequest(
			env.Ctx,
			workspaceAddress,
			"keyring1k6vc6vhp6e6l3rxalue9v4ux",
			"ecdsa",
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(hash).ToNot(BeEmpty())
		r, err := env.Tx.GetTx(env.Ctx, hash)
		Expect(err).ToNot(HaveOccurred())
		Expect(r.TxResponse).ToNot(BeNil())
		Expect(r.TxResponse.RawLog).To(BeEmpty())

		var ecdsaRequestIDStr string
		for _, event := range r.TxResponse.Events {
			if event.Type == "new_key_request" {
				for _, attr := range event.Attributes {
					if attr.Key == "request_id" {
						ecdsaRequestIDStr = attr.Value
						break
					}
				}
			}
		}

		ecdsaRequestID, err = strconv.ParseUint(ecdsaRequestIDStr, 10, 64)
		Expect(err).ToNot(HaveOccurred())
		Expect(ecdsaRequestID).ToNot(BeNil())
		GinkgoWriter.Printf("ECDSA Key Request created: %d\n", ecdsaRequestID)
	})

	// Waiting for the ecdsa key request to be included onchain
	It("fetches the ecdsa key request within 5 seconds", func() {
		Eventually(func() (uint64, error) {
			req, err := env.Query.GetKeyRequest(env.Ctx, ecdsaRequestID)
			if err != nil {
				return 0, err
			}
			return req.Id, nil
		}, "5s", "1s").Should(Equal(ecdsaRequestID))
		GinkgoWriter.Printf("ECDSA Key Request fetched: %d\n", ecdsaRequestID)
	})

	// Waiting for the ecdsa key request to be fulfilled
	It("gets fulfilled within 15 seconds", func() {
		Eventually(func() (string, error) {
			req, err := env.Query.GetKeyRequest(env.Ctx, ecdsaRequestID)
			if err != nil {
				return "", err
			}
			return req.Status, nil
		}, "15s", "1s").Should(Equal(types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED.String()))
		GinkgoWriter.Printf("ECDSA Key Request fulfilled: %d\n", ecdsaRequestID)
	})

	// Fetching the ecdsa zrchain address from the returned public key
	It("fetches the ecdsa zrchain address", func() {
		req, err := env.Query.GetKey(env.Ctx, ecdsaRequestID)
		Expect(err).ToNot(HaveOccurred())
		for _, w := range req.Wallets {
			if w.Type == "WALLET_TYPE_NATIVE" {
				rockAddress = w.Address
			}
		}
		Expect(rockAddress).ToNot(BeEmpty())
		GinkgoWriter.Printf("zrchain address: %s\n", rockAddress)
	})

	// Depositing 100000 ROCK on Zrchain to convert to btc
	It("deposits 100000 ROCK on Zrchain", func() {
		r, err := env.Tx.BankTxClient.SendCoins(env.Ctx, env.Ident.Address.String(), rockAddress, 100000000000, "urock")
		Expect(err).ToNot(HaveOccurred())
		Expect(r).ToNot(BeEmpty())
		GinkgoWriter.Printf("response docker cmd: %s\n", r)
	})

	// Creating a new bitcoin key request
	It("creates a new bitcoin key request", func() {
		hash, err := env.Tx.NewKeyRequest(
			env.Ctx,
			workspaceAddress,
			"keyring1k6vc6vhp6e6l3rxalue9v4ux",
			"bitcoin",
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
		btcRequestID, err = strconv.ParseUint(requestIDStr, 10, 64)
		Expect(err).ToNot(HaveOccurred())
		Expect(btcRequestID).ToNot(BeNil())
		GinkgoWriter.Printf("Bitcoin Key Request created: %d\n", btcRequestID)
	})

	// Waiting for the bitcoin key request to be included onchain
	It("fetches the bitcoin key request within 5 seconds", func() {
		Eventually(func() (uint64, error) {
			req, err := env.Query.GetKeyRequest(env.Ctx, btcRequestID)
			if err != nil {
				return 0, err
			}
			return req.Id, nil
		}, "5s", "1s").Should(Equal(btcRequestID))
		GinkgoWriter.Printf("Bitcoin Key Request fetched: %d\n", btcRequestID)
	})

	// Waiting for the bitcoin key request to be fulfilled
	It("gets fulfilled within 15 seconds", func() {
		Eventually(func() (string, error) {
			req, err := env.Query.GetKeyRequest(env.Ctx, btcRequestID)
			if err != nil {
				return "", err
			}
			return req.Status, nil
		}, "15s", "1s").Should(Equal(types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED.String()))
		GinkgoWriter.Printf("Bitcoin Key Request fulfilled: %d\n", btcRequestID)
	})

	// Fetching the bitcoin REGNET key address from the returned public key
	It("fetches the bitcoin REGNET key address", func() {
		req, err := env.Query.GetKey(env.Ctx, btcRequestID)
		Expect(err).ToNot(HaveOccurred())
		for _, w := range req.Wallets {
			if w.Type == "WALLET_TYPE_BTC_REGNET" {
				bitcoinAddress = w.Address
			}
		}
		Expect(bitcoinAddress).ToNot(BeEmpty())
		GinkgoWriter.Printf("Bitcoin REGNET address: %s\n", bitcoinAddress)
	})

	// Depositing 1 BTC on Bitcoin
	It("deposits 1 BTC on Bitcoin", func() {
		r, err := env.Docker.Exec("bitcoin", []string{"/app/send.sh", "1.0", bitcoinAddress})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("response docker cmd: %s\n", r)
	})

	// Verifying BTC deposit arrived
	It("verifies BTC deposit arrived", func() {
		out, err := env.Docker.Exec("bitcoin", []string{"/app/balance.sh", bitcoinAddress})
		Expect(err).ToNot(HaveOccurred())
		balance, err := extractBTCBalance(out)
		Expect(err).ToNot(HaveOccurred())
		Expect(balance).To(BeNumerically(">=", 1.0))
		GinkgoWriter.Printf("BTC balance on address %s: %f\n", bitcoinAddress, balance)
	})
	// Creating a new bitcoin key request for the poolkeyid
	It("creates a new bitcoin key request for the pool keyid ", func() {

		hash, err := env.Tx.NewKeyRequest(
			env.Ctx,
			workspaceAddress,
			"keyring1k6vc6vhp6e6l3rxalue9v4ux",
			"bitcoin",
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

		btcPoolKeyRequestId, err = strconv.ParseUint(requestIDStr, 10, 64)
		Expect(err).ToNot(HaveOccurred())
		Expect(btcRequestID).ToNot(BeNil())
		GinkgoWriter.Printf("Bitcoin Key Request for pool address created: %d\n", btcPoolKeyRequestId)
	})

	// Waiting for the bitcoin key request to be included onchain
	It("fetches the bitcoin key request within 5 seconds", func() {
		Eventually(func() (uint64, error) {
			req, err := env.Query.GetKeyRequest(env.Ctx, btcPoolKeyRequestId)
			if err != nil {
				return 0, err
			}
			return req.Id, nil
		}, "5s", "1s").Should(Equal(btcPoolKeyRequestId))
		GinkgoWriter.Printf("Bitcoin Key Request fetched: %d\n", btcPoolKeyRequestId)
	})

	// Waiting for the bitcoin key request to be fulfilled
	It("gets fulfilled within 15 seconds", func() {
		Eventually(func() (string, error) {
			req, err := env.Query.GetKeyRequest(env.Ctx, btcPoolKeyRequestId)
			if err != nil {
				return "", err
			}
			return req.Status, nil
		}, "15s", "1s").Should(Equal(types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED.String()))
		GinkgoWriter.Printf("Bitcoin Key Request fulfilled: %d\n", btcPoolKeyRequestId)
	})

	// Fetching the bitcoin REGNET key address from the returned public key
	It("fetches the bitcoin REGNET key address", func() {
		req, err := env.Query.GetKey(env.Ctx, btcPoolKeyRequestId)
		Expect(err).ToNot(HaveOccurred())
		for _, w := range req.Wallets {
			if w.Type == "WALLET_TYPE_BTC_REGNET" {
				bitcoinPoolKeyAddress = w.Address
			}
		}
		Expect(bitcoinPoolKeyAddress).ToNot(BeEmpty())
		GinkgoWriter.Printf("Bitcoin REGNET address: %s\n", bitcoinPoolKeyAddress)
	})

	// Creating a zenex rockbtc swap request on zrchain
	It("creates a zenex rockbtc swap request on zrchain", func() {
		amountRockIn = 30000000000 // 30000 ROCK

		// Create a zenex swap request transaction
		hash, err := env.Tx.NewMsgSwapRequest(
			env.Ctx,
			amountRockIn,
			ecdsaRequestID,
			btcRequestID,
			zenextypes.TradePair_TRADE_PAIR_ROCK_BTC,
			workspaceAddress,
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(hash).ToNot(BeEmpty())

		r, err := env.Tx.GetTx(env.Ctx, hash)
		Expect(err).ToNot(HaveOccurred())
		Expect(r.TxResponse).ToNot(BeNil())
		Expect(r.TxResponse.RawLog).To(BeEmpty())

		// Extract swap ID from transaction events
		for _, event := range r.TxResponse.Events {
			if event.Type == "swap_request" {
				for _, attr := range event.Attributes {
					if attr.Key == "swap_id" {
						swapID, err = strconv.ParseUint(attr.Value, 10, 64)
						Expect(err).ToNot(HaveOccurred())
						break
					}
				}
			}
		}
		GinkgoWriter.Printf("Zenex swap request created with swap ID: %d, tx hash: %s\n", swapID, hash)
	})

	It("checks the swap status to be SWAP_STATUS_INITIATED", func() {
		Expect(swapID).ToNot(BeZero())
		swaps, err := env.Query.Swaps(env.Ctx, "", "", "", swapID, zenextypes.SwapStatus_SWAP_STATUS_UNSPECIFIED, zenextypes.TradePair_TRADE_PAIR_ROCK_BTC)
		Expect(err).ToNot(HaveOccurred())
		Expect(swaps.Swaps).ToNot(BeEmpty())
		Expect(swaps.Swaps[0].SwapId).To(Equal(swapID))
		Expect(swaps.Swaps[0].Status).To(Equal(zenextypes.SwapStatus_SWAP_STATUS_INITIATED))
	})

	// checks the swap amount is in the module account
	It("checks the swap amount is in the module account", func() {
		rockPool, err := env.Query.RockPool(env.Ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(rockPool.RockBalance).To(BeNumerically("==", amountRockIn))
		GinkgoWriter.Printf("Rock balance in module account: %d\n", rockPool.RockBalance)
	})

	// Waiting for the swap transfer request from the zenex proxy to be included onchain
	It("waits for the transfer request of the zenex proxy within 5 seconds", func() {
		Eventually(func() (zenextypes.SwapStatus, error) {
			swaps, err := env.Query.Swaps(env.Ctx, "", "", "", swapID, zenextypes.SwapStatus_SWAP_STATUS_REQUESTED, zenextypes.TradePair_TRADE_PAIR_ROCK_BTC)
			if err != nil {
				return 0, err
			}
			return swaps.Swaps[0].Status, nil
		}, "5s", "1s").Should(Equal(zenextypes.SwapStatus_SWAP_STATUS_REQUESTED.String()))
	})

	// Waiting for the ecdsa key request to be included onchain
	It("waits for the transfer request of the zenex proxy within 5 seconds", func() {
		Eventually(func() (zenextypes.SwapStatus, error) {
			swaps, err := env.Query.Swaps(env.Ctx, "", "", "", swapID, zenextypes.SwapStatus_SWAP_STATUS_COMPLETED, zenextypes.TradePair_TRADE_PAIR_ROCK_BTC)
			if err != nil {
				return 0, err
			}
			return swaps.Swaps[0].Status, nil
		}, "5s", "1s").Should(Equal(zenextypes.SwapStatus_SWAP_STATUS_COMPLETED.String()))
	})
})
