package types

import (
	"crypto/ed25519"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bin "github.com/gagliardetto/binary"
	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

type SolanaWallet struct {
	key *ed25519.PublicKey
}

var (
	_ Wallet   = &SolanaWallet{}
	_ TxParser = &SolanaWallet{}
)

func NewSolanaWallet(k *Key) (*SolanaWallet, error) {
	pubkey, err := k.ToEdDSAEd25519()
	if err != nil {
		return nil, err
	}
	return &SolanaWallet{key: pubkey}, nil
}

// Address returns a Solana address for the wallet - a string representation of the public key in base 58
// TODO: are the nil checks necessary?
func (w *SolanaWallet) Address() string {
	if *w.key == nil {
		panic("key is not set")
	}

	pk := solana.PublicKeyFromBytes(*w.key)

	if pk.IsZero() {
		panic("public key not parsed to solana pk")
	}

	return pk.String()
}

// ParseTx parses data from the raw bytes of an unsigned Solana transaction.
// A recent Solana blockhash is included in the serialized txBytes. This is valid
// for 150 blocks - approximately 1 minute after creation of the unsigned transaction.
func (*SolanaWallet) ParseTx(rawTx []byte, _ Metadata) (Transfer, error) {
	tx := &solana.Transaction{
		Message: solana.Message{},
	}
	if err := tx.Message.UnmarshalWithDecoder(bin.NewBinDecoder(rawTx)); err != nil {
		return Transfer{}, err
	}

	solTransfer, err := GetTransferFromInstruction(tx.Message)
	if err != nil {
		return Transfer{}, err
	}
	amount := new(big.Int)
	if solTransfer.Lamports != nil {
		amount.SetUint64(*solTransfer.Lamports)
	}

	to := solTransfer.GetRecipientAccount()
	receiverAddress := to.PublicKey

	coinIdentifier := []byte("SOL/")

	return Transfer{
		To:             []byte(receiverAddress.String()),
		Amount:         amount,
		CoinIdentifier: coinIdentifier,
		DataForSigning: rawTx,
	}, nil
}

// ParseSignedTx parses data from the raw bytes of a signed Solana transaction.
func (*SolanaWallet) ParseSignedTx(txBytes []byte, _ Metadata) (Transfer, error) {
	decodedTx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(txBytes))
	if err != nil {
		return Transfer{}, err
	}

	solTransfer, err := GetTransferFromInstruction(decodedTx.Message)
	if err != nil {
		return Transfer{}, err
	}

	amount := new(big.Int)
	if solTransfer.Lamports != nil {
		amount.SetUint64(*solTransfer.Lamports)
	}

	to := solTransfer.GetRecipientAccount()
	receiverAddress := to.PublicKey

	coinIdentifier := []byte("SOL/")

	return Transfer{
		To:             []byte(receiverAddress.String()),
		Amount:         amount,
		CoinIdentifier: coinIdentifier,
		DataForSigning: txBytes,
	}, nil
}

// GetTransferFromInstruction for a given solana.Message decodes the instruction and returns system.Transfer which contains from, to, amount
func GetTransferFromInstruction(msg solana.Message) (*system.Transfer, error) {
	for _, inst := range msg.Instructions {
		accounts, err := inst.ResolveInstructionAccounts(&msg)
		if err != nil {
			return nil, err
		}
		instruction, err := system.DecodeInstruction(accounts, inst.Data)
		if err != nil {
			return nil, err
		}
		if st, ok := instruction.Impl.(*system.Transfer); ok {
			return st, nil
		}
	}
	return nil, fmt.Errorf("no transfer instruction found")
}

// SolanaTransfer possibly not needed
// TODO Check
type SolanaTransfer struct {
	To             *common.Address
	Amount         *big.Int
	Contract       *common.Address
	DataForSigning []byte
}

// DecodeUnsignedSolanaPayload
func DecodeUnsignedSolanaPayload(msg []byte) (types.TxData, error) {
	panic("Not implemented")
}

// ParseSolanaTransaction is a placeholder
func ParseSolanaTransaction(b []byte, chainID *big.Int) (*SolanaTransfer, error) {
	panic("Not implemented")
}
