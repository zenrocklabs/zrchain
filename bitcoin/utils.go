package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func CalculateTXID(rawtx string) *chainhash.Hash {
	rawTxBytes, err := hex.DecodeString(rawtx)
	if err != nil {
		log.Fatalf("Failed to decode transaction hex: %v", err)
	}

	// Use a bytes.Reader to read the transaction bytes
	reader := bytes.NewReader(rawTxBytes)

	// Parse the transaction bytes into a wire.MsgTx object
	var msgTx wire.MsgTx
	err = msgTx.Deserialize(reader)
	if err != nil {
		log.Fatalf("Failed to deserialize transaction: %v", err)
	}

	var buf bytes.Buffer

	BlankWitnessData(&msgTx)

	err = msgTx.Serialize(&buf)
	if err != nil {
		log.Fatalf("Failed to serialize transaction: %v", err)
	}

	// Compute the double SHA-256 hash of the serialized transaction
	firstHash := sha256.Sum256(buf.Bytes())
	secondHash := sha256.Sum256(firstHash[:])

	// The txid is the double SHA-256 hash in reverse order (little-endian)
	txid := ReverseBytes(secondHash[:])

	// Print the txid
	fmt.Printf("Transaction ID (txid): %x\n", txid)

	calculatedTXid, _ := chainhash.NewHash(txid)
	return calculatedTXid
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

func DecodeOutputs(rawTx string, chainName string) ([]TXOutputs, error) {
	rawTxBytes, err := hex.DecodeString(rawTx)
	// Use a bytes.Reader to read the transaction bytes
	reader := bytes.NewReader(rawTxBytes)

	// Parse the transaction bytes into a wire.MsgTx object
	var msgTx wire.MsgTx
	err = msgTx.Deserialize(reader)
	if err != nil {
		log.Fatalf("Failed to deserialize transaction: %v", err)
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
			address = addrs[0].String()
		}

		outputs = append(outputs, TXOutputs{
			OutputIndex: uint(i),
			Amount:      uint64(out.Value),
			Address:     address,
		})
	}
	return outputs, nil
}

func ChainFromString(chainName string) *chaincfg.Params {
	chainName = strings.ToLower(chainName)
	switch chainName {
	case "mainnet":
		return &chaincfg.MainNetParams
	case "testnet", "testnet3":
		return &chaincfg.TestNet3Params
	case "regtest", "regnet":
		return &chaincfg.RegressionNetParams
	case "testnet4":
		return &chaincfg.TestNet3Params //TestNet4Params not available yet (22/7/24)

	default:
		return nil
	}
}
