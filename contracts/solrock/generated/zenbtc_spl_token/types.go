// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package zenbtc_spl_token

import (
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

type AddFeeAuthorityArgs struct {
	NewFeeAuthority ag_solanago.PublicKey
}

func (obj AddFeeAuthorityArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `NewFeeAuthority` param:
	err = encoder.Encode(obj.NewFeeAuthority)
	if err != nil {
		return err
	}
	return nil
}

func (obj *AddFeeAuthorityArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `NewFeeAuthority`:
	err = decoder.Decode(&obj.NewFeeAuthority)
	if err != nil {
		return err
	}
	return nil
}

type AddMintAuthorityArgs struct {
	NewMintAuthority ag_solanago.PublicKey
}

func (obj AddMintAuthorityArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `NewMintAuthority` param:
	err = encoder.Encode(obj.NewMintAuthority)
	if err != nil {
		return err
	}
	return nil
}

func (obj *AddMintAuthorityArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `NewMintAuthority`:
	err = decoder.Decode(&obj.NewMintAuthority)
	if err != nil {
		return err
	}
	return nil
}

type DeployTokenArgs struct {
	MintAuthorities []ag_solanago.PublicKey
	FeeAuthorities  []ag_solanago.PublicKey
	FeeWallet       ag_solanago.PublicKey
	BurnFeeBps      uint64
	TokenName       string
	TokenSymbol     string
	TokenDecimals   uint8
	TokenUri        string
}

func (obj DeployTokenArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `MintAuthorities` param:
	err = encoder.Encode(obj.MintAuthorities)
	if err != nil {
		return err
	}
	// Serialize `FeeAuthorities` param:
	err = encoder.Encode(obj.FeeAuthorities)
	if err != nil {
		return err
	}
	// Serialize `FeeWallet` param:
	err = encoder.Encode(obj.FeeWallet)
	if err != nil {
		return err
	}
	// Serialize `BurnFeeBps` param:
	err = encoder.Encode(obj.BurnFeeBps)
	if err != nil {
		return err
	}
	// Serialize `TokenName` param:
	err = encoder.Encode(obj.TokenName)
	if err != nil {
		return err
	}
	// Serialize `TokenSymbol` param:
	err = encoder.Encode(obj.TokenSymbol)
	if err != nil {
		return err
	}
	// Serialize `TokenDecimals` param:
	err = encoder.Encode(obj.TokenDecimals)
	if err != nil {
		return err
	}
	// Serialize `TokenUri` param:
	err = encoder.Encode(obj.TokenUri)
	if err != nil {
		return err
	}
	return nil
}

func (obj *DeployTokenArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `MintAuthorities`:
	err = decoder.Decode(&obj.MintAuthorities)
	if err != nil {
		return err
	}
	// Deserialize `FeeAuthorities`:
	err = decoder.Decode(&obj.FeeAuthorities)
	if err != nil {
		return err
	}
	// Deserialize `FeeWallet`:
	err = decoder.Decode(&obj.FeeWallet)
	if err != nil {
		return err
	}
	// Deserialize `BurnFeeBps`:
	err = decoder.Decode(&obj.BurnFeeBps)
	if err != nil {
		return err
	}
	// Deserialize `TokenName`:
	err = decoder.Decode(&obj.TokenName)
	if err != nil {
		return err
	}
	// Deserialize `TokenSymbol`:
	err = decoder.Decode(&obj.TokenSymbol)
	if err != nil {
		return err
	}
	// Deserialize `TokenDecimals`:
	err = decoder.Decode(&obj.TokenDecimals)
	if err != nil {
		return err
	}
	// Deserialize `TokenUri`:
	err = decoder.Decode(&obj.TokenUri)
	if err != nil {
		return err
	}
	return nil
}

type GlobalConfig struct {
	// The authority which can update the config
	GlobalAuthority ag_solanago.PublicKey

	// The bump
	Bump uint8
}

func (obj GlobalConfig) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `GlobalAuthority` param:
	err = encoder.Encode(obj.GlobalAuthority)
	if err != nil {
		return err
	}
	// Serialize `Bump` param:
	err = encoder.Encode(obj.Bump)
	if err != nil {
		return err
	}
	return nil
}

