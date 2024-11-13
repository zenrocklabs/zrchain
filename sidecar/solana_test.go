package main_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/near/borsh-go"
)

// TODO: Keep this for now as a reference but remove this file later or add useful tests

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
