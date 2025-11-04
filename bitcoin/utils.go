package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	bitcoinecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// Zcash encrypted ciphertext boundaries (ZIP-244)
// encCiphertext is structured as: [compact (52 bytes)][memo (512 bytes)][non-compact (remainder)]
const (
	encCiphertextCompactSize    = 52  // Size of compact encrypted data
	encCiphertextMemoSize       = 512 // Size of memo field
	encCiphertextNonCompactStart = encCiphertextCompactSize + encCiphertextMemoSize // 564: start of non-compact data
)

type saplingSpend struct {
	cv        []byte
	nullifier []byte
	rk        []byte
}

type saplingOutput struct {
	cv            []byte
	cmu           []byte
	ephemeralKey  []byte
	encCiphertext []byte
	outCiphertext []byte
}

type orchardAction struct {
	cv            []byte
	nullifier     []byte
	rk            []byte
	cmx           []byte
	ephemeralKey  []byte
	encCiphertext []byte
	outCiphertext []byte
}

func CalculateTXID(rawtx string, chainName string) (*chainhash.Hash, error) {
	rawTxBytes, err := hex.DecodeString(rawtx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction hex: %w", err)
	}

	if isZcashChain(chainName) {
		return calculateZcashTxID(rawTxBytes)
	}

	reader := bytes.NewReader(rawTxBytes)
	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(reader); err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	BlankWitnessData(&msgTx)

	var buf bytes.Buffer
	if err := msgTx.Serialize(&buf); err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	return hashSerializedTransaction(buf.Bytes())
}

func BlankWitnessData(tx *wire.MsgTx) {
	for i := range tx.TxIn {
		tx.TxIn[i].Witness = nil // or tx.TxIn[i].Witness = [][]byte{}
	}
}

func ReverseBytes(data []byte) []byte {
	for i := 0; i < len(data)/2; i++ {
		data[i], data[len(data)-1-i] = data[len(data)-1-i], data[i]
	}
	return data
}

func ReverseHex(hexStr string) string {
	n := len(hexStr)
	if n%2 != 0 {
		// Ensure the hex string length is even
		return "invalid hex string"
	}
	result := make([]byte, n)
	for i := 0; i < n; i += 2 {
		// Copy two characters (one byte) at a time from the end to the beginning
		result[n-i-2], result[n-i-1] = hexStr[i], hexStr[i+1]
	}
	return string(result)
}

func MergeHashes(left, right *chainhash.Hash) *chainhash.Hash {
	var buffer bytes.Buffer
	buffer.Write(left.CloneBytes())
	buffer.Write(right.CloneBytes())
	mergedHash := chainhash.DoubleHashH(buffer.Bytes())
	return &mergedHash
}

func DecodeTX(rawTx []byte) (*wire.MsgTx, error) {
	reader := bytes.NewReader(rawTx)
	// Parse the transaction bytes into a wire.MsgTx object
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}
	return &msgTx, nil
}

func DecodeOutputs(rawTx string, chainName string) ([]TXOutputs, error) {
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction hex: %w", err)
	}

	// Check if this is a Zcash transaction
	if isZcashChain(chainName) {
		// For Zcash v5 transactions, use custom parsing
		return decodeZcashOutputs(rawTxBytes, chainName)
	}

	// Use a bytes.Reader to read the transaction bytes
	reader := bytes.NewReader(rawTxBytes)

	// Parse the transaction bytes into a wire.MsgTx object
	var msgTx wire.MsgTx
	err = msgTx.Deserialize(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	chain := ChainFromString(chainName)

	var outputs []TXOutputs
	//
	for i, out := range msgTx.TxOut {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, chain)
		if err != nil {
			fmt.Println("Failed to decode address from output script for output", i, err)
			continue
		}

		//if len(addrs) > 0 {
		//	fmt.Println("Output:", i, "Amount:", out.Value, "Address:", addrs[0])
		//}
		var address string
		if len(addrs) > 0 {
			if converted, ok := convertAddressForChain(addrs[0], chainName); ok {
				address = converted
			} else {
				address = addrs[0].String()
			}
		}

		outputs = append(outputs, TXOutputs{
			OutputIndex: uint(i),
			Amount:      uint64(out.Value),
			Address:     address,
		})
	}
	return outputs, nil
}

