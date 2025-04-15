// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package rock_spl_token

import (
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

type AddFeeAuthorityArgs struct {
	// Public key of the new fee authority to add
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
	// Public key of the new mint authority to add
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

type GlobalConfig struct {
	// The authority which can update the config and manage other authorities
	GlobalAuthority ag_solanago.PublicKey

	// List of authorities that can mint tokens
	MintAuthorities []ag_solanago.PublicKey

	// List of authorities that can manage fees
	FeeAuthorities []ag_solanago.PublicKey

	// Address of the token mint controlled by this config
	MintAddress ag_solanago.PublicKey

	// Wallet that receives all fees
	FeeWallet ag_solanago.PublicKey

	// Fee percentage in basis points (1 bp = 0.01%)
	BurnFeeBps uint64

	// Counter for total mint operations
	MintCounter ag_binary.Uint128

	// Counter for total redemption operations
	RedemptionCounter ag_binary.Uint128

	// PDA bump seed
	Bump uint8
}

func (obj GlobalConfig) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `GlobalAuthority` param:
	err = encoder.Encode(obj.GlobalAuthority)
	if err != nil {
		return err
	}
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
	// Serialize `MintAddress` param:
	err = encoder.Encode(obj.MintAddress)
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

func (obj *GlobalConfig) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `GlobalAuthority`:
	err = decoder.Decode(&obj.GlobalAuthority)
	if err != nil {
		return err
	}
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
	// Deserialize `MintAddress`:
	err = decoder.Decode(&obj.MintAddress)
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

type InitializeArgs struct {
	// Authority that can manage the program configuration
	GlobalAuthority ag_solanago.PublicKey

	// Initial list of authorities that can mint tokens
	MintAuthorities []ag_solanago.PublicKey

	// Initial list of authorities that can manage fees
	FeeAuthorities []ag_solanago.PublicKey

	// Wallet that will receive all fees
	FeeWallet ag_solanago.PublicKey

	// Initial burn fee in basis points (1 bp = 0.01%)
	BurnFeeBps uint64

	// Name of the token
	TokenName string

	// Symbol/ticker of the token
	TokenSymbol string

	// Number of decimal places
	TokenDecimals uint8

	// URI for token metadata
	TokenUri string
}

func (obj InitializeArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `GlobalAuthority` param:
	err = encoder.Encode(obj.GlobalAuthority)
	if err != nil {
		return err
	}
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

func (obj *InitializeArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `GlobalAuthority`:
	err = decoder.Decode(&obj.GlobalAuthority)
	if err != nil {
		return err
	}
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

type RemoveFeeAuthorityArgs struct {
	// Public key of the fee authority to remove
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
	// Public key of the mint authority to remove
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

type TokenRedemption struct {
	// Address burning the tokens
	Redeemer ag_solanago.PublicKey

	// Amount of tokens burned (after fee)
	Value uint64

	// Destination address on external chain
	DestAddr [25]uint8

	// Fee amount taken
	Fee uint64

	// Token mint address
	Mint ag_solanago.PublicKey

	// Unique identifier for this redemption
	Id ag_binary.Uint128
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
	// Address receiving the tokens
	Recipient ag_solanago.PublicKey

	// Amount of tokens minted (after fee)
	Value uint64

	// Fee amount taken
	Fee uint64

	// Token mint address
	Mint ag_solanago.PublicKey

	// Unique identifier for this mint
	Id ag_binary.Uint128
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
	// Amount of tokens to burn
	Value uint64

	// Destination address on external chain
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
	// New burn fee in basis points (1 bp = 0.01%)
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
	// New wallet address that will receive fees
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

type UpdateGlobalAuthorityArgs struct {
	// Public key of the new global authority
	NewGlobalAuthority ag_solanago.PublicKey
}

func (obj UpdateGlobalAuthorityArgs) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `NewGlobalAuthority` param:
	err = encoder.Encode(obj.NewGlobalAuthority)
	if err != nil {
		return err
	}
	return nil
}

func (obj *UpdateGlobalAuthorityArgs) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `NewGlobalAuthority`:
	err = decoder.Decode(&obj.NewGlobalAuthority)
	if err != nil {
		return err
	}
	return nil
}

type WrapArgs struct {
	// Total amount to mint
	Value uint64

	// Fee amount to take from total
	Fee uint64
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
