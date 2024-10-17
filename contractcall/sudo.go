package contractcall

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SudoMessage struct {
	ParseInput ParseInput `json:"parse_input"`
}

type ParseInput struct {
	Input []byte `json:"input"`
}

type SudoResponse struct {
	ParseInputResponse parseInputResponse `json:"parse_input_response"`
}

type parseInputResponse struct {
	Value string `json:"value"`
}

func Sudo(wasm *wasmkeeper.PermissionedKeeper, ctx sdk.Context, address sdk.AccAddress, msg SudoMessage) (*SudoResponse, error) {
	bz, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	res, err := wasm.Sudo(ctx, address, bz)
	if err != nil {
		return nil, err
	}

	resp := SudoResponse{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
