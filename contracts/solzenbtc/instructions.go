package solzenbtc

import (
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc/generated/zenbtc_spl_token"
	"github.com/gagliardetto/solana-go"
)

func Initialize(
	programID solana.PublicKey,
	args zenbtc_spl_token.InitializeArgs,
	signer solana.PublicKey,
	mint solana.PublicKey,

) *zenbtc_spl_token.Instruction {
	zenbtc_spl_token.SetProgramID(programID)

	globalConfigPDA, _ := GetGlobalConfigPDA(programID)
	wrappedMetadataPDA, _ := GetMetadataPDA(mint)

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

	return instruction
}

func Wrap(
	programID solana.PublicKey,
	args zenbtc_spl_token.WrapArgs,
	signer solana.PublicKey,
	mint solana.PublicKey,
	multisigKey solana.PublicKey,
	feeWallet solana.PublicKey,
	feeWalletAta solana.PublicKey,
	receiver solana.PublicKey,
	receiverAta solana.PublicKey,
) *zenbtc_spl_token.Instruction {
	zenbtc_spl_token.SetProgramID(programID)

	globalConfigPDA, _ := GetGlobalConfigPDA(programID)

	instruction := zenbtc_spl_token.NewWrapInstruction(
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
	).Build()

	return instruction
}

func Unwrap(
	programID solana.PublicKey,
	args zenbtc_spl_token.UnwrapArgs,
	signer solana.PublicKey,
	mint solana.PublicKey,
	multisigKey solana.PublicKey,
	feeWallet solana.PublicKey,
) *zenbtc_spl_token.Instruction {
	zenbtc_spl_token.SetProgramID(programID)

	globalConfigPDA, _ := GetGlobalConfigPDA(programID)
	signerAta, _, _ := solana.FindAssociatedTokenAddress(signer, mint)
	feeWalletAta, _, _ := solana.FindAssociatedTokenAddress(feeWallet, mint)

	instruction := zenbtc_spl_token.NewUnwrapInstruction(
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
	).Build()

	return instruction
}
