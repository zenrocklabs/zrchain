package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type EthereumWallet struct {
	key *ecdsa.PublicKey
}

var _ Wallet = &EthereumWallet{}
var _ TxParser = &EthereumWallet{}

func NewEthereumWallet(k *Key) (*EthereumWallet, error) {
	pubkey, err := k.ToECDSASecp256k1()
	if err != nil {
		return nil, err
	}
	return &EthereumWallet{key: pubkey}, nil
}

func (w *EthereumWallet) Address() string {
	addr := crypto.PubkeyToAddress(*w.key)
	return addr.Hex()
}

func (*EthereumWallet) ParseTx(b []byte, m Metadata) (Transfer, error) {
	meta, ok := m.(*MetadataEthereum)
	if !ok || meta == nil {
		return Transfer{}, fmt.Errorf("invalid metadata field, expected *MetadataEthereum, got %T", m)
	}

	tx, err := ParseEthereumTransaction(b, big.NewInt(int64(meta.ChainId)))
	if err != nil {
		return Transfer{}, err
	}

	coinIdentifier := []byte("ETH/")
	if tx.Contract != nil {
		coinIdentifier = append(coinIdentifier, tx.Contract.Bytes()...)
	}

	return Transfer{
		To:             tx.To.Bytes(),
		Amount:         tx.Amount,
		CoinIdentifier: coinIdentifier,
		DataForSigning: []byte(hex.EncodeToString(tx.DataForSigning)),
	}, nil
}

// EthereumTransfer represents an ETH transfer or an ERC-20 transfer on the
// Ethereum blockchain.
type EthereumTransfer struct {
	// To is the destination of the transfer.
	To *common.Address

	// Amount is the amount being transferred.
	Amount *big.Int

	// Contract is nil if the native currency (ETH) is being transferred,
	// or is the address of the contract if a ERC-20 token is being
	// transferred.
	Contract *common.Address

	DataForSigning []byte
}

// ParseEthereumTransaction parses an unsigned transaction that can be an ETH
// transfer or a ERC-20 transfer.
func ParseEthereumTransaction(b []byte, chainID *big.Int) (*EthereumTransfer, error) {
	txData, err := shared.DecodeUnsignedPayload(b)
	if err != nil {
		return nil, err
	}
	// create new types Transaction from input fields
	tx := types.NewTx(txData)

	value := tx.Value()

	// Use latest signer for the supplied chainID
	signer := types.LatestSignerForChainID(chainID)

	hash := signer.Hash(tx)

	transfer := &EthereumTransfer{
		To:             tx.To(),
		Amount:         value,
		DataForSigning: hash.Bytes(),
	}

	if len(tx.Data()) > 0 {
		// a contract call is being made
		transfer.Contract = tx.To()
		callMsg, parsed, err := parseCallData(tx.Data()) // - TODO we should refactor this so that value can be extracted from all known contract calls
		if err != nil {
			return nil, err
		}
		if !parsed {
			// Most contract calls will fall into this category. Over time parseCallData must be improved so that
			// asset value movements can be tracked over an increasing set of contract types.
			return transfer, nil
		}
		transfer.To = callMsg.To
		transfer.Amount = callMsg.Value
	}

	return transfer, nil
}

func parseCallData(txData []byte) (call *ethereum.CallMsg, parsed bool, err error) {
	if len(txData) < 4 {
		return nil, false, fmt.Errorf("invalid contract call")
	}

	switch {
	case checkMethodId(txData, transferMethodID):
		// 32 bytes - recipient address
		// 32 bytes - amount
		to, amt, err := rawUnpackERC20Transfer(txData)
		if err != nil {
			return nil, false, err
		}
		return &ethereum.CallMsg{To: to, Value: amt}, true, nil

	case checkMethodId(txData, transferFromMethodID):
		// 32 bytes - sender address
		// 32 bytes - recipient address
		// 32 bytes - amount
		to, amt, err := rawUnpackERC20TransferFrom(txData)
		if err != nil {
			return nil, false, err
		}
		return &ethereum.CallMsg{To: to, Value: amt}, true, nil

	default:
		return nil, false, nil
	}
}

var (
	methodIdLen          = 4
	transferMethodID     = crypto.Keccak256Hash([]byte("transfer(address,uint256)")).Bytes()[:methodIdLen]
	transferFromMethodID = crypto.Keccak256Hash([]byte("transferFrom(address,address,uint256)")).Bytes()[:methodIdLen]
	addrPrefix           = hexutil.MustDecode("0x000000000000000000000000")
)

// rawUnpackERC20Transfer Unpack without use of Go ABI package. Assumes correct ERC20 payload formatting
func rawUnpackERC20Transfer(txData []byte) (to *common.Address, amount *big.Int, err error) {
	if !checkMethodId(txData, transferMethodID) {
		return nil, nil, fmt.Errorf("wrong method id for transfer")
	}
	to, err = getAddress(txData, 1)
	if err != nil {
		return nil, nil, errors.New("invalid ERC-20 transfer: recipient address is not 20 bytes")
	}
	amount, err = getAmount(txData, 2)
	return
}

// rawUnpackERC20TransferFrom Unpack without use of Go ABI package. Assumes correct ERC20 payload formatting
func rawUnpackERC20TransferFrom(txData []byte) (to *common.Address, amount *big.Int, err error) {
	if !checkMethodId(txData, transferFromMethodID) {
		return nil, nil, fmt.Errorf("wrong method id for transferFrom")
	}
	to, err = getAddress(txData, 2)
	if err != nil {
		return nil, nil, errors.New("invalid ERC-20 transfer from: recipient address is not 20 bytes")
	}
	amount, err = getAmount(txData, 3)
	return
}

func checkMethodId(txData []byte, methodId []byte) bool {
	return bytes.Equal(txData[:methodIdLen], methodId)
}

func getAddress(txData []byte, field int) (*common.Address, error) {
	addrBytes, err := getField(txData, field)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(addrBytes[:12], addrPrefix) {
		return nil, errors.New("invalid ERC20 address")
	}

	addr := common.BytesToAddress(addrBytes[12:])
	return &addr, nil
}

func getAmount(txData []byte, field int) (*big.Int, error) {
	amountBytes, err := getField(txData, field)
	if err != nil {
		return nil, err
	}

	amount := new(big.Int).SetBytes(amountBytes)
	return amount, nil
}

func getField(txData []byte, field int) ([]byte, error) {
	start := 4 + (32 * (field - 1))
	end := start + 32

	if start > len(txData) || end > len(txData) {
		return nil, errors.New("invalid data length")
	}

	return txData[start:end], nil
}