// decodeZcashOutputs parses Zcash v5 transaction outputs
func decodeZcashOutputs(rawTxBytes []byte, chainName string) ([]TXOutputs, error) {
	// Zcash v5 transaction format is complex with Sapling/Orchard components
	// For now, we'll parse just the transparent outputs which follow a similar format to Bitcoin
	// but we need to skip the v5-specific header

	if len(rawTxBytes) < 4 {
		return nil, fmt.Errorf("transaction too short")
	}

	// Check version (first 4 bytes)
	version := uint32(rawTxBytes[0]) | uint32(rawTxBytes[1])<<8 | uint32(rawTxBytes[2])<<16 | uint32(rawTxBytes[3])<<24

	// Version 5 has different structure
	if version == 5 || (rawTxBytes[0] == 0x05 && rawTxBytes[1] == 0x00 && rawTxBytes[2] == 0x00 && rawTxBytes[3] == 0x80) {
		return decodeZcashV5Outputs(rawTxBytes, chainName)
	}

	// Check if it's Zcash v2, v3, or v4 overwintered (has overwintered flag set in version)
	isOverwintered := (rawTxBytes[3] & 0x80) != 0
	if isOverwintered {
		if version == 4 || rawTxBytes[0] == 0x04 {
			return decodeZcashV4Outputs(rawTxBytes, chainName)
		}
		if version == 3 || rawTxBytes[0] == 0x03 {
			// V3 (Sapling) has same structure as V4 for transparent outputs
			return decodeZcashV4Outputs(rawTxBytes, chainName)
		}
		if version == 2 || rawTxBytes[0] == 0x02 {
			// V2 (Overwinter) has same structure as V4 for transparent outputs
			return decodeZcashV4Outputs(rawTxBytes, chainName)
		}
	}

	// For non-overwintered transactions (v1, v2), try Bitcoin-style parsing
	reader := bytes.NewReader(rawTxBytes)
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize Zcash transaction: %w", err)
	}

	chain := ChainFromString(chainName)
	var outputs []TXOutputs

	for i, out := range msgTx.TxOut {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, chain)
		if err != nil {
			fmt.Println("Failed to decode address from output script for output", i, err)
			continue
		}

		var address string
		if len(addrs) > 0 {
			if converted, ok := convertAddressForChain(addrs[0], chainName); ok {
				address = converted
			} else {
				address = addrs[0].String()
			}
		}

		outputs = append(outputs, TXOutputs{
			OutputIndex: uint(i),
			Amount:      uint64(out.Value),
			Address:     address,
		})
	}
	return outputs, nil
}

// decodeZcashV5Outputs parses outputs from a Zcash v5 transaction
func decodeZcashV5Outputs(rawTxBytes []byte, chainName string) ([]TXOutputs, error) {
	// Zcash v5 transaction structure:
	// - Header (4 bytes): version
	// - Header (4 bytes): version group id
	// - Header (4 bytes): consensus branch id
	// - Header (4 bytes): lock time
	// - Header (4 bytes): expiry height
	// - Transparent inputs (varint count + inputs)
	// - Transparent outputs (varint count + outputs) <- we want this
	// - Sapling/Orchard data...

	offset := 4 + 4 + 4 + 4 + 4 // Skip header (20 bytes total)

	if len(rawTxBytes) < offset {
		return nil, fmt.Errorf("zcash v5 transaction too short")
	}

	// Read transparent input count (varint)
	inputCount, bytesRead := readVarInt(rawTxBytes[offset:])
	offset += bytesRead

	// Skip all transparent inputs
	for i := uint64(0); i < inputCount; i++ {
		// Previous output (36 bytes: 32 byte hash + 4 byte index)
		offset += 36
		if offset > len(rawTxBytes) {
			return nil, fmt.Errorf("transaction data truncated at input %d", i)
		}

		// Script length (varint)
		scriptLen, bytesRead := readVarInt(rawTxBytes[offset:])
		offset += bytesRead

		// Script bytes
		offset += int(scriptLen)

		// Sequence (4 bytes)
		offset += 4

		if offset > len(rawTxBytes) {
			return nil, fmt.Errorf("transaction data truncated at input %d", i)
		}
	}

	// Read transparent output count (varint)
	outputCount, bytesRead := readVarInt(rawTxBytes[offset:])
	offset += bytesRead

	chain := ChainFromString(chainName)
	var outputs []TXOutputs

	// Parse transparent outputs
	for i := uint64(0); i < outputCount; i++ {
		if offset+8 > len(rawTxBytes) {
			return nil, fmt.Errorf("transaction data truncated at output %d", i)
		}

		// Amount (8 bytes, little-endian)
		amount := uint64(rawTxBytes[offset]) |
			uint64(rawTxBytes[offset+1])<<8 |
			uint64(rawTxBytes[offset+2])<<16 |
			uint64(rawTxBytes[offset+3])<<24 |
			uint64(rawTxBytes[offset+4])<<32 |
			uint64(rawTxBytes[offset+5])<<40 |
			uint64(rawTxBytes[offset+6])<<48 |
			uint64(rawTxBytes[offset+7])<<56
		offset += 8

		// Script length (varint)
		scriptLen, bytesRead := readVarInt(rawTxBytes[offset:])
		offset += bytesRead

		if offset+int(scriptLen) > len(rawTxBytes) {
			return nil, fmt.Errorf("transaction data truncated at output %d script", i)
		}

		// Script bytes
		pkScript := rawTxBytes[offset : offset+int(scriptLen)]
		offset += int(scriptLen)

		// Extract address from script
		var address string
		if len(pkScript) > 0 && chain != nil {
			_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript, chain)
			if err == nil && len(addrs) > 0 {
				if converted, ok := convertAddressForChain(addrs[0], chainName); ok {
					address = converted
				} else {
					address = addrs[0].String()
				}
			}
		}

		outputs = append(outputs, TXOutputs{
			OutputIndex: uint(i),
			Amount:      amount,
			Address:     address,
		})
	}

	return outputs, nil
}

