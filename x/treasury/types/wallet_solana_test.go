package types

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/stretchr/testify/require"
)

// TestSolanaWalletAddress tests Solana address derivation by building a Solana wallet and checking the solanaWallet.Address
// Solana addresses are a simple base58 representation of a public key. The test seed is used to deterministically generate
// an Ed25519 public key.
func TestSolanaWalletAddress(t *testing.T) {
	testCases := []struct {
		desc            string
		seed            string
		expectedAddress string
	}{
		{
			desc:            "wallet address test 0",
			seed:            "example seed",
			expectedAddress: "4y4Hs9PQNWMnG8WJAQMQDh6crkqZngbNKY97BGbX29i4",
		},
		{
			desc:            "wallet address test 1",
			seed:            "example seed 2",
			expectedAddress: "AbeXScYJEq9Ece3fUsFbZsH29qJ45qEVbBh39uLkD6bK",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			wallet := solanaWallet(t, tC.seed)
			require.Equal(t, tC.expectedAddress, wallet.Address())
		})
	}
}

// TestSolanaParseTx
func TestSolanaParseTx(t *testing.T) {
	testCases := []struct {
		name        string
		to          string
		amount      int64
		txHexstring string
	}{
		{
			name:        "Test Case 0",
			to:          "FYJ5gRsqYAwWcxL7LsxsZuunWTysMF7if9Sp5nPADQYF",
			amount:      10_042,
			txHexstring: "0x01000103e01763c7d59132d8423e17a8e285d86d0ae555bae8ce098ae14b5fbc709bb238d807ef6c24d679daf65e8097a455a26dae192524f343b27d1c998a9407c05cc00000000000000000000000000000000000000000000000000000000000000000f0857f2a581eea6588796d8d39e204dd354ea2d8bad46e9fc7591bc5f00c534d01020200010c020000003a27000000000000",
		},
		{
			name:        "Test Case 1",
			to:          "HpzusjfWgokpwuz6D8GhyCELJM83e6FC7KeAvzbXtF6R",
			amount:      1_000_013,
			txHexstring: "0x01000103ec683a77d4e09795bd0ee0828a553fc7acdd3c8d537f95945579b6780c429d37fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a0000000000000000000000000000000000000000000000000000000000000000a1e2e5e4e30c35c282579e0338eeb9739f34f06d3db52e024e8767fcdb4dd47501020200010c020000004d420f0000000000",
		},
		{
			name:        "Test Case 2",
			to:          "HpzusjfWgokpwuz6D8GhyCELJM83e6FC7KeAvzbXtF6R",
			amount:      1_000_000_042,
			txHexstring: "0x01000103ec683a77d4e09795bd0ee0828a553fc7acdd3c8d537f95945579b6780c429d37fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a0000000000000000000000000000000000000000000000000000000000000000bd349395c46ca6be0be09c0989649ecb2def6488a15bd6bdf44d36f9ca5aabaa01020200010c020000002aca9a3b00000000",
		},
		{
			name:        "Test Case 3",
			to:          "5Wvu3L3vVkQi2RPA12sTbzFTPgzgim5kQc3iRHFVw6zZ",
			amount:      10_042,
			txHexstring: "0x01000103ec683a77d4e09795bd0ee0828a553fc7acdd3c8d537f95945579b6780c429d3743198871268c38c8e04a2aaa78618bf49d085c39c2631e37d3321d10c8b8c0ae0000000000000000000000000000000000000000000000000000000000000000b25330dd9a65d49b61bb4993123e10f813b0f8a2740d7a2043a246d866ade95901020200010c020000003a27000000000000",
		},
		{
			name:        "Debug JS app",
			to:          "HpzusjfWgokpwuz6D8GhyCELJM83e6FC7KeAvzbXtF6R",
			amount:      10_042,
			txHexstring: "0x01000103ec683a77d4e09795bd0ee0828a553fc7acdd3c8d537f95945579b6780c429d37fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a0000000000000000000000000000000000000000000000000000000000000000530441528372542b2f23866d8cc58a6b8d764aee09f9c26fb263513e7fed3d8d01020200010c020000003a27000000000000",
		},
		{
			name:        "Debug JS app - FINAL",
			to:          "HpzusjfWgokpwuz6D8GhyCELJM83e6FC7KeAvzbXtF6R",
			amount:      10_042,
			txHexstring: "0x01000103ec683a77d4e09795bd0ee0828a553fc7acdd3c8d537f95945579b6780c429d37fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a00000000000000000000000000000000000000000000000000000000000000006bda0ef1adffccbcf2f2425b1af761b9db4122dce6a7889e4ddb2520baff79e801020200010c020000003a27000000000000",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wallet := solanaWallet(t, "example seed")
			txBytes := hexutil.MustDecode(tc.txHexstring)
			transfer, err := wallet.ParseTx(txBytes, nil)
			expectedDataForSigning := []byte(hex.EncodeToString(txBytes))
			require.NoError(t, err)
			require.Equal(t, []byte(tc.to), transfer.To, "to address mismatch")
			require.Equal(t, tc.amount, transfer.Amount.Int64(), "amount mismatch")
			require.Equal(t, hex.EncodeToString(expectedDataForSigning), hex.EncodeToString(transfer.DataForSigning), "data for signing mismatch")
		})
	}
}

func solanaWallet(t *testing.T, seed string) *SolanaWallet {
	t.Helper()
	hashedSeed := sha256.Sum256([]byte(seed))
	publicKey, _, err := ed25519.GenerateKey(bytes.NewReader(hashedSeed[:]))
	if err != nil {
		log.Fatal("Error generating ed25519 key pair:", err)
	}

	k := &Key{
		Id:            0,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_EDDSA_ED25519,
		PublicKey:     []byte(publicKey),
	}

	wallet, err := NewSolanaWallet(k)
	if err != nil {
		t.Fatal(err)
	}
	return wallet
}

// TestGetTransferFromInstruction
// TODO complete this test
func TestGetTransferFromInstruction(t *testing.T) {
	t.Skip()
	type args struct {
		msg solana.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *system.Transfer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTransferFromInstruction(tt.args.msg)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}
