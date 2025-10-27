package keeper_test

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationkeeper "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	validationtestutil "github.com/Zenrock-Foundation/zrchain/v6/x/validation/testutil"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupSolanaDCTMintKeeper(t *testing.T) (*validationkeeper.Keeper, sdk.Context, *validationtestutil.MockDCTKeeper, *validationtestutil.MockTreasuryKeeper, *gomock.Controller) {
	t.Helper()

	ctrl := gomock.NewController(t)

	storeKey := storetypes.NewKVStoreKey(validationtypes.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewTestLogger(t), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewTestLogger(t))

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	mockAccountKeeper := validationtestutil.NewMockAccountKeeper(ctrl)
	mockBankKeeper := validationtestutil.NewMockBankKeeper(ctrl)
	mockTreasuryKeeper := validationtestutil.NewMockTreasuryKeeper(ctrl)
	mockDCTKeeper := validationtestutil.NewMockDCTKeeper(ctrl)
	mockZentpKeeper := validationtestutil.NewMockZentpKeeper(ctrl)
	mockSlashingKeeper := validationtestutil.NewMockSlashingKeeper(ctrl)

	mockAccountKeeper.EXPECT().GetModuleAddress(validationtypes.BondedPoolName).Return(sdk.AccAddress{}).AnyTimes()
	mockAccountKeeper.EXPECT().GetModuleAddress(validationtypes.NotBondedPoolName).Return(sdk.AccAddress{}).AnyTimes()
	mockAccountKeeper.EXPECT().AddressCodec().Return(addresscodec.NewBech32Codec("zen")).AnyTimes()

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	mockAccountKeeper.EXPECT().GetModuleAddress(govtypes.ModuleName).Return(authority).AnyTimes()

	k := validationkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		mockAccountKeeper,
		mockBankKeeper,
		authority.String(),
		nil,
		nil,
		mockTreasuryKeeper,
		nil,
		mockDCTKeeper,
		mockZentpKeeper,
		mockSlashingKeeper,
		addresscodec.NewBech32Codec("zenvaloper"),
		addresscodec.NewBech32Codec("zenvalcons"),
	)
	require.NotNil(t, k)

	return k, ctx, mockDCTKeeper, mockTreasuryKeeper, ctrl
}

func newMintEvent(id uint64, coin sidecarapitypes.Coin, recipient, mint solana.PublicKey, amount uint64) sidecarapitypes.SolanaMintEvent {
	sigBytes := make([]byte, 16)
	binary.LittleEndian.PutUint64(sigBytes, id)
	return sidecarapitypes.SolanaMintEvent{
		Coint:     coin,
		SigHash:   sigBytes,
		Recipient: recipient[:],
		Value:     amount,
		Mint:      mint[:],
		TxSig:     hex.EncodeToString(sigBytes),
	}
}

func recordProcessedEvent(t *testing.T, k *validationkeeper.Keeper, ctx sdk.Context, asset dcttypes.Asset, event sidecarapitypes.SolanaMintEvent) {
	t.Helper()

	eventHash := base64.StdEncoding.EncodeToString(event.SigHash)
	key := collections.Join(asset.String(), eventHash)
	require.NoError(t, k.ProcessedSolanaMintEvents.Set(ctx, key, true))
}

