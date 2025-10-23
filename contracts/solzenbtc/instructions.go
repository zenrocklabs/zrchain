package solzenbtc

import (
	"fmt"
	"math/big"

	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc/generated/zenbtc_spl_token"
	"github.com/gagliardetto/solana-go"
)

func Initialize(
	programID solana.PublicKey,
	args zenbtc_spl_token.InitializeArgs,
	signer solana.PublicKey,
	mint solana.PublicKey,
) (*zenbtc_spl_token.Instruction, error) {
	zenbtc_spl_token.SetProgramID(programID)

	globalConfigPDA, err := GetGlobalConfigPDA(programID)
	if err != nil {
		return nil, err
	}
	wrappedMetadataPDA, err := GetMetadataPDA(mint)
	if err != nil {
		return nil, err
	}

	instruction := zenbtc_spl_token.NewInitializeInstruction(
		args,
		signer,
		globalConfigPDA,
		mint,
		solana.SystemProgramID,
		solana.TokenProgramID,
		wrappedMetadataPDA,
		solana.TokenMetadataProgramID,
		solana.SysVarRentPubkey,
	).Build()

	return instruction, nil
}

func Wrap(
	programID solana.PublicKey,
	eventStoreProgramID solana.PublicKey,
	eventID *big.Int,
	args zenbtc_spl_token.WrapArgs,
	signer solana.PublicKey,
	mint solana.PublicKey,
	multisigKey solana.PublicKey,
	feeWallet solana.PublicKey,
	feeWalletAta solana.PublicKey,
	receiver solana.PublicKey,
	receiverAta solana.PublicKey,
) (*zenbtc_spl_token.Instruction, error) {
	zenbtc_spl_token.SetProgramID(programID)

	if eventID == nil {
		return nil, fmt.Errorf("eventID is required")
	}

	globalConfigPDA, err := GetGlobalConfigPDA(programID)
	if err != nil {
		return nil, err
	}

	eventStoreGlobalConfig, err := GetEventStoreGlobalConfigPDA(eventStoreProgramID)
	if err != nil {
		return nil, err
	}

	zenbtcWrapShard, err := GetEventStoreZenbtcWrapShardPDA(eventStoreProgramID, eventID)
	if err != nil {
		return nil, err
	}

	builder := zenbtc_spl_token.NewWrapInstruction(
		args,
		signer,
		globalConfigPDA,
		multisigKey,
		mint,
		feeWallet,
		feeWalletAta,
		receiver,
		receiverAta,
		solana.SystemProgramID,
		solana.TokenProgramID,
		solana.SPLAssociatedTokenAccountProgramID,
		eventStoreProgramID,
		eventStoreGlobalConfig,
		programID,
		zenbtcWrapShard,
	)

	if account := builder.GetZenbtcWrapShardAccount(); account != nil {
		account.WRITE()
	}

	instruction, err := builder.ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	return instruction, nil
}

func Unwrap(
	programID solana.PublicKey,
	eventStoreProgramID solana.PublicKey,
	eventID *big.Int,
	args zenbtc_spl_token.UnwrapArgs,
	signer solana.PublicKey,
	mint solana.PublicKey,
	multisigKey solana.PublicKey,
	feeWallet solana.PublicKey,
) (*zenbtc_spl_token.Instruction, error) {
	zenbtc_spl_token.SetProgramID(programID)

	if eventID == nil {
		return nil, fmt.Errorf("eventID is required")
	}

	globalConfigPDA, err := GetGlobalConfigPDA(programID)
	if err != nil {
		return nil, err
	}
	eventStoreGlobalConfig, err := GetEventStoreGlobalConfigPDA(eventStoreProgramID)
	if err != nil {
		return nil, err
	}

	zenbtcUnwrapShard, err := GetEventStoreZenbtcUnwrapShardPDA(eventStoreProgramID, eventID)
	if err != nil {
		return nil, err
	}

	signerAta, _, _ := solana.FindAssociatedTokenAddress(signer, mint)
	feeWalletAta, _, _ := solana.FindAssociatedTokenAddress(feeWallet, mint)

	builder := zenbtc_spl_token.NewUnwrapInstruction(
		args,
		signer,
		globalConfigPDA,
		multisigKey,
		mint,
		signerAta,
		feeWallet,
		feeWalletAta,
		solana.SystemProgramID,
		solana.TokenProgramID,
		solana.SPLAssociatedTokenAccountProgramID,
		eventStoreProgramID,
		eventStoreGlobalConfig,
		programID,
		zenbtcUnwrapShard,
	)

	if account := builder.GetZenbtcUnwrapShardAccount(); account != nil {
		account.WRITE()
	}

	instruction, err := builder.ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	return instruction, nil
}