// decodeZcashV4Outputs parses outputs from a Zcash v4 overwintered transaction
func decodeZcashV4Outputs(rawTxBytes []byte, chainName string) ([]TXOutputs, error) {
	// Zcash v4 overwintered transaction structure:
	// - Header (4 bytes): version (with overwintered flag)
	// - Header (4 bytes): version group id
	// - Transparent inputs (varint count + inputs)
	// - Transparent outputs (varint count + outputs) <- we want this
	// - Lock time (4 bytes)
	// - Expiry height (4 bytes)
	// - Value balance (8 bytes)
	// - Shielded spends/outputs...

	offset := 4 + 4 // Skip version and version group id (8 bytes total)

	if len(rawTxBytes) < offset {
		return nil, fmt.Errorf("zcash v4 transaction too short")
	}

	// Read transparent input count (varint)
	inputCount, bytesRead := readVarInt(rawTxBytes[offset:])
	offset += bytesRead

	// Skip all transparent inputs
	for i := uint64(0); i < inputCount; i++ {
		// Previous output (36 bytes: 32 byte hash + 4 byte index)
		offset += 36
		if offset > len(rawTxBytes) {
			return nil, fmt.Errorf("transaction data truncated at input %d", i)
		}

		// Script length (varint)
		scriptLen, bytesRead := readVarInt(rawTxBytes[offset:])
		offset += bytesRead

		// Script bytes
		offset += int(scriptLen)

		// Sequence (4 bytes)
		offset += 4

		if offset > len(rawTxBytes) {
			return nil, fmt.Errorf("transaction data truncated at input %d", i)
		}
	}

	// Read transparent output count (varint)
	outputCount, bytesRead := readVarInt(rawTxBytes[offset:])
	offset += bytesRead

	chain := ChainFromString(chainName)
	var outputs []TXOutputs

	// Parse transparent outputs
	for i := uint64(0); i < outputCount; i++ {
		if offset+8 > len(rawTxBytes) {
			return nil, fmt.Errorf("transaction data truncated at output %d", i)
		}

		// Amount (8 bytes, little-endian)
		amount := uint64(rawTxBytes[offset]) |
			uint64(rawTxBytes[offset+1])<<8 |
			uint64(rawTxBytes[offset+2])<<16 |
			uint64(rawTxBytes[offset+3])<<24 |
			uint64(rawTxBytes[offset+4])<<32 |
			uint64(rawTxBytes[offset+5])<<40 |
			uint64(rawTxBytes[offset+6])<<48 |
			uint64(rawTxBytes[offset+7])<<56
		offset += 8

		// Script length (varint)
		scriptLen, bytesRead := readVarInt(rawTxBytes[offset:])
		offset += bytesRead

		if offset+int(scriptLen) > len(rawTxBytes) {
			return nil, fmt.Errorf("transaction data truncated at output %d script", i)
		}

		// Script bytes
		pkScript := rawTxBytes[offset : offset+int(scriptLen)]
		offset += int(scriptLen)

		// Extract address from script
		var address string
		if len(pkScript) > 0 && chain != nil {
			_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript, chain)
			if err == nil && len(addrs) > 0 {
				if converted, ok := convertAddressForChain(addrs[0], chainName); ok {
					address = converted
				} else {
					address = addrs[0].String()
				}
			}
		}

		outputs = append(outputs, TXOutputs{
			OutputIndex: uint(i),
			Amount:      amount,
			Address:     address,
		})
	}

	return outputs, nil
}