func TestProcessSolanaDCTMintEvents_SkipsAlreadyProcessedEvent(t *testing.T) {
	k, ctx, dctKeeper, treasuryKeeper, ctrl := setupSolanaDCTMintKeeper(t)
	t.Cleanup(ctrl.Finish)

	asset := dcttypes.Asset_ASSET_ZENZEC
	recipient := solana.MustPublicKeyFromBase58("2RoRSPwFmMcFDd3DLgoomMRFkhUSaVNSDo7XyH2tw6de")
	mint := solana.MustPublicKeyFromBase58("3YMe7Bbus2rZiDR7ijRBhT6hNFvNwFnBgGEwxkw3L71g")
	amount := uint64(30_000_000)

	pending := dcttypes.PendingMintTransaction{
		Id:               1,
		Asset:            asset,
		Amount:           amount,
		RecipientAddress: recipient.String(),
		ZrchainTxId:      42,
		Status:           dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
	}

	event1 := newMintEvent(1, sidecarapitypes.Coin_ZENZEC, recipient, mint, amount)

	require.NoError(t, k.SolanaCounters.Set(ctx, asset.String(), validationtypes.SolanaCounters{MintCounter: 1}))
	recordProcessedEvent(t, k, ctx, asset, event1)

	signTxReq := &treasurytypes.SignTransactionRequest{SignRequestId: 100}
	mainSignReq := &treasurytypes.SignRequest{ChildReqIds: []uint64{101}}
	signedChild := &treasurytypes.SignRequest{SignedData: []*treasurytypes.SignedDataWithID{{SignedData: []byte{0x01}}}}

	dctKeeper.EXPECT().ListSupportedAssets(gomock.Any()).Return([]dcttypes.Asset{asset}, nil)
	dctKeeper.EXPECT().GetSolanaParams(gomock.Any(), asset).Return(&dcttypes.Solana{MintAddress: mint.String()}, nil)
	dctKeeper.EXPECT().GetFirstPendingSolMintTransaction(gomock.Any(), asset).Return(pending.Id, nil)
	dctKeeper.EXPECT().GetPendingMintTransaction(gomock.Any(), asset, pending.Id).Return(pending, nil)
	dctKeeper.EXPECT().GetSupply(gomock.Any(), asset).Times(0)
	dctKeeper.EXPECT().SetSupply(gomock.Any(), gomock.Any()).Times(0)
	dctKeeper.EXPECT().SetPendingMintTransaction(gomock.Any(), gomock.Any()).Times(0)
	dctKeeper.EXPECT().SetFirstPendingSolMintTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	treasuryKeeper.EXPECT().GetSignTransactionRequest(gomock.Any(), pending.ZrchainTxId).Return(signTxReq, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), signTxReq.SignRequestId).Return(mainSignReq, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), uint64(101)).Return(signedChild, nil)

	oracleData := validationkeeper.OracleData{SolanaMintEvents: []sidecarapitypes.SolanaMintEvent{event1}}
	k.ProcessSolanaDCTMintEventsTestHelper(ctx, oracleData)

	counters, err := k.SolanaCounters.Get(ctx, asset.String())
	require.NoError(t, err)
	require.Equal(t, uint64(1), counters.MintCounter)

	eventHash := base64.StdEncoding.EncodeToString(event1.SigHash)
	key := collections.Join(asset.String(), eventHash)
	processed, err := k.ProcessedSolanaMintEvents.Get(ctx, key)
	require.NoError(t, err)
	require.True(t, processed)
}

func TestProcessSolanaDCTMintEvents_ProcessesNextEventID(t *testing.T) {
	k, ctx, dctKeeper, treasuryKeeper, ctrl := setupSolanaDCTMintKeeper(t)
	t.Cleanup(ctrl.Finish)

	asset := dcttypes.Asset_ASSET_ZENZEC
	recipient := solana.MustPublicKeyFromBase58("2RoRSPwFmMcFDd3DLgoomMRFkhUSaVNSDo7XyH2tw6de")
	mint := solana.MustPublicKeyFromBase58("3YMe7Bbus2rZiDR7ijRBhT6hNFvNwFnBgGEwxkw3L71g")
	amount := uint64(30_000_000)

	pending := dcttypes.PendingMintTransaction{
		Id:               2,
		Asset:            asset,
		Amount:           amount,
		RecipientAddress: recipient.String(),
		ZrchainTxId:      4242,
		Status:           dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
	}

	event1 := newMintEvent(1, sidecarapitypes.Coin_ZENZEC, recipient, mint, amount)
	event2 := newMintEvent(2, sidecarapitypes.Coin_ZENZEC, recipient, mint, amount)

	require.NoError(t, k.SolanaCounters.Set(ctx, asset.String(), validationtypes.SolanaCounters{MintCounter: 1}))
	recordProcessedEvent(t, k, ctx, asset, event1)

	signTxReq := &treasurytypes.SignTransactionRequest{SignRequestId: 200}
	mainSignReq := &treasurytypes.SignRequest{ChildReqIds: []uint64{301, 302}}
	signedChild := &treasurytypes.SignRequest{SignedData: []*treasurytypes.SignedDataWithID{{SignedData: []byte{0xAA}}}}

	dctKeeper.EXPECT().ListSupportedAssets(gomock.Any()).Return([]dcttypes.Asset{asset}, nil)
	dctKeeper.EXPECT().GetSolanaParams(gomock.Any(), asset).Return(&dcttypes.Solana{MintAddress: mint.String()}, nil)
	dctKeeper.EXPECT().GetFirstPendingSolMintTransaction(gomock.Any(), asset).Return(pending.Id, nil)
	dctKeeper.EXPECT().GetPendingMintTransaction(gomock.Any(), asset, pending.Id).Return(pending, nil)
	dctKeeper.EXPECT().GetSupply(gomock.Any(), asset).Return(dcttypes.Supply{Asset: asset, PendingAmount: amount}, nil)
	dctKeeper.EXPECT().SetSupply(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, supply dcttypes.Supply) error {
		require.Equal(t, asset, supply.Asset)
		require.Equal(t, uint64(0), supply.PendingAmount)
		require.Equal(t, amount, supply.MintedAmount)
		return nil
	})
	dctKeeper.EXPECT().SetPendingMintTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, tx dcttypes.PendingMintTransaction) error {
		require.Equal(t, asset, tx.Asset)
		require.Equal(t, dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED, tx.Status)
		require.Equal(t, event2.TxSig, tx.TxHash)
		return nil
	})
	dctKeeper.EXPECT().SetFirstPendingSolMintTransaction(gomock.Any(), asset, uint64(0)).Return(nil)
	dctKeeper.EXPECT().GetPendingMintTransaction(gomock.Any(), asset, uint64(3)).Return(dcttypes.PendingMintTransaction{}, collections.ErrNotFound)

	treasuryKeeper.EXPECT().GetSignTransactionRequest(gomock.Any(), pending.ZrchainTxId).Return(signTxReq, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), signTxReq.SignRequestId).Return(mainSignReq, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), uint64(301)).Return(signedChild, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), uint64(302)).Return(signedChild, nil)

	oracleData := validationkeeper.OracleData{SolanaMintEvents: []sidecarapitypes.SolanaMintEvent{event2}}
	k.ProcessSolanaDCTMintEventsTestHelper(ctx, oracleData)

	counters, err := k.SolanaCounters.Get(ctx, asset.String())
	require.NoError(t, err)
	require.Equal(t, uint64(2), counters.MintCounter)

	eventHash := base64.StdEncoding.EncodeToString(event2.SigHash)
	key := collections.Join(asset.String(), eventHash)
	processed, err := k.ProcessedSolanaMintEvents.Get(ctx, key)
	require.NoError(t, err)
	require.True(t, processed)
}

