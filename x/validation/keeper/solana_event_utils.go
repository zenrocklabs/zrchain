package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"

	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
)

func eventIDFromSolanaMintEvent(event sidecarapitypes.SolanaMintEvent) (*big.Int, error) {
	if len(event.SigHash) != 16 {
		return nil, fmt.Errorf("unsupported Solana event sig hash length: %d", len(event.SigHash))
	}
	return eventIDFromLEBytes(event.SigHash), nil
}

func eventIDFromLEBytes(b []byte) *big.Int {
	var idBytes [16]byte
	copy(idBytes[:], b)

	hi := binary.LittleEndian.Uint64(idBytes[8:])
	lo := binary.LittleEndian.Uint64(idBytes[:8])

	result := new(big.Int).SetUint64(hi)
	result.Lsh(result, 64)
	result.Add(result, new(big.Int).SetUint64(lo))
	return result
}
