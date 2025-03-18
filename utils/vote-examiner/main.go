package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

// ConsensusVote represents a vote from a validator
type ConsensusVote struct {
	Validator struct {
		Address string `json:"address"`
		Power   int    `json:"power"`
	} `json:"validator"`
	VoteExtension      string `json:"vote_extension"`
	ExtensionSignature string `json:"extension_signature"`
	BlockIDFlag        int    `json:"block_id_flag"`
}

// BlockData represents the block data structure
type BlockData struct {
	ConsensusData struct {
		Votes []ConsensusVote `json:"votes"`
	} `json:"ConsensusData"`
}

// ValidatorInfo stores information about a validator
type ValidatorInfo struct {
	Address string
	Power   int
	Index   int
	HasVote bool
}

const NO_VOTE = "<no vote extension>"

func main() {
	// Command line flags
	useFilePtr := flag.String("file", "", "Use file instead of executing command (optional)")
	rpcNodePtr := flag.String("node", "https://rpc.diamond.zenrocklabs.io:443", "RPC node URL")
	blockHeightPtr := flag.String("height", "", "Block height (default: latest)")
	flag.Parse()

	var input []byte
	var err error

	// If a file is provided, read from it
	if *useFilePtr != "" {
		input, err = os.ReadFile(*useFilePtr)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
	} else {
		// Otherwise execute the command
		fmt.Println("Querying latest block data...")

		input, err = executeZenrockdCommand(*rpcNodePtr, *blockHeightPtr)
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			return
		}
	}

	// Convert to string for processing
	inputStr := string(input)

	// Extract the JSON object - find text between { and the last }
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	match := re.FindString(inputStr)
	if match == "" {
		fmt.Println("Could not find a JSON object in the input")
		return
	}

	// Parse the JSON data
	var blockData BlockData
	err = json.Unmarshal([]byte(match), &blockData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Calculate total validator power
	totalPower := 0
	for _, vote := range blockData.ConsensusData.Votes {
		totalPower += vote.Validator.Power
	}

	// Store all decoded extensions for comparison
	allExtensions := make([]map[string]interface{}, 0)
	allValidators := make([]ValidatorInfo, 0)
	validatorsWithExtensions := 0
	totalVotingPower := 0
	totalVotedPower := 0

	// Process each vote
	fmt.Printf("Found %d validators in the block\n", len(blockData.ConsensusData.Votes))
	for i, vote := range blockData.ConsensusData.Votes {
		validatorHasVote := false

		validatorInfo := ValidatorInfo{
			Address: vote.Validator.Address,
			Power:   vote.Validator.Power,
			Index:   i + 1,
			HasVote: false,
		}

		totalVotingPower += vote.Validator.Power

		fmt.Printf("\n=== Validator %d (Power: %d) ===\n", i+1, vote.Validator.Power)

		// Skip if vote extension is empty
		if vote.VoteExtension == "" {
			fmt.Println("No vote extension found")
		} else {
			validatorHasVote = true
			validatorsWithExtensions++
			totalVotedPower += vote.Validator.Power

			// Decode base64 data
			decodedBytes, err := base64.StdEncoding.DecodeString(vote.VoteExtension)
			if err != nil {
				fmt.Printf("Error decoding base64: %v\n", err)
			} else {
				// Parse and pretty print the JSON
				var prettyJSON map[string]interface{}
				err = json.Unmarshal(decodedBytes, &prettyJSON)
				if err != nil {
					fmt.Printf("Error parsing decoded JSON: %v\n", err)
					fmt.Println("Raw decoded data:", string(decodedBytes))
				} else {
					// Format the JSON with indentation
					prettyJSONBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
					if err != nil {
						fmt.Printf("Error formatting JSON: %v\n", err)
					} else {
						fmt.Println("Decoded vote extension:")
						fmt.Println(string(prettyJSONBytes))

						// Store extension for comparison
						allExtensions = append(allExtensions, prettyJSON)
					}
				}
			}
		}

		validatorInfo.HasVote = validatorHasVote
		allValidators = append(allValidators, validatorInfo)

		// If no extension, add a nil map to maintain index alignment
		if !validatorHasVote {
			allExtensions = append(allExtensions, nil)
		}
	}

	// Show vote extension participation stats
	fmt.Printf("\n=== VOTE EXTENSION PARTICIPATION ===\n")
	fmt.Printf("Validators with vote extensions: %d of %d (%.2f%%)\n",
		validatorsWithExtensions,
		len(blockData.ConsensusData.Votes),
		float64(validatorsWithExtensions)/float64(len(blockData.ConsensusData.Votes))*100)
	fmt.Printf("Voting power with extensions: %d of %d (%.2f%%)\n",
		totalVotedPower,
		totalVotingPower,
		float64(totalVotedPower)/float64(totalVotingPower)*100)

	// Compare extensions to find differences
	if len(allExtensions) > 0 {
		fmt.Println("\n\n=== DIFFERENCES BETWEEN VOTE EXTENSIONS ===")
		findDifferences(allExtensions, allValidators, totalVotingPower)
	}
}

