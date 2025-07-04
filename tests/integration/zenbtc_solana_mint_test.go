package integration_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	zentype "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

var _ = Describe("ZenBTC Solana mint:", func() {
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
			"HhnPD1MTBBZ6h15pdJdRySXi7QW4Dvi65UNYRZTCFAWN",
			types.WalletType_WALLET_TYPE_SOLANA,
			"solana:HK8b7Skns2TX3FvXQxm2mPQbY2nVY8GD",
			"HhnPD1MTBBZ6h15pdJdRySXi7QW4Dvi65UNYRZTCFAWN",
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

	It("gets fulfilled within 15 seconds", func() {
		Eventually(func() string {
			req, err := env.Query.GetKeyRequest(env.Ctx, requestID)
			Expect(err).ToNot(HaveOccurred())
			return req.Status
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
		r, err := env.Docker.Exec("bitcoin", []string{"/app/mine.sh", bitcoinAddress})
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("response docker cmd: %s\n", r)
	})

	It("creates a mint transaction on zrchain", func() {
		var lastCount int
		var err error
		// Initial count of pending mint transactions
		initialResp, err := env.Query.PendingMintTransactions(env.Ctx, 1)
		Expect(err).ToNot(HaveOccurred())
		lastCount = len(initialResp.PendingMintTransactions)

		var newTx zentype.PendingMintTransaction

		Eventually(func() int {
			resp, err := env.Query.PendingMintTransactions(env.Ctx, 1)
			Expect(err).ToNot(HaveOccurred())

			if len(resp.PendingMintTransactions) > lastCount {
				newTx = *resp.PendingMintTransactions[len(resp.PendingMintTransactions)-1]
			}

			return len(resp.PendingMintTransactions)
		}, "30s", "2s").Should(BeNumerically(">", lastCount))

		Expect(newTx.Status).To(Equal(zentype.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED))
		GinkgoWriter.Printf("Mint transaction created with ID %d\n", newTx.Id)
	})

	It("mint gets staked", func() {
		Eventually(func() zentype.MintTransactionStatus {
			resp, err := env.Query.PendingMintTransactions(env.Ctx, 1)
			Expect(err).ToNot(HaveOccurred())
			lastTx := *resp.PendingMintTransactions[len(resp.PendingMintTransactions)-1]

			return lastTx.Status
		}, "90s", "5s").Should(Equal(zentype.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED))
		GinkgoWriter.Printf("Mint transaction moved to staked \n")
	})

	It("mint gets minted", func() {
		Eventually(func() zentype.MintTransactionStatus {
			resp, err := env.Query.PendingMintTransactions(env.Ctx, 1)
			Expect(err).ToNot(HaveOccurred())
			lastTx := *resp.PendingMintTransactions[len(resp.PendingMintTransactions)-1]

			return lastTx.Status
		}, "150s", "5s").Should(Equal(zentype.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED))
		GinkgoWriter.Printf("Mint transaction moved to minted \n")
	})
})