// readVarInt reads a Bitcoin/Zcash variable length integer
func readVarInt(data []byte) (uint64, int) {
	if len(data) == 0 {
		return 0, 0
	}

	first := data[0]
	if first < 0xfd {
		return uint64(first), 1
	}

	if first == 0xfd {
		if len(data) < 3 {
			return 0, 0
		}
		return uint64(data[1]) | uint64(data[2])<<8, 3
	}

	if first == 0xfe {
		if len(data) < 5 {
			return 0, 0
		}
		return uint64(data[1]) | uint64(data[2])<<8 | uint64(data[3])<<16 | uint64(data[4])<<24, 5
	}

	// first == 0xff
	if len(data) < 9 {
		return 0, 0
	}
	return uint64(data[1]) | uint64(data[2])<<8 | uint64(data[3])<<16 | uint64(data[4])<<24 |
		uint64(data[5])<<32 | uint64(data[6])<<40 | uint64(data[7])<<48 | uint64(data[8])<<56, 9
}

func calculateZcashTxID(raw []byte) (*chainhash.Hash, error) {
	if len(raw) < 20 {
		return nil, fmt.Errorf("zcash transaction too short")
	}

	offset := 0

	versionBytes, err := readBytes(raw, &offset, 4)
	if err != nil {
		return nil, err
	}
	version := binary.LittleEndian.Uint32(versionBytes)
	effectiveVersion := version & 0x7fffffff

	// Versions prior to 5 use legacy txid calculation.
	if effectiveVersion < 5 {
		return hashSerializedTransaction(raw)
	}
	if effectiveVersion != 5 {
		return nil, fmt.Errorf("unsupported zcash transaction version %d", effectiveVersion)
	}

	versionGroupIDBytes, err := readBytes(raw, &offset, 4)
	if err != nil {
		return nil, err
	}
	branchIDBytes, err := readBytes(raw, &offset, 4)
	if err != nil {
		return nil, err
	}
	lockTimeBytes, err := readBytes(raw, &offset, 4)
	if err != nil {
		return nil, err
	}
	expiryHeightBytes, err := readBytes(raw, &offset, 4)
	if err != nil {
		return nil, err
	}

	headerDigest, err := blake2bHashPersonalString(
		"ZTxIdHeadersHash",
		versionBytes,
		versionGroupIDBytes,
		branchIDBytes,
		lockTimeBytes,
		expiryHeightBytes,
	)
	if err != nil {
		return nil, err
	}

	prevoutsBuf := bytes.Buffer{}
	sequenceBuf := bytes.Buffer{}
	txInCount, _, err := readCompactSize(raw, &offset)
	if err != nil {
		return nil, err
	}
	for i := uint64(0); i < txInCount; i++ {
		outpoint, err := readBytes(raw, &offset, 36)
		if err != nil {
			return nil, err
		}
		prevoutsBuf.Write(outpoint)

		scriptLen, _, err := readCompactSize(raw, &offset)
		if err != nil {
			return nil, err
		}
		if scriptLen > uint64(len(raw)-offset) {
			return nil, io.ErrUnexpectedEOF
		}
		if scriptLen > 0 {
			if _, err := readBytes(raw, &offset, int(scriptLen)); err != nil {
				return nil, err
			}
		}

		sequence, err := readBytes(raw, &offset, 4)
		if err != nil {
			return nil, err
		}
		sequenceBuf.Write(sequence)
	}

	outputsBuf := bytes.Buffer{}
	txOutCount, _, err := readCompactSize(raw, &offset)
	if err != nil {
		return nil, err
	}
	for i := uint64(0); i < txOutCount; i++ {
		amountBytes, err := readBytes(raw, &offset, 8)
		if err != nil {
			return nil, err
		}
		outputsBuf.Write(amountBytes)

		scriptLen, scriptLenBytes, err := readCompactSize(raw, &offset)
		if err != nil {
			return nil, err
		}
		outputsBuf.Write(scriptLenBytes)
		if scriptLen > uint64(len(raw)-offset) {
			return nil, io.ErrUnexpectedEOF
		}
		if scriptLen > 0 {
			scriptBytes, err := readBytes(raw, &offset, int(scriptLen))
			if err != nil {
				return nil, err
			}
			outputsBuf.Write(scriptBytes)
		}
	}

	nSpendsSapling, _, err := readCompactSize(raw, &offset)
	if err != nil {
		return nil, err
	}
	nOutputsSapling, _, err := readCompactSize(raw, &offset)
	if err != nil {
		return nil, err
	}

	saplingSpends := make([]saplingSpend, nSpendsSapling)
	for i := uint64(0); i < nSpendsSapling; i++ {
		cv, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		nullifier, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		rk, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		saplingSpends[i] = saplingSpend{
			cv:        append([]byte(nil), cv...),
			nullifier: append([]byte(nil), nullifier...),
			rk:        append([]byte(nil), rk...),
		}
	}

	saplingOutputs := make([]saplingOutput, nOutputsSapling)
	for i := uint64(0); i < nOutputsSapling; i++ {
		cv, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		cmu, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		ephemeralKey, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		encCiphertext, err := readBytes(raw, &offset, 580)
		if err != nil {
			return nil, err
		}
		outCiphertext, err := readBytes(raw, &offset, 80)
		if err != nil {
			return nil, err
		}
		saplingOutputs[i] = saplingOutput{
			cv:            append([]byte(nil), cv...),
			cmu:           append([]byte(nil), cmu...),
			ephemeralKey:  append([]byte(nil), ephemeralKey...),
			encCiphertext: append([]byte(nil), encCiphertext...),
			outCiphertext: append([]byte(nil), outCiphertext...),
		}
	}

	saplingPresent := nSpendsSapling > 0 || nOutputsSapling > 0

	var valueBalanceSaplingBytes []byte
	if saplingPresent {
		valueBalanceSaplingBytes, err = readBytes(raw, &offset, 8)
		if err != nil {
			return nil, err
		}
	}

	var anchorSapling []byte
	if nSpendsSapling > 0 {
		anchorSapling, err = readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
	}

	if nSpendsSapling > 0 {
		if err := skipBytes(raw, &offset, 192, nSpendsSapling); err != nil {
			return nil, err
		}
		if err := skipBytes(raw, &offset, 64, nSpendsSapling); err != nil {
			return nil, err
		}
	}

	if nOutputsSapling > 0 {
		if err := skipBytes(raw, &offset, 192, nOutputsSapling); err != nil {
			return nil, err
		}
	}

	if saplingPresent {
		if _, err := readBytes(raw, &offset, 64); err != nil {
			return nil, err
		}
	}

	nActionsOrchard, _, err := readCompactSize(raw, &offset)
	if err != nil {
		return nil, err
	}

	orchardActions := make([]orchardAction, nActionsOrchard)
	for i := uint64(0); i < nActionsOrchard; i++ {
		cv, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		nullifier, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		rk, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		cmx, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		ephemeralKey, err := readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}
		encCiphertext, err := readBytes(raw, &offset, 580)
		if err != nil {
			return nil, err
		}
		outCiphertext, err := readBytes(raw, &offset, 80)
		if err != nil {
			return nil, err
		}
		orchardActions[i] = orchardAction{
			cv:            append([]byte(nil), cv...),
			nullifier:     append([]byte(nil), nullifier...),
			rk:            append([]byte(nil), rk...),
			cmx:           append([]byte(nil), cmx...),
			ephemeralKey:  append([]byte(nil), ephemeralKey...),
			encCiphertext: append([]byte(nil), encCiphertext...),
			outCiphertext: append([]byte(nil), outCiphertext...),
		}
	}

	var (
		flagsOrchard             byte
		valueBalanceOrchardBytes []byte
		anchorOrchard            []byte
	)

	if nActionsOrchard > 0 {
		flags, err := readBytes(raw, &offset, 1)
		if err != nil {
			return nil, err
		}
		flagsOrchard = flags[0]

		valueBalanceOrchardBytes, err = readBytes(raw, &offset, 8)
		if err != nil {
			return nil, err
		}

		anchorOrchard, err = readBytes(raw, &offset, 32)
		if err != nil {
			return nil, err
		}

		sizeProofs, _, err := readCompactSize(raw, &offset)
		if err != nil {
			return nil, err
		}
		if sizeProofs > 0 {
			if err := skipBytes(raw, &offset, 1, sizeProofs); err != nil {
				return nil, err
			}
		}

		if err := skipBytes(raw, &offset, 64, nActionsOrchard); err != nil {
			return nil, err
		}

		if _, err := readBytes(raw, &offset, 64); err != nil {
			return nil, err
		}
	}

	// After Sapling/Orchard components, no further fields should remain.
	if offset != len(raw) {
		return nil, fmt.Errorf("unexpected extra data (%d bytes) in zcash transaction", len(raw)-offset)
	}

	prevoutsDigest, err := digestWithEmptyFallback("ZTxIdPrevoutHash", prevoutsBuf.Bytes(), txInCount > 0)
	if err != nil {
		return nil, err
	}
	sequenceDigest, err := digestWithEmptyFallback("ZTxIdSequencHash", sequenceBuf.Bytes(), txInCount > 0)
	if err != nil {
		return nil, err
	}
	outputsDigest, err := digestWithEmptyFallback("ZTxIdOutputsHash", outputsBuf.Bytes(), txOutCount > 0)
	if err != nil {
		return nil, err
	}

	transparentDigest, err := blake2bHashPersonalString(
		"ZTxIdTranspaHash",
		prevoutsDigest,
		sequenceDigest,
		outputsDigest,
	)
	if err != nil {
		return nil, err
	}

	saplingDigest, err := computeSaplingDigest(saplingSpends, saplingOutputs, anchorSapling, valueBalanceSaplingBytes)
	if err != nil {
		return nil, err
	}

	orchardDigest, err := computeOrchardDigest(orchardActions, flagsOrchard, valueBalanceOrchardBytes, anchorOrchard)
	if err != nil {
		return nil, err
	}

	topPersonal := make([]byte, 0, 16)
	topPersonal = append(topPersonal, []byte("ZcashTxHash_")...)
	topPersonal = append(topPersonal, branchIDBytes...)

	txidDigest, err := blake2bHashPersonal(
		topPersonal,
		headerDigest,
		transparentDigest,
		saplingDigest,
		orchardDigest,
	)
	if err != nil {
		return nil, err
	}

	txidLittleEndian := ReverseBytes(txidDigest)
	return chainhash.NewHash(txidLittleEndian)
}

