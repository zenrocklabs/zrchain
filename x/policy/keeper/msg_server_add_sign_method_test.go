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

func Test_msgServer_AddSignMethod(t *testing.T) {
	empty_config, err := codectypes.NewAnyWithValue(&types.SignMethodPasskey{})
	require.NoError(t, err)

	unknown_config_type, err := codectypes.NewAnyWithValue(&types.MsgAddSignMethod{})
	require.NoError(t, err)

	invalid_config, err := codectypes.NewAnyWithValue(&types.SignMethodPasskey{
		RawId:             []byte("some-id"),
		AttestationObject: []byte{},
		ClientDataJson:    []byte{},
	})
	require.NoError(t, err)
	invalid_clientdata_challenge, _ := base64.StdEncoding.DecodeString("eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoiQ2lFRFIwTEhVVVJaWXpuSWgwbXZqM3F1dDVXUlRmczBhdEM5TThTY3dSVzZETFkiLCJvcmlnaW4iOiJodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZyIsImNyb3NzT3JpZ2luIjpmYWxzZX0=")

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

	valid_sign_method.Active = true
	stored_valid_config, err := codectypes.NewAnyWithValue(valid_sign_method)
	require.NoError(t, err)

	invalid_config_challenge, err := codectypes.NewAnyWithValue(&types.SignMethodPasskey{
		RawId:             valid_rawid,
		AttestationObject: valid_attestation,
		ClientDataJson:    invalid_clientdata_challenge,
	})
	require.NoError(t, err)

	type args struct {
		creator  string
		config   *codectypes.Any
		tx_bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgAddSignMethodResponse
		wantErr string
	}{
		{
			name: "FAIL: invalid config",
			args: args{
				creator: "some-creator",
				config:  empty_config,
			},
			want:    nil,
			wantErr: "invalid config: ",
		},
		{
			name: "FAIL: invalid config",
			args: args{
				creator: "some-creator",
				config:  unknown_config_type,
			},
			want:    nil,
			wantErr: "no concrete type registered",
		},
		{
			name: "FAIL: invalid config",
			args: args{
				creator: "some-creator",
				config:  invalid_config,
			},
			want:    nil,
			wantErr: "invalid config: ",
		},
		{
			name: "FAIL: invalid config",
			args: args{
				creator: "some-creator",
				config:  invalid_config,
			},
			want:    nil,
			wantErr: "invalid config: ",
		},
		{
			name: "FAIL: invalid tx",
			args: args{
				creator: "some-creator",
				config:  valid_config,
			},
			want:    nil,
			wantErr: "invalid config: no signers in tx",
		},
		{
			name: "FAIL: invalid config challenge",
			args: args{
				creator:  "some-creator",
				config:   invalid_config_challenge,
				tx_bytes: valid_tx_bytes,
			},
			want:    nil,
			wantErr: "invalid config: invalid challenge",
		},
		{
			name: "PASS: valid config",
			args: args{
				creator:  "some-creator",
				config:   valid_config,
				tx_bytes: valid_tx_bytes,
			},
			want:    &types.MsgAddSignMethodResponse{},
			wantErr: "",
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

			// Act
			res, err := msgSer.AddSignMethod(ctx, &types.MsgAddSignMethod{
				Creator: tt.args.creator,
				Config:  tt.args.config,
			})

			// Assert
			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}
			require.Contains(t, errMsg, tt.wantErr)
			require.Equal(t, tt.want, res)

			if tt.wantErr == "" {
				getRes, getErr := pk.SignMethodsByAddress(ctx, &types.QuerySignMethodsByAddressRequest{
					Address: tt.args.creator,
				})
				require.Nil(t, getErr)
				require.NotNil(t, getRes)
				require.Len(t, getRes.Config, 1)
				require.Equal(t, stored_valid_config.Compare(getRes.Config[0]), 0)
			}
		})
	}
}
