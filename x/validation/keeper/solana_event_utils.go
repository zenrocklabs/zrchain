package keeper

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
)

func eventIDFromSolanaMintEvent(event sidecarapitypes.SolanaMintEvent) (*big.Int, error) {
	if len(event.SigHash) > 0 && len(event.SigHash) <= 16 {
		return eventIDFromLEBytes(event.SigHash), nil
	}
	return eventIDFromTxSig(event.TxSig)
}

func eventIDFromTxSig(txSig string) (*big.Int, error) {
	decoded, err := hex.DecodeString(txSig)
	if err != nil {
		return nil, fmt.Errorf("invalid tx_sig hex: %w", err)
	}
	if len(decoded) > 16 {
		return nil, fmt.Errorf("tx_sig length %d exceeds expected 16 bytes", len(decoded))
	}
	return eventIDFromLEBytes(decoded), nil
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
