package wasmbinding

import (
	"encoding/json"

	cosmoserrors "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	identitykeeper "github.com/Zenrock-Foundation/zrchain/v4/x/identity/keeper"
	identitytypes "github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	treasurykeeper "github.com/Zenrock-Foundation/zrchain/v4/x/treasury/keeper"
)

type keyringByAddressQuery struct {
	KeyringAddr string `json:"keyring_addr"`
}

type zenrockQueries struct {
	KeyringByAddressQuery *keyringByAddressQuery `json:"keyring_by_address_query,omitempty"`
}

type CustomQueryPlugin struct {
	policy   *policykeeper.Keeper
	identity *identitykeeper.Keeper
	treasury *treasurykeeper.Keeper
}

func NewQueryPlugin(
	policy *policykeeper.Keeper,
	identity *identitykeeper.Keeper,
	treasury *treasurykeeper.Keeper,
) *CustomQueryPlugin {
	return &CustomQueryPlugin{
		policy:   policy,
		identity: identity,
		treasury: treasury,
	}
}

func CustomQuerier(qp *CustomQueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		msg, err := request.MarshalJSON()
		if err != nil {
			return nil, cosmoserrors.Wrap(err, "CustomQuerier - marshal")
		}
		ctx.Logger().Info(string(msg))

		var query zenrockQueries
		if err := json.Unmarshal(request, &query); err != nil {
			return nil, cosmoserrors.Wrap(err, "unmarshal supported messages")
		}

		if query.KeyringByAddressQuery != nil {
			return qp.keyringByAddress(ctx, query.KeyringByAddressQuery)
		}

		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown query"}
	}
}

func (qp *CustomQueryPlugin) keyringByAddress(ctx sdk.Context, query *keyringByAddressQuery) ([]byte, error) {
	internalQuery := &identitytypes.QueryKeyringByAddressRequest{
		KeyringAddr: query.KeyringAddr,
	}

	res, err := identitykeeper.Keeper.KeyringByAddress(*qp.identity, ctx, internalQuery)
	if err != nil {
		return []byte{}, cosmoserrors.Wrap(err, "KeyringByAddress")
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return []byte{}, cosmoserrors.Wrap(err, "KeyringByAddress - marshal")
	}

	return resBytes, nil
}
