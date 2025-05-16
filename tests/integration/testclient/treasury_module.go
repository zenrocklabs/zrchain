package testclient

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/gogoproto/proto"
)

func (c *TestClient) CreateKey(ctx context.Context, workspace, keyring, keytype string, signPolicyId uint64) (*treasurytypes.MsgNewKeyRequestResponse, []abcitypes.Event, error) {
	msg := &treasurytypes.MsgNewKeyRequest{
		Creator:       c.IdentitySigner.Address.String(),
		WorkspaceAddr: workspace,
		KeyringAddr:   keyring,
		KeyType:       keytype,
		Btl:           1000,
		SignPolicyId:  signPolicyId,
	}

	txres, evts, err := c.executeTxWithIdentity(ctx, c.IdentitySigner, msg)
	if err != nil {
		return nil, nil, err
	}

	res := &treasurytypes.MsgNewKeyRequestResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, evts, err
}

func (c *TestClient) CreateSignatureRequest(ctx context.Context, keyId uint64, dataForSigning string) (*treasurytypes.MsgNewSignatureRequestResponse, []abcitypes.Event, error) {
	msg := &treasurytypes.MsgNewSignatureRequest{
		Creator:        c.IdentitySigner.Address.String(),
		KeyIds:         []uint64{keyId},
		DataForSigning: dataForSigning,
	}

	txres, evts, err := c.executeTxWithIdentity(ctx, c.IdentitySigner, msg)
	if err != nil {
		return nil, nil, err
	}

	res := &treasurytypes.MsgNewSignatureRequestResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, evts, err
}

func (c *TestClient) CreateSignETHTransactionRequest(ctx context.Context, keyId uint64, unsignedTx []byte, chainId uint64) (*treasurytypes.MsgNewSignTransactionRequestResponse, []abcitypes.Event, error) {
	meta := &treasurytypes.MetadataEthereum{
		ChainId: chainId,
	}

	anyMeta, err := codectypes.NewAnyWithValue(meta)
	if err != nil {
		return nil, nil, err
	}

	msg := &treasurytypes.MsgNewSignTransactionRequest{
		Creator:             c.IdentitySigner.Address.String(),
		KeyId:               keyId,
		WalletType:          treasurytypes.WalletType_WALLET_TYPE_EVM,
		UnsignedTransaction: unsignedTx,
		Metadata:            anyMeta,
	}

	txres, evts, err := c.executeTxWithIdentity(ctx, c.IdentitySigner, msg)
	if err != nil {
		return nil, nil, err
	}

	res := &treasurytypes.MsgNewSignTransactionRequestResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, evts, err
}

func (c *TestClient) CreateICATransactionRequest(ctx context.Context) (any, error) {
	return nil, nil
}

func (c *TestClient) TransferFromKeyring(ctx context.Context, identity client.Identity, keyring, recipient, denom string, amount uint64) (*treasurytypes.MsgTransferFromKeyringResponse, error) {
	msg := &treasurytypes.MsgTransferFromKeyring{
		Creator:   identity.Address.String(),
		Keyring:   keyring,
		Recipient: recipient,
		Amount:    amount,
		Denom:     denom,
	}

	txres, _, err := c.executeTxWithIdentity(ctx, c.IdentitySigner, msg)
	if err != nil {
		return nil, err
	}

	res := &treasurytypes.MsgTransferFromKeyringResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) UpdateKeyPolicy(ctx context.Context, identity client.Identity, keyId, signPolicyId uint64) (*treasurytypes.MsgUpdateKeyPolicyResponse, error) {
	msg := &treasurytypes.MsgUpdateKeyPolicy{
		Creator:      c.IdentitySigner.Address.String(),
		KeyId:        keyId,
		SignPolicyId: signPolicyId,
	}

	txres, _, err := c.executeTxWithIdentity(ctx, identity, msg)
	if err != nil {
		return nil, err
	}

	res := &treasurytypes.MsgUpdateKeyPolicyResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) FulfilKeyRequest(ctx context.Context, identity client.Identity, requestId uint64, status treasurytypes.KeyRequestStatus, sig []byte, pubkey []byte) (*treasurytypes.MsgFulfilKeyRequestResponse, []abcitypes.Event, error) {
	msg := &treasurytypes.MsgFulfilKeyRequest{
		Creator:               identity.Address.String(),
		RequestId:             requestId,
		KeyringPartySignature: sig,
		Status:                status,
		Result: &treasurytypes.MsgFulfilKeyRequest_Key{
			Key: &treasurytypes.MsgNewKey{
				PublicKey: pubkey,
			},
		},
	}

	txres, evts, err := c.executeTxWithIdentity(ctx, identity, msg)
	if err != nil {
		return nil, nil, err
	}

	res := &treasurytypes.MsgFulfilKeyRequestResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, evts, err
}

func (c *TestClient) FulfilSignatureRequest(ctx context.Context) (any, error) {
	return nil, nil
}

func (c *TestClient) FulfilICATransactionRequest(ctx context.Context) (any, error) {
	return nil, nil
}

func (c *TestClient) GetKeyCount(ctx context.Context) (uint64, error) {
	res, err := c.tqc.Keys(ctx, &treasurytypes.QueryKeysRequest{})
	if err != nil {
		return 0, err
	}
	return res.Pagination.Total, nil
}

func (c *TestClient) GetKeys(ctx context.Context, workspaceAddress string, walletType treasurytypes.WalletType, prefixes []string) ([]*treasurytypes.KeyAndWalletResponse, error) {
	res, err := c.tqc.Keys(ctx, &treasurytypes.QueryKeysRequest{
		WorkspaceAddr: workspaceAddress,
		WalletType:    walletType,
		Prefixes:      prefixes,
	})
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

func (c *TestClient) GetKeyRequestCount(ctx context.Context) (uint64, error) {
	res, err := c.tqc.KeyRequests(ctx, &treasurytypes.QueryKeyRequestsRequest{})
	if err != nil {
		return 0, err
	}
	return res.Pagination.Total, nil
}

func (c *TestClient) GetKeyRequest(ctx context.Context, keyreqId uint64) (*treasurytypes.KeyReqResponse, error) {
	res, err := c.tqc.KeyRequestByID(ctx, &treasurytypes.QueryKeyRequestByIDRequest{
		Id: keyreqId,
	})
	if err != nil {
		return nil, err
	}

	return res.KeyRequest, nil
}

func (c *TestClient) GetSignTransactionRequestCount(ctx context.Context) (uint64, error) {
	res, err := c.tqc.SignTransactionRequests(ctx, &treasurytypes.QuerySignTransactionRequestsRequest{})
	if err != nil {
		return 0, err
	}
	return res.Pagination.Total, nil
}

func (c *TestClient) GetSignTransactionRequest(ctx context.Context, id uint64) (*treasurytypes.SignTxReqResponse, error) {
	res, err := c.tqc.SignTransactionRequestByID(ctx, &treasurytypes.QuerySignTransactionRequestByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return res.SignTransactionRequest, nil
}

func (c *TestClient) GetSignatureRequestCount(ctx context.Context) (uint64, error) {
	res, err := c.tqc.SignatureRequests(ctx, &treasurytypes.QuerySignatureRequestsRequest{})
	if err != nil {
		return 0, err
	}
	return res.Pagination.Total, nil
}

func (c *TestClient) GetSignatureRequest(ctx context.Context, id uint64) (*treasurytypes.SignReqResponse, error) {
	res, err := c.tqc.SignatureRequestByID(ctx, &treasurytypes.QuerySignatureRequestByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return res.SignRequest, nil
}
