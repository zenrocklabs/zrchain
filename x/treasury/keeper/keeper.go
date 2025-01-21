package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/app/params"
	shared "github.com/Zenrock-Foundation/zrchain/v5/shared"
	idtypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
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

	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	policy "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
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
	mintKeeper         types.MintKeeper
	memStore           store.MemoryStoreService
	zenBTCKeeper       shared.ZenBTCKeeper
}

type ValidationKeeper interface {
	GetBitcoinProxyCreatorID(ctx context.Context) string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	bankKeeper types.BankKeeper,
	identityKeeper identity.Keeper,
	policyKeeper policy.Keeper,
	mintKeeper types.MintKeeper,
	memStore store.MemoryStoreService,
	zenBTCKeeper shared.ZenBTCKeeper,
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
		mintKeeper:     mintKeeper,
		zenBTCKeeper:   zenBTCKeeper,

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
		memStore:                    memStore,
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
		err := k.EscrowKeyringFee(ctx, msg.Creator, keyring.KeyReqFee)
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

	timeout := keyring.MpcDefaultTimeout
	if msg.MpcTimeout > keyring.MpcMinimumTimeout {
		timeout = msg.MpcTimeout
	}
	btl := k.secondsToBTL(ctx, timeout)

	req := &types.KeyRequest{
		Creator:        msg.Creator,
		WorkspaceAddr:  msg.WorkspaceAddr,
		KeyringAddr:    msg.KeyringAddr,
		KeyType:        keyType,
		Status:         types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
		Index:          msg.Index,
		SignPolicyId:   msg.SignPolicyId,
		ZenbtcMetadata: msg.ZenbtcMetadata,
		Btl:            uint64(ctx.BlockHeight()) + btl,
		Fee:            keyring.KeyReqFee,
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

func (k *Keeper) HandleSignatureRequest(ctx sdk.Context, msg *types.MsgNewSignatureRequest) (*types.MsgNewSignatureRequestResponse, error) {
	dataForSigning, err := dataForSigning(msg.DataForSigning)
	if err != nil {
		return nil, err
	}

	// Verify the number of key IDs matches the number of data elements
	if len(dataForSigning) != len(msg.KeyIds) {
		return nil, fmt.Errorf("number of key IDs (%d) does not match number of data elements (%d)",
			len(msg.KeyIds), len(dataForSigning))
	}

	verified, err := VerifyDataForSigning(dataForSigning, msg.VerifySigningData, msg.VerifySigningDataVersion)
	if verified == types.Verification_Failed {
		return nil, fmt.Errorf("transaction & hash verification transaction did not verify")
	}
	if err != nil {
		return nil, fmt.Errorf("error whilst verifying transaction & hashes %s", err.Error())
	}

	id, err := k.processSignatureRequests(ctx, dataForSigning, msg.KeyIds, &types.SignRequest{
		Creator:        msg.Creator,
		DataForSigning: dataForSigning,
		KeyIds:         msg.KeyIds,
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
		CacheId:        msg.CacheId,
	}, msg.MpcTimeout)
	if err != nil {
		return nil, errorsmod.Wrap(err, "processSignatureRequests")
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventNewSignRequest,
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(id, 10)),
		),
	})

	return &types.MsgNewSignatureRequestResponse{SigReqId: id}, nil
}

