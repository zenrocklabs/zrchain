// Code generated by MockGen. DO NOT EDIT.
// Source: x/identity/types/expected_keepers.go

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	policy "github.com/Zenrock-Foundation/zrchain/v5/policy"
	types "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	codec "github.com/cosmos/cosmos-sdk/codec"
	types0 "github.com/cosmos/cosmos-sdk/codec/types"
	types1 "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

// MockAccountKeeper is a mock of AccountKeeper interface.
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper.
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance.
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockAccountKeeper) GetAccount(arg0 context.Context, arg1 types1.AccAddress) types1.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(types1.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountKeeperMockRecorder) GetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), arg0, arg1)
}

// MockPolicyKeeper is a mock of PolicyKeeper interface.
type MockPolicyKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyKeeperMockRecorder
}

// MockPolicyKeeperMockRecorder is the mock recorder for MockPolicyKeeper.
type MockPolicyKeeperMockRecorder struct {
	mock *MockPolicyKeeper
}

// NewMockPolicyKeeper creates a new mock instance.
func NewMockPolicyKeeper(ctrl *gomock.Controller) *MockPolicyKeeper {
	mock := &MockPolicyKeeper{ctrl: ctrl}
	mock.recorder = &MockPolicyKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPolicyKeeper) EXPECT() *MockPolicyKeeperMockRecorder {
	return m.recorder
}

