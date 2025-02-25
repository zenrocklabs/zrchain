package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	"github.com/Zenrock-Foundation/zrchain/v5/policy"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
)

type ExportedKeeper interface {
	AddAction(ctx sdk.Context, creator string, msg sdk.Msg, policyID, btl uint64, policyData map[string][]byte) (*types.Action, error)
	Codec() codec.BinaryCodec
	GeneratorHandler(reqType string) (func(sdk.Context, *cdctypes.Any) (policy.Policy, error), bool)
	RegisterPolicyGeneratorHandler(reqType string, f func(sdk.Context, *cdctypes.Any) (policy.Policy, error))
	ActionHandler(actionType string) (func(sdk.Context, *types.Action) (any, error), bool)
	RegisterActionHandler(actionType string, f func(sdk.Context, *types.Action) (any, error))
	GetPolicyParticipants(ctx context.Context, policyId uint64) (map[string]struct{}, error)
	PolicyMembersAreOwners(ctx context.Context, policyId uint64, workspaceOwners []string)
}

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	Schema          collections.Schema
	ParamStore      collections.Item[types.Params]
	ActionStore     collections.Map[uint64, types.Action]
	ActionCount     collections.Item[uint64]
	PolicyStore     collections.Map[uint64, types.Policy]
	PolicyCount     collections.Item[uint64]
	SignMethodStore collections.Map[collections.Pair[string, string], cdctypes.Any]

	ibcKeeperFn             func() *ibckeeper.Keeper
	capabilityScopedFn      func(string) capabilitykeeper.ScopedKeeper
	scopedKeeper            exported.ScopedKeeper
	actionHandlers          map[string]func(sdk.Context, *types.Action) (any, error)
	policyGeneratorHandlers map[string]func(sdk.Context, *cdctypes.Any) (policy.Policy, error)
	authzKeeper             types.AuthzKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	authzKeeper types.AuthzKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s - %v", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
		authzKeeper:  authzKeeper,

		ParamStore:      collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Params](cdc)),
		ActionStore:     collections.NewMap(sb, types.ActionsKey, types.ActionsIndex, collections.Uint64Key, codec.CollValue[types.Action](cdc)),
		ActionCount:     collections.NewItem(sb, types.ActionCountKey, types.ActionCountIndex, collections.Uint64Value),
		PolicyStore:     collections.NewMap(sb, types.PoliciesKey, types.PoliciesIndex, collections.Uint64Key, codec.CollValue[types.Policy](cdc)),
		PolicyCount:     collections.NewItem(sb, types.PolicyCountKey, types.PolicyCountIndex, collections.Uint64Value),
		SignMethodStore: collections.NewMap(sb, types.SignMethodsKey, types.SignMethodsIndex, collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[cdctypes.Any](cdc)),

		actionHandlers:          make(map[string]func(sdk.Context, *types.Action) (any, error)),
		policyGeneratorHandlers: make(map[string]func(sdk.Context, *cdctypes.Any) (policy.Policy, error)),
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

func (k Keeper) Codec() codec.BinaryCodec {
	return k.cdc
}

func (k Keeper) GeneratorHandler(url string) (func(sdk.Context, *cdctypes.Any) (policy.Policy, error), bool) {
	f, found := k.policyGeneratorHandlers[url]
	return f, found
}

func (k Keeper) RegisterPolicyGeneratorHandler(reqType string, f func(sdk.Context, *cdctypes.Any) (policy.Policy, error)) {
	k.policyGeneratorHandlers[reqType] = f
}

func (k Keeper) ActionHandler(actionType string) (func(sdk.Context, *types.Action) (any, error), bool) {
	f, found := k.actionHandlers[actionType]
	return f, found
}

func (k Keeper) RegisterActionHandler(actionType string, f func(sdk.Context, *types.Action) (any, error)) {
	k.actionHandlers[actionType] = f
}

func (k Keeper) GetPolicy(ctx sdk.Context, policyId uint64) (*types.Policy, error) {
	policy, err := k.PolicyStore.Get(ctx, policyId)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (k Keeper) Unpack(policyPb *types.Policy) (policy.Policy, error) {
	var p policy.Policy
	err := k.Codec().UnpackAny(policyPb.Policy, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (k Keeper) SetAction(ctx sdk.Context, action *types.Action) error {
	return k.ActionStore.Set(ctx, action.Id, *action)
}
