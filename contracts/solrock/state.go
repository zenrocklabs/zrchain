package solrock

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/contracts/solrock/generated/zenbtc_spl_token"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetGlobalConfig(context context.Context, client *rpc.Client, programID solana.PublicKey) (zenbtc_spl_token.GlobalConfigAccount, error) {
	globalConfigPDA, err := GetGlobalConfigPDA(programID)
	if err != nil {
		return zenbtc_spl_token.GlobalConfigAccount{}, err
	}

	accountInfo, err := GetAccountInfo(context, client, globalConfigPDA)
	if err != nil {
		return zenbtc_spl_token.GlobalConfigAccount{}, err
	}

	data := accountInfo.Value.Data.GetBinary()

	globalConfig := new(zenbtc_spl_token.GlobalConfigAccount)
	decoder := bin.NewBorshDecoder(data)

	err = globalConfig.UnmarshalWithDecoder(decoder)
	if err != nil {
		return zenbtc_spl_token.GlobalConfigAccount{}, err
	}

	return *globalConfig, nil

}

func GetTokenConfig(context context.Context, client *rpc.Client, programID solana.PublicKey, mintPubkey solana.PublicKey) (zenbtc_spl_token.TokenConfigAccount, error) {
	tokenConfigPDA, err := GetTokenConfigPDA(programID, mintPubkey)
	if err != nil {
		return zenbtc_spl_token.TokenConfigAccount{}, err
	}

	accountInfo, err := GetAccountInfo(context, client, tokenConfigPDA)
	if err != nil {
		return zenbtc_spl_token.TokenConfigAccount{}, err
	}

	data := accountInfo.Value.Data.GetBinary()

	tokenConfig := new(zenbtc_spl_token.TokenConfigAccount)
	decoder := bin.NewBorshDecoder(data)

	err = tokenConfig.UnmarshalWithDecoder(decoder)
	if err != nil {
		return zenbtc_spl_token.TokenConfigAccount{}, err
	}

	return *tokenConfig, nil
}

func GetMint(context context.Context, client *rpc.Client, mintPubkey solana.PublicKey) (token.Mint, error) {
	accountInfo, err := GetAccountInfo(context, client, mintPubkey)
	if err != nil {
		return token.Mint{}, err
	}

	data := accountInfo.Value.Data.GetBinary()

	var mint token.Mint

	err = bin.NewBorshDecoder(data).Decode(&mint)
	if err != nil {
		return token.Mint{}, err
	}

	return mint, nil
}

func GetTokenAccount(context context.Context, client *rpc.Client, tokenAccountPubkey solana.PublicKey) (token.Account, error) {
	accountInfo, err := GetAccountInfo(context, client, tokenAccountPubkey)
	if err != nil {
		return token.Account{}, err
	}

	data := accountInfo.Value.Data.GetBinary()

	tokenAccount := new(token.Account)
	decoder := bin.NewBorshDecoder(data)

	err = tokenAccount.UnmarshalWithDecoder(decoder)
	if err != nil {
		return token.Account{}, err
	}

	return *tokenAccount, nil
}

func GetNonceAccount(context context.Context, client *rpc.Client, nonceAccountPubkey solana.PublicKey) (system.NonceAccount, error) {
	accountInfo, err := GetAccountInfo(context, client, nonceAccountPubkey)
	if err != nil {
		return system.NonceAccount{}, err
	}

	data := accountInfo.Value.Data.GetBinary()

	nonceAccount := new(system.NonceAccount)
	decoder := bin.NewBorshDecoder(data)

	err = nonceAccount.UnmarshalWithDecoder(decoder)
	if err != nil {
		return system.NonceAccount{}, err
	}

	return *nonceAccount, nil
}