func (obj *GlobalConfig) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `GlobalAuthority`:
	err = decoder.Decode(&obj.GlobalAuthority)
	if err != nil {
		return err
	}
	// Deserialize `Bump`:
	err = decoder.Decode(&obj.Bump)
	if err != nil {
		return err
	}
	return nil
}

type InitializeArgs struct {
	GlobalAuthority ag_solanago.PublicKey
}

func (obj InitializeArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `GlobalAuthority` param:
	err = encoder.Encode(obj.GlobalAuthority)
	if err != nil {
		return err
	}
	return nil
}

func (obj *InitializeArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `GlobalAuthority`:
	err = decoder.Decode(&obj.GlobalAuthority)
	if err != nil {
		return err
	}
	return nil
}

type RemoveFeeAuthorityArgs struct {
	FeeAuthorityToRemove ag_solanago.PublicKey
}

func (obj RemoveFeeAuthorityArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `FeeAuthorityToRemove` param:
	err = encoder.Encode(obj.FeeAuthorityToRemove)
	if err != nil {
		return err
	}
	return nil
}

func (obj *RemoveFeeAuthorityArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `FeeAuthorityToRemove`:
	err = decoder.Decode(&obj.FeeAuthorityToRemove)
	if err != nil {
		return err
	}
	return nil
}

type RemoveMintAuthorityArgs struct {
	MintAuthorityToRemove ag_solanago.PublicKey
}

func (obj RemoveMintAuthorityArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `MintAuthorityToRemove` param:
	err = encoder.Encode(obj.MintAuthorityToRemove)
	if err != nil {
		return err
	}
	return nil
}

func (obj *RemoveMintAuthorityArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `MintAuthorityToRemove`:
	err = decoder.Decode(&obj.MintAuthorityToRemove)
	if err != nil {
		return err
	}
	return nil
}

type TokenConfig struct {
	MintAuthorities   []ag_solanago.PublicKey
	FeeAuthorities    []ag_solanago.PublicKey
	FeeWallet         ag_solanago.PublicKey
	BurnFeeBps        uint64
	MintCounter       ag_binary.Uint128
	RedemptionCounter ag_binary.Uint128

	// The bump
	Bump uint8
}

func (obj TokenConfig) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `MintAuthorities` param:
	err = encoder.Encode(obj.MintAuthorities)
	if err != nil {
		return err
	}
	// Serialize `FeeAuthorities` param:
	err = encoder.Encode(obj.FeeAuthorities)
	if err != nil {
		return err
	}
	// Serialize `FeeWallet` param:
	err = encoder.Encode(obj.FeeWallet)
	if err != nil {
		return err
	}
	// Serialize `BurnFeeBps` param:
	err = encoder.Encode(obj.BurnFeeBps)
	if err != nil {
		return err
	}
	// Serialize `MintCounter` param:
	err = encoder.Encode(obj.MintCounter)
	if err != nil {
		return err
	}
	// Serialize `RedemptionCounter` param:
	err = encoder.Encode(obj.RedemptionCounter)
	if err != nil {
		return err
	}
	// Serialize `Bump` param:
	err = encoder.Encode(obj.Bump)
	if err != nil {
		return err
	}
	return nil
}

