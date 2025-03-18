package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
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

	// Process each vote
	fmt.Printf("Found %d votes to decode\n", len(blockData.ConsensusData.Votes))
	for i, vote := range blockData.ConsensusData.Votes {
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
	}
}
