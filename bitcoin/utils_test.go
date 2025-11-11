package bitcoin

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeTX_ZcashSaplingTransparent(t *testing.T) {
	raw := mustDecodeBase64(t, "BAAAAIUgL4kF83d0mzxzgASNFmR47PF4zGUQvGU2pRvsSv1a7shF4BsAAAAAAP////91WorSLUlMFDgW5eaL62c6/ER2/rwuk48K6OEeqgT2fgAAAAAA/////3jxbqz4HqYRvJ7sb0RrB6FXvHeu2echTIY3x/pHUy/WAAAAAAD/////1QAEnZG1hifi6VUNqjOatNnQUXWePf/2vXI8hkIao7wAAAAAAP////9jXhLiUycT5HTpQZadyCJuYvyUmbc3gEyoy4M/DAe2cAAAAAAA/////wQVecMIAAAAABl2qRQGvCVYPtJPB2l/cl0CwCmBAehlroisMBsPAAAAAAAZdqkUdEZtXvZWSok9W9JHwrYFhoWF3RaIrDAbDwAAAAAAGXapFHRGbV72VkqJPVvSR8K2BYaFhd0WiKwwGw8AAAAAABl2qRR0Rm1e9lZKiT1b0kfCtgWGhYXdFoisAAAAAAAAAAAAAAAAAAAAAAAA")

	tx, err := DecodeTX(raw, "zcashtestnet")
	require.NoError(t, err)

	require.Len(t, tx.TxIn, 5, "expected 5 transparent inputs")
	require.Len(t, tx.TxOut, 4, "expected 4 transparent outputs")

	require.EqualValues(t, 147028245, tx.TxOut[0].Value)
	require.EqualValues(t, 990000, tx.TxOut[1].Value)
	require.EqualValues(t, 990000, tx.TxOut[2].Value)
	require.EqualValues(t, 990000, tx.TxOut[3].Value)

	for _, out := range tx.TxOut {
		require.NotZero(t, len(out.PkScript))
	}
}

func mustDecodeBase64(t *testing.T, in string) []byte {
	t.Helper()
	out, err := base64.StdEncoding.DecodeString(in)
	require.NoError(t, err)
	return out
}
