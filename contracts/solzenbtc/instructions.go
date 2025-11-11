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

// WrapZenzec creates a wrap instruction for ZenZEC (uses zenbtc_wrap seed)
func WrapZenzec(
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

	wrapShard, err := GetEventStoreZenzecWrapShardPDA(eventStoreProgramID, eventID)
	if err != nil {
		return nil, err
	}

	return buildWrapInstruction(programID, eventStoreProgramID, args, signer, globalConfigPDA, multisigKey, mint, feeWallet, feeWalletAta, receiver, receiverAta, eventStoreGlobalConfig, wrapShard)
}

// WrapZenbtc2 creates a wrap instruction for ZenBTC (uses zenbtc2_wrap seed)
func WrapZenbtc2(
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

	wrapShard, err := GetEventStoreZenbtc2WrapShardPDA(eventStoreProgramID, eventID)
	if err != nil {
		return nil, err
	}

	return buildWrapInstruction(programID, eventStoreProgramID, args, signer, globalConfigPDA, multisigKey, mint, feeWallet, feeWalletAta, receiver, receiverAta, eventStoreGlobalConfig, wrapShard)
}

// WrapRock creates a wrap instruction for ROCK (uses rock_wrap seed)
func WrapRock(
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

	wrapShard, err := GetEventStoreRockWrapShardPDA(eventStoreProgramID, eventID)
	if err != nil {
		return nil, err
	}

	return buildWrapInstruction(programID, eventStoreProgramID, args, signer, globalConfigPDA, multisigKey, mint, feeWallet, feeWalletAta, receiver, receiverAta, eventStoreGlobalConfig, wrapShard)
}

// buildWrapInstruction is a helper that builds the actual wrap instruction
func buildWrapInstruction(
	programID solana.PublicKey,
	eventStoreProgramID solana.PublicKey,
	args zenbtc_spl_token.WrapArgs,
	signer solana.PublicKey,
	globalConfigPDA solana.PublicKey,
	multisigKey solana.PublicKey,
	mint solana.PublicKey,
	feeWallet solana.PublicKey,
	feeWalletAta solana.PublicKey,
	receiver solana.PublicKey,
	receiverAta solana.PublicKey,
	eventStoreGlobalConfig solana.PublicKey,
	wrapShard solana.PublicKey,
) (*zenbtc_spl_token.Instruction, error) {

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
		wrapShard,
	)

	if account := builder.GetGlobalConfigAccount(); account != nil {
		account.WRITE()
	}
	if account := builder.GetZenbtcWrapShardAccount(); account != nil {
		account.WRITE()
	}

	instruction, err := builder.ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	return instruction, nil
}

// UnwrapZenzec creates an unwrap instruction for ZenZEC (uses zenbtc_unwrap seed)
func UnwrapZenzec(
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

	unwrapShard, err := GetEventStoreZenzecUnwrapShardPDA(eventStoreProgramID, eventID)
	if err != nil {
		return nil, err
	}

	return buildUnwrapInstruction(programID, eventStoreProgramID, args, signer, globalConfigPDA, multisigKey, mint, feeWallet, eventStoreGlobalConfig, unwrapShard)
}

// UnwrapZenbtc2 creates an unwrap instruction for ZenBTC (uses zenbtc2_unwrap seed)
func UnwrapZenbtc2(
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

	unwrapShard, err := GetEventStoreZenbtc2UnwrapShardPDA(eventStoreProgramID, eventID)
	if err != nil {
		return nil, err
	}

	return buildUnwrapInstruction(programID, eventStoreProgramID, args, signer, globalConfigPDA, multisigKey, mint, feeWallet, eventStoreGlobalConfig, unwrapShard)
}

// UnwrapRock creates an unwrap instruction for ROCK (uses rock_unwrap seed)
func UnwrapRock(
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

	unwrapShard, err := GetEventStoreRockUnwrapShardPDA(eventStoreProgramID, eventID)
	if err != nil {
		return nil, err
	}

	return buildUnwrapInstruction(programID, eventStoreProgramID, args, signer, globalConfigPDA, multisigKey, mint, feeWallet, eventStoreGlobalConfig, unwrapShard)
}

// buildUnwrapInstruction is a helper that builds the actual unwrap instruction
func buildUnwrapInstruction(
	programID solana.PublicKey,
	eventStoreProgramID solana.PublicKey,
	args zenbtc_spl_token.UnwrapArgs,
	signer solana.PublicKey,
	globalConfigPDA solana.PublicKey,
	multisigKey solana.PublicKey,
	mint solana.PublicKey,
	feeWallet solana.PublicKey,
	eventStoreGlobalConfig solana.PublicKey,
	unwrapShard solana.PublicKey,
) (*zenbtc_spl_token.Instruction, error) {

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
		unwrapShard,
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
