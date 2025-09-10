package e2e

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	zenextypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

var _ = Describe("Zenex BTC flow:", func() {
	var env *TestEnv
	var shouldSkip bool
	var btcRequestID uint64
	var ecdsaRequestID uint64
	var rockAddress string
	var bitcoinAddress string

	BeforeEach(func() {
		env = setupTestEnv(GinkgoT())
		if shouldSkip {
			Skip("Previous test failed, skipping the rest")
		}
	})

	It("creates a new ecdsa key request", func() {

		hash, err := env.Tx.NewKeyRequest(
			env.Ctx,
			"workspace1mphgzyhncnzyggfxmv4nmh",
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

	It("deposits 100000 ROCK on Zrchain", func() {
		r, err := env.Docker.Exec("zenrockd", []string{"/app/send.sh", "100000", rockAddress})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("response docker cmd: %s\n", r)
	})

	It("creates a new bitcoin key request", func() {

		hash, err := env.Tx.NewKeyRequest(
			env.Ctx,
			"workspace1mphgzyhncnzyggfxmv4nmh",
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

	It("deposits 1 BTC on Bitcoin", func() {
		r, err := env.Docker.Exec("bitcoin", []string{"/app/send.sh", "1.0", bitcoinAddress})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("response docker cmd: %s\n", r)
	})

	It("verifies BTC deposit arrived", func() {
		out, err := env.Docker.Exec("bitcoin", []string{"/app/balance.sh", bitcoinAddress})
		Expect(err).ToNot(HaveOccurred())
		balance, err := extractBTCBalance(out)
		Expect(err).ToNot(HaveOccurred())
		Expect(balance).To(BeNumerically(">=", 1.0))
		GinkgoWriter.Printf("BTC balance on address %s: %f\n", bitcoinAddress, balance)
	})

	It("creates a zenex rockbtc swap request on zrchain", func() {
		// Create a zenex swap request transaction
		hash, err := env.Tx.NewMsgSwapRequest(
			env.Ctx,
			10000000000, // 100000 ROCK
			ecdsaRequestID,
			btcRequestID,
			"rockbtc",
			"workspace1mphgzyhncnzyggfxmv4nmh",
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(hash).ToNot(BeEmpty())

		// Verify the transaction was successful
		r, err := env.Tx.GetTx(env.Ctx, hash)
		Expect(err).ToNot(HaveOccurred())
		Expect(r.TxResponse).ToNot(BeNil())
		Expect(r.TxResponse.RawLog).To(BeEmpty())

		// Extract swap ID from transaction events
		var swapID uint64
		for _, event := range r.TxResponse.Events {
			if event.Type == "zenex_swap_request" {
				for _, attr := range event.Attributes {
					if attr.Key == "swap_id" {
						swapID, err = strconv.ParseUint(attr.Value, 10, 64)
						Expect(err).ToNot(HaveOccurred())
						break
					}
				}
			}
		}
		Expect(swapID).ToNot(BeZero())
		swaps, err := env.Query.Swaps(env.Ctx, "", "", "", swapID, zenextypes.SwapStatus_SWAP_STATUS_UNSPECIFIED)
		Expect(err).ToNot(HaveOccurred())
		Expect(swaps.Swaps).ToNot(BeEmpty())
		Expect(swaps.Swaps[0].SwapId).To(Equal(swapID))
		Expect(swaps.Swaps[0].Status).To(Equal(zenextypes.SwapStatus_SWAP_STATUS_REQUESTED))

		GinkgoWriter.Printf("Zenex swap request created with swap ID: %d, tx hash: %s\n", swapID, hash)
	})

})
