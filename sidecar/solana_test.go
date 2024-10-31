package main_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	"github.com/stretchr/testify/require"

	"github.com/near/borsh-go"
)

func TestSolanaDecodeInstructionData(t *testing.T) {
	base64InstructionData := "SY2S0J4ZvbGKN12w/5pWlX2ZpLrVfIBdKhSoAPMGys6reiieVXHQ7w=="
	// Decode the instruction data
	decodedData, err := decodeInstructionData(base64InstructionData)
	if err != nil {
		log.Fatalf("Error decoding instruction data: %v", err)
	}

	// Deserialize the ZrKeyRequest instruction
	zrKeyRequest, err := decodeZrKeyRequestInstruction(decodedData)
	if err != nil {
		log.Fatalf("Error decoding ZrKeyRequest instruction: %v", err)
	}

	fmt.Printf("%+v\n", zrKeyRequest)
}

func TestSolanaDecodeTransaction(t *testing.T) {
	oracle := initTestOracle()
	oracleService := sidecar.NewOracleService(oracle)
	out, err := oracleService.DebugGetSolanaTransaction(context.Background(), &api.SolanaTransactionRequest{
		TxSignature: "4pym3sTe9dNcCM4ZCxLuwBZk46QicTWcCjaDNkh4pXS1cu4j3TMVUnx4AAVBPsoKjTmBM51oehyMFUMwLQbVK1JU",
	})
	require.NoError(t, err)
	require.NotNil(t, out.Tx)
	require.NotNil(t, out.Tx.Transaction)

	transaction, err := out.Tx.Transaction.GetTransaction()
	require.NoError(t, err)
	require.NotNil(t, transaction)

	message := transaction.Message
	require.NotNil(t, message)

	// Iterate over the instructions in the message
	for _, instruction := range message.Instructions {
		// Get the program ID to identify the target program
		programID := message.AccountKeys[instruction.ProgramIDIndex]

		data := instruction.Data
		fmt.Printf("Program ID: %s\n", programID.String())
		fmt.Printf("Instruction Data (base64): %s\n", data)
	}
}

type ZrKeyRequestArgs struct {
	WalletTypeID string
}

type ZrKeyRequestInstruction struct {
	WalletTypeID [32]byte // 32-byte array for wallet_type_id
}

func decodeInstructionData(base64Data string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}
	return decodedData, nil
}

func decodeZrKeyRequestInstruction(data []byte) (*ZrKeyRequestInstruction, error) {
	var instruction ZrKeyRequestInstruction
	err := borsh.Deserialize(&instruction, data)
	if err != nil {
		return nil, err
	}
	return &instruction, nil
}
