package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gagliardetto/solana-go"
)

type eventStoreWrapType int

const (
	eventStoreWrapTypeZenbtc eventStoreWrapType = iota + 1
	eventStoreWrapTypeRock
)

var (
	eventStoreGlobalConfigSeed        = []byte("global_config")
	eventStoreZenbtcWrapSeed          = []byte("zenbtc_wrap")
	eventStoreRockWrapSeed            = []byte("rock_wrap")
	storeZenbtcWrapEventDiscriminator = []byte{89, 219, 17, 100, 222, 151, 27, 199}
	storeRockWrapEventDiscriminator   = []byte{196, 31, 228, 113, 200, 144, 183, 76}

	eventStoreZenbtcShardCount uint16 = 10
	eventStoreRockShardCount   uint16 = 10
)

type eventStoreRequest struct {
	ProgramID string
	WrapType  eventStoreWrapType
	EventID   [16]byte
}

func newEventStoreEventID(prefix uint64, id uint64) [16]byte {
	var out [16]byte
	binary.LittleEndian.PutUint64(out[:8], id)
	binary.LittleEndian.PutUint64(out[8:], prefix)
	return out
}

func computeShardIndex(eventID [16]byte, shardCount uint16) uint16 {
	low := binary.LittleEndian.Uint64(eventID[:8])
	high := binary.LittleEndian.Uint64(eventID[8:])

	var value big.Int
	var highInt big.Int
	highInt.SetUint64(high)
	value.Lsh(&highInt, 64)
	var lowInt big.Int
	lowInt.SetUint64(low)
	value.Add(&value, &lowInt)

	var modulus big.Int
	modulus.Mod(&value, big.NewInt(int64(shardCount)))
	return uint16(modulus.Uint64())
}

type tokensMintedWithFee struct {
	Recipient solana.PublicKey
	Value     uint64
	Fee       uint64
	Mint      solana.PublicKey
	ID        [16]byte
}

func encodeTokensMintedWithFee(event tokensMintedWithFee) []byte {
	buf := make([]byte, 0, 32+8+8+32+16)
	buf = append(buf, event.Recipient[:]...)

	var uintBuf [8]byte
	binary.LittleEndian.PutUint64(uintBuf[:], event.Value)
	buf = append(buf, uintBuf[:]...)

	binary.LittleEndian.PutUint64(uintBuf[:], event.Fee)
	buf = append(buf, uintBuf[:]...)

	buf = append(buf, event.Mint[:]...)
	buf = append(buf, event.ID[:]...)
	return buf
}

func encodeStoreZenbtcWrapEventArgs(event tokensMintedWithFee) []byte {
	return encodeTokensMintedWithFee(event)
}

func (k Keeper) buildEventStoreWrapInstruction(
	ctx sdk.Context,
	req *solanaMintTxRequest,
	eventStore *eventStoreRequest,
	signerPubKey solana.PublicKey,
	mintKey solana.PublicKey,
	recipient solana.PublicKey,
) (solana.Instruction, error) {
	if eventStore == nil {
		return nil, fmt.Errorf("event store request is nil")
	}

	programID, err := solana.PublicKeyFromBase58(eventStore.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("invalid event store program id: %w", err)
	}

	globalConfig, _, err := solana.FindProgramAddress(
		[][]byte{eventStoreGlobalConfigSeed},
		programID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive event store global config PDA: %w", err)
	}

	var (
		discriminator []byte
		shardSeed     []byte
		shardCount    uint16
	)

	switch eventStore.WrapType {
	case eventStoreWrapTypeZenbtc:
		discriminator = storeZenbtcWrapEventDiscriminator
		shardSeed = eventStoreZenbtcWrapSeed
		shardCount = eventStoreZenbtcShardCount
	case eventStoreWrapTypeRock:
		discriminator = storeRockWrapEventDiscriminator
		shardSeed = eventStoreRockWrapSeed
		shardCount = eventStoreRockShardCount
	default:
		return nil, fmt.Errorf("unsupported event store wrap type: %d", eventStore.WrapType)
	}

	shardIndex := computeShardIndex(eventStore.EventID, shardCount)
	var indexBytes [2]byte
	binary.LittleEndian.PutUint16(indexBytes[:], shardIndex)

	shardAccount, _, err := solana.FindProgramAddress(
		[][]byte{shardSeed, indexBytes[:]},
		programID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive event store shard PDA: %w", err)
	}

	event := tokensMintedWithFee{
		Recipient: recipient,
		Value:     req.amount,
		Fee:       req.fee,
		Mint:      mintKey,
		ID:        eventStore.EventID,
	}

	argsData := encodeStoreZenbtcWrapEventArgs(event)

	data := make([]byte, len(discriminator)+len(argsData))
	copy(data, discriminator)
	copy(data[len(discriminator):], argsData)

	accounts := solana.AccountMetaSlice{
		solana.Meta(globalConfig),
		solana.Meta(signerPubKey).SIGNER(),
		solana.Meta(shardAccount).WRITE(),
	}

	k.Logger(ctx).Info("Added event store instruction to Solana mint transaction",
		"program_id", programID.String(),
		"shard_index", shardIndex,
	)

	return solana.NewInstruction(
		programID,
		accounts,
		data,
	), nil
}
