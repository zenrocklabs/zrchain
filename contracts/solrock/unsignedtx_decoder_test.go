package solrock

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
)

// TransactionJSON represents the decoded transaction in JSON format
type TransactionJSON struct {
	RecentBlockhash string            `json:"recent_blockhash"`
	FeePayer        string            `json:"fee_payer"`
	NumInstructions int               `json:"num_instructions"`
	NumAccounts     int               `json:"num_accounts"`
	Accounts        []string          `json:"accounts"`
	Instructions    []InstructionJSON `json:"instructions"`
}

// InstructionJSON represents a decoded instruction in JSON format
type InstructionJSON struct {
	Index     int                    `json:"index"`
	ProgramID string                 `json:"program_id"`
	Data      string                 `json:"data"` // hex encoded
	Accounts  []string               `json:"accounts"`
	Type      string                 `json:"type"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// decodeTransactionToJSON decodes a transaction and returns it as JSON
func decodeTransactionToJSON(t *testing.T, base64Tx string) (*TransactionJSON, error) {
	// Decode base64
	txBytes, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Parse the transaction message (unsigned transaction)
	var msg solana.Message
	err = msg.UnmarshalWithDecoder(bin.NewBinDecoder(txBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction message: %w", err)
	}

	// Convert to JSON structure
	txJSON := &TransactionJSON{
		RecentBlockhash: msg.RecentBlockhash.String(),
		FeePayer:        msg.AccountKeys[msg.Header.NumRequiredSignatures-1].String(),
		NumInstructions: len(msg.Instructions),
		NumAccounts:     len(msg.AccountKeys),
		Accounts:        make([]string, len(msg.AccountKeys)),
		Instructions:    make([]InstructionJSON, len(msg.Instructions)),
	}

	// Convert accounts
	for i, acc := range msg.AccountKeys {
		txJSON.Accounts[i] = acc.String()
	}

	// Convert instructions
	for i, inst := range msg.Instructions {
		programID := msg.AccountKeys[inst.ProgramIDIndex]

		// Get account addresses
		accounts := make([]string, len(inst.Accounts))
		for j, accIdx := range inst.Accounts {
			accounts[j] = msg.AccountKeys[accIdx].String()
		}

		instructionJSON := InstructionJSON{
			Index:     i,
			ProgramID: programID.String(),
			Data:      hex.EncodeToString(inst.Data),
			Accounts:  accounts,
			Type:      "unknown",
			Details:   make(map[string]interface{}),
		}

		// Identify instruction type and add details
		switch programID {
		case system.ProgramID:
			if len(inst.Data) > 0 {
				discriminator := inst.Data[0]
				switch discriminator {
				case 4:
					instructionJSON.Type = "system_advance_nonce"
					if len(accounts) >= 3 {
						instructionJSON.Details["nonce_account"] = accounts[0]
						instructionJSON.Details["recent_blockhash"] = accounts[1]
						instructionJSON.Details["nonce_authority"] = accounts[2]
					}
				default:
					instructionJSON.Type = "system_unknown"
					instructionJSON.Details["discriminator"] = discriminator
				}
			}
		case token.ProgramID:
			instructionJSON.Type = "token_program"
		case ata.ProgramID:
			if len(inst.Data) > 0 {
				discriminator := inst.Data[0]
				switch discriminator {
				case 1:
					instructionJSON.Type = "ata_create"
					if len(accounts) >= 4 {
						instructionJSON.Details["payer"] = accounts[0]
						instructionJSON.Details["owner"] = accounts[1]
						instructionJSON.Details["mint"] = accounts[2]
						instructionJSON.Details["ata"] = accounts[3]
					}
				default:
					instructionJSON.Type = "ata_unknown"
					instructionJSON.Details["discriminator"] = discriminator
				}
			}
		default:
			// Check if it's a solrock instruction
			if len(inst.Data) >= 8 {
				discriminator := inst.Data[:8]
				wrapDiscriminator := []byte{178, 40, 10, 189, 228, 129, 186, 140} // Instruction_Wrap

				if string(discriminator) == string(wrapDiscriminator) {
					instructionJSON.Type = "solrock_wrap"

					// Try to decode wrap args
					if len(inst.Data) >= 8 {
						decoder := bin.NewBinDecoder(inst.Data[8:])
						var value uint64
						var fee uint64

						if err := decoder.Decode(&value); err == nil {
							instructionJSON.Details["value"] = value
						}
						if err := decoder.Decode(&fee); err == nil {
							instructionJSON.Details["fee"] = fee
						}
					}

					// Add account details
					if len(accounts) >= 10 {
						instructionJSON.Details["signer"] = accounts[0]
						instructionJSON.Details["global_config"] = accounts[1]
						instructionJSON.Details["mint"] = accounts[2]
						instructionJSON.Details["fee_wallet"] = accounts[3]
						instructionJSON.Details["fee_wallet_ata"] = accounts[4]
						instructionJSON.Details["receiver"] = accounts[5]
						instructionJSON.Details["receiver_ata"] = accounts[6]
						instructionJSON.Details["system_program"] = accounts[7]
						instructionJSON.Details["token_program"] = accounts[8]
						instructionJSON.Details["associated_token_program"] = accounts[9]
					}
				} else {
					instructionJSON.Type = "unknown_program"
					instructionJSON.Details["discriminator"] = hex.EncodeToString(discriminator)
				}
			} else {
				instructionJSON.Type = "unknown_program"
			}
		}

		txJSON.Instructions[i] = instructionJSON
	}

	return txJSON, nil
}

// Simple decoder for testing
func decodeTransaction(t *testing.T, base64Tx string) {
	// Decode base64
	txBytes, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	fmt.Printf("Transaction bytes length: %d\n", len(txBytes))
	fmt.Printf("Transaction bytes (hex): %s\n", hex.EncodeToString(txBytes))

	// Parse the transaction message (unsigned transaction)
	var msg solana.Message
	err = msg.UnmarshalWithDecoder(bin.NewBinDecoder(txBytes))
	if err != nil {
		t.Fatalf("Failed to parse transaction message: %v", err)
	}

	fmt.Printf("\n=== Transaction Details ===\n")
	fmt.Printf("Recent Blockhash: %s\n", msg.RecentBlockhash.String())
	fmt.Printf("Fee Payer: %s\n", msg.AccountKeys[msg.Header.NumRequiredSignatures-1].String())
	fmt.Printf("Number of Instructions: %d\n", len(msg.Instructions))
	fmt.Printf("Number of Accounts: %d\n", len(msg.AccountKeys))

	// Print all accounts
	fmt.Printf("\n=== All Accounts ===\n")
	for i, acc := range msg.AccountKeys {
		fmt.Printf("[%d] %s\n", i, acc.String())
	}

	// Decode each instruction
	fmt.Printf("\n=== Instructions ===\n")
	for i, inst := range msg.Instructions {
		fmt.Printf("\nInstruction %d:\n", i)
		fmt.Printf("  Program ID: %s\n", msg.AccountKeys[inst.ProgramIDIndex].String())
		fmt.Printf("  Data (hex): %s\n", hex.EncodeToString(inst.Data))
		fmt.Printf("  Data length: %d\n", len(inst.Data))

		// Get account addresses
		fmt.Printf("  Accounts:\n")
		for j, accIdx := range inst.Accounts {
			fmt.Printf("    [%d] %s\n", j, msg.AccountKeys[accIdx].String())
		}

		// Try to identify instruction type
		programID := msg.AccountKeys[inst.ProgramIDIndex]
		switch programID {
		case system.ProgramID:
			if len(inst.Data) > 0 {
				discriminator := inst.Data[0]
				switch discriminator {
				case 4:
					fmt.Printf("  Type: System - Advance Nonce Account\n")
				default:
					fmt.Printf("  Type: System - Unknown (discriminator: %d)\n", discriminator)
				}
			}
		case token.ProgramID:
			fmt.Printf("  Type: Token Program\n")
		case ata.ProgramID:
			if len(inst.Data) > 0 {
				discriminator := inst.Data[0]
				switch discriminator {
				case 1:
					fmt.Printf("  Type: ATA - Create\n")
				default:
					fmt.Printf("  Type: ATA - Unknown (discriminator: %d)\n", discriminator)
				}
			}
		default:
			// Check if it's a solrock instruction
			if len(inst.Data) >= 8 {
				discriminator := inst.Data[:8]
				wrapDiscriminator := []byte{178, 40, 10, 189, 228, 129, 186, 140} // Instruction_Wrap

				if string(discriminator) == string(wrapDiscriminator) {
					fmt.Printf("  Type: Solrock - Wrap\n")

					// Try to decode wrap args
					if len(inst.Data) >= 8 {
						decoder := bin.NewBinDecoder(inst.Data[8:])
						var value uint64
						var fee uint64

						if err := decoder.Decode(&value); err == nil {
							fmt.Printf("    Value: %d\n", value)
						}
						if err := decoder.Decode(&fee); err == nil {
							fmt.Printf("    Fee: %d\n", fee)
						}
					}
				} else {
					fmt.Printf("  Type: Unknown Program - %s\n", programID.String())
					fmt.Printf("    Discriminator: %s\n", hex.EncodeToString(discriminator))
				}
			} else {
				fmt.Printf("  Type: Unknown Program - %s\n", programID.String())
			}
		}
	}
}

func TestDecodeUnsignedTx(t *testing.T) {
	// The unsigned transaction from the user
	base64Tx := "AgEHD0Khtw04VpiJjlQSixKwyHIX9ZlSFK78hJKUsHHY2ghCXaXGvr5ZNUkv/FBXvHfzVr1CUOToih8xQuiXa6BW/5CFLrOPaN3y7Bz1/SKzLyfnWqedMsgzkMVFqLGJovCjrF662Yr4MTFK1cJ6SBti+Z5uasr+2CSDciH+bzHdQ/pO2FDFWsTkA1QPZ3On31P6QdIQJbL7OEv7PqCZK/6DfPhENuh06HvQLjxa9a4DxAo9ff/Keux5qFBJ11abAKzAZFF1mvHQF/cT/OcJUI6K9O2uJIBUb8lzN1V+Js94B5LC7UveGJCnRSPrgV8HK9YRW408Mi2vEIJfukdJJ/T2g1IGp9UXGSxWjuCKhF9z0peIzwNcMUWyGrNE2AYuqUAAABU2HXi28k+dQwmEPRbXuQ4qq0pKULDxnnnaXLiOCLBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAG3fbh12Whk9nL4UbO63msHLSF7V9bN5E6jPWFfv8AqQan1RcZLFxRIYzJTD1K8X9Y2u4Im6H9ROPb2YoAAAAAjJclj04kifG7PRApFI4NgwtaE5na/xCEBI572Nvp+Fk5AedvJ0aNSNvw/Fsllvmbo06Xv1VeYtxIS+w7wwqxUP9B0pe3Cg4wuwPdLmQviOJtC6g1FJPq/lRjp2VCAfT8AwoDAggBBAQAAAANBwADCQUKCwwADgoABAUGBwkDCgsNGLIoCr3kgbqMgJaYAAAAAAAAAAAAAAAAAA=="

	// Decode to JSON
	txJSON, err := decodeTransactionToJSON(t, base64Tx)
	if err != nil {
		t.Fatalf("Failed to decode transaction to JSON: %v", err)
	}

	// Output JSON
	jsonBytes, err := json.MarshalIndent(txJSON, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	fmt.Printf("\n=== Transaction JSON ===\n")
	fmt.Println(string(jsonBytes))

	// Also run the original decoder for comparison
	decodeTransaction(t, base64Tx)
}
