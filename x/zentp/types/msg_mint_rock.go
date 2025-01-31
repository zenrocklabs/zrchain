package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgMintRock{}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgMintRock) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if m.Amount == 0 {
		return errors.New("amount must be greater than zero")
	}

	return nil
}