func (obj *TokenConfig) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `MintAuthorities`:
	err = decoder.Decode(&obj.MintAuthorities)
	if err != nil {
		return err
	}
	// Deserialize `FeeAuthorities`:
	err = decoder.Decode(&obj.FeeAuthorities)
	if err != nil {
		return err
	}
	// Deserialize `FeeWallet`:
	err = decoder.Decode(&obj.FeeWallet)
	if err != nil {
		return err
	}
	// Deserialize `BurnFeeBps`:
	err = decoder.Decode(&obj.BurnFeeBps)
	if err != nil {
		return err
	}
	// Deserialize `MintCounter`:
	err = decoder.Decode(&obj.MintCounter)
	if err != nil {
		return err
	}
	// Deserialize `RedemptionCounter`:
	err = decoder.Decode(&obj.RedemptionCounter)
	if err != nil {
		return err
	}
	// Deserialize `Bump`:
	err = decoder.Decode(&obj.Bump)
	if err != nil {
		return err
	}
	return nil
}

type TokenRedemption struct {
	Redeemer ag_solanago.PublicKey
	Value    uint64
	DestAddr [25]uint8
	Fee      uint64
	Mint     ag_solanago.PublicKey
	Id       ag_binary.Uint128
}

func (obj TokenRedemption) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Redeemer` param:
	err = encoder.Encode(obj.Redeemer)
	if err != nil {
		return err
	}
	// Serialize `Value` param:
	err = encoder.Encode(obj.Value)
	if err != nil {
		return err
	}
	// Serialize `DestAddr` param:
	err = encoder.Encode(obj.DestAddr)
	if err != nil {
		return err
	}
	// Serialize `Fee` param:
	err = encoder.Encode(obj.Fee)
	if err != nil {
		return err
	}
	// Serialize `Mint` param:
	err = encoder.Encode(obj.Mint)
	if err != nil {
		return err
	}
	// Serialize `Id` param:
	err = encoder.Encode(obj.Id)
	if err != nil {
		return err
	}
	return nil
}

func (obj *TokenRedemption) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Redeemer`:
	err = decoder.Decode(&obj.Redeemer)
	if err != nil {
		return err
	}
	// Deserialize `Value`:
	err = decoder.Decode(&obj.Value)
	if err != nil {
		return err
	}
	// Deserialize `DestAddr`:
	err = decoder.Decode(&obj.DestAddr)
	if err != nil {
		return err
	}
	// Deserialize `Fee`:
	err = decoder.Decode(&obj.Fee)
	if err != nil {
		return err
	}
	// Deserialize `Mint`:
	err = decoder.Decode(&obj.Mint)
	if err != nil {
		return err
	}
	// Deserialize `Id`:
	err = decoder.Decode(&obj.Id)
	if err != nil {
		return err
	}
	return nil
}

type TokensMintedWithFee struct {
	Recipient ag_solanago.PublicKey
	Value     uint64
	Fee       uint64
	Mint      ag_solanago.PublicKey
	Id        ag_binary.Uint128
}

func (obj TokensMintedWithFee) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Recipient` param:
	err = encoder.Encode(obj.Recipient)
	if err != nil {
		return err
	}
	// Serialize `Value` param:
	err = encoder.Encode(obj.Value)
	if err != nil {
		return err
	}
	// Serialize `Fee` param:
	err = encoder.Encode(obj.Fee)
	if err != nil {
		return err
	}
	// Serialize `Mint` param:
	err = encoder.Encode(obj.Mint)
	if err != nil {
		return err
	}
	// Serialize `Id` param:
	err = encoder.Encode(obj.Id)
	if err != nil {
		return err
	}
	return nil
}

func (obj *TokensMintedWithFee) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Recipient`:
	err = decoder.Decode(&obj.Recipient)
	if err != nil {
		return err
	}
	// Deserialize `Value`:
	err = decoder.Decode(&obj.Value)
	if err != nil {
		return err
	}
	// Deserialize `Fee`:
	err = decoder.Decode(&obj.Fee)
	if err != nil {
		return err
	}
	// Deserialize `Mint`:
	err = decoder.Decode(&obj.Mint)
	if err != nil {
		return err
	}
	// Deserialize `Id`:
	err = decoder.Decode(&obj.Id)
	if err != nil {
		return err
	}
	return nil
}

