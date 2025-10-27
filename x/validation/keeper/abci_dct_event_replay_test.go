package keeper_test

import (
	"context"
	"encoding/base64"
	"testing"

	"cosmossdk.io/collections"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationkeeper "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type sidecarMint struct {
	sigHash   string
	recipient string
	mint      string
	txSig     string
	value     uint64
}

func decodeSidecarMint(t *testing.T, raw sidecarMint) sidecarapitypes.SolanaMintEvent {
	t.Helper()

	sigHash, err := base64.StdEncoding.DecodeString(raw.sigHash)
	require.NoError(t, err)
	recipient, err := base64.StdEncoding.DecodeString(raw.recipient)
	require.NoError(t, err)
	mint, err := base64.StdEncoding.DecodeString(raw.mint)
	require.NoError(t, err)
	return sidecarapitypes.SolanaMintEvent{
		Coint:     sidecarapitypes.Coin_ZENZEC,
		SigHash:   sigHash,
		Recipient: recipient,
		Value:     raw.value,
		Mint:      mint,
		TxSig:     raw.txSig,
	}
}

func TestProcessSolanaDCTMintEvents_SidecarSequentialEvents(t *testing.T) {
	k, ctx, dctKeeper, treasuryKeeper, ctrl := setupSolanaDCTMintKeeper(t)
	t.Cleanup(ctrl.Finish)

	asset := dcttypes.Asset_ASSET_ZENZEC
	amount := uint64(30_000_000)

	rawEvents := []sidecarMint{
		{"AQAAAAAAAAAAAAAAAAAAAA==", "FTYdeLbyT51DCYQ9Fte5DiqrSkpQsPGeedpcuI4IsEE=", "Jb+3Zaz9AGEurTPMi20cPXN8TLwB+duKzvQ54/oftBc=", "01000000000000000000000000000000", amount},
		{"AgAAAAAAAAAAAAAAAAAAAA==", "L9Yz0VLzXC6okhPDrwowdbezh3WxjTvUbqh5ltN8AA0=", "Jb+3Zaz9AGEurTPMi20cPXN8TLwB+duKzvQ54/oftBc=", "02000000000000000000000000000000", amount},
		{"AwAAAAAAAAAAAAAAAAAAAA==", "FTYdeLbyT51DCYQ9Fte5DiqrSkpQsPGeedpcuI4IsEE=", "Jb+3Zaz9AGEurTPMi20cPXN8TLwB+duKzvQ54/oftBc=", "03000000000000000000000000000000", amount},
		{"BAAAAAAAAAAAAAAAAAAAAA==", "FTYdeLbyT51DCYQ9Fte5DiqrSkpQsPGeedpcuI4IsEE=", "Jb+3Zaz9AGEurTPMi20cPXN8TLwB+duKzvQ54/oftBc=", "04000000000000000000000000000000", amount},
		{"BQAAAAAAAAAAAAAAAAAAAA==", "FTYdeLbyT51DCYQ9Fte5DiqrSkpQsPGeedpcuI4IsEE=", "Jb+3Zaz9AGEurTPMi20cPXN8TLwB+duKzvQ54/oftBc=", "05000000000000000000000000000000", amount},
		{"BgAAAAAAAAAAAAAAAAAAAA==", "FTYdeLbyT51DCYQ9Fte5DiqrSkpQsPGeedpcuI4IsEE=", "Jb+3Zaz9AGEurTPMi20cPXN8TLwB+duKzvQ54/oftBc=", "06000000000000000000000000000000", amount},
	}

	events := make([]sidecarapitypes.SolanaMintEvent, len(rawEvents))
	for i, raw := range rawEvents {
		events[i] = decodeSidecarMint(t, raw)
	}

	require.NoError(t, k.SolanaCounters.Set(ctx, asset.String(), validationtypes.SolanaCounters{MintCounter: 4}))
	for i := 0; i < 4; i++ {
		recordProcessedEvent(t, k, ctx, asset, events[i])
	}

	targetEvent := events[4]
	recipientPub := solana.PublicKeyFromBytes(targetEvent.Recipient)
	mintPub := solana.PublicKeyFromBytes(targetEvent.Mint)

	pending := dcttypes.PendingMintTransaction{
		Id:               22,
		Asset:            asset,
		Amount:           amount,
		RecipientAddress: recipientPub.String(),
		ZrchainTxId:      20741,
		Status:           dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
	}

	signTxReq := &treasurytypes.SignTransactionRequest{SignRequestId: 9001}
	mainSignReq := &treasurytypes.SignRequest{ChildReqIds: []uint64{9101}}
	signedChild := &treasurytypes.SignRequest{SignedData: []*treasurytypes.SignedDataWithID{{SignedData: []byte{0xAA}}}}

	dctKeeper.EXPECT().ListSupportedAssets(gomock.Any()).Return([]dcttypes.Asset{asset}, nil)
	dctKeeper.EXPECT().GetSolanaParams(gomock.Any(), asset).Return(&dcttypes.Solana{
		MintAddress: mintPub.String(),
	}, nil)
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
		require.Equal(t, targetEvent.TxSig, tx.TxHash)
		return nil
	})
	dctKeeper.EXPECT().SetFirstPendingSolMintTransaction(gomock.Any(), asset, uint64(0)).Return(nil)
	dctKeeper.EXPECT().GetPendingMintTransaction(gomock.Any(), asset, pending.Id+1).Return(dcttypes.PendingMintTransaction{}, collections.ErrNotFound)

	treasuryKeeper.EXPECT().GetSignTransactionRequest(gomock.Any(), pending.ZrchainTxId).Return(signTxReq, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), signTxReq.SignRequestId).Return(mainSignReq, nil)
	treasuryKeeper.EXPECT().GetSignRequest(gomock.Any(), uint64(9101)).Return(signedChild, nil)

	oracleData := validationkeeper.OracleData{SolanaMintEvents: events}
	k.ProcessSolanaDCTMintEventsTestHelper(ctx, oracleData)

	counters, err := k.SolanaCounters.Get(ctx, asset.String())
	require.NoError(t, err)
	require.Equal(t, uint64(5), counters.MintCounter)

	eventHash := base64.StdEncoding.EncodeToString(targetEvent.SigHash)
	key := collections.Join(asset.String(), eventHash)
	processed, err := k.ProcessedSolanaMintEvents.Get(ctx, key)
	require.NoError(t, err)
	require.True(t, processed)
}
