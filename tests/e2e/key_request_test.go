package e2e

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

var _ = Describe("MPC Key Requests:", func() {
	var env *TestEnv
	var requestID uint64

	testKeyType := func(keyType string) {
		BeforeEach(func() {
			env = setupTestEnv(GinkgoT())
		})

		It("creates a new request", func() {
			hash, err := env.Tx.NewKeyRequest(
				env.Ctx,
				"workspace1mphgzyhncnzyggfxmv4nmh",
				"keyring1k6vc6vhp6e6l3rxalue9v4ux",
				keyType,
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
			GinkgoWriter.Printf("ECDSA Key Request created: %d\n", requestID)
		})

		It("fetches the request within 5 seconds", func() {
			Eventually(func() (uint64, error) {
				req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
				if err != nil {
					return 0, err
				}
				return req.Id, nil
			}, "5s", "1s").Should(Equal(requestID))
			GinkgoWriter.Printf("ECDSA Key Request fetched: %d\n", requestID)
		})

		It("gets fulfilled within 15 seconds", func() {
			Eventually(func() (string, error) {
				req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
				if err != nil {
					return "", err
				}
				return req.Status, nil
			}, "15s", "1s").Should(Equal(types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED.String()))
			GinkgoWriter.Printf("ECDSA Key Request fulfilled: %d\n", requestID)
		})
	}

	Describe("ECDSA", func() {
		testKeyType("ecdsa")
	})

	Describe("EDDSA", func() {
		testKeyType("eddsa")
	})

	Describe("BITCOIN", func() {
		testKeyType("bitcoin")
	})
})