// executeZenrockdCommand runs the zenrockd command and returns the output
func executeZenrockdCommand(rpcNode, blockHeight string) ([]byte, error) {
	var cmd *exec.Cmd

	if blockHeight == "" {
		// Use bash to execute the full command pipeline
		fullCmd := fmt.Sprintf(`zenrockd --node=%s query block --type=height $(zenrockd --node=%s status | jq -r '.sync_info.latest_block_height') --output=json | jq -r '.data.txs[0]' | base64 --decode`, rpcNode, rpcNode)
		cmd = exec.Command("bash", "-c", fullCmd)
	} else {
		// Use specific block height
		fullCmd := fmt.Sprintf(`zenrockd --node=%s query block --type=height %s --output=json | jq -r '.data.txs[0]' | base64 --decode`, rpcNode, blockHeight)
		cmd = exec.Command("bash", "-c", fullCmd)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("command execution failed: %s\nStderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

func findDifferences(extensions []map[string]interface{}, validators []ValidatorInfo, totalPower int) {
	// Map to track differences by field
	differences := make(map[string]map[string][]ValidatorInfo)

	// Get all possible keys from all extensions
	allKeys := make(map[string]bool)
	for _, ext := range extensions {
		if ext == nil {
			continue
		}
		for k := range ext {
			allKeys[k] = true
		}
	}

	// Check each key for differences
	for key := range allKeys {
		valueMap := make(map[string][]ValidatorInfo)

		for i, ext := range extensions {
			// Skip if index is out of range
			if i >= len(validators) {
				continue
			}

			var valueStr string
			if ext == nil {
				valueStr = NO_VOTE
			} else if val, ok := ext[key]; ok {
				// Convert value to string representation
				bytes, err := json.Marshal(val)
				if err == nil {
					valueStr = string(bytes)
				} else {
					valueStr = fmt.Sprintf("%v", val)
				}
			} else {
				valueStr = "<missing field>"
			}

			// Add validator to the list for this value
			valueMap[valueStr] = append(valueMap[valueStr], validators[i])
		}

		// If there's more than one unique value, we have a difference
		if len(valueMap) > 1 {
			differences[key] = valueMap
		}
	}

	// Display differences in a nice format
	if len(differences) == 0 {
		fmt.Println("No differences found between vote extensions")
		return
	}

	// Sort keys for consistent output
	sortedKeys := make([]string, 0, len(differences))
	for k := range differences {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		fmt.Printf("\nDifferences in field: %s\n", key)
		fmt.Printf("-------------------------\n")

		valueMap := differences[key]

		// Sort values by number of validators (descending) and power
		type ValueCount struct {
			Value      string
			Validators []ValidatorInfo
			TotalPower int
		}

		valueCounts := make([]ValueCount, 0, len(valueMap))
		for v, validators := range valueMap {
			thisPower := 0
			for _, val := range validators {
				thisPower += val.Power
			}

			valueCounts = append(valueCounts, ValueCount{
				Value:      v,
				Validators: validators,
				TotalPower: thisPower,
			})
		}

		// Sort by power (descending)
		sort.Slice(valueCounts, func(i, j int) bool {
			return valueCounts[i].TotalPower > valueCounts[j].TotalPower
		})

		// Format for display
		for i, vc := range valueCounts {
			displayValue := vc.Value
			if displayValue == NO_VOTE {
				displayValue = "NO VOTE EXTENSION SUBMITTED"
			} else if len(displayValue) > 100 {
				displayValue = displayValue[:97] + "..."
			}

			valueLabel := "MAJORITY"
			if i > 0 {
				valueLabel = "MINORITY"
			}

			powerPercentage := float64(vc.TotalPower) / float64(totalPower) * 100

			fmt.Printf("%s VALUE (Power: %d, %.2f%% of total)\n", valueLabel, vc.TotalPower, powerPercentage)
			fmt.Printf("  %s\n", displayValue)
			fmt.Printf("  Supported by %d validators (%.2f%% of validators)\n",
				len(vc.Validators),
				float64(len(vc.Validators))/float64(len(validators))*100)

			// Count validators with and without votes
			votedCount := 0
			for _, v := range vc.Validators {
				if v.HasVote {
					votedCount++
				}
			}

			noVoteCount := len(vc.Validators) - votedCount
			if displayValue != "NO VOTE EXTENSION SUBMITTED" && noVoteCount > 0 {
				fmt.Printf("  Note: %d validators in this group did not submit a vote extension\n", noVoteCount)
			}

			fmt.Printf("  Validators: %s\n\n", formatValidatorList(vc.Validators, 5))
		}
	}
}

// Helper function to format a validator list nicely
func formatValidatorList(validators []ValidatorInfo, maxItems int) string {
	if len(validators) <= maxItems {
		var result []string
		for _, v := range validators {
			tag := ""
			if !v.HasVote {
				tag = " (no vote)"
			}
			result = append(result, fmt.Sprintf("#%d (power: %d)%s", v.Index, v.Power, tag))
		}
		return strings.Join(result, ", ")
	}

	var result []string
	for _, v := range validators[:maxItems] {
		tag := ""
		if !v.HasVote {
			tag = " (no vote)"
		}
		result = append(result, fmt.Sprintf("#%d (power: %d)%s", v.Index, v.Power, tag))
	}
	return strings.Join(result, ", ") + fmt.Sprintf(" and %d more...", len(validators)-maxItems)
}
