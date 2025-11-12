package solzenbtc

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

const (
	eventStoreGlobalConfigSeed   = "global_config"
	eventStoreZenzecWrapSeed     = "zenbtc_wrap"      // ZenZEC uses original "zenbtc_wrap" seed
	eventStoreZenzecUnwrapSeed   = "zenbtc_unwrap"    // ZenZEC uses original "zenbtc_unwrap" seed
	eventStoreZenbtc2WrapSeed    = "zenbtc2_wrap"     // ZenBTC uses new "zenbtc2_wrap" seed
	eventStoreZenbtc2UnwrapSeed  = "zenbtc2_unwrap"   // ZenBTC uses new "zenbtc2_unwrap" seed
	eventStoreRockWrapSeed       = "rock_wrap"        // ROCK uses "rock_wrap" seed
	eventStoreRockUnwrapSeed     = "rock_unwrap"      // ROCK uses "rock_unwrap" seed
	zenbtcWrapShardCount         = 10
	zenbtcUnwrapShardCount       = 17
	rockWrapShardCount           = 10
	rockUnwrapShardCount         = 17
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

// GetEventStoreZenzecWrapShardPDA returns the PDA for ZenZEC wrap events (uses "zenbtc_wrap" seed)
func GetEventStoreZenzecWrapShardPDA(eventStoreProgramID solana.PublicKey, eventID *big.Int) (solana.PublicKey, error) {
	if eventID == nil {
		return solana.PublicKey{}, fmt.Errorf("eventID cannot be nil")
	}

	shardIndexBig := new(big.Int).Mod(eventID, big.NewInt(zenbtcWrapShardCount))
	index := uint16(shardIndexBig.Uint64())

	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], index)

	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(eventStoreZenzecWrapSeed),
			indexBytes[:],
		},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

// GetEventStoreZenzecUnwrapShardPDA returns the PDA for ZenZEC unwrap events (uses "zenbtc_unwrap" seed)
func GetEventStoreZenzecUnwrapShardPDA(eventStoreProgramID solana.PublicKey, eventID *big.Int) (solana.PublicKey, error) {
	if eventID == nil {
		return solana.PublicKey{}, fmt.Errorf("eventID cannot be nil")
	}

	shardIndexBig := new(big.Int).Mod(eventID, big.NewInt(zenbtcUnwrapShardCount))
	index := uint16(shardIndexBig.Uint64())

	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], index)

	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(eventStoreZenzecUnwrapSeed),
			indexBytes[:],
		},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

// GetEventStoreZenbtc2WrapShardPDA returns the PDA for ZenBTC wrap events (uses "zenbtc2_wrap" seed)
func GetEventStoreZenbtc2WrapShardPDA(eventStoreProgramID solana.PublicKey, eventID *big.Int) (solana.PublicKey, error) {
	if eventID == nil {
		return solana.PublicKey{}, fmt.Errorf("eventID cannot be nil")
	}

	shardIndexBig := new(big.Int).Mod(eventID, big.NewInt(zenbtcWrapShardCount))
	index := uint16(shardIndexBig.Uint64())

	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], index)

	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(eventStoreZenbtc2WrapSeed),
			indexBytes[:],
		},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

// GetEventStoreZenbtc2UnwrapShardPDA returns the PDA for ZenBTC unwrap events (uses "zenbtc2_unwrap" seed)
func GetEventStoreZenbtc2UnwrapShardPDA(eventStoreProgramID solana.PublicKey, eventID *big.Int) (solana.PublicKey, error) {
	if eventID == nil {
		return solana.PublicKey{}, fmt.Errorf("eventID cannot be nil")
	}

	shardIndexBig := new(big.Int).Mod(eventID, big.NewInt(zenbtcUnwrapShardCount))
	index := uint16(shardIndexBig.Uint64())

	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], index)

	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(eventStoreZenbtc2UnwrapSeed),
			indexBytes[:],
		},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

// GetEventStoreRockWrapShardPDA returns the PDA for ROCK wrap events (uses "rock_wrap" seed)
func GetEventStoreRockWrapShardPDA(eventStoreProgramID solana.PublicKey, eventID *big.Int) (solana.PublicKey, error) {
	if eventID == nil {
		return solana.PublicKey{}, fmt.Errorf("eventID cannot be nil")
	}

	shardIndexBig := new(big.Int).Mod(eventID, big.NewInt(rockWrapShardCount))
	index := uint16(shardIndexBig.Uint64())

	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], index)

	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(eventStoreRockWrapSeed),
			indexBytes[:],
		},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}

// GetEventStoreRockUnwrapShardPDA returns the PDA for ROCK unwrap events (uses "rock_unwrap" seed)
func GetEventStoreRockUnwrapShardPDA(eventStoreProgramID solana.PublicKey, eventID *big.Int) (solana.PublicKey, error) {
	if eventID == nil {
		return solana.PublicKey{}, fmt.Errorf("eventID cannot be nil")
	}

	shardIndexBig := new(big.Int).Mod(eventID, big.NewInt(rockUnwrapShardCount))
	index := uint16(shardIndexBig.Uint64())

	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], index)

	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(eventStoreRockUnwrapSeed),
			indexBytes[:],
		},
		eventStoreProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}
