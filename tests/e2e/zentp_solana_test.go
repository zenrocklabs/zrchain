package integration_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

var _ = Describe("ZenTP Solana:", func() {
	var env *TestEnv
	var requestID uint64
	var solanaAddress string

	BeforeEach(func() {
		env = setupTestEnv(GinkgoT())
	})

	It("creates a new EDDSA key request", func() {
		hash, err := env.Tx.NewKeyRequest(
			env.Ctx,
			"workspace1mphgzyhncnzyggfxmv4nmh",
			"keyring1k6vc6vhp6e6l3rxalue9v4ux",
			"eddsa",
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
		GinkgoWriter.Printf("EDDSA Key Request created: %d\n", requestID)
	})

	It("fetches the EDDSA key request within 5 seconds", func() {
		Eventually(func() uint64 {
			req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
			Expect(err).ToNot(HaveOccurred())
			return req.Id
		}, "5s", "1s").Should(Equal(requestID))
		GinkgoWriter.Printf("EDDSA Key Request fetched: %d\n", requestID)
	})

	It("gets fulfilled within 15 seconds", func() {
		Eventually(func() string {
			req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
			Expect(err).ToNot(HaveOccurred())
			return req.Status
		}, "15s", "1s").Should(Equal(types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED.String()))
		GinkgoWriter.Printf("EDDSA Key Request fulfilled: %d\n", requestID)
	})

	It("fetches the solana address", func() {
		req, err := env.Query.GetKey(env.Ctx, requestID)
		Expect(err).ToNot(HaveOccurred())
		for _, w := range req.Wallets {
			if w.Type == "WALLET_TYPE_SOLANA" {
				solanaAddress = w.Address
			}
		}
		Expect(solanaAddress).ToNot(BeEmpty())
		GinkgoWriter.Printf("Solana address: %s\n", solanaAddress)
	})

	It("creates a bridge transaction on zrchain", func() {
		var lastCount int
		var err error
		var newTx zentptypes.Bridge

		// Initial count of mints
		initialResp, err := env.Query.Mints(env.Ctx, "", "", zentptypes.BridgeStatus_BRIDGE_STATUS_UNSPECIFIED)
		Expect(err).ToNot(HaveOccurred())
		lastCount = len(initialResp.Mints)

		// Create a new mint tx
		_, err = env.Tx.NewMsgBridge(
			env.Ctx,
			"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
			100000000,
			"urock",
			"solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD",
			solanaAddress,
		)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() int {
			resp, err := env.Query.Mints(env.Ctx, "", "", zentptypes.BridgeStatus_BRIDGE_STATUS_UNSPECIFIED)
			Expect(err).ToNot(HaveOccurred())

			if len(resp.Mints) > lastCount {
				newTx = *resp.Mints[len(resp.Mints)-1]
			}

			return len(resp.Mints)
		}, "30s", "2s").Should(BeNumerically(">", lastCount))

		Expect(newTx.State).To(Equal(zentptypes.BridgeStatus_BRIDGE_STATUS_PENDING))
		GinkgoWriter.Printf("Bridge mint transaction created with ID %d\n", newTx.Id)
	})

	It("bridge tx gets completed", func() {
		Eventually(func() zentptypes.BridgeStatus {
			resp, err := env.Query.Mints(env.Ctx, "", "", zentptypes.BridgeStatus_BRIDGE_STATUS_UNSPECIFIED)
			Expect(err).ToNot(HaveOccurred())
			lastTx := *resp.Mints[len(resp.Mints)-1]

			return lastTx.State
		}, "150s", "5s").Should(Equal(zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED))
		GinkgoWriter.Printf("Bridge mint transaction completed \n")
	})
})
