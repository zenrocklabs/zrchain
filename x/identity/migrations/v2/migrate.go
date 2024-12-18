package v2

import (
	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "identity"
)

var ParamsKey = []byte{0x01}

// MigrateKeyrings migrates the x/identity to v3. It migrates the keyring fields
// sig_req_fee -> fees.signature.rock_amount
// key_req_fee -> fees.key.rock_amount
func MigrateKeyrings(
	ctx sdk.Context,
	ks collections.Map[string, types.Keyring],
	cdc codec.BinaryCodec,
) error {

	it, err := ks.Iterate(ctx, nil)
	if err != nil {
		return errors.New("failed to iterate keyring")
	}

	for ; it.Valid(); it.Next() {
		key, err := it.Key()
		if err != nil {
			return err
		}
		value, err := it.Value()
		if err != nil {
			return err
		}

		value.Fees = &types.KeyringFees{
			Key: &types.KeyringFee{
				RockAmount: value.KeyReqFee,
			},
			Signature: &types.KeyringFee{
				RockAmount: value.SigReqFee,
			},
		}

		ks.Set(ctx, string(key), value)
	}
	return nil
}
