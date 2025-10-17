package bitcoin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeOutputs_ZcashV5(t *testing.T) {
	rawTx := "050000800a27a726f04dec4d00000000f4493700010000000000000000000000000000000000000000000000000000000000000000ffffffff0603f449370106ffffffff0240597307000000001976a91402b7b5b3afa00d56eb1a2e76b8db889c935c3f1088ac20bcbe000000000017a9147a86d6c7eb12ce0aa309d7391a6f338eba3c242b87000000"

	outputs, err := DecodeOutputs(rawTx, "zcashtestnet")
	require.NoError(t, err)
	require.Len(t, outputs, 2)

	require.Equal(t, uint(0), outputs[0].OutputIndex)
	require.Equal(t, uint64(125000000), outputs[0].Amount)
	require.Equal(t, "tm9xikw8UhJXxUTsWphMfGvGPcjL1PA3CHz", outputs[0].Address)

	require.Equal(t, uint(1), outputs[1].OutputIndex)
	require.Equal(t, uint64(12500000), outputs[1].Amount)
	require.Equal(t, "t2HifwjUj9uyxr9bknR8LFuQbc98c3vkXtu", outputs[1].Address)
}

func TestCalculateTXID_ZcashV5(t *testing.T) {
	rawTx := "050000800a27a726f04dec4d00000000f4493700010000000000000000000000000000000000000000000000000000000000000000ffffffff0603f449370106ffffffff0240597307000000001976a91402b7b5b3afa00d56eb1a2e76b8db889c935c3f1088ac20bcbe000000000017a9147a86d6c7eb12ce0aa309d7391a6f338eba3c242b87000000"
	expectedTxID := "dc7cfca537414cbf45b568f7afad4ccdd402b8455fbdf9b5fed26e96b3fe4a83"

	txID, err := CalculateTXID(rawTx, "zcashtestnet")
	require.NoError(t, err)

	require.Equal(t, expectedTxID, ReverseHex(txID.String()))
}