func digestWithEmptyFallback(personal string, data []byte, hasData bool) ([]byte, error) {
	if hasData && len(data) > 0 {
		return blake2bHashPersonalString(personal, data)
	}
	return blake2bHashPersonalString(personal)
}

func computeSaplingDigest(
	spends []saplingSpend,
	outputs []saplingOutput,
	anchorSapling []byte,
	valueBalanceBytes []byte,
) ([]byte, error) {
	saplingPresent := len(spends) > 0 || len(outputs) > 0
	if !saplingPresent {
		return blake2bHashPersonalString("ZTxIdSaplingHash")
	}

	spendCompactBuf := bytes.Buffer{}
	for _, spend := range spends {
		spendCompactBuf.Write(spend.nullifier)
	}
	spendCompactDigest, err := digestWithEmptyFallback("ZTxIdSSpendCHash", spendCompactBuf.Bytes(), len(spends) > 0)
	if err != nil {
		return nil, err
	}

	if len(spends) > 0 && len(anchorSapling) != 32 {
		return nil, fmt.Errorf("invalid sapling anchor length %d", len(anchorSapling))
	}

	spendNonCompactBuf := bytes.Buffer{}
	for _, spend := range spends {
		spendNonCompactBuf.Write(spend.cv)
		spendNonCompactBuf.Write(anchorSapling)
		spendNonCompactBuf.Write(spend.rk)
	}
	spendNonCompactDigest, err := digestWithEmptyFallback("ZTxIdSSpendNHash", spendNonCompactBuf.Bytes(), len(spends) > 0)
	if err != nil {
		return nil, err
	}

	var spendDigest []byte
	if len(spends) > 0 {
		spendDigest, err = blake2bHashPersonalString(
			"ZTxIdSSpendsHash",
			spendCompactDigest,
			spendNonCompactDigest,
		)
		if err != nil {
			return nil, err
		}
	} else {
		spendDigest, err = blake2bHashPersonalString("ZTxIdSSpendsHash")
		if err != nil {
			return nil, err
		}
	}

	outputCompactBuf := bytes.Buffer{}
	outputMemosBuf := bytes.Buffer{}
	outputNonCompactBuf := bytes.Buffer{}
	for _, output := range outputs {
		outputCompactBuf.Write(output.cmu)
		outputCompactBuf.Write(output.ephemeralKey)
		outputCompactBuf.Write(output.encCiphertext[:encCiphertextCompactSize])

		outputMemosBuf.Write(output.encCiphertext[encCiphertextCompactSize:encCiphertextNonCompactStart])

		outputNonCompactBuf.Write(output.cv)
		outputNonCompactBuf.Write(output.encCiphertext[encCiphertextNonCompactStart:])
		outputNonCompactBuf.Write(output.outCiphertext)
	}

	outputCompactDigest, err := digestWithEmptyFallback("ZTxIdSOutC__Hash", outputCompactBuf.Bytes(), len(outputs) > 0)
	if err != nil {
		return nil, err
	}
	outputMemosDigest, err := digestWithEmptyFallback("ZTxIdSOutM__Hash", outputMemosBuf.Bytes(), len(outputs) > 0)
	if err != nil {
		return nil, err
	}
	outputNonCompactDigest, err := digestWithEmptyFallback("ZTxIdSOutN__Hash", outputNonCompactBuf.Bytes(), len(outputs) > 0)
	if err != nil {
		return nil, err
	}

	var outputDigest []byte
	if len(outputs) > 0 {
		outputDigest, err = blake2bHashPersonalString(
			"ZTxIdSOutputHash",
			outputCompactDigest,
			outputMemosDigest,
			outputNonCompactDigest,
		)
		if err != nil {
			return nil, err
		}
	} else {
		outputDigest, err = blake2bHashPersonalString("ZTxIdSOutputHash")
		if err != nil {
			return nil, err
		}
	}

	if len(valueBalanceBytes) != 8 {
		return nil, fmt.Errorf("invalid sapling value balance length %d", len(valueBalanceBytes))
	}

	return blake2bHashPersonalString(
		"ZTxIdSaplingHash",
		spendDigest,
		outputDigest,
		valueBalanceBytes,
	)
}

