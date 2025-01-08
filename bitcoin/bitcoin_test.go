package bitcoin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ECDSA_convert_to_Bitcoin_SigdeECDSA(t *testing.T) {
	ecdsaSig1 := "e4899cc47ab30e84aed02fe91cac71fa4e89110e81158b612c3e21b02ca37ea656a5dcfd089c9b8e65c386c79433f4f2e66f0f7b4804c5dd6a6c002ab29c9d9c00"
	bitcoinSig1 := "3045022100e4899cc47ab30e84aed02fe91cac71fa4e89110e81158b612c3e21b02ca37ea6022056a5dcfd089c9b8e65c386c79433f4f2e66f0f7b4804c5dd6a6c002ab29c9d9c"
	calculatedBitcoinSig1, err := ConvertECDSASigtoBitcoinSig(ecdsaSig1)
	require.NoError(t, err)
	require.Equal(t, calculatedBitcoinSig1, bitcoinSig1)

	ecdsaSig2 := "d129c4e978e0633304246bec0acfb760991ac68e2169c1eea6edcb8ece708ee864d3fe6d5fca57e5d3f86046413669a66bc2845e9abe9821e062ea13cb3de6bc"
	bitcoinSig2 := "3045022100d129c4e978e0633304246bec0acfb760991ac68e2169c1eea6edcb8ece708ee8022064d3fe6d5fca57e5d3f86046413669a66bc2845e9abe9821e062ea13cb3de6bc"
	calculatedBitcoinSig2, err := ConvertECDSASigtoBitcoinSig(ecdsaSig2)
	require.NoError(t, err)
	require.Equal(t, calculatedBitcoinSig2, bitcoinSig2)
}
