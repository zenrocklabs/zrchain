package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
}

func main() {
	var input []byte
	var err error

	// If a file is provided, read from it
	if len(os.Args) > 1 {
		input, err = os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
	} else {
		// Otherwise read from stdin
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Error reading from stdin: %v\n", err)
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

	// Store all decoded extensions for comparison
	allExtensions := make([]map[string]interface{}, 0)
	validatorInfos := make([]ValidatorInfo, 0)

	// Process each vote
	fmt.Printf("Found %d votes to decode\n", len(blockData.ConsensusData.Votes))
	for i, vote := range blockData.ConsensusData.Votes {
		validatorInfo := ValidatorInfo{
			Address: vote.Validator.Address,
			Power:   vote.Validator.Power,
			Index:   i + 1,
		}

		fmt.Printf("\n=== Validator %d (Power: %d) ===\n", i+1, vote.Validator.Power)

		// Skip if vote extension is empty
		if vote.VoteExtension == "" {
			fmt.Println("No vote extension found")
			continue
		}

		// Decode base64 data
		decodedBytes, err := base64.StdEncoding.DecodeString(vote.VoteExtension)
		if err != nil {
			fmt.Printf("Error decoding base64: %v\n", err)
			continue
		}

		// Parse and pretty print the JSON
		var prettyJSON map[string]interface{}
		err = json.Unmarshal(decodedBytes, &prettyJSON)
		if err != nil {
			fmt.Printf("Error parsing decoded JSON: %v\n", err)
			fmt.Println("Raw decoded data:", string(decodedBytes))
			continue
		}

		// Format the JSON with indentation
		prettyJSONBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			continue
		}

		fmt.Println("Decoded vote extension:")
		fmt.Println(string(prettyJSONBytes))

		// Store for later comparison
		allExtensions = append(allExtensions, prettyJSON)
		validatorInfos = append(validatorInfos, validatorInfo)
	}

	// Compare extensions to find differences
	if len(allExtensions) > 1 {
		fmt.Println("\n\n=== DIFFERENCES BETWEEN VOTE EXTENSIONS ===")
		findDifferences(allExtensions, validatorInfos)
	}
}

func findDifferences(extensions []map[string]interface{}, validators []ValidatorInfo) {
	// Map to track differences by field
	differences := make(map[string]map[string][]ValidatorInfo)

	// Get all possible keys from all extensions
	allKeys := make(map[string]bool)
	for _, ext := range extensions {
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
			if val, ok := ext[key]; ok {
				// Convert value to string representation
				bytes, err := json.Marshal(val)
				if err == nil {
					valueStr = string(bytes)
				} else {
					valueStr = fmt.Sprintf("%v", val)
				}
			} else {
				valueStr = "<missing>"
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
			totalPower := 0
			for _, val := range validators {
				totalPower += val.Power
			}

			valueCounts = append(valueCounts, ValueCount{
				Value:      v,
				Validators: validators,
				TotalPower: totalPower,
			})
		}

		// Sort by power (descending)
		sort.Slice(valueCounts, func(i, j int) bool {
			return valueCounts[i].TotalPower > valueCounts[j].TotalPower
		})

		// Format for display
		for i, vc := range valueCounts {
			displayValue := vc.Value
			if len(displayValue) > 100 {
				displayValue = displayValue[:97] + "..."
			}

			valueLabel := "MAJORITY"
			if i > 0 {
				valueLabel = "MINORITY"
			}

			fmt.Printf("%s VALUE (Power: %d):\n", valueLabel, vc.TotalPower)
			fmt.Printf("  %s\n", displayValue)
			fmt.Printf("  Used by %d validators (%.2f%% of votes)\n",
				len(vc.Validators),
				float64(len(vc.Validators))/float64(len(validators))*100)

			fmt.Printf("  Validators: %s\n\n", formatValidatorList(vc.Validators, 3))
		}
	}
}

// Helper function to format a validator list nicely
func formatValidatorList(validators []ValidatorInfo, maxItems int) string {
	if len(validators) <= maxItems {
		var result []string
		for _, v := range validators {
			result = append(result, fmt.Sprintf("#%d (power: %d)", v.Index, v.Power))
		}
		return strings.Join(result, ", ")
	}

	var result []string
	for _, v := range validators[:maxItems] {
		result = append(result, fmt.Sprintf("#%d (power: %d)", v.Index, v.Power))
	}
	return strings.Join(result, ", ") + fmt.Sprintf(" and %d more...", len(validators)-maxItems)
}
