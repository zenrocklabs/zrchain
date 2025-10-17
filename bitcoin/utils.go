package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func CalculateTXID(rawtx string, chainName string) (*chainhash.Hash, error) {
	rawTxBytes, err := hex.DecodeString(rawtx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction hex: %w", err)
	}

	if isZcashChain(chainName) {
		return hashSerializedTransaction(rawTxBytes)
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
	isZcash := strings.Contains(strings.ToLower(chainName), "zcash")

	if isZcash {
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

	// For older Zcash versions (v4 and below), try Bitcoin-style parsing
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