func TestProcessSolanaDCTMintEvents_SkipsUnexpectedEventID(t *testing.T) {
	k, ctx, dctKeeper, treasuryKeeper, ctrl := setupSolanaDCTMintKeeper(t)
	t.Cleanup(ctrl.Finish)

	asset := dcttypes.Asset_ASSET_ZENZEC
	recipient := solana.MustPublicKeyFromBase58("2RoRSPwFmMcFDd3DLgoomMRFkhUSaVNSDo7XyH2tw6de")
	mint := solana.MustPublicKeyFromBase58("3YMe7Bbus2rZiDR7ijRBhT6hNFvNwFnBgGEwxkw3L71g")
	amount := uint64(30_000_000)

	pending := dcttypes.PendingMintTransaction{
		Id:               3,
		Asset:            asset,
		Amount:           amount,
		RecipientAddress: recipient.String(),
		ZrchainTxId:      5151,
		Status:           dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
	}

	eventOld := newMintEvent(1, sidecarapitypes.Coin_ZENZEC, recipient, mint, amount)

	require.NoError(t, k.SolanaCounters.Set(ctx, asset.String(), validationtypes.SolanaCounters{MintCounter: 2}))

	signTxReq := &treasurytypes.SignTransactionRequest{SignRequestId: 400}
	mainSignReq := &treasurytypes.SignRequest{ChildReqIds: []uint64{501}}
	signedChild := &treasurytypes.SignRequest{SignedData: []*treasurytypes.SignedDataWithID{{SignedData: []byte{0xBB}}}}

	dctKeeper.EXPECT().ListSupportedAssets(gomock.Any()).Return([]dcttypes.Asset{asset}, nil)
	dctKeeper.EXPECT().GetSolanaParams(gomock.Any(), asset).Return(&dcttypes.Solana{MintAddress: mint.String()}, nil)
	dctKeeper.EXPECT().GetFirstPendingSolMintTransaction(gomock.Any(), asset).Return(pending.Id, nil)
	dctKeeper.EXPECT().GetPendingMintTransaction(gomock.Any(), asset, pending.Id).Return(pending, nil)
	dctKeeper.EXPECT().GetSupply(gomock.Any(), asset).Times(0)
	dctKeeper.EXPECT().SetSupply(gomock.Any(), gomock.Any()).Times(0)
	dctKeeper.EXPECT().SetPendingMintTransaction(gomock.Any(), gomock.Any()).Times(0)
	dctKeeper.EXPECT().SetFirstPendingSolMintTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	treasuryKeeper.EXPECT().GetSignTransactionRequest(gomock.Any(), pending.ZrchainTxId).Return(signTxReq, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), signTxReq.SignRequestId).Return(mainSignReq, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), uint64(501)).Return(signedChild, nil)

	oracleData := validationkeeper.OracleData{SolanaMintEvents: []sidecarapitypes.SolanaMintEvent{eventOld}}
	k.ProcessSolanaDCTMintEventsTestHelper(ctx, oracleData)

	counters, err := k.SolanaCounters.Get(ctx, asset.String())
	require.NoError(t, err)
	require.Equal(t, uint64(2), counters.MintCounter)

	eventHash := base64.StdEncoding.EncodeToString(eventOld.SigHash)
	key := collections.Join(asset.String(), eventHash)
	_, err = k.ProcessedSolanaMintEvents.Get(ctx, key)
	require.ErrorIs(t, err, collections.ErrNotFound)
}
