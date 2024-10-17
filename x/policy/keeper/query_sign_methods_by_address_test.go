package keeper_test

import (
	"encoding/base64"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SignMethodsByAddress(t *testing.T) {
	valid_tx_bytes, _ := base64.StdEncoding.DecodeString("CoEECv4DCiAvemVucm9jay5wb2xpY3kuTXNnQWRkU2lnbk1ldGhvZBLZAwoqemVuMTN5M3RtNjhnbXU5a250Y3h3dm11ZTgycDZha2FjbnB0MnY3bnR5EqoDCiEvemVucm9jay5wb2xpY3kuU2lnbk1ldGhvZFBhc3NrZXkShAMKG1RiZlJGUDUzYTh5SlBrU0hJeXh3eWN2VFhYWRIUTbfRFP53a8yJPkSHIyxwycvTXXYatgGjY2ZtdGRub25lZ2F0dFN0bXSgaGF1dGhEYXRhWJgRfYLbodlevvhdzK1ieUAGBKTRmNw28dlJczoSN+bjLV0AAAAA+/wwBxVOTsyMC24CBVfXvQAUTbfRFP53a8yJPkSHIyxwycvTXXalAQIDJiABIVggf6JpOflG9El/S3+/YvEBT69317zGkfjG2XcvBluSMBYiWCDqkyO8H1ULfUL6lLUz7iQTC0Ilqu7Kwa8yS4nnQDsv5CKVAXsidHlwZSI6IndlYmF1dGhuLmNyZWF0ZSIsImNoYWxsZW5nZSI6IkNpRURSMExIVVVSWll6bkloMG12ajNxdXQ1V1JUZnMwYXRDOU04U2N3Ulc2RExVIiwib3JpZ2luIjoiaHR0cHM6Ly9kZXZlbG9wZXIubW96aWxsYS5vcmciLCJjcm9zc09yaWdpbiI6ZmFsc2V9ElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQNHQsdRRFljOciHSa+Peq63lZFN+zRq0L0zxJzBFboMtRIECgIIARgBEgQQwJoMGkD3sYFTIYM+izWXmSKQ+6ouuUwHfqPT0IWqG1mFnlmzhmiyk0MMGEhvlDDX2kMvCey/B7NS1+f5V1z9uQLN1S5R")
	valid_clientdata, _ := base64.StdEncoding.DecodeString("eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoiQ2lFRFIwTEhVVVJaWXpuSWgwbXZqM3F1dDVXUlRmczBhdEM5TThTY3dSVzZETFUiLCJvcmlnaW4iOiJodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZyIsImNyb3NzT3JpZ2luIjpmYWxzZX0=")
	valid_attestation, _ := base64.StdEncoding.DecodeString("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYEX2C26HZXr74XcytYnlABgSk0ZjcNvHZSXM6Ejfm4y1dAAAAAPv8MAcVTk7MjAtuAgVX170AFE230RT+d2vMiT5EhyMscMnL0112pQECAyYgASFYIH+iaTn5RvRJf0t/v2LxAU+vd9e8xpH4xtl3LwZbkjAWIlgg6pMjvB9VC31C+pS1M+4kEwtCJaruysGvMkuJ50A7L+Q=")
	valid_sign_method := &types.SignMethodPasskey{
		RawId:             []byte("some-passkey-id"),
		AttestationObject: valid_attestation,
		ClientDataJson:    valid_clientdata,
	}
	valid_config, err := codectypes.NewAnyWithValue(valid_sign_method)
	if err != nil {
		t.Fatalf("error encoding valid config, err %v", err)
	}

	keepers := keepertest.NewTest(t)
	ctx := keepers.Ctx.WithTxBytes(valid_tx_bytes)
	pk := keepers.PolicyKeeper
	msgSer := keeper.NewMsgServerImpl(*pk)

	_, err = msgSer.AddSignMethod(ctx, &types.MsgAddSignMethod{
		Creator: "owner1",
		Config:  valid_config,
	})
	assert.Nil(t, err)

	_, err = msgSer.AddSignMethod(ctx, &types.MsgAddSignMethod{
		Creator: "owner2",
		Config:  valid_config,
	})
	assert.Nil(t, err)

	res, err := pk.SignMethodsByAddress(ctx, &types.QuerySignMethodsByAddressRequest{
		Address: "owner2",
	})

	require.Nil(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Config)
	assert.Len(t, res.Config, 1)
}
