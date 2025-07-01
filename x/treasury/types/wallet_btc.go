package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
)

type BitcoinWallet struct {
	key   *ecdsa.PublicKey
	chain *chaincfg.Params
}

var _ Wallet = &BitcoinWallet{}
var _ TxParser = &BitcoinWallet{}

func NewBTCWallet(k *Key, chain *chaincfg.Params) (*BitcoinWallet, error) {
	pubkey, err := k.ToBitcoinSecp256k1()
	if err != nil {
		return nil, err
	}
	return &BitcoinWallet{key: pubkey, chain: chain}, nil
}

//P2SH is deprecated in favour of P2WPKH
//func (w *BitcoinWallet) P2SHAddress() string {
//	var pubkey secp256k1.PubKey
//	pubkey.Key = crypto.CompressPubkey(w.key)
//	chain := w.chain
//	witnessProg := btcutil.Hash160(pubkey.Key)
//	p2wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chain)
//	if err != nil {
//		return ""
//	}
//	p2wpkhScript, err := txscript.PayToAddrScript(p2wpkhAddress)
//	if err != nil {
//		return ""
//	}
//	p2shAddress, err := btcutil.NewAddressScriptHash(p2wpkhScript, chain)
//	if err != nil {
//		return ""
//	}
//	return p2shAddress.EncodeAddress()
//}

func (w *BitcoinWallet) Address() string {
	//Generate a P2WPKH address
	var pubkey secp256k1.PubKey
	chain := w.chain
	pubkey.Key = crypto.CompressPubkey(w.key)
	publicKey, err := btcec.ParsePubKey(pubkey.Key)
	if err != nil {
		return ""
	}
	// Generate P2WPKH address from the public key
	witnessProg := btcutil.Hash160(publicKey.SerializeCompressed())
	address, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chain)
	if err != nil {
		return ""
	}
	return address.EncodeAddress()
}

// // ParseTx implements TxParser.
func (w *BitcoinWallet) ParseTx(b []byte, m Metadata) (Transfer, error) {
	// Although the pre-created hashes are passed in the unsigned transaction as additional fields in the witness data
	// We need to re-calculate them to ensure that they are hashes of the transaction sent & approved by the user and
	// not some hashes injected by an attacker who has compromised the zr-bitcoin-proxy
	hashes, err := w.SigHashes(b)

	if err != nil {
		return Transfer{}, err
	}
	var dataForSigning []string
	for _, hash := range hashes {
		dataForSigning = append(dataForSigning, hex.EncodeToString(hash))
	}
	return Transfer{
		SigHashes:      hashes,
		DataForSigning: []byte(strings.Join(dataForSigning, ",")),
	}, nil
}

func (w *BitcoinWallet) SigHashes(b []byte) (hashes [][]byte, err error) {
	//Sighashes re-calculates the hash for the transaction to ensure they have not been tampered with
	msgTx, err := DeserializeTransaction(b)
	_ = msgTx
	if err != nil {
		return nil, err
	}
	inputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	for _, txin := range msgTx.TxIn {
		inputFetcher.AddPrevOut(txin.PreviousOutPoint, &wire.TxOut{})
	}
	for index, txin := range msgTx.TxIn {
		//Collect data from TX Witness
		amount := BytesToInt64(txin.Witness[0])
		pkscript := txin.Witness[2]

		//Recalculate the Hash
		prevOutFetcher := txscript.NewCannedPrevOutputFetcher(pkscript, amount)
		sigHashes := txscript.NewTxSigHashes(msgTx, prevOutFetcher)
		hash, err := txscript.CalcWitnessSigHash(pkscript, sigHashes, txscript.SigHashAll, msgTx, index, amount)
		if err != nil {
			return nil, fmt.Errorf("error CalcWitnessSigHash: %w", err)
		}
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

func DeserializeTransaction(txBytes []byte) (*wire.MsgTx, error) {
	// Initialize a new empty MsgTx.
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// Create a bytes.Buffer from the transaction bytes.
	rbuf := bytes.NewReader(txBytes)

	// Deserialize the transaction using BtcDecode into the empty MsgTx.
	err := msgTx.BtcDecode(rbuf, wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		return nil, err
	}
	// Return the deserialized transaction.
	return msgTx, nil
}

func BytesToInt64(buf []byte) int64 {
	buffer := bytes.NewBuffer(buf)
	var i int64
	err := binary.Read(buffer, binary.LittleEndian, &i)
	if err != nil {
		fmt.Println("Error in binary.Read", err)
	}
	return i
}

func VerifyBitcoinSigHashes(dataForSigning [][]byte, tx []byte) (status VerificationStatus, err error) {
	//Do not check non existent data
	if tx == nil {
		return Verification_NotVerified, nil
	}
	if len(dataForSigning) == 0 {
		return Verification_NotVerified, nil
	}
	msgTx, err := DeserializeTransaction(tx)
	if err != nil {
		return Verification_NotVerified, fmt.Errorf("error deserializing transaction %s", tx)
	}

	if len(msgTx.TxOut) != len(dataForSigning) {
		return Verification_NotVerified, fmt.Errorf("dataforsigning (%d) should be the same count as UTXO inputs (%d)", len(dataForSigning), len(msgTx.TxOut))
	}

	const ScriptBytes = 2
	const Amount = 0

	for index, txin := range msgTx.TxIn {
		pkscript := txin.Witness[ScriptBytes]
		amount := BytesToInt64(txin.Witness[Amount])
		prevOutFetcher := txscript.NewCannedPrevOutputFetcher(pkscript, amount)
		sigHashes := txscript.NewTxSigHashes(msgTx, prevOutFetcher)

		hash, err := txscript.CalcWitnessSigHash(pkscript, sigHashes, txscript.SigHashAll, msgTx, index, amount)
		if err != nil {
			return Verification_NotVerified, fmt.Errorf("error CalcWitnessSigHash: %w", err)
		}
		if bytes.Compare(hash, dataForSigning[index]) != 0 {
			return Verification_Failed, fmt.Errorf("hash index %d is invalid", index)
		}
	}
	return Verification_Suceeded, nil
}