func (k *Keeper) processSignatureRequests(ctx sdk.Context, dataForSigning [][]byte, keyIds []uint64, req *types.SignRequest, mpcTimeout uint64) (uint64, error) {
	// Verify all keys exist and collect keyring fees
	bh := ctx.BlockHeight()
	var sigReqs []*types.SignRequest
	for _, keyID := range keyIds {
		key, err := k.KeyStore.Get(ctx, keyID)
		if err != nil {
			return 0, fmt.Errorf("key %v not found", keyID)
		}

		if err := k.validateZenBTCSignRequest(ctx, *req, key); err != nil {
			return 0, err
		}

		keyring, err := k.identityKeeper.KeyringStore.Get(ctx, key.KeyringAddr)
		if err != nil {
			return 0, fmt.Errorf("keyring %s not found", key.KeyringAddr)
		}
		timeout := keyring.MpcDefaultTimeout
		if mpcTimeout > keyring.MpcMinimumTimeout {
			timeout = mpcTimeout
		}

		sigReqs = append(sigReqs, &types.SignRequest{
			Btl: uint64(bh) + k.secondsToBTL(ctx, timeout),
			Fee: keyring.SigReqFee,
		})

		// Accumulate fees per keyring
		if keyring.SigReqFee > 0 {
			err := k.EscrowKeyringFee(ctx, req.Creator, keyring.SigReqFee)
			if err != nil {
				return 0, err
			}
		}
		req.KeyType = key.GetType()
	}
	req.Btl = sigReqs[0].Btl
	req.Fee = sigReqs[0].Fee

	// Create parent request
	parentID, err := k.CreateSignRequest(ctx, req)
	if err != nil {
		return 0, err
	}

	var childIDs []uint64

	// Create child requests if there are multiple data elements
	if len(dataForSigning) > 1 {
		for i, data := range dataForSigning {
			sigReqs[i].Creator = req.Creator
			sigReqs[i].KeyIds = []uint64{keyIds[i]} // Use keyId corresponding to the data (hash)
			sigReqs[i].KeyType = req.KeyType
			sigReqs[i].DataForSigning = [][]byte{data} // only first element is used for data + keyId on child req
			sigReqs[i].Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING
			sigReqs[i].ParentReqId = parentID
			sigReqs[i].CacheId = req.CacheId
			sigReqs[i].Btl = req.Btl
			sigReqs[i].Fee = req.Fee

			childID, err := k.CreateSignRequest(ctx, sigReqs[i])
			if err != nil {
				return 0, err
			}
			childIDs = append(childIDs, childID)
		}

		// Update parent with child IDs
		req.ChildReqIds = childIDs
		k.SignRequestStore.Set(ctx, parentID, *req)
	}

	return parentID, nil
}

func (k *Keeper) HandleSignTransactionRequest(ctx sdk.Context, msg *types.MsgNewSignTransactionRequest, data []byte) (*types.MsgNewSignTransactionRequestResponse, error) {
	dataForSigning, err := dataForSigning(string(data))
	if err != nil {
		return nil, err
	}

	keyIDs := []uint64{msg.KeyId}
	id, err := k.processSignatureRequests(ctx, dataForSigning, keyIDs, &types.SignRequest{
		Creator:        msg.Creator,
		KeyIds:         keyIDs,
		DataForSigning: dataForSigning,
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
		Metadata:       msg.Metadata,
		CacheId:        msg.CacheId,
	}, msg.MpcTimeout)
	if err != nil {
		return nil, err
	}

	tID, err := k.CreateSignTransactionRequest(ctx, &types.SignTransactionRequest{
		Creator:             msg.Creator,
		SignRequestId:       id,
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
			sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(tID, 10)),
		),
	})
	return &types.MsgNewSignTransactionRequestResponse{Id: id, SignatureRequestId: id}, nil
}

func (k *Keeper) validateZenBTCSignRequest(ctx context.Context, req types.SignRequest, key types.Key) error {
	if key.ZenbtcMetadata != nil && key.ZenbtcMetadata.RecipientAddr != "" &&
		req.Creator != k.zenBTCKeeper.GetBitcoinProxyAddress(ctx) {
		return fmt.Errorf("only the Bitcoin proxy service can request signatures from zenBTC deposit keys")
	}
	return nil
}

func dataForSigning(data string) ([][]byte, error) {
	var dataForSigning [][]byte
	payload := strings.Split(data, ",")
	for _, p := range payload {
		data, err := hex.DecodeString(p)
		if err != nil {
			return nil, err
		}
		dataForSigning = append(dataForSigning, data)
	}
	return dataForSigning, nil
}

func (k *Keeper) EscrowKeyringFee(ctx context.Context, from string, fee uint64) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(from),
		types.KeyringEscrowName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(fee))),
	)
	return err
}

