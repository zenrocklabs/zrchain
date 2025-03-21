package keeper_test

import (
	"encoding/base64"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/keeper"
	policy "github.com/Zenrock-Foundation/zrchain/v6/x/policy/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_RemoveSignMethod(t *testing.T) {
	valid_tx_bytes, _ := base64.StdEncoding.DecodeString("CoEECv4DCiAvemVucm9jay5wb2xpY3kuTXNnQWRkU2lnbk1ldGhvZBLZAwoqemVuMTN5M3RtNjhnbXU5a250Y3h3dm11ZTgycDZha2FjbnB0MnY3bnR5EqoDCiEvemVucm9jay5wb2xpY3kuU2lnbk1ldGhvZFBhc3NrZXkShAMKG1RiZlJGUDUzYTh5SlBrU0hJeXh3eWN2VFhYWRIUTbfRFP53a8yJPkSHIyxwycvTXXYatgGjY2ZtdGRub25lZ2F0dFN0bXSgaGF1dGhEYXRhWJgRfYLbodlevvhdzK1ieUAGBKTRmNw28dlJczoSN+bjLV0AAAAA+/wwBxVOTsyMC24CBVfXvQAUTbfRFP53a8yJPkSHIyxwycvTXXalAQIDJiABIVggf6JpOflG9El/S3+/YvEBT69317zGkfjG2XcvBluSMBYiWCDqkyO8H1ULfUL6lLUz7iQTC0Ilqu7Kwa8yS4nnQDsv5CKVAXsidHlwZSI6IndlYmF1dGhuLmNyZWF0ZSIsImNoYWxsZW5nZSI6IkNpRURSMExIVVVSWll6bkloMG12ajNxdXQ1V1JUZnMwYXRDOU04U2N3Ulc2RExVIiwib3JpZ2luIjoiaHR0cHM6Ly9kZXZlbG9wZXIubW96aWxsYS5vcmciLCJjcm9zc09yaWdpbiI6ZmFsc2V9ElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQNHQsdRRFljOciHSa+Peq63lZFN+zRq0L0zxJzBFboMtRIECgIIARgBEgQQwJoMGkD3sYFTIYM+izWXmSKQ+6ouuUwHfqPT0IWqG1mFnlmzhmiyk0MMGEhvlDDX2kMvCey/B7NS1+f5V1z9uQLN1S5R")
	valid_rawid, _ := base64.StdEncoding.DecodeString("TbfRFP53a8yJPkSHIyxwycvTXXY=")
	valid_clientdata, _ := base64.StdEncoding.DecodeString("eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoiQ2lFRFIwTEhVVVJaWXpuSWgwbXZqM3F1dDVXUlRmczBhdEM5TThTY3dSVzZETFUiLCJvcmlnaW4iOiJodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZyIsImNyb3NzT3JpZ2luIjpmYWxzZX0=")
	valid_attestation, _ := base64.StdEncoding.DecodeString("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYEX2C26HZXr74XcytYnlABgSk0ZjcNvHZSXM6Ejfm4y1dAAAAAPv8MAcVTk7MjAtuAgVX170AFE230RT+d2vMiT5EhyMscMnL0112pQECAyYgASFYIH+iaTn5RvRJf0t/v2LxAU+vd9e8xpH4xtl3LwZbkjAWIlgg6pMjvB9VC31C+pS1M+4kEwtCJaruysGvMkuJ50A7L+Q=")
	valid_sign_method := &types.SignMethodPasskey{
		RawId:             valid_rawid,
		AttestationObject: valid_attestation,
		ClientDataJson:    valid_clientdata,
	}
	valid_config, err := codectypes.NewAnyWithValue(valid_sign_method)
	require.NoError(t, err)

	type args struct {
		creator  string
		config   *codectypes.Any
		tx_bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgRemoveSignMethodResponse
		wantErr string
	}{
		{
			name:    "FAIL: id not found",
			want:    nil,
			wantErr: "not found",
		},
		{
			name:    "PASS: deavtivated",
			want:    &types.MsgRemoveSignMethodResponse{},
			wantErr: "",
			args: args{
				creator:  "some-creator",
				config:   valid_config,
				tx_bytes: valid_tx_bytes,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			keepers := keepertest.NewTest(t)
			ctx := keepers.Ctx.WithTxBytes(tt.args.tx_bytes)
			pk := keepers.PolicyKeeper
			msgSer := keeper.NewMsgServerImpl(*pk)
			polGenesis := types.GenesisState{
				Params:   types.Params{},
				PortId:   "42",
				Policies: []types.Policy{},
				Actions:  []types.Action{},
			}
			policy.InitGenesis(ctx, *pk, polGenesis)

			if tt.args.config != nil {
				_, err := msgSer.AddSignMethod(ctx, &types.MsgAddSignMethod{
					Creator: tt.args.creator,
					Config:  tt.args.config,
				})
				require.NoError(t, err)
			}

			// Act
			res, err := msgSer.RemoveSignMethod(ctx, &types.MsgRemoveSignMethod{
				Creator: tt.args.creator,
				Id:      valid_sign_method.GetConfigId(),
			})

			// Assert
			if tt.wantErr != "" {
				require.Contains(t, err.Error(), tt.wantErr)
				require.Equal(t, tt.want, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, res)

				getRes, err := pk.SignMethodsByAddress(ctx, &types.QuerySignMethodsByAddressRequest{
					Address: tt.args.creator,
				})
				require.NoError(t, err)
				require.NotNil(t, getRes)
				require.Len(t, getRes.Config, 1)
				require.Equal(t, valid_config.Compare(getRes.Config[0]), 0)
			}
		})
	}
}
