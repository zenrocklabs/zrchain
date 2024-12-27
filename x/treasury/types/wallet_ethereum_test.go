package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func Test_EthereumWallet_Address(t *testing.T) {
	wallet := ethereumWallet(t)
	require.Equal(t, "0xdD1d3fF09C5EdfF1bE7d466cA614cB1cF3f78738", wallet.Address())
}

func ethereumWallet(t *testing.T) *EthereumWallet {
	t.Helper()
	k := &Key{
		Id:            0,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		Type:          KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     hexutil.MustDecode("0x025cd45a6614df5348692ea4d0f7c16255b75a6b6f67bea5013621fe84af8031f0"),
	}

	wallet, err := NewEthereumWallet(k)
	require.NoError(t, err)
	return wallet
}

func Test_ParseEthereumTransaction(t *testing.T) {
	tests := []struct {
		name         string
		b            []byte
		chainID      *big.Int
		wantTo       string
		wantAmount   *big.Int
		wantContract string
		wantErr      bool
	}{
		// The following two txs are LegacyTxs that are no longer supported. Leaving the test cases here in case we want to support them again in the future.
		{
			name:         "ETH transfer",
			b:            hexutil.MustDecode("0xeb80843b9aca0082520894ea223ca8968ca59e0bc79ba331c2f6f636a3fb82880de0b6b3a764000080808080"),
			wantTo:       "0xeA223Ca8968Ca59e0Bc79Ba331c2F6f636A3fB82",
			wantAmount:   big.NewInt(1000000000000000000),
			wantContract: "",
			wantErr:      false,
		},
		{
			name:         "ERC-20 transfer",
			b:            hexutil.MustDecode("0xf8aa80850b68a0aa0083010d6b94a0b86991c6218b36c1d19d4a2e9eb0ce3606eb4880b844a9059cbb00000000000000000000000048c04ed5691981c42154c6167398f95e8f38a7ff00000000000000000000000000000000000000000000000000000000017d784026a01ad4a933da06f76b08a20784661b2ccc55a8f25492bbc6b66d4b84ef97e1db47a01ab2ff0cd0fb01e2990dd4196412baf173484d91c7f836727d554cdf1cd70c64"),
			wantTo:       "0x48c04ed5691981C42154C6167398f95e8f38a7fF",
			wantAmount:   big.NewInt(25000000),
			wantContract: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			wantErr:      false,
		},
		{
			name:         "ERC-20 transferFrom",
			b:            hexutil.MustDecode("0xf8ca80850b68a0aa0083010d6b94a0b86991c6218b36c1d19d4a2e9eb0ce3606eb4880b86423b872dd00000000000000000000000048c04ed5691981c42154c6167398f95e8f38a7ff00000000000000000000000048c04ed5691981c42154c6167398f95e8f38a7aa00000000000000000000000000000000000000000000000000000000017d784026a01ad4a933da06f76b08a20784661b2ccc55a8f25492bbc6b66d4b84ef97e1db47a01ab2ff0cd0fb01e2990dd4196412baf173484d91c7f836727d554cdf1cd70c64"),
			wantTo:       "0x48C04Ed5691981c42154c6167398F95e8f38A7aA",
			wantAmount:   big.NewInt(25000000),
			wantContract: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			wantErr:      false,
		},
		// {
		// 	name:         "DynamicFeeTx",
		// 	b:            hexutil.MustDecode("0x02f902b583aa36a7040385042d03bb3d8302a43b943fc91a3afd70395cd496c647d5a6cc9d4b2b7fad872386f26fc10000b902843593564c000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000006595a2b000000000000000000000000000000000000000000000000000000000000000020b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000002386f26fc10000000000000000000000000000000000000000000000000000001925fd93f197ab00000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002bfff9976782d46cc05630d1f6ebab18b2324d6b14000bb81f9840a85d5af5bf1d1762f925bdaddc4201f984000000000000000000000000000000000000000000c0"),
		// 	chainID:      big.NewInt(11155111),
		// 	wantTo:       "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
		// 	wantAmount:   big.NewInt(10000000000000000),
		// 	wantContract: "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
		// 	wantErr:      false,
		// },
		// {
		// 	name:         "Uniswap custom tx",
		// 	b:            hexutil.MustDecode("0x02f86c83aa36a70203850703deeb8b82b78b941f9840a85d5af5bf1d1762f925bdaddc4201f98480b844095ea7b3000000000000000000000000000000000022d473030f116ddee9f6b43ac78ba3ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0"),
		// 	chainID:      big.NewInt(11155111),
		// 	wantTo:       "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984",
		// 	wantAmount:   big.NewInt(0),
		// 	wantContract: "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984",
		// 	wantErr:      false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := ParseEthereumTransaction(tt.b, tt.chainID)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantTo, tx.To.Hex())
			require.Equal(t, tt.wantAmount, tx.Amount)
			if len(tt.wantContract) == 0 {
				require.Nil(t, tx.Contract)
			} else {
				require.Equal(t, tt.wantContract, tx.Contract.Hex())
			}
		})
	}
}
