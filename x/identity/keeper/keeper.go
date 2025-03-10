package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	Schema         collections.Schema
	ParamStore     collections.Item[types.Params]
	KeyringStore   collections.Map[string, types.Keyring]
	KeyringCount   collections.Item[uint64]
	WorkspaceStore collections.Map[string, types.Workspace]
	WorkspaceCount collections.Item[uint64]

	ibcKeeperFn        func() *ibckeeper.Keeper
	capabilityScopedFn func(string) capabilitykeeper.ScopedKeeper
	scopedKeeper       exported.ScopedKeeper
	bankKeeper         types.BankKeeper
	policyKeeper       types.PolicyKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	bankKeeper types.BankKeeper,
	policyKeeper types.PolicyKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
		bankKeeper:   bankKeeper,
		policyKeeper: policyKeeper,

		ParamStore:     collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Params](cdc)),
		KeyringStore:   collections.NewMap(sb, types.KeyringsKey, types.KeyringsIndex, collections.StringKey, codec.CollValue[types.Keyring](cdc)),
		KeyringCount:   collections.NewItem(sb, types.KeyringCountKey, types.KeyringCountIndex, collections.Uint64Value),
		WorkspaceStore: collections.NewMap(sb, types.WorkspacesKey, types.WorkspacesIndex, collections.StringKey, codec.CollValue[types.Workspace](cdc)),
		WorkspaceCount: collections.NewItem(sb, types.WorkspaceCountKey, types.WorkspaceCountIndex, collections.Uint64Value),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ----------------------------------------------------------------------------
// IBC Keeper Logic
// ----------------------------------------------------------------------------

// ChanCloseInit defines a wrapper function for the channel Keeper's function.
func (k *Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := host.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.ScopedKeeper().GetCapability(ctx, capName)
	if !ok {
		return errorsmod.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.ibcKeeperFn().ChannelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

// ShouldBound checks if the IBC app module can be bound to the desired port
func (k *Keeper) ShouldBound(ctx sdk.Context, portID string) bool {
	scopedKeeper := k.ScopedKeeper()
	if scopedKeeper == nil {
		return false
	}
	_, ok := scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return !ok
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to module's InitGenesis function
func (k *Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.ibcKeeperFn().PortKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the IBC app module. Used in ExportGenesis
func (k *Keeper) GetPort(ctx sdk.Context) string {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the IBC app module. Used in InitGenesis
func (k *Keeper) SetPort(ctx sdk.Context, portID string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(types.PortKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k *Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.ScopedKeeper().AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the IBC app module to claim a capability that core IBC
// passes to it
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.ScopedKeeper().ClaimCapability(ctx, cap, name)
}

// ScopedKeeper returns the ScopedKeeper
func (k *Keeper) ScopedKeeper() exported.ScopedKeeper {
	if k.scopedKeeper == nil && k.capabilityScopedFn != nil {
		k.scopedKeeper = k.capabilityScopedFn(types.ModuleName)
	}
	return k.scopedKeeper
}

func (k Keeper) GetZrSignWorkspace(goCtx context.Context, ethAddress string, walletType uint64) (string, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	walletTypeStr := strconv.FormatUint(walletType, 10)
	ws, _, err := query.CollectionFilteredPaginate(
		goCtx,
		k.WorkspaceStore,
		nil,
		func(key string, value types.Workspace) (bool, error) {
			return ethAddress == value.Creator, nil
		},
		func(key string, value types.Workspace) (*types.Workspace, error) {
			return &value, nil
		},
	)
	if err != nil {
		return "", err
	}

	var (
		parentWsID string
		childID    string
	)
	if len(ws) == 0 {
		resp, err := k.CreateWorkspace(ctx, &types.Workspace{
			Creator:       ethAddress,
			AdminPolicyId: 0,
			SignPolicyId:  0,
			Owners:        []string{ethAddress},
		})
		if err != nil {
			return "", errors.Wrap(err, "create workspace")
		}
		parentWsID = resp
	} else {
		parentWsID = ws[0].Address // default to first if there aren't any with children
		for _, w := range ws {
			if len(w.ChildWorkspaces) > 0 {
				parentWsID = w.Address
				break
			}
		}
	}

	parentW, err := k.WorkspaceStore.Get(ctx, parentWsID)
	if err != nil {
		return "", errors.Wrapf(types.ErrNotFound, "parent workspace %s not found", parentWsID)
	}

	for _, wID := range parentW.ChildWorkspaces {
		child, err := k.WorkspaceStore.Get(ctx, wID)
		if err != nil {
			return "", errors.Wrapf(types.ErrNotFound, "child workspace %s not found", parentWsID)
		}

		if child.GetAlias() == walletTypeStr {
			childID = child.Address
			break
		}
	}

	if len(childID) == 0 {
		resp, err := k.storeChildWorkspace(ctx, &parentW, &types.Workspace{
			Creator:       ethAddress,
			Owners:        []string{ethAddress},
			AdminPolicyId: parentW.AdminPolicyId,
			SignPolicyId:  parentW.SignPolicyId,
			Alias:         walletTypeStr,
		})
		if err != nil {
			return "", errors.Wrap(err, "create new child workspace")
		}
		childID = resp.Address
	}

	return childID, nil
}

func (k Keeper) GetZrSignWorkspaces(goCtx context.Context, ethAddress, walletType string) (map[string]string, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ws, _, err := query.CollectionFilteredPaginate(
		goCtx,
		k.WorkspaceStore,
		nil,
		func(key string, value types.Workspace) (bool, error) {
			return ethAddress == value.Creator && len(value.ChildWorkspaces) > 0, nil
		},
		func(key string, value types.Workspace) (*types.Workspace, error) {
			return &value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if len(ws) == 0 {
		return nil, fmt.Errorf("no workspaces")
	}

	parentWsID := ws[0].Address // default to first if there aren't any with children

	parentW, err := k.WorkspaceStore.Get(ctx, parentWsID)
	if err != nil {
		return nil, errors.Wrapf(types.ErrNotFound, "parent workspace %s not found", parentWsID)
	}

	workspaceList := map[string]string{}
	for _, wID := range parentW.ChildWorkspaces {
		child, err := k.WorkspaceStore.Get(ctx, wID)
		if err != nil {
			return nil, errors.Wrapf(types.ErrNotFound, "child workspace %s not found", parentWsID)
		}

		if walletType != "" && child.GetAlias() != walletType {
			continue
		}
		workspaceList[child.Alias] = child.Address
	}

	return workspaceList, nil
}

func (k Keeper) GetKeyring(ctx sdk.Context, id string) (*types.Keyring, error) {
	keyring, err := k.KeyringStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &keyring, nil
}

func (k Keeper) GetWorkspace(ctx sdk.Context, id string) (*types.Workspace, error) {
	workspace, err := k.WorkspaceStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}