func computeOrchardDigest(
	actions []orchardAction,
	flags byte,
	valueBalanceBytes []byte,
	anchor []byte,
) ([]byte, error) {
	if len(actions) == 0 {
		return blake2bHashPersonalString("ZTxIdOrchardHash")
	}

	compactBuf := bytes.Buffer{}
	memosBuf := bytes.Buffer{}
	nonCompactBuf := bytes.Buffer{}
	for _, action := range actions {
		compactBuf.Write(action.nullifier)
		compactBuf.Write(action.cmx)
		compactBuf.Write(action.ephemeralKey)
		compactBuf.Write(action.encCiphertext[:encCiphertextCompactSize])

		memosBuf.Write(action.encCiphertext[encCiphertextCompactSize:encCiphertextNonCompactStart])

		nonCompactBuf.Write(action.cv)
		nonCompactBuf.Write(action.rk)
		nonCompactBuf.Write(action.encCiphertext[encCiphertextNonCompactStart:])
		nonCompactBuf.Write(action.outCiphertext)
	}

	compactDigest, err := blake2bHashPersonalString("ZTxIdOrcActCHash", compactBuf.Bytes())
	if err != nil {
		return nil, err
	}
	memosDigest, err := blake2bHashPersonalString("ZTxIdOrcActMHash", memosBuf.Bytes())
	if err != nil {
		return nil, err
	}
	nonCompactDigest, err := blake2bHashPersonalString("ZTxIdOrcActNHash", nonCompactBuf.Bytes())
	if err != nil {
		return nil, err
	}

	if len(valueBalanceBytes) != 8 {
		return nil, fmt.Errorf("invalid orchard value balance length %d", len(valueBalanceBytes))
	}
	if len(anchor) != 32 {
		return nil, fmt.Errorf("invalid orchard anchor length %d", len(anchor))
	}

	return blake2bHashPersonalString(
		"ZTxIdOrchardHash",
		compactDigest,
		memosDigest,
		nonCompactDigest,
		[]byte{flags},
		valueBalanceBytes,
		anchor,
	)
}

