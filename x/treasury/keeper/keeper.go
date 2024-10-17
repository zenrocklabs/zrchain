package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v4/app/params"
	shared "github.com/Zenrock-Foundation/zrchain/v4/shared"
	"github.com/cosmos/cosmos-sdk/types/query"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	identity "github.com/Zenrock-Foundation/zrchain/v4/x/identity/keeper"
	policy "github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	Schema                      collections.Schema
	ParamStore                  collections.Item[types.Params]
	KeyStore                    collections.Map[uint64, types.Key]
	KeyRequestStore             collections.Map[uint64, types.KeyRequest]
	KeyRequestCount             collections.Item[uint64]
	SignRequestStore            collections.Map[uint64, types.SignRequest]
	SignRequestCount            collections.Item[uint64]
	SignTransactionRequestStore collections.Map[uint64, types.SignTransactionRequest]
	SignTransactionRequestCount collections.Item[uint64]
	ICATransactionRequestStore  collections.Map[uint64, types.ICATransactionRequest]
	ICATransactionRequestCount  collections.Item[uint64]

	ibcKeeperFn        func() *ibckeeper.Keeper
	capabilityScopedFn func(string) capabilitykeeper.ScopedKeeper
	scopedKeeper       exported.ScopedKeeper
	bankKeeper         types.BankKeeper
	identityKeeper     identity.Keeper
	policyKeeper       policy.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	bankKeeper types.BankKeeper,
	identityKeeper identity.Keeper,
	policyKeeper policy.Keeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:            cdc,
		storeService:   storeService,
		authority:      authority,
		logger:         logger,
		bankKeeper:     bankKeeper,
		identityKeeper: identityKeeper,
		policyKeeper:   policyKeeper,

		ParamStore:                  collections.NewItem(sb, types.ParamsKey, types.ParamsIndex, codec.CollValue[types.Params](cdc)),
		KeyStore:                    collections.NewMap(sb, types.KeysKey, types.KeysIndex, collections.Uint64Key, codec.CollValue[types.Key](cdc)),
		KeyRequestStore:             collections.NewMap(sb, types.KeyRequestsKey, types.KeyRequestsIndex, collections.Uint64Key, codec.CollValue[types.KeyRequest](cdc)),
		KeyRequestCount:             collections.NewItem(sb, types.KeyRequestCountKey, types.KeyRequestCountIndex, collections.Uint64Value),
		SignRequestStore:            collections.NewMap(sb, types.SignRequestsKey, types.SignRequestsIndex, collections.Uint64Key, codec.CollValue[types.SignRequest](cdc)),
		SignRequestCount:            collections.NewItem(sb, types.SignRequestCountKey, types.SignRequestCountIndex, collections.Uint64Value),
		SignTransactionRequestStore: collections.NewMap(sb, types.SignTransactionRequestsKey, types.SignTransactionRequestsIndex, collections.Uint64Key, codec.CollValue[types.SignTransactionRequest](cdc)),
		SignTransactionRequestCount: collections.NewItem(sb, types.SignTransactionRequestCountKey, types.SignTransactionRequestCountIndex, collections.Uint64Value),
		ICATransactionRequestStore:  collections.NewMap(sb, types.ICATransactionRequestsKey, types.ICATransactionRequestsIndex, collections.Uint64Key, codec.CollValue[types.ICATransactionRequest](cdc)),
		ICATransactionRequestCount:  collections.NewItem(sb, types.ICATransactionRequestCountKey, types.ICATransactionRequestCountIndex, collections.Uint64Value),
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

