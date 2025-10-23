package solzenbtc

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

const (
	eventStoreGlobalConfigSeed   = "global_config"
	eventStoreZenbtcWrapSeed     = "zenbtc_wrap"
	eventStoreZenbtcUnwrapSeed   = "zenbtc_unwrap"
	zenbtcWrapShardCount         = 10
	zenbtcUnwrapShardCount       = 17
)

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

func GetMintAddress(programID solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("wrapped_mint"),
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

func GetEventStoreGlobalConfigPDA(eventStoreProgramID solana.PublicKey) (solana.PublicKey, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{[]byte(eventStoreGlobalConfigSeed)},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

func GetEventStoreZenbtcWrapShardPDA(eventStoreProgramID solana.PublicKey, eventID *big.Int) (solana.PublicKey, error) {
	if eventID == nil {
		return solana.PublicKey{}, fmt.Errorf("eventID cannot be nil")
	}

	shardIndexBig := new(big.Int).Mod(eventID, big.NewInt(zenbtcWrapShardCount))
	index := uint16(shardIndexBig.Uint64())

	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], index)

	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(eventStoreZenbtcWrapSeed),
			indexBytes[:],
		},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

func GetEventStoreZenbtcUnwrapShardPDA(eventStoreProgramID solana.PublicKey, eventID *big.Int) (solana.PublicKey, error) {
	if eventID == nil {
		return solana.PublicKey{}, fmt.Errorf("eventID cannot be nil")
	}

	shardIndexBig := new(big.Int).Mod(eventID, big.NewInt(zenbtcUnwrapShardCount))
	index := uint16(shardIndexBig.Uint64())

	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], index)

	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(eventStoreZenbtcUnwrapSeed),
			indexBytes[:],
		},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}