func (k *Keeper) SplitKeyringFee(ctx context.Context, from, to string, fee uint64) error {
	prms, err := k.ParamStore.Get(ctx)
	if err != nil {
		return err
	}

	zenrockFee := uint64(math.Round(float64(fee) * (float64(prms.KeyringCommission) / 100.0)))
	keyringFee := fee - zenrockFee

	if err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(from),
		types.KeyringCollectorName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(zenrockFee))),
	); err != nil {
		return err
	}

	if to == types.KeyringCollectorName {
		err = k.bankKeeper.SendCoinsFromAccountToModule(
			ctx,
			sdk.MustAccAddressFromBech32(from),
			types.KeyringCollectorName,
			sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(keyringFee))),
		)
	} else {
		err = k.bankKeeper.SendCoins(
			ctx,
			sdk.MustAccAddressFromBech32(from),
			sdk.MustAccAddressFromBech32(to),
			sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(keyringFee))),
		)
	}

	return err
}

func (k Keeper) CheckForKeyMPCTimeouts(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	requests, _, err := query.CollectionFilteredPaginate[uint64, types.KeyRequest, collections.Map[uint64, types.KeyRequest], *types.KeyRequest](
		goCtx,
		k.KeyRequestStore,
		nil,
		func(key uint64, value types.KeyRequest) (bool, error) {
			return value.Status == types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING, nil
		},
		func(key uint64, value types.KeyRequest) (*types.KeyRequest, error) {
			return &value, nil
		},
	)
	if err != nil {
		return err
	}

	blockHeight := ctx.BlockHeight()
	for _, req := range requests {
		if req.Btl < uint64(blockHeight) {
			if req.Fee > 0 {
				err = k.bankKeeper.SendCoinsFromModuleToAccount(
					goCtx,
					types.KeyringEscrowName,
					sdk.MustAccAddressFromBech32(req.Creator),
					sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(req.Fee))),
				)
				if err != nil {
					return err
				}
			}
			req.Status = types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED
			if err := k.KeyRequestStore.Set(ctx, req.Id, *req); err != nil {
				return err
			}
			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventKeyRequestRejected,
					sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(req.GetId(), 10)),
				),
			})
		}
	}

	return nil
}

func (k Keeper) CheckForSignatureMPCTimeouts(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	keyrings := map[string]*idtypes.Keyring{}
	requests, _, err := query.CollectionFilteredPaginate[uint64, types.SignRequest, collections.Map[uint64, types.SignRequest], *types.SignRequest](
		goCtx,
		k.SignRequestStore,
		nil,
		func(key uint64, value types.SignRequest) (bool, error) {
			return value.Status == types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING, nil
		},
		func(key uint64, value types.SignRequest) (*types.SignRequest, error) {
			return &value, nil
		},
	)
	if err != nil {
		return err
	}

	blockHeight := ctx.BlockHeight()
	for _, req := range requests {
		if req.Btl < uint64(blockHeight) {
			var keyring idtypes.Keyring
			keyID := req.KeyId
			if len(req.KeyIds) > 0 {
				keyID = req.KeyIds[0]
			}
			key, err := k.KeyStore.Get(goCtx, keyID)
			if err != nil {
				return fmt.Errorf("key %v not found", keyID)
			}
			if kr, ok := keyrings[key.KeyringAddr]; ok {
				keyring = *kr
			} else {

				kr, err := k.identityKeeper.KeyringStore.Get(goCtx, key.KeyringAddr)
				if err != nil {
					return fmt.Errorf("keyring %s not found", key.KeyringAddr)
				}
				keyring = kr
				keyrings[key.KeyringAddr] = &kr
			}
			if keyring.SigReqFee > 0 {
				err = k.bankKeeper.SendCoinsFromModuleToAccount(
					goCtx,
					types.KeyringEscrowName,
					sdk.MustAccAddressFromBech32(req.Creator),
					sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(req.Fee))),
				)
				if err != nil {
					return err
				}
			}
			req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
			if err := k.SignRequestStore.Set(ctx, req.Id, *req); err != nil {
				return err
			}
			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventSignRequestRejected,
					sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(req.GetId(), 10)),
				),
			})
		}
	}

	return nil
}

func (k Keeper) secondsToBTL(ctx sdk.Context, seconds uint64) uint64 {
	avgBT := uint64(k.AverageBlockTime(ctx))

	if seconds <= avgBT {
		return 1
	}

	return seconds / avgBT
}
