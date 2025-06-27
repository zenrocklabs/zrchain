package integration_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

var _ = Describe("ZenBTC ETH mint:", func() {
	var env *TestEnv
	var requestID uint64
	var bitcoinAddress string

	BeforeEach(func() {
		env = setupTestEnv(GinkgoT())
	})

	It("creates a new bitcoin key request", func() {
		hash, err := env.Tx.NewZenBTCKeyRequest(
			env.Ctx,
			"workspace1mphgzyhncnzyggfxmv4nmh",
			"keyring1k6vc6vhp6e6l3rxalue9v4ux",
			"bitcoin",
			"0xC50279996508f4562aB0B3f48D98653CE34a1667",
			types.WalletType_WALLET_TYPE_EVM,
			"eip155:17000",
			"0xC50279996508f4562aB0B3f48D98653CE34a1667",
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
		Eventually(func() uint64 {
			req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
			Expect(err).ToNot(HaveOccurred())
			return req.Id
		}, "5s", "1s").Should(Equal(requestID))
		GinkgoWriter.Printf("Bitcoin Key Request fetched: %d\n", requestID)
	})

	It("gets fulfilled within 10 seconds", func() {
		Eventually(func() string {
			req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
			Expect(err).ToNot(HaveOccurred())
			return req.Status
		}, "10s", "1s").Should(Equal(types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED.String()))
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
})
