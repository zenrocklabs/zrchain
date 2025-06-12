package types

import (
	"errors"
	"slices"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBridge{}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgBridge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if m.Amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

func IsValidChain(ctx sdk.Context, chain string) bool {
	chainID := ctx.ChainID()
	if chainID == "" {
		chainID = "zenrock"
	}

	chainMap := map[string][]string{
		"zenrock": {"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1"},
		"amber":   {"solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1"},
		"gardia":  {"solana:4uhcVJyU9pJkvQyS88uRDiswHXSCkY3z", "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1"},
		"diamond": {"solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp"},
	}

	var allowedChains []string
	for prefix, chains := range chainMap {
		if strings.HasPrefix(chainID, prefix) {
			allowedChains = chains
			break
		}
	}

	if allowedChains == nil {
		return false
	}

	return slices.Contains(allowedChains, chain)
}
