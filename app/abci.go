package app

// Ref: https://github.com/cosmos/cosmos-sdk/blob/c64d1010800d60677cc25e2fca5b3d8c37b683cc/baseapp/abci_utils.go

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	validation "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
)

func (app *ZenrockApp) ExtendVoteHandler(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
	return app.ValidationKeeper.ExtendVoteHandler(ctx, req)
}

func (app *ZenrockApp) VerifyVoteExtensionHandler(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	return app.ValidationKeeper.VerifyVoteExtensionHandler(ctx, req)
}

func (app *ZenrockApp) PrepareProposalHandler(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
	txSelector := newTxSelectorWithVoteExt()

	if err := app.prepareVoteExtension(ctx, req, txSelector); err != nil {
		return nil, err
	}

	if err := app.prepareTransactions(req, maxGasPerBlock(ctx), txSelector); err != nil {
		return nil, err
	}

	return &abci.ResponsePrepareProposal{Txs: txSelector.SelectedTxs}, nil
}

func (app *ZenrockApp) ProcessProposalHandler(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	if res, err := app.ValidationKeeper.ProcessProposal(ctx, req); res.Status == abci.ResponseProcessProposal_REJECT {
		return res, err
	}

	if err := app.processTransactions(ctx, req, maxGasPerBlock(ctx)); err != nil {
		return validation.REJECT_PROPOSAL, nil
	}

	return validation.ACCEPT_PROPOSAL, nil
}

func (app *ZenrockApp) prepareVoteExtension(ctx sdk.Context, req *abci.RequestPrepareProposal, txSelector *TxSelectorWithVoteExt) error {
	voteExtension, err := app.ValidationKeeper.PrepareProposal(ctx, req)
	if err != nil {
		return err
	}

	if voteExtension != nil {
		txSelector.insertVoteExtension(voteExtension)
	}

	return nil
}

func (app *ZenrockApp) prepareTransactions(req *abci.RequestPrepareProposal, maxGasPerBlock uint64, txSelector *TxSelectorWithVoteExt) error {
	maxBytesPerTx := uint64(req.MaxTxBytes)

	for _, txBytes := range req.Txs {
		transaction, err := app.TxDecode(txBytes)
		if err != nil {
			return err
		}

		txSelector.insertTransaction(maxBytesPerTx, maxGasPerBlock, transaction, txBytes)
		if txSelector.capacityReached(maxBytesPerTx, maxGasPerBlock) {
			break
		}
	}
	return nil
}

func (app *ZenrockApp) processTransactions(ctx sdk.Context, req *abci.RequestProcessProposal, maxGasPerBlock uint64) error {
	gasUsed := uint64(0)

	for _, tx := range extractTxsFromRequest(ctx, req, app.txConfig.TxDecoder()) {
		verifiedTx, err := app.ProcessProposalVerifyTx(tx)
		if err != nil || !isGasWithinLimit(verifiedTx, maxGasPerBlock, &gasUsed) {
			return err
		}
	}

	return nil
}

func extractTxsFromRequest(ctx sdk.Context, req *abci.RequestProcessProposal, txUnmarshaler sdk.TxDecoder) [][]byte {
	if req.Txs == nil || len(req.Txs) == 0 {
		return nil
	}

	if validation.VoteExtensionsEnabled(ctx) && validation.ContainsVoteExtension(req.Txs[0], txUnmarshaler) {
		return req.Txs[1:] // the first tx is the vote extension so we don't include it with the actual txs
	}

	return req.Txs
}