func blake2bHashPersonalString(personal string, parts ...[]byte) ([]byte, error) {
	return blake2bHashPersonal([]byte(personal), parts...)
}

func blake2bHashPersonal(personal []byte, parts ...[]byte) ([]byte, error) {
	if len(personal) != 16 {
		return nil, fmt.Errorf("invalid personalization length %d", len(personal))
	}
	return blake2bPersonalHash(personal, parts...)
}

func readBytes(data []byte, offset *int, length int) ([]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("negative length requested")
	}
	if *offset+length > len(data) {
		return nil, io.ErrUnexpectedEOF
	}
	bytes := data[*offset : *offset+length]
	*offset += length
	return bytes, nil
}

func skipBytes(data []byte, offset *int, elementSize int, count uint64) error {
	if count == 0 {
		return nil
	}
	if elementSize <= 0 {
		return fmt.Errorf("invalid element size %d", elementSize)
	}
	if count > uint64(math.MaxInt/elementSize) {
		return fmt.Errorf("element count too large")
	}
	total := int(count) * elementSize
	_, err := readBytes(data, offset, total)
	return err
}

func readCompactSize(data []byte, offset *int) (uint64, []byte, error) {
	if *offset >= len(data) {
		return 0, nil, io.ErrUnexpectedEOF
	}

	start := *offset
	prefix := data[*offset]
	*offset++

	switch prefix {
	case 0xfd:
		if *offset+2 > len(data) {
			return 0, nil, io.ErrUnexpectedEOF
		}
		value := binary.LittleEndian.Uint16(data[*offset : *offset+2])
		*offset += 2
		return uint64(value), data[start:*offset], nil
	case 0xfe:
		if *offset+4 > len(data) {
			return 0, nil, io.ErrUnexpectedEOF
		}
		value := binary.LittleEndian.Uint32(data[*offset : *offset+4])
		*offset += 4
		return uint64(value), data[start:*offset], nil
	case 0xff:
		if *offset+8 > len(data) {
			return 0, nil, io.ErrUnexpectedEOF
		}
		value := binary.LittleEndian.Uint64(data[*offset : *offset+8])
		*offset += 8
		return value, data[start:*offset], nil
	default:
		return uint64(prefix), data[start:*offset], nil
	}
}