type UnwrapArgs struct {
	Value    uint64
	DestAddr [25]uint8
}

func (obj UnwrapArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Value` param:
	err = encoder.Encode(obj.Value)
	if err != nil {
		return err
	}
	// Serialize `DestAddr` param:
	err = encoder.Encode(obj.DestAddr)
	if err != nil {
		return err
	}
	return nil
}

func (obj *UnwrapArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Value`:
	err = decoder.Decode(&obj.Value)
	if err != nil {
		return err
	}
	// Deserialize `DestAddr`:
	err = decoder.Decode(&obj.DestAddr)
	if err != nil {
		return err
	}
	return nil
}

type UpdateBurnFeeBpsArgs struct {
	NewBurnFeeBps uint64
}

func (obj UpdateBurnFeeBpsArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `NewBurnFeeBps` param:
	err = encoder.Encode(obj.NewBurnFeeBps)
	if err != nil {
		return err
	}
	return nil
}

func (obj *UpdateBurnFeeBpsArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `NewBurnFeeBps`:
	err = decoder.Decode(&obj.NewBurnFeeBps)
	if err != nil {
		return err
	}
	return nil
}

type UpdateFeeWalletArgs struct {
	NewFeeWallet ag_solanago.PublicKey
}

func (obj UpdateFeeWalletArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `NewFeeWallet` param:
	err = encoder.Encode(obj.NewFeeWallet)
	if err != nil {
		return err
	}
	return nil
}

func (obj *UpdateFeeWalletArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `NewFeeWallet`:
	err = decoder.Decode(&obj.NewFeeWallet)
	if err != nil {
		return err
	}
	return nil
}

type WhitelistedWallet struct {
	// The authority which can update the config
	Wallet            ag_solanago.PublicKey
	IsWhitelisted     bool
	MintCounter       uint64
	RedemptionCounter uint64

	// The bump
	Bump uint8
}

func (obj WhitelistedWallet) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Wallet` param:
	err = encoder.Encode(obj.Wallet)
	if err != nil {
		return err
	}
	// Serialize `IsWhitelisted` param:
	err = encoder.Encode(obj.IsWhitelisted)
	if err != nil {
		return err
	}
	// Serialize `MintCounter` param:
	err = encoder.Encode(obj.MintCounter)
	if err != nil {
		return err
	}
	// Serialize `RedemptionCounter` param:
	err = encoder.Encode(obj.RedemptionCounter)
	if err != nil {
		return err
	}
	// Serialize `Bump` param:
	err = encoder.Encode(obj.Bump)
	if err != nil {
		return err
	}
	return nil
}

func (obj *WhitelistedWallet) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Wallet`:
	err = decoder.Decode(&obj.Wallet)
	if err != nil {
		return err
	}
	// Deserialize `IsWhitelisted`:
	err = decoder.Decode(&obj.IsWhitelisted)
	if err != nil {
		return err
	}
	// Deserialize `MintCounter`:
	err = decoder.Decode(&obj.MintCounter)
	if err != nil {
		return err
	}
	// Deserialize `RedemptionCounter`:
	err = decoder.Decode(&obj.RedemptionCounter)
	if err != nil {
		return err
	}
	// Deserialize `Bump`:
	err = decoder.Decode(&obj.Bump)
	if err != nil {
		return err
	}
	return nil
}

type WrapArgs struct {
	Value uint64
	Fee   uint64
}

func (obj WrapArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Value` param:
	err = encoder.Encode(obj.Value)
	if err != nil {
		return err
	}
	// Serialize `Fee` param:
	err = encoder.Encode(obj.Fee)
	if err != nil {
		return err
	}
	return nil
}

func (obj *WrapArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Value`:
	err = decoder.Decode(&obj.Value)
	if err != nil {
		return err
	}
	// Deserialize `Fee`:
	err = decoder.Decode(&obj.Fee)
	if err != nil {
		return err
	}
	return nil
}
