package wasmbinding

import (
	"encoding/json"
	"strconv"

	cosmoserrors "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// custom keepers
	identitykeeper "github.com/Zenrock-Foundation/zrchain/v4/x/identity/keeper"
	identitytypes "github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	treasurykeeper "github.com/Zenrock-Foundation/zrchain/v4/x/treasury/keeper"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
)

type CustomMessageInterceptor struct {
	wrapped  wasmkeeper.Messenger
	policy   *policykeeper.Keeper
	identity *identitykeeper.Keeper
	treasury *treasurykeeper.Keeper
}

type newWorkspaceRequest struct {
	Creator       string `json:"creator"`
	AdminPolicyId string `json:"admin_policy_id"`
	SignPolicyId  string `json:"sign_policy_id"`
}

type newKeyRequest struct {
	Creator       string `json:"creator"`
	WorkspaceAddr string `json:"workspace_addr"`
	KeyringAddr   string `json:"keyring_addr"`
	KeyType       string `json:"key_type"`
}

type signDataRequest struct {
	Creator        string `json:"creator"`
	KeyId          string `json:"key_id"`
	DataForSigning string `json:"data_for_signing"`
}

type transactionMetadata struct {
	TypeUrl string `json:"type_url"`
	Value   []byte `json:"value"`
}

type signTransactionRequest struct {
	Creator             string               `json:"creator"`
	KeyId               string               `json:"key_id"`
	WalletType          string               `json:"wallet_type"`
	UnsignedTransaction []byte               `json:"unsigned_transaction"`
	Metadata            *transactionMetadata `json:"metadata"`
}

type addWorkspaceOwnerRequest struct {
	Creator       string `json:"creator"`
	WorkspaceAddr string `json:"workspace_addr"`
	NewOwner      string `json:"new_owner"`
}

type zenrockMessages struct {
	NewWorkspaceRequest       *newWorkspaceRequest      `json:"new_workspace_request,omitempty"`
	NewKeyRequest             *newKeyRequest            `json:"new_key_request,omitempty"`
	NewSignDataRequest        *signDataRequest          `json:"new_sign_data_request,omitempty"`
	NewSignTransactionRequest *signTransactionRequest   `json:"new_sign_transaction_request,omitempty"`
	AddWorkspaceOwnerRequest  *addWorkspaceOwnerRequest `json:"add_workspace_owner_request,omitempty"`
}

func NewCustomMessageDecorator(
	policy *policykeeper.Keeper,
	identity *identitykeeper.Keeper,
	treasury *treasurykeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessageInterceptor{
			wrapped:  old,
			policy:   policy,
			identity: identity,
			treasury: treasury,
		}
	}
}

func (m *CustomMessageInterceptor) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, responses [][]*codectypes.Any, err error) {
	if msg.Custom != nil {
		ctx.Logger().Info(string(msg.Custom))

		var message zenrockMessages

		if err := json.Unmarshal(msg.Custom, &message); err != nil {
			return nil, nil, nil, cosmoserrors.Wrap(err, "unmarshal supported messages")
		}

		if message.NewWorkspaceRequest != nil {
			return nil, nil, nil, m.newWorkspace(ctx, message.NewWorkspaceRequest)
		}
		if message.NewKeyRequest != nil {
			return nil, nil, nil, m.newKey(ctx, message.NewKeyRequest)
		}
		if message.NewSignDataRequest != nil {
			return nil, nil, nil, m.signData(ctx, message.NewSignDataRequest)
		}
		if message.NewSignTransactionRequest != nil {
			return nil, nil, nil, m.signTx(ctx, message.NewSignTransactionRequest)
		}
		if message.AddWorkspaceOwnerRequest != nil {
			return nil, nil, nil, m.addWorkspaceOwner(ctx, message.AddWorkspaceOwnerRequest)
		}
	}

	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessageInterceptor) newWorkspace(ctx sdk.Context, req *newWorkspaceRequest) (err error) {
	internalReq := &identitytypes.MsgNewWorkspace{
		Creator: req.Creator,
	}

	internalReq.AdminPolicyId, err = strconv.ParseUint(req.AdminPolicyId, 10, 64)
	if err != nil {
		return cosmoserrors.Wrap(err, "adminpolicyid")
	}
	internalReq.SignPolicyId, err = strconv.ParseUint(req.SignPolicyId, 10, 64)
	if err != nil {
		return cosmoserrors.Wrap(err, "signpolicyid")
	}

	identityMsgServer := identitykeeper.NewMsgServerImpl(*m.identity)
	_, err = identityMsgServer.NewWorkspace(ctx, internalReq)
	if err != nil {
		return cosmoserrors.Wrap(err, "NewWorkspace")
	}

	return nil
}

