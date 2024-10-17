package shared

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type DynamicFeeTxWithoutSignature struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int
	Data       []byte
	AccessList types.AccessList
}

type AccessListTxWithoutSignature struct {
	ChainID    *big.Int         // destination chain ID
	Nonce      uint64           // nonce of sender account
	GasPrice   *big.Int         // wei per gas
	Gas        uint64           // gas limit
	To         *common.Address  `rlp:"nil"` // nil means contract creation
	Value      *big.Int         // wei amount
	Data       []byte           // contract invocation input data
	AccessList types.AccessList // EIP-2930 access list
}

// The following code doesn't work for unsigned transactions:
//
//	var tx types.Transaction
//	tx.UnmarshalBinary(b)
//
// This function is a workaround taken from https://github.com/ethereum/go-ethereum/issues/26236.
func DecodeUnsignedPayload(msg []byte) (types.TxData, error) {
	if len(msg) > 0 && msg[0] > 0x7f {
		// Legacy transaction
		var res types.LegacyTx
		err := rlp.DecodeBytes(msg, &res)
		return &res, err
	}
	switch msg[0] {
	case types.AccessListTxType:
		var res AccessListTxWithoutSignature
		err := rlp.DecodeBytes(msg[1:], &res)
		return &types.AccessListTx{
			ChainID:    res.ChainID,
			Nonce:      res.Nonce,
			GasPrice:   res.GasPrice,
			Gas:        res.Gas,
			To:         res.To,
			Value:      res.Value,
			Data:       res.Data,
			AccessList: res.AccessList,
		}, err
	case types.DynamicFeeTxType:
		var res DynamicFeeTxWithoutSignature
		err := rlp.DecodeBytes(msg[1:], &res)
		return &types.DynamicFeeTx{
			ChainID:    res.ChainID,
			Nonce:      res.Nonce,
			GasTipCap:  res.GasTipCap,
			GasFeeCap:  res.GasFeeCap,
			Gas:        res.Gas,
			To:         res.To,
			Value:      res.Value,
			Data:       res.Data,
			AccessList: res.AccessList,
		}, err
	default:
		return nil, fmt.Errorf("unsupported transaction type: %v", msg[0])
	}
}
