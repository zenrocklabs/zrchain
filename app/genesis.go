package app

import (
	"encoding/json"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GenesisState of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

var _ authtypes.GenesisAccount = (*GenesisAccount)(nil)

// GenesisAccount defines a type that implements the GenesisAccount interface
// to be used for simulation accounts in the genesis state.
type GenesisAccount struct {
	*authtypes.BaseAccount

	// vesting account fields
	OriginalVesting  sdk.Coins `json:"original_vesting" yaml:"original_vesting"`   // total vesting coins upon initialization
	DelegatedFree    sdk.Coins `json:"delegated_free" yaml:"delegated_free"`       // delegated vested coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting" yaml:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        int64     `json:"start_time" yaml:"start_time"`               // vesting start time (UNIX Epoch time)
	EndTime          int64     `json:"end_time" yaml:"end_time"`                   // vesting end time (UNIX Epoch time)

	// module account fields
	ModuleName        string   `json:"module_name" yaml:"module_name"`               // name of the module account
	ModulePermissions []string `json:"module_permissions" yaml:"module_permissions"` // permissions of module account
}

// Validate checks for errors on the vesting and module account parameters
func (sga GenesisAccount) Validate() error {
	if !sga.OriginalVesting.IsZero() {
		if sga.StartTime >= sga.EndTime {
			return errors.New("vesting start-time cannot be before end-time")
		}
	}

	if sga.ModuleName != "" {
		ma := authtypes.ModuleAccount{
			BaseAccount: sga.BaseAccount, Name: sga.ModuleName, Permissions: sga.ModulePermissions,
		}
		if err := ma.Validate(); err != nil {
			return err
		}
	}

	return sga.BaseAccount.Validate()
}
