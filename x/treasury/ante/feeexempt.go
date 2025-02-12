package ante

import (
	"context"

	"github.com/cockroachdb/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FeeExempt is an interface that defines the methods for a fee exempt.
type FeeExempt interface {
	IsAllowed(ctx context.Context, typeURL string) (bool, error)
}

// FeeExemptDecorator is an AnteDecorator that checks if the transaction type is allowed to skip transaction fees
type FeeExemptDecorator struct {
	feeExemptKeeper FeeExempt
}

func NewFeeExemptDecorator(fk FeeExempt) FeeExemptDecorator {
	return FeeExemptDecorator{
		feeExemptKeeper: fk,
	}
}

func (fed FeeExemptDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// loop through all the messages and check if the message type is allowed
	for _, msg := range tx.GetMsgs() {
		isAllowed, err := fed.feeExemptKeeper.IsAllowed(ctx, sdk.MsgTypeURL(msg))
		if err != nil {
			return ctx, err
		}

		if !isAllowed {
			return ctx, errors.New("tx type not allowed")
		}
	}

	return next(ctx, tx, simulate)
}
