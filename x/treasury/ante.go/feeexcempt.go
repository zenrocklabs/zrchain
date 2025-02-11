package ante

import (
	"context"

	"github.com/cockroachdb/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FeeExcempt is an interface that defines the methods for a fee excempt.
type FeeExcempt interface {
	IsAllowed(ctx context.Context, typeURL string) (bool, error)
}

// FeeExcemptDecorator is an AnteDecorator that checks if the transaction type is allowed to skip transaction fees
type FeeExcemptDecorator struct {
	feeExcemptKeeper FeeExcempt
}

func NewFeeExcemptDecorator(fk FeeExcempt) FeeExcemptDecorator {
	return FeeExcemptDecorator{
		feeExcemptKeeper: fk,
	}
}

func (fed FeeExcemptDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// loop through all the messages and check if the message type is allowed
	for _, msg := range tx.GetMsgs() {
		isAllowed, err := fed.feeExcemptKeeper.IsAllowed(ctx, sdk.MsgTypeURL(msg))
		if err != nil {
			return ctx, err
		}

		if !isAllowed {
			return ctx, errors.New("tx type not allowed")
		}
	}

	return next(ctx, tx, simulate)
}
