package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	dbm "github.com/cosmos/cosmos-db"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type authzMock struct {
	GrantCalled  int
	RevokeCalled int

	Grants  []*authztypes.MsgGrant
	Revokes []*authztypes.MsgRevoke
}

func newAuthzMock() *authzMock {
	return &authzMock{
		Grants:  []*authztypes.MsgGrant{},
		Revokes: []*authztypes.MsgRevoke{},
	}
}

func (m *authzMock) Grant(_ context.Context, msg *authztypes.MsgGrant) (*authztypes.MsgGrantResponse, error) {
	m.GrantCalled++
	m.Grants = append(m.Grants, msg)
	return &authztypes.MsgGrantResponse{}, nil
}

func (m *authzMock) Revoke(_ context.Context, msg *authztypes.MsgRevoke) (*authztypes.MsgRevokeResponse, error) {
	m.RevokeCalled++
	m.Revokes = append(m.Revokes, msg)
	return &authztypes.MsgRevokeResponse{}, nil
}

func Test_msgServer_AddMultiGrant(t *testing.T) {
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
			name: "PASS: Add 2 messages to grantee",
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
			name: "FAIL: invalid grantee address",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := dbm.NewMemDB()
			stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
			authzmock := newAuthzMock()
			pk, ctx := keepertest.PolicyKeeper(t, db, stateStore, authzmock)

			msgser := keeper.NewMsgServerImpl(pk)
			res, err := msgser.AddMultiGrant(ctx, &types.MsgAddMultiGrant{
				Creator: tt.args.granter,
				Grantee: tt.args.grantee,
				Msgs:    tt.args.msgs,
			})

			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.NotNil(t, res)
				assert.Equal(t, 2, authzmock.GrantCalled)
				require.Equal(t, 2, len(authzmock.Grants))

				for i, _ := range tt.args.msgs {
					assert.Equal(t, tt.args.granter, authzmock.Grants[i].Granter)
					assert.Equal(t, tt.args.grantee, authzmock.Grants[i].Grantee)
					auth := &authztypes.GenericAuthorization{
						Msg: tt.args.msgs[i],
					}
					authAny, err := codectypes.NewAnyWithValue(auth)
					require.Nil(t, err)
					assert.Equal(t, authAny.TypeUrl, authzmock.Grants[i].Grant.Authorization.TypeUrl)
					assert.Equal(t, authAny.Value, authzmock.Grants[i].Grant.Authorization.Value)
				}
			}
		})
	}
}
