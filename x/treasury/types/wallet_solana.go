package types

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bin "github.com/gagliardetto/binary"
	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
)

type transfer struct {
	amount    *big.Int
	recipient string
}
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
func (*SolanaWallet) ParseTx(rawTx []byte, md Metadata) (Transfer, error) {
	tx := &solana.Transaction{
		Message: solana.Message{},
	}
	if err := tx.Message.UnmarshalWithDecoder(bin.NewBinDecoder(rawTx)); err != nil {
		return Transfer{}, err
	}

	solanaTx, err := extractTransferFromMessage(tx.Message)
	if err != nil {
		return Transfer{}, err
	}

	meta, ok := md.(*MetadataSolana)
	if !ok || meta == nil {
		return Transfer{}, fmt.Errorf("invalid metadata field, expected *MetadataSolana, got %T", md)
	}

	coinIdentifier := []byte(fmt.Sprintf("SOL/%s", meta.MintAddress))

	return Transfer{
		To:             []byte(solanaTx.recipient),
		Amount:         solanaTx.amount,
		CoinIdentifier: coinIdentifier,
		DataForSigning: []byte(hex.EncodeToString(rawTx)),
	}, nil
}

// ParseSignedTx parses data from the raw bytes of a signed Solana transaction.
func (*SolanaWallet) ParseSignedTx(txBytes []byte, md Metadata) (Transfer, error) {
	decodedTx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(txBytes))
	if err != nil {
		return Transfer{}, err
	}

	solanaTx, err := extractTransferFromMessage(decodedTx.Message)
	if err != nil {
		return Transfer{}, err
	}

	meta, ok := md.(*MetadataSolana)
	if !ok || meta == nil {
		return Transfer{}, fmt.Errorf("invalid metadata field, expected *MetadataSolana, got %T", md)
	}

	coinIdentifier := []byte(fmt.Sprintf("SOL/%s", meta.MintAddress))

	return Transfer{
		To:             []byte(solanaTx.recipient),
		Amount:         solanaTx.amount,
		CoinIdentifier: coinIdentifier,
		DataForSigning: txBytes,
	}, nil
}

// extractTransferFromMessage examines each instruction in a Solana message and
// returns a transfer (either a system or token transfer) that contains the recipient
// and the amount. If any transfer instruction is processed (indicated by
// amountAndRecipientRequired being true) but both recipient and amount are empty,
// an error is returned. This allows non-transfer transactions to pass without error.
func extractTransferFromMessage(msg solana.Message) (*transfer, error) {
	tx := &transfer{
		amount: new(big.Int),
	}

	amountAndRecipientRequired := false

	for _, inst := range msg.Instructions {
		accounts, err := inst.ResolveInstructionAccounts(&msg)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve instruction accounts: %w", err)
		}

		if int(inst.ProgramIDIndex) >= len(msg.AccountKeys) {
			continue
		}
		programID := msg.AccountKeys[inst.ProgramIDIndex]

		switch {
		case programID.Equals(solana.SystemProgramID):
			// Attempt to decode a system transfer.
			instruction, err := system.DecodeInstruction(accounts, inst.Data)
			if err != nil {
				continue
			}

			sysTransfer, ok := instruction.Impl.(*system.Transfer)
			if !ok {
				continue
			}

			if sysTransfer.Lamports != nil {
				tx.amount.SetUint64(*sysTransfer.Lamports)
			}
			tx.recipient = sysTransfer.GetRecipientAccount().PublicKey.String()
			amountAndRecipientRequired = true

		case programID.Equals(solana.TokenProgramID):
			// Attempt to decode a token transfer.
			instruction, err := token.DecodeInstruction(accounts, inst.Data)
			if err != nil {
				continue
			}

			tokenTransfer, ok := instruction.Impl.(*token.Transfer)
			if !ok {
				continue
			}

			tx.amount.SetUint64(*tokenTransfer.Amount)

			// For token transfers, the recipient is typically found at the second account index.
			if len(inst.Accounts) < 2 || int(inst.Accounts[1]) >= len(msg.AccountKeys) {
				continue
			}
			tx.recipient = msg.AccountKeys[inst.Accounts[1]].String()
			amountAndRecipientRequired = true

		default:
			continue
		}
	}

	if amountAndRecipientRequired {
		if len(tx.recipient) == 0 && tx.amount.Uint64() == 0 {
			return nil, fmt.Errorf("this transaction is a transfer with incomplete data")
		}
	}

	return tx, nil
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
