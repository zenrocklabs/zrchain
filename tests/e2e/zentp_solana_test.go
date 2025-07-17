package e2e

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	solrock "github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock/generated/rock_spl_token"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var _ = Describe("ZenTP Solana:", func() {
	var env *TestEnv
	var solanaAccount string
	var burnTxHash solana.Signature

	BeforeEach(func() {
		env = setupTestEnv(GinkgoT())
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

	It("creates a mint transaction on zrchain", func() {
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
			solanaAccount,
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

	It("mint tx gets completed", func() {
		Eventually(func() zentptypes.BridgeStatus {
			resp, err := env.Query.Mints(env.Ctx, "", "", zentptypes.BridgeStatus_BRIDGE_STATUS_UNSPECIFIED)
			Expect(err).ToNot(HaveOccurred())
			lastTx := *resp.Mints[len(resp.Mints)-1]

			return lastTx.State
		}, "150s", "5s").Should(Equal(zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED))
		GinkgoWriter.Printf("Bridge mint transaction completed \n")
	})

	It("creates a burn", func() {
		var burnTx zentptypes.Bridge
		var client = rpc.New("http://localhost:8899")
		// Read private key from solana container
		out, err := env.Docker.Exec("solana", []string{
			"cat", "/root/.config/solana/id.json",
		})
		Expect(err).ToNot(HaveOccurred())

		signerWallet, err := LoadSolanaPrivateKeyFromJSON(out)
		Expect(err).ToNot(HaveOccurred())
		aliceBytes, err := sdk.GetFromBech32(ZENROCK_ALICE_ADDRESS, "zen")
		Expect(err).ToNot(HaveOccurred())
		aliceAddress := [25]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
		copy(aliceAddress[:], aliceBytes)
		programId, err := solana.PublicKeyFromBase58(SOLANA_ROCK_PROGRAM_ID)
		Expect(err).ToNot(HaveOccurred())
		signer, err := solana.PublicKeyFromBase58(solanaAccount)
		Expect(err).ToNot(HaveOccurred())
		feeWallet, err := solana.PublicKeyFromBase58(SOLANA_FEE_WALLET)
		Expect(err).ToNot(HaveOccurred())
		mintAddress, err := solana.PublicKeyFromBase58(SOLANA_TOKEN_ADDRESS)
		Expect(err).ToNot(HaveOccurred())

		latest, err := client.GetLatestBlockhash(env.Ctx, rpc.CommitmentProcessed)
		Expect(err).ToNot(HaveOccurred())

		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				solrock.Unwrap(
					programId,
					rock_spl_token.UnwrapArgs{
						Value:    uint64(10000000),
						DestAddr: aliceAddress,
					},
					signer,
					mintAddress,
					feeWallet,
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
		Eventually(func() int {
			resp, err := env.Query.Burns(env.Ctx, "", burnTxHash.String())
			Expect(err).ToNot(HaveOccurred())
			if len(resp.Burns) > 0 {
				burnTx = *resp.Burns[0]
			}

			return len(resp.Burns)
		}, "150s", "5s").Should(BeNumerically("==", 1))
		GinkgoWriter.Printf("Burn tx created: %v\n", burnTx)
	})

	It("burn tx gets completed", func() {
		var burnTx zentptypes.Bridge

		Eventually(func() zentptypes.BridgeStatus {
			resp, err := env.Query.Burns(env.Ctx, "", burnTxHash.String())

			Expect(err).ToNot(HaveOccurred())
			if len(resp.Burns) > 0 {
				burnTx = *resp.Burns[0]
			}

			return burnTx.State
		}, "15s", "5s").Should(Equal(zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED))
		GinkgoWriter.Printf("Burn transaction completed: %v\n", burnTx)
	})
})
