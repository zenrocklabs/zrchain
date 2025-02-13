package solrock

import (
	"github.com/Zenrock-Foundation/zrchain/v5/contracts/solrock/generated/zenbtc_spl_token"
	"github.com/gagliardetto/solana-go"
)

func DeployToken(
	args zenbtc_spl_token.DeployTokenArgs,
	programID solana.PublicKey,
	mint solana.PublicKey,
	signer solana.PublicKey,
) *zenbtc_spl_token.Instruction {
	zenbtc_spl_token.SetProgramID(programID)

	globalConfigPDA, _ := GetGlobalConfigPDA(programID)
	tokenConfigPDA, _ := GetTokenConfigPDA(programID, mint)
	metadataPDA, _ := GetMetadataPDA(mint)

	instruction := zenbtc_spl_token.NewDeployTokenInstruction(
		args,
		signer,
		globalConfigPDA,
		tokenConfigPDA,
		mint,
		solana.SystemProgramID,
		solana.TokenProgramID,
		metadataPDA,
		solana.TokenMetadataProgramID,
		solana.SysVarRentPubkey,
	).Build()

	return instruction
}

func Initialize(
	args zenbtc_spl_token.InitializeArgs,
	programID solana.PublicKey,
	signer solana.PublicKey,
) *zenbtc_spl_token.Instruction {
	zenbtc_spl_token.SetProgramID(programID)

	globalConfigPDA, _ := GetGlobalConfigPDA(programID)

	instruction := zenbtc_spl_token.NewInitializeInstruction(
		args,
		signer,
		globalConfigPDA,
		solana.SystemProgramID,
	).Build()

	return instruction
}

func Wrap(
	args zenbtc_spl_token.WrapArgs,
	programID solana.PublicKey,
	mint solana.PublicKey,
	signer solana.PublicKey,
	feeWallet solana.PublicKey,
	feeWalletAta solana.PublicKey,
	receiver solana.PublicKey,
	receiverAta solana.PublicKey,
) *zenbtc_spl_token.Instruction {
	zenbtc_spl_token.SetProgramID(programID)

	globalConfigPDA, _ := GetGlobalConfigPDA(programID)
	tokenConfigPDA, _ := GetTokenConfigPDA(programID, mint)
	whitelistedWalletPDA, _ := GetWhitelistedWalletPDA(programID, receiver)

	instruction := zenbtc_spl_token.NewWrapInstruction(
		args,
		signer,
		globalConfigPDA,
		tokenConfigPDA,
		mint,
		feeWallet,
		feeWalletAta,
		receiver,
		receiverAta,
		whitelistedWalletPDA,
		solana.SystemProgramID,
		solana.TokenProgramID,
		solana.SPLAssociatedTokenAccountProgramID,
	).Build()

	return instruction
}

func Unwrap(
	args zenbtc_spl_token.UnwrapArgs,
	programID solana.PublicKey,
	mint solana.PublicKey,
	signer solana.PublicKey,
	feeWallet solana.PublicKey,

) *zenbtc_spl_token.Instruction {
	zenbtc_spl_token.SetProgramID(programID)

	globalConfigPDA, _ := GetGlobalConfigPDA(programID)
	tokenConfigPDA, _ := GetTokenConfigPDA(programID, mint)
	whitelistedWalletPDA, _ := GetWhitelistedWalletPDA(programID, signer)
	signerAta, _, _ := solana.FindAssociatedTokenAddress(signer, mint)
	feeWalletAta, _, _ := solana.FindAssociatedTokenAddress(feeWallet, mint)

	instruction := zenbtc_spl_token.NewUnwrapInstruction(
		args,
		signer,
		globalConfigPDA,
		tokenConfigPDA,
		mint,
		signerAta,
		whitelistedWalletPDA,
		feeWallet,
		feeWalletAta,
		solana.SystemProgramID,
		solana.TokenProgramID,
		solana.SPLAssociatedTokenAccountProgramID,
	).Build()

	return instruction
}