func (k *Keeper) zrSignKeyRequest(goCtx context.Context, msg *types.MsgNewKeyRequest) (*types.MsgNewKeyRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := k.ParamStore.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "get params")
	}

	ConvertedKeyType, err := shared.WalletTypeToKeyType(msg.ExtKeyType)
	if err != nil {
		return nil, errorsmod.Wrap(err, "covert wallet type to key type")
	}

	workspace, err := k.identityKeeper.GetZrSignWorkspace(goCtx, msg.ExtRequester, msg.ExtKeyType)
	if err != nil {
		return nil, errorsmod.Wrap(err, "get zr sign workspaces")
	}

	allKeys, _, err := query.CollectionFilteredPaginate[uint64, types.Key, collections.Map[uint64, types.Key], *types.Key](
		goCtx,
		k.KeyStore,
		nil,
		func(key uint64, value types.Key) (bool, error) {
			return value.WorkspaceAddr == workspace, nil
		},
		func(key uint64, value types.Key) (*types.Key, error) {
			return &value, nil
		},
	)

	unfulfilledRequests, _, err := query.CollectionFilteredPaginate[uint64, types.KeyRequest, collections.Map[uint64, types.KeyRequest], *types.KeyReqResponse](
		goCtx,
		k.KeyRequestStore,
		nil,
		func(key uint64, value types.KeyRequest) (bool, error) {
			return workspace == value.WorkspaceAddr &&
				value.Status != types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED, nil
		},
		func(key uint64, value types.KeyRequest) (*types.KeyReqResponse, error) {
			return &types.KeyReqResponse{
				Id:                     value.Id,
				Creator:                value.Creator,
				WorkspaceAddr:          value.WorkspaceAddr,
				KeyringAddr:            value.KeyringAddr,
				KeyType:                value.KeyType.String(),
				Status:                 value.Status.String(),
				KeyringPartySignatures: value.KeyringPartySignatures,
				RejectReason:           value.RejectReason,
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// take into account unfulfilled keys when calculating the latest index
	var index = uint64(len(allKeys)) + uint64(len(unfulfilledRequests))

	result, err := k.newKeyRequest(ctx, &types.MsgNewKeyRequest{
		Creator:       msg.ExtRequester,
		WorkspaceAddr: workspace,
		KeyringAddr:   params.MpcKeyring,
		KeyType:       ConvertedKeyType,
		Index:         index,
	})
	if err != nil {
		return nil, errorsmod.Wrap(err, "newKeyRequest failed")
	}

	return result, nil
}

func (k *Keeper) newKeyRequest(ctx sdk.Context, msg *types.MsgNewKeyRequest) (*types.MsgNewKeyRequestResponse, error) {
	if _, err := k.identityKeeper.WorkspaceStore.Get(ctx, msg.WorkspaceAddr); err != nil {
		return nil, fmt.Errorf("workspace %s not found", msg.WorkspaceAddr)
	}

	keyring, err := k.identityKeeper.KeyringStore.Get(ctx, msg.KeyringAddr)
	if err != nil {
		return nil, fmt.Errorf("keyring %s not found", msg.KeyringAddr)
	}

	if keyring.KeyReqFee > 0 {
		err := k.SplitKeyringFee(ctx, msg.Creator, keyring.Address, keyring.KeyReqFee)
		if err != nil {
			return nil, err
		}
	}

	var keyType types.KeyType
	typeStr := strings.ToLower(msg.KeyType)

	switch {
	case strings.Contains(typeStr, "ecdsa"):
		keyType = types.KeyType_KEY_TYPE_ECDSA_SECP256K1
	case strings.Contains(typeStr, "ed25519") || strings.Contains(typeStr, "eddsa"):
		keyType = types.KeyType_KEY_TYPE_EDDSA_ED25519
	case strings.Contains(typeStr, "bitcoin") || strings.Contains(typeStr, "btc"):
		keyType = types.KeyType_KEY_TYPE_BITCOIN_SECP256K1

	default:
		return nil, fmt.Errorf("unknown key type: %s", msg.KeyType)
	}

	req := &types.KeyRequest{
		Creator:        msg.Creator,
		WorkspaceAddr:  msg.WorkspaceAddr,
		KeyringAddr:    msg.KeyringAddr,
		KeyType:        keyType,
		Status:         types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
		Index:          msg.Index,
		SignPolicyId:   msg.SignPolicyId,
		ZenbtcMetadata: msg.ZenbtcMetadata,
	}

	id, err := k.AppendKeyRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventNewKeyRequest,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(id, 10)),
		),
	})

	return &types.MsgNewKeyRequestResponse{
		KeyReqId: id,
	}, nil
}

