package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	keepertest "github.com/Zenrock-Foundation/zrchain/v4/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/testutil/sample"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_RemoveMultiGrant(t *testing.T) {
	type args struct {
		granter string
		grantee string
		msgs    []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "PASS: Remove 2 messages from grantee",
			args: args{
				granter: sample.AccAddress(),
				grantee: sample.AccAddress(),
				msgs: []string{
					"some-msg-1",
					"some-msg-2",
				},
			},
			wantErr: false,
		},
		{
			name: "FAIL: invalid granter address",
			args: args{
				granter: "invalid-address",
				grantee: sample.AccAddress(),
				msgs: []string{
					"some-msg-1",
					"some-msg-2",
				},
			},
			wantErr: true,
		},
		{
			name: "PASS: invalid grantee address",
			args: args{
				granter: sample.AccAddress(),
				grantee: "invalid-address",
				msgs: []string{
					"some-msg-1",
					"some-msg-2",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := dbm.NewMemDB()
			stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
			authzmock := newAuthzMock()
			pk, ctx := keepertest.PolicyKeeper(t, db, stateStore, authzmock)

			msgser := keeper.NewMsgServerImpl(pk)
			res, err := msgser.RemoveMultiGrant(ctx, &types.MsgRemoveMultiGrant{
				Creator: tt.args.granter,
				Grantee: tt.args.grantee,
				Msgs:    tt.args.msgs,
			})

			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.NotNil(t, res)
				assert.Equal(t, 2, authzmock.RevokeCalled)
				require.Equal(t, 2, len(authzmock.Revokes))

				for i, _ := range tt.args.msgs {
					assert.Equal(t, tt.args.granter, authzmock.Revokes[i].Granter)
					assert.Equal(t, tt.args.grantee, authzmock.Revokes[i].Grantee)
					assert.Equal(t, tt.args.msgs[i], authzmock.Revokes[i].MsgTypeUrl)
				}
			}
		})
	}
}