// ActionHandler mocks base method.
func (m *MockPolicyKeeper) ActionHandler(actionType string) (func(types1.Context, *types.Action) (any, error), bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActionHandler", actionType)
	ret0, _ := ret[0].(func(types1.Context, *types.Action) (any, error))
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ActionHandler indicates an expected call of ActionHandler.
func (mr *MockPolicyKeeperMockRecorder) ActionHandler(actionType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionHandler", reflect.TypeOf((*MockPolicyKeeper)(nil).ActionHandler), actionType)
}

// AddAction mocks base method.
func (m *MockPolicyKeeper) AddAction(ctx types1.Context, creator string, msg types1.Msg, policyID, btl uint64, policyData map[string][]byte) (*types.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAction", ctx, creator, msg, policyID, btl, policyData)
	ret0, _ := ret[0].(*types.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAction indicates an expected call of AddAction.
func (mr *MockPolicyKeeperMockRecorder) AddAction(ctx, creator, msg, policyID, btl, policyData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAction", reflect.TypeOf((*MockPolicyKeeper)(nil).AddAction), ctx, creator, msg, policyID, btl, policyData)
}

// Codec mocks base method.
func (m *MockPolicyKeeper) Codec() codec.BinaryCodec {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Codec")
	ret0, _ := ret[0].(codec.BinaryCodec)
	return ret0
}

// Codec indicates an expected call of Codec.
func (mr *MockPolicyKeeperMockRecorder) Codec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Codec", reflect.TypeOf((*MockPolicyKeeper)(nil).Codec))
}

// GeneratorHandler mocks base method.
func (m *MockPolicyKeeper) GeneratorHandler(reqType string) (func(types1.Context, *types0.Any) (policy.Policy, error), bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneratorHandler", reqType)
	ret0, _ := ret[0].(func(types1.Context, *types0.Any) (policy.Policy, error))
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GeneratorHandler indicates an expected call of GeneratorHandler.
func (mr *MockPolicyKeeperMockRecorder) GeneratorHandler(reqType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneratorHandler", reflect.TypeOf((*MockPolicyKeeper)(nil).GeneratorHandler), reqType)
}

// GetPolicy mocks base method.
func (m *MockPolicyKeeper) GetPolicy(ctx types1.Context, policyId uint64) (*types.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicy", ctx, policyId)
	ret0, _ := ret[0].(*types.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPolicy indicates an expected call of GetPolicy.
func (mr *MockPolicyKeeperMockRecorder) GetPolicy(ctx, policyId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicy", reflect.TypeOf((*MockPolicyKeeper)(nil).GetPolicy), ctx, policyId)
}

// GetPolicyParticipants mocks base method.
func (m *MockPolicyKeeper) GetPolicyParticipants(ctx context.Context, policyId uint64) (map[string]struct{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicyParticipants", ctx, policyId)
	ret0, _ := ret[0].(map[string]struct{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPolicyParticipants indicates an expected call of GetPolicyParticipants.
func (mr *MockPolicyKeeperMockRecorder) GetPolicyParticipants(ctx, policyId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicyParticipants", reflect.TypeOf((*MockPolicyKeeper)(nil).GetPolicyParticipants), ctx, policyId)
}

// PolicyForAction mocks base method.
func (m *MockPolicyKeeper) PolicyForAction(ctx types1.Context, act *types.Action) (policy.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PolicyForAction", ctx, act)
	ret0, _ := ret[0].(policy.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PolicyForAction indicates an expected call of PolicyForAction.
func (mr *MockPolicyKeeperMockRecorder) PolicyForAction(ctx, act interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PolicyForAction", reflect.TypeOf((*MockPolicyKeeper)(nil).PolicyForAction), ctx, act)
}

// PolicyMembersAreOwners mocks base method.
func (m *MockPolicyKeeper) PolicyMembersAreOwners(ctx context.Context, policyId uint64, wsOwners []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PolicyMembersAreOwners", ctx, policyId, wsOwners)
	ret0, _ := ret[0].(error)
	return ret0
}

// PolicyMembersAreOwners indicates an expected call of PolicyMembersAreOwners.
func (mr *MockPolicyKeeperMockRecorder) PolicyMembersAreOwners(ctx, policyId, wsOwners interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PolicyMembersAreOwners", reflect.TypeOf((*MockPolicyKeeper)(nil).PolicyMembersAreOwners), ctx, policyId, wsOwners)
}

// RegisterActionHandler mocks base method.
func (m *MockPolicyKeeper) RegisterActionHandler(actionType string, f func(types1.Context, *types.Action) (any, error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterActionHandler", actionType, f)
}

// RegisterActionHandler indicates an expected call of RegisterActionHandler.
func (mr *MockPolicyKeeperMockRecorder) RegisterActionHandler(actionType, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterActionHandler", reflect.TypeOf((*MockPolicyKeeper)(nil).RegisterActionHandler), actionType, f)
}

// RegisterPolicyGeneratorHandler mocks base method.
func (m *MockPolicyKeeper) RegisterPolicyGeneratorHandler(reqType string, f func(types1.Context, *types0.Any) (policy.Policy, error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterPolicyGeneratorHandler", reqType, f)
}

// RegisterPolicyGeneratorHandler indicates an expected call of RegisterPolicyGeneratorHandler.
func (mr *MockPolicyKeeperMockRecorder) RegisterPolicyGeneratorHandler(reqType, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPolicyGeneratorHandler", reflect.TypeOf((*MockPolicyKeeper)(nil).RegisterPolicyGeneratorHandler), reqType, f)
}

// SetAction mocks base method.
func (m *MockPolicyKeeper) SetAction(ctx types1.Context, action *types.Action) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAction", ctx, action)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAction indicates an expected call of SetAction.
func (mr *MockPolicyKeeperMockRecorder) SetAction(ctx, action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAction", reflect.TypeOf((*MockPolicyKeeper)(nil).SetAction), ctx, action)
}

// Unpack mocks base method.
func (m *MockPolicyKeeper) Unpack(policyPb *types.Policy) (policy.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unpack", policyPb)
	ret0, _ := ret[0].(policy.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unpack indicates an expected call of Unpack.
func (mr *MockPolicyKeeperMockRecorder) Unpack(policyPb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpack", reflect.TypeOf((*MockPolicyKeeper)(nil).Unpack), policyPb)
}

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// SendCoinsFromAccountToModule mocks base method.
func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr types1.AccAddress, recipientModule string, amt types1.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromAccountToModule", ctx, senderAddr, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromAccountToModule indicates an expected call of SendCoinsFromAccountToModule.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromAccountToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromAccountToModule), ctx, senderAddr, recipientModule, amt)
}

// SpendableCoins mocks base method.
func (m *MockBankKeeper) SpendableCoins(arg0 context.Context, arg1 types1.AccAddress) types1.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendableCoins", arg0, arg1)
	ret0, _ := ret[0].(types1.Coins)
	return ret0
}

// SpendableCoins indicates an expected call of SpendableCoins.
func (mr *MockBankKeeperMockRecorder) SpendableCoins(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendableCoins", reflect.TypeOf((*MockBankKeeper)(nil).SpendableCoins), arg0, arg1)
}

// MockParamSubspace is a mock of ParamSubspace interface.
type MockParamSubspace struct {
	ctrl     *gomock.Controller
	recorder *MockParamSubspaceMockRecorder
}

// MockParamSubspaceMockRecorder is the mock recorder for MockParamSubspace.
type MockParamSubspaceMockRecorder struct {
	mock *MockParamSubspace
}

// NewMockParamSubspace creates a new mock instance.
func NewMockParamSubspace(ctrl *gomock.Controller) *MockParamSubspace {
	mock := &MockParamSubspace{ctrl: ctrl}
	mock.recorder = &MockParamSubspaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParamSubspace) EXPECT() *MockParamSubspaceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockParamSubspace) Get(arg0 context.Context, arg1 []byte, arg2 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Get", arg0, arg1, arg2)
}

// Get indicates an expected call of Get.
func (mr *MockParamSubspaceMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockParamSubspace)(nil).Get), arg0, arg1, arg2)
}

// Set mocks base method.
func (m *MockParamSubspace) Set(arg0 context.Context, arg1 []byte, arg2 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1, arg2)
}

// Set indicates an expected call of Set.
func (mr *MockParamSubspaceMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockParamSubspace)(nil).Set), arg0, arg1, arg2)
}
