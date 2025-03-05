package solrock

import "github.com/gagliardetto/solana-go"

func GetGlobalConfigPDA(programID solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("global_config"),
	}
	addr, _, err := solana.FindProgramAddress(seeds, programID)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

func GetTokenConfigPDA(programID solana.PublicKey, mintAddress solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("token_config"),
		mintAddress.Bytes(),
	}
	addr, _, err := solana.FindProgramAddress(seeds, programID)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

func GetMintAddress(programID solana.PublicKey, tokenSymbol string) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("wrapped_mint"),
		[]byte(tokenSymbol),
	}
	addr, _, err := solana.FindProgramAddress(seeds, programID)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

func GetMetadataPDA(mint solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("metadata"),
		solana.TokenMetadataProgramID.Bytes(),
		mint.Bytes(),
	}
	addr, _, err := solana.FindProgramAddress(seeds, solana.TokenMetadataProgramID)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

func GetWhitelistedWalletPDA(programID solana.PublicKey, receiver solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("whitelisted_wallet"),
		receiver.Bytes(),
	}
	addr, _, err := solana.FindProgramAddress(seeds, programID)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}
