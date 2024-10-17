package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TxSelectorWithVoteExt struct {
	BytesUsed   uint64
	GasUsed     uint64
	SelectedTxs [][]byte
}

func newTxSelectorWithVoteExt() *TxSelectorWithVoteExt {
	return &TxSelectorWithVoteExt{}
}

func (ts *TxSelectorWithVoteExt) insertVoteExtension(txBytes []byte) {
	ts.BytesUsed += uint64(len(txBytes))
	ts.SelectedTxs = append(ts.SelectedTxs, txBytes)
}

func (ts *TxSelectorWithVoteExt) insertTransaction(maxBytesPerTx, maxGasPerBlock uint64, memTx sdk.Tx, txBytes []byte) {
	txSize := uint64(len(txBytes))
	txGasLimit := txGasLimit(memTx)

	if ts.canAddTransaction(txSize, maxBytesPerTx, txGasLimit, maxGasPerBlock) {
		ts.addTransaction(txSize, txGasLimit, txBytes)
	}
}

func (ts *TxSelectorWithVoteExt) canAddTransaction(txSize, maxBytesPerTx, txGasLimit, maxGasPerBlock uint64) bool {
	if (txSize + ts.BytesUsed) > maxBytesPerTx {
		return false
	}

	if maxGasPerBlock > 0 && (txGasLimit+ts.GasUsed) > maxGasPerBlock {
		return false
	}

	return true
}

func (ts *TxSelectorWithVoteExt) addTransaction(txSize, txGasLimit uint64, txBytes []byte) {
	ts.GasUsed += txGasLimit
	ts.BytesUsed += txSize
	ts.SelectedTxs = append(ts.SelectedTxs, txBytes)
}

func (ts *TxSelectorWithVoteExt) capacityReached(maxBytesPerTx, maxGasPerBlock uint64) bool {
	return ts.BytesUsed >= maxBytesPerTx || (maxGasPerBlock > 0 && ts.GasUsed >= maxGasPerBlock)
}

func maxGasPerBlock(ctx sdk.Context) uint64 {
	if b := ctx.ConsensusParams().Block; b != nil {
		return uint64(b.MaxGas)
	}

	return 0
}

func txGasLimit(memTx sdk.Tx) uint64 {
	if tx, ok := memTx.(baseapp.GasTx); ok {
		return tx.GetGas()
	}

	return 0
}

func isGasWithinLimit(tx sdk.Tx, maxGasPerBlock uint64, gasUsed *uint64) bool {
	if maxGasPerBlock > 0 {
		if gasTx, ok := tx.(baseapp.GasTx); ok {
			*gasUsed += gasTx.GetGas()
		}

		if *gasUsed > maxGasPerBlock {
			return false
		}
	}

	return true
}