func (k *Keeper) signatureRequest(ctx sdk.Context, msg *types.MsgNewSignatureRequest) (*types.MsgNewSignatureRequestResponse, error) {
	key, err := k.KeyStore.Get(ctx, msg.KeyId)
	if err != nil {
		return nil, fmt.Errorf("key %v not found", msg.KeyId)
	}

	if _, err := k.identityKeeper.WorkspaceStore.Get(ctx, key.WorkspaceAddr); err != nil {
		return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
	}

	keyring, err := k.identityKeeper.KeyringStore.Get(ctx, key.KeyringAddr)
	if err != nil {
		return nil, fmt.Errorf("keyring %s not found", key.KeyringAddr)
	}

	var dataForSigning [][]byte
	payload := strings.Split(msg.DataForSigning, ",")
	for _, p := range payload {
		data, err := hex.DecodeString(p)
		if err != nil {
			return nil, err
		}
		dataForSigning = append(dataForSigning, data)
	}

	verified, err := VerifyDataForSigning(dataForSigning, msg.VerifySigningData, msg.VerifySigningDataVersion)
	if verified == types.Verification_Failed {
		return nil, fmt.Errorf("transaction & hash verfication transaction did not verify")
	}
	if err != nil {
		return nil, fmt.Errorf("error whilst verifying transaction & hashes %s", err.Error())
	}

	if keyring.SigReqFee > 0 {
		err := k.SplitKeyringFee(ctx, msg.Creator, keyring.Address, keyring.SigReqFee)
		if err != nil {
			return nil, err
		}
	}

	parentReq := &types.SignRequest{
		Creator:        msg.Creator,
		KeyId:          msg.KeyId,
		KeyType:        key.Type,
		DataForSigning: dataForSigning,
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
		CacheId:        msg.CacheId,
	}
	parentID, err := k.CreateSignRequest(ctx, parentReq)
	if err != nil {
		return nil, err
	}

	var childIDs []uint64

	for _, data := range dataForSigning {
		if len(dataForSigning) < 2 {
			break
		}

		req := &types.SignRequest{
			Creator:        msg.Creator,
			KeyId:          msg.KeyId,
			KeyType:        key.Type,
			DataForSigning: [][]byte{data},
			Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
			ParentReqId:    parentID,
		}

		id, err := k.CreateSignRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		childIDs = append(childIDs, id)
	}

	parentReq.ChildReqIds = childIDs
	k.SignRequestStore.Set(ctx, parentID, *parentReq)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventNewSignRequest,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(parentID, 10)),
		),
	})

	return &types.MsgNewSignatureRequestResponse{SigReqId: parentID}, nil
}

func (k *Keeper) HandleSignTransactionRequest(ctx sdk.Context, msg *types.MsgNewSignTransactionRequest, data []byte) (*types.MsgNewSignTransactionRequestResponse, error) {
	key, err := k.KeyStore.Get(ctx, msg.KeyId)
	if err != nil {
		return nil, fmt.Errorf("key not found")
	}

	keyring, err := k.identityKeeper.KeyringStore.Get(ctx, key.KeyringAddr)
	if err != nil {
		return nil, fmt.Errorf("keyring not found")
	}

	if keyring.SigReqFee > 0 {
		err := k.SplitKeyringFee(ctx, msg.Creator, keyring.Address, keyring.SigReqFee)
		if err != nil {
			return nil, err
		}
	}

	// generate signature request
	signatureRequest := &types.SignRequest{
		Creator:        msg.Creator,
		KeyId:          msg.KeyId,
		KeyType:        key.Type,
		DataForSigning: [][]byte{data},
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
		Metadata:       msg.Metadata,
		CacheId:        msg.CacheId,
	}

	signRequestID, err := k.CreateSignRequest(ctx, signatureRequest)
	if err != nil {
		return nil, err
	}

	id, err := k.CreateSignTransactionRequest(ctx, &types.SignTransactionRequest{
		Creator:             msg.Creator,
		SignRequestId:       signRequestID,
		KeyId:               msg.KeyId,
		WalletType:          msg.WalletType,
		UnsignedTransaction: msg.UnsignedTransaction,
		NoBroadcast:         msg.NoBroadcast,
	})
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventNewSignRequest,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(id, 10)),
		),
	})

	return &types.MsgNewSignTransactionRequestResponse{Id: id, SignatureRequestId: signRequestID}, nil
}

func (k *Keeper) SplitKeyringFee(ctx context.Context, from, to string, fee uint64) error {
	prms, err := k.ParamStore.Get(ctx)
	if err != nil {
		return err
	}

	zenrockFee := uint64(math.Round(float64(fee) * (float64(prms.KeyringCommission) / 100.0)))
	keyringFee := fee - zenrockFee

	dest := prms.KeyringCommissionDestination
	err = k.bankKeeper.SendCoins(
		ctx,
		sdk.MustAccAddressFromBech32(from),
		sdk.MustAccAddressFromBech32(dest),
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(zenrockFee))),
	)

	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoins(
		ctx,
		sdk.MustAccAddressFromBech32(from),
		sdk.MustAccAddressFromBech32(to),
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(keyringFee))),
	)

	return err
}
