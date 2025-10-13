package types

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
)

type Wallet interface {
	// Address returns a human readable version of the address.
	Address() string
}

var ErrUnknownWalletType = fmt.Errorf("error in NewWallet: unknown wallet type")

func NewWallet(k *Key, w WalletType) (Wallet, error) {
	switch w {
	case WalletType_WALLET_TYPE_NATIVE:
		return NewNativeWallet(k)
	case WalletType_WALLET_TYPE_EVM:
		return NewEthereumWallet(k)
	case WalletType_WALLET_TYPE_BTC_TESTNET:
		return NewBTCWallet(k, &chaincfg.TestNet3Params)
	case WalletType_WALLET_TYPE_BTC_MAINNET:
		return NewBTCWallet(k, &chaincfg.MainNetParams)
	case WalletType_WALLET_TYPE_BTC_REGNET:
		return NewBTCWallet(k, &chaincfg.RegressionNetParams)
	case WalletType_WALLET_TYPE_SOLANA:
		return NewSolanaWallet(k)
	case WalletType_WALLET_TYPE_ZCASH_MAINNET:
		return NewZCashWallet(k, "mainnet")
	case WalletType_WALLET_TYPE_ZCASH_TESTNET:
		return NewZCashWallet(k, "testnet")
	case WalletType_WALLET_TYPE_ZCASH_REGNET:
		return NewZCashWallet(k, "regtest")
	}
	return nil, ErrUnknownWalletType
}

// Transfer represents a generic transfer of tokens on a layer 1 blockchain.
type Transfer struct {
	// To uniquely identifies the recipient of the transfer.
	To []byte

	// Amount is the amount being transferred.
	Amount *big.Int

	// CoinIdentifier uniquely identifies the coin being transferred.
	CoinIdentifier []byte

	// DataForSigning is the data that will be signed by the key.
	DataForSigning []byte

	// DataForSigning (Hashes) when the Transaction requires multiple signatures (eg. each Bitcoin UTXO)
	SigHashes [][]byte
}

// TxParser can be implemented by wallets that are able to parse unsigned
// transactions into the common Layer1Tx format.
//
// By doing that, wallets can expose more functionalities
type TxParser interface {
	ParseTx(b []byte, m Metadata) (Transfer, error)
}

type Metadata any