func (m *CustomMessageInterceptor) newKey(ctx sdk.Context, req *newKeyRequest) (err error) {
	internalReq := &treasurytypes.MsgNewKeyRequest{
		Creator:       req.Creator,
		WorkspaceAddr: req.WorkspaceAddr,
		KeyringAddr:   req.KeyringAddr,
		KeyType:       req.KeyType,
	}

	treasuryMsgServer := treasurykeeper.NewMsgServerImpl(*m.treasury)
	if _, err = treasuryMsgServer.NewKeyRequest(ctx, internalReq); err != nil {
		return cosmoserrors.Wrap(err, "NewKey")
	}

	return nil
}

func (m *CustomMessageInterceptor) signData(ctx sdk.Context, req *signDataRequest) (err error) {
	internalReq := &treasurytypes.MsgNewSignatureRequest{
		Creator:        req.Creator,
		DataForSigning: req.DataForSigning,
	}

	internalReq.KeyId, err = strconv.ParseUint(req.KeyId, 10, 64)
	if err != nil {
		return cosmoserrors.Wrap(err, "keyId")
	}

	treasuryMsgServer := treasurykeeper.NewMsgServerImpl(*m.treasury)
	if _, err = treasuryMsgServer.NewSignatureRequest(ctx, internalReq); err != nil {
		return cosmoserrors.Wrap(err, "SignData")
	}

	return nil
}

func (m *CustomMessageInterceptor) signTx(ctx sdk.Context, req *signTransactionRequest) (err error) {
	internalReq := &treasurytypes.MsgNewSignTransactionRequest{
		Creator:             req.Creator,
		UnsignedTransaction: req.UnsignedTransaction,
	}

	internalReq.KeyId, err = strconv.ParseUint(req.KeyId, 10, 64)
	if err != nil {
		return cosmoserrors.Wrap(err, "keyId")
	}

	if req.Metadata != nil {
		internalReq.Metadata = &codectypes.Any{
			TypeUrl: req.Metadata.TypeUrl,
			Value:   req.Metadata.Value,
		}
	}

	switch req.WalletType {
	case "unspecified":
		internalReq.WalletType = treasurytypes.WalletType_WALLET_TYPE_UNSPECIFIED
	case "native":
		internalReq.WalletType = treasurytypes.WalletType_WALLET_TYPE_NATIVE
	case "evm":
		internalReq.WalletType = treasurytypes.WalletType_WALLET_TYPE_EVM
	case "btc_testnet":
		internalReq.WalletType = treasurytypes.WalletType_WALLET_TYPE_BTC_TESTNET
	case "btc_mainnet":
		internalReq.WalletType = treasurytypes.WalletType_WALLET_TYPE_BTC_MAINNET
	case "btc_regnet":
		internalReq.WalletType = treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET
	case "solana":
		internalReq.WalletType = treasurytypes.WalletType_WALLET_TYPE_SOLANA
	}

	treasuryMsgServer := treasurykeeper.NewMsgServerImpl(*m.treasury)
	if _, err = treasuryMsgServer.NewSignTransactionRequest(ctx, internalReq); err != nil {
		return cosmoserrors.Wrap(err, "SignTransaction")
	}

	return nil
}

func (m *CustomMessageInterceptor) addWorkspaceOwner(ctx sdk.Context, req *addWorkspaceOwnerRequest) (err error) {
	internalReq := &identitytypes.MsgAddWorkspaceOwner{
		Creator:       req.Creator,
		WorkspaceAddr: req.WorkspaceAddr,
		NewOwner:      req.NewOwner,
	}

	identityMsgServer := identitykeeper.NewMsgServerImpl(*m.identity)
	_, err = identityMsgServer.AddWorkspaceOwner(ctx, internalReq)
	if err != nil {
		return cosmoserrors.Wrap(err, "AddWorkspaceOwner")
	}

	return nil
}
