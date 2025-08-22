package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	eventstore "github.com/yourorg/eventstore-sdk"
)

func main() {
	// Create RPC client (using devnet for example)
	client := rpc.New(rpc.DevNet_RPC)

	// Create EventStore client with default program ID
	// For custom program ID, use: eventstore.NewClient(client, &customProgramID)
	esClient := eventstore.NewClient(client, nil)

	ctx := context.Background()

	fmt.Println("🚀 EventStore Go SDK Example")
	fmt.Println("=============================")

	// Example 1: Get all events in a single call
	fmt.Println("\n📊 Fetching all events...")
	allEvents, err := esClient.GetAllEvents(ctx)
	if err != nil {
		log.Printf("Error fetching all events: %v", err)
	} else {
		fmt.Printf("✅ Successfully fetched events:")
		fmt.Printf("   • ZenBTC Wrap Events: %d\n", len(allEvents.ZenbtcWrapEvents))
		fmt.Printf("   • ZenBTC Unwrap Events: %d\n", len(allEvents.ZenbtcUnwrapEvents))
		fmt.Printf("   • Rock Wrap Events: %d\n", len(allEvents.RockWrapEvents))
		fmt.Printf("   • Rock Unwrap Events: %d\n", len(allEvents.RockUnwrapEvents))

		// Show some sample events
		if len(allEvents.ZenbtcWrapEvents) > 0 {
			fmt.Printf("\n🔍 Sample ZenBTC Wrap Event:\n")
			printWrapEvent(allEvents.ZenbtcWrapEvents[0])
		}

		if len(allEvents.ZenbtcUnwrapEvents) > 0 {
			fmt.Printf("\n🔍 Sample ZenBTC Unwrap Event:\n")
			printUnwrapEvent(allEvents.ZenbtcUnwrapEvents[0])
		}
	}

	// Example 2: Get event counts
	fmt.Println("\n📈 Event Statistics...")
	counts, err := esClient.GetEventCounts(ctx)
	if err != nil {
		log.Printf("Error getting event counts: %v", err)
	} else {
		fmt.Println("Event counts:", counts)
	}

	// Example 3: Get only ZenBTC unwrap events with Bitcoin address analysis
	fmt.Println("\n🪙 Analyzing ZenBTC Unwrap Events...")
	zenbtcUnwraps, err := esClient.GetZenbtcUnwrapEvents(ctx)
	if err != nil {
		log.Printf("Error fetching ZenBTC unwrap events: %v", err)
	} else {
		fmt.Printf("Found %d ZenBTC unwrap events\n", len(zenbtcUnwraps))

		// Analyze Bitcoin address types
		addressTypes := make(map[string]int)
		for _, event := range zenbtcUnwraps {
			addr := event.GetBitcoinAddress()
			addrType := eventstore.GetBitcoinAddressType(addr)
			addressTypes[addrType]++
		}

		fmt.Println("Bitcoin address types distribution:")
		for addrType, count := range addressTypes {
			fmt.Printf("   • %s: %d addresses\n", addrType, count)
		}

		// Show detailed info for first few events
		for i, event := range zenbtcUnwraps {
			if i >= 3 { // Show only first 3
				break
			}
			fmt.Printf("\n   Event #%d:\n", i+1)
			fmt.Printf("     ID: %d\n", event.GetID())
			fmt.Printf("     Value: %d sats\n", event.Value)
			fmt.Printf("     Fee: %d sats\n", event.Fee)
			fmt.Printf("     Bitcoin Address: %s\n", event.GetBitcoinAddress())
			fmt.Printf("     Address Type: %s\n", eventstore.GetBitcoinAddressType(event.GetBitcoinAddress()))
			fmt.Printf("     Redeemer: %s\n", event.Redeemer.String())
		}
	}

	// Example 4: Get only Rock events
	fmt.Println("\n🪨 Fetching Rock Events...")
	rockWraps, err := esClient.GetRockWrapEvents(ctx)
	if err != nil {
		log.Printf("Error fetching Rock wrap events: %v", err)
	} else {
		fmt.Printf("Rock wrap events: %d\n", len(rockWraps))
	}

	rockUnwraps, err := esClient.GetRockUnwrapEvents(ctx)
	if err != nil {
		log.Printf("Error fetching Rock unwrap events: %v", err)
	} else {
		fmt.Printf("Rock unwrap events: %d\n", len(rockUnwraps))
	}

	// Example 5: JSON export
	fmt.Println("\n💾 JSON Export Example...")
	if allEvents != nil {
		jsonData, err := json.MarshalIndent(allEvents, "", "  ")
		if err != nil {
			log.Printf("Error marshaling to JSON: %v", err)
		} else {
			fmt.Printf("JSON export size: %d bytes\n", len(jsonData))

			// Save to file (optional)
			// os.WriteFile("events.json", jsonData, 0644)
			// fmt.Println("✅ Events saved to events.json")
		}
	}

	// Example 6: Real-time monitoring (polling)
	fmt.Println("\n🔄 Starting real-time monitoring (5 second intervals)...")
	fmt.Println("Press Ctrl+C to stop...")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lastCounts := make(map[string]int)

	for i := 0; i < 3; i++ { // Run for 3 iterations as example
		select {
		case <-ticker.C:
			currentCounts, err := esClient.GetEventCounts(ctx)
			if err != nil {
				log.Printf("Error in monitoring: %v", err)
				continue
			}

			fmt.Printf("⏰ %s - Current counts: ", time.Now().Format("15:04:05"))
			hasChanges := false

			for eventType, count := range currentCounts {
				if eventType == "total" {
					continue
				}

				lastCount := lastCounts[eventType]
				if count != lastCount {
					fmt.Printf("%s: %d (+%d) ", eventType, count, count-lastCount)
					hasChanges = true
				} else {
					fmt.Printf("%s: %d ", eventType, count)
				}
			}

			if !hasChanges {
				fmt.Printf("(no changes)")
			}
			fmt.Println()

			lastCounts = currentCounts
		}
	}

	fmt.Println("\n🎉 Example completed successfully!")
	fmt.Println("\n💡 Pro Tips:")
	fmt.Println("   • Use GetAllEvents() for efficient single-call data retrieval")
	fmt.Println("   • Bitcoin addresses are automatically decoded from FlexibleAddress")
	fmt.Println("   • The SDK handles all 1000+ events across multiple shards seamlessly")
	fmt.Println("   • Event IDs provide chronological ordering across shards")
	fmt.Println("   • All Bitcoin address formats (Legacy, P2SH, Bech32, Taproot) are supported")
}

func printWrapEvent(event eventstore.TokensMintedWithFee) {
	fmt.Printf("   Event ID: %d\n", event.GetID())
	fmt.Printf("   Recipient: %s\n", event.Recipient.String())
	fmt.Printf("   Value: %d\n", event.Value)
	fmt.Printf("   Fee: %d\n", event.Fee)
	fmt.Printf("   Mint: %s\n", event.Mint.String())
}

func printUnwrapEvent(event eventstore.ZenbtcTokenRedemption) {
	fmt.Printf("   Event ID: %d\n", event.GetID())
	fmt.Printf("   Redeemer: %s\n", event.Redeemer.String())
	fmt.Printf("   Value: %d\n", event.Value)
	fmt.Printf("   Fee: %d\n", event.Fee)
	fmt.Printf("   Bitcoin Address: %s\n", event.GetBitcoinAddress())
	fmt.Printf("   Address Type: %s\n", eventstore.GetBitcoinAddressType(event.GetBitcoinAddress()))
	fmt.Printf("   Mint: %s\n", event.Mint.String())
}