func hashSerializedTransaction(serialized []byte) (*chainhash.Hash, error) {
	firstHash := sha256.Sum256(serialized)
	secondHash := sha256.Sum256(firstHash[:])
	txid := ReverseBytes(secondHash[:])
	calculatedTXID, err := chainhash.NewHash(txid)
	if err != nil {
		return nil, fmt.Errorf("failed to derive txid hash: %w", err)
	}
	return calculatedTXID, nil
}

func ChainFromString(chainName string) *chaincfg.Params {
	chainName = strings.ToLower(chainName)
	switch chainName {
	// Bitcoin chains
	case "mainnet":
		return &chaincfg.MainNetParams
	case "testnet", "testnet3":
		return &chaincfg.TestNet3Params
	case "regtest", "regnet":
		return &chaincfg.RegressionNetParams
	case "testnet4":
		return &chaincfg.TestNet3Params //TestNet4Params not available yet (22/7/24)

	// Zcash chains - use Bitcoin params as base since Zcash is a Bitcoin fork
	// The actual Zcash-specific parameters (like Sapling, etc.) are handled by the Zcash proxy
	case "zcash-mainnet", "zcashmainnet", "zcash":
		return &chaincfg.MainNetParams
	case "zcash-testnet", "zcashtestnet":
		return &chaincfg.TestNet3Params
	case "zcash-regtest", "zcash-regnet", "zcashregtest":
		return &chaincfg.RegressionNetParams

	default:
		return nil
	}
}

func isZcashChain(chainName string) bool {
	return strings.Contains(strings.ToLower(chainName), "zcash")
}

func ConvertBigIntToModNScalar(b *big.Int) *btcec.ModNScalar {
	modNScalar := new(btcec.ModNScalar)
	bytes := b.Bytes()
	if len(bytes) > 32 {
		bytes = bytes[:32] // Truncate if necessary
	}
	modNScalar.SetBytes(padTo32Bytes(bytes))
	return modNScalar
}

func padTo32Bytes(b []byte) *[32]byte {
	var padded [32]byte
	copy(padded[32-len(b):], b)
	return &padded
}

func ConvertECDSASigtoBitcoinSig(ecdsaSig string) (string, error) {
	if len(ecdsaSig) >= 132 {
		return "", fmt.Errorf("ConvertECDSASigtoBitcoinSig - invalid ecdsa signature")
	}
	r := ecdsaSig[:64]
	s := ecdsaSig[64:128]
	rBig, _ := new(big.Int).SetString(r, 16)
	sBig, _ := new(big.Int).SetString(s, 16)
	rawsig := bitcoinecdsa.NewSignature(ConvertBigIntToModNScalar(rBig), ConvertBigIntToModNScalar(sBig))
	return hex.EncodeToString(rawsig.Serialize()), nil
}

func convertAddressForChain(addr btcutil.Address, chainName string) (string, bool) {
	if !isZcashChain(chainName) {
		return "", false
	}

	switch a := addr.(type) {
	case *btcutil.AddressPubKeyHash:
		version, ok := zcashVersionBytes(chainName, false)
		if !ok {
			return "", false
		}
		hash := a.Hash160()
		return encodeZcashBase58(version, hash[:]), true
	case *btcutil.AddressScriptHash:
		version, ok := zcashVersionBytes(chainName, true)
		if !ok {
			return "", false
		}
		hash := a.Hash160()
		// convert pointer to slice
		return encodeZcashBase58(version, hash[:]), true
	default:
		return "", false
	}
}

func zcashVersionBytes(chainName string, isScriptHash bool) ([]byte, bool) {
	net := strings.ToLower(chainName)
	net = strings.TrimPrefix(net, "zcash-")
	net = strings.TrimPrefix(net, "zcash")

	switch {
	case net == "" || strings.Contains(net, "main"):
		if isScriptHash {
			return []byte{0x1c, 0xbd}, true
		}
		return []byte{0x1c, 0xb8}, true
	case strings.Contains(net, "test") || strings.Contains(net, "reg"):
		if isScriptHash {
			return []byte{0x1c, 0xba}, true
		}
		return []byte{0x1d, 0x25}, true
	default:
		return nil, false
	}
}

func encodeZcashBase58(version []byte, payload []byte) string {
	buf := make([]byte, 0, len(version)+len(payload)+4)
	buf = append(buf, version...)
	buf = append(buf, payload...)
	checksum := doubleSha256(buf)
	buf = append(buf, checksum[:4]...)
	return base58.Encode(buf)
}

func doubleSha256(data []byte) [32]byte {
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}
