package eventstore

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

const (
	// EventStore program constants (must match on-chain program)
	// Separate targets because wrap capacity (10*100 = 1000) differs from unwrap (17*60 = 1020).
	TARGET_WRAP_EVENTS        = ZENBTC_WRAP_SHARD_COUNT * SHARD_SIZE_WRAP     // 1000
	TARGET_UNWRAP_EVENTS      = ZENBTC_UNWRAP_SHARD_COUNT * SHARD_SIZE_UNWRAP // 1020
	SHARD_SIZE_WRAP           = 100
	SHARD_SIZE_UNWRAP         = 60
	ZENBTC_WRAP_SHARD_COUNT   = 10
	ZENBTC_UNWRAP_SHARD_COUNT = 17
	ROCK_WRAP_SHARD_COUNT     = 10
	ROCK_UNWRAP_SHARD_COUNT   = 17

	// PDA seeds
	GLOBAL_CONFIG_SEED       = "global_config"
	ZENBTC_WRAP_SHARD_SEED   = "zenbtc_wrap"
	ZENBTC_UNWRAP_SHARD_SEED = "zenbtc_unwrap"
	ROCK_WRAP_SHARD_SEED     = "rock_wrap"
	ROCK_UNWRAP_SHARD_SEED   = "rock_unwrap"

	// Default EventStore program ID
	DEFAULT_PROGRAM_ID = "4KFjSTnBjbJbWXAiwpWjBCCfAKhjqMp3yfXYpoR3eVis"
)

// FlexibleAddress represents a variable-length Bitcoin address
type FlexibleAddress struct {
	Len  uint8     `borsh_struct:"true"`
	Data [63]uint8 `borsh_struct:"true"`
}

// String returns the Bitcoin address as a string
func (fa *FlexibleAddress) String() string {
	if fa.Len == 0 {
		return ""
	}
	return string(fa.Data[:fa.Len])
}

// TokensMintedWithFee represents a wrap event (shared by zenbtc and rock)
type TokensMintedWithFee struct {
	Recipient solana.PublicKey `borsh_struct:"true"`
	Value     uint64           `borsh_struct:"true"`
	Fee       uint64           `borsh_struct:"true"`
	Mint      solana.PublicKey `borsh_struct:"true"`
	ID        [16]uint8        `borsh_struct:"true"` // u128 as bytes
}

// GetID returns the event ID as uint64 (lower 64 bits)
func (t *TokensMintedWithFee) GetID() uint64 {
	return binary.LittleEndian.Uint64(t.ID[:8])
}

// ZenbtcTokenRedemption represents a zenbtc unwrap event
type ZenbtcTokenRedemption struct {
	Redeemer solana.PublicKey `borsh_struct:"true"`
	Value    uint64           `borsh_struct:"true"`
	DestAddr FlexibleAddress  `borsh_struct:"true"`
	Fee      uint64           `borsh_struct:"true"`
	Mint     solana.PublicKey `borsh_struct:"true"`
	ID       [16]uint8        `borsh_struct:"true"` // u128 as bytes
}

// GetID returns the event ID as uint64
func (z *ZenbtcTokenRedemption) GetID() uint64 {
	return binary.LittleEndian.Uint64(z.ID[:8])
}

// GetBitcoinAddress returns the destination Bitcoin address
func (z *ZenbtcTokenRedemption) GetBitcoinAddress() string {
	return z.DestAddr.String()
}

// RockTokenRedemption represents a rock unwrap event
type RockTokenRedemption struct {
	Redeemer solana.PublicKey `borsh_struct:"true"`
	Value    uint64           `borsh_struct:"true"`
	DestAddr FlexibleAddress  `borsh_struct:"true"`
	Fee      uint64           `borsh_struct:"true"`
	Mint     solana.PublicKey `borsh_struct:"true"`
	ID       [16]uint8        `borsh_struct:"true"` // u128 as bytes
}

// GetID returns the event ID as uint64
func (r *RockTokenRedemption) GetID() uint64 {
	return binary.LittleEndian.Uint64(r.ID[:8])
}

// GetBitcoinAddress returns the destination Bitcoin address
func (r *RockTokenRedemption) GetBitcoinAddress() string {
	return r.DestAddr.String()
}

// Shard account structures
type ZenbtcWrapShard struct {
	ShardIndex   uint16                `borsh_struct:"true"`
	EventsStored uint64                `borsh_struct:"true"`
	Events       []TokensMintedWithFee `borsh_struct:"true"`
}

type ZenbtcUnwrapShard struct {
	ShardIndex   uint16                  `borsh_struct:"true"`
	EventsStored uint64                  `borsh_struct:"true"`
	Events       []ZenbtcTokenRedemption `borsh_struct:"true"`
}

type RockWrapShard struct {
	ShardIndex   uint16                `borsh_struct:"true"`
	EventsStored uint64                `borsh_struct:"true"`
	Events       []TokensMintedWithFee `borsh_struct:"true"`
}

type RockUnwrapShard struct {
	ShardIndex   uint16                `borsh_struct:"true"`
	EventsStored uint64                `borsh_struct:"true"`
	Events       []RockTokenRedemption `borsh_struct:"true"`
}

// AllEvents represents all events from all shards
type AllEvents struct {
	ZenbtcWrapEvents   []TokensMintedWithFee   `json:"zenbtc_wrap_events"`
	ZenbtcUnwrapEvents []ZenbtcTokenRedemption `json:"zenbtc_unwrap_events"`
	RockWrapEvents     []TokensMintedWithFee   `json:"rock_wrap_events"`
	RockUnwrapEvents   []RockTokenRedemption   `json:"rock_unwrap_events"`
}

// Client provides access to the EventStore program
type Client struct {
	rpcClient *rpc.Client
	programID solana.PublicKey
}

// NewClient creates a new EventStore client
func NewClient(rpcClient *rpc.Client, programID *solana.PublicKey) *Client {
	var pid solana.PublicKey
	if programID != nil {
		pid = *programID
	} else {
		// Use default program ID
		pid = solana.MustPublicKeyFromBase58(DEFAULT_PROGRAM_ID)
	}

	return &Client{
		rpcClient: rpcClient,
		programID: pid,
	}
}

// DeriveGlobalConfigPDA returns the global config PDA for a given EventStore program ID.
func DeriveGlobalConfigPDA(programID solana.PublicKey) (solana.PublicKey, error) {
	pda, _, err := solana.FindProgramAddress([][]byte{[]byte(GLOBAL_CONFIG_SEED)}, programID)
	return pda, err
}

// DeriveWrapShardPDA deterministically derives the wrap shard PDA for a given event ID.
// eventID MUST be the canonical u128 (lower 64 bits used here) prior to submission.
// shardCount is the total number of wrap shards for the asset type (zenbtc or rock).
// seed must be one of the *_WRAP_SHARD_SEED constants.
func DeriveWrapShardPDA(programID solana.PublicKey, seed string, eventID uint64, shardCount uint16) (solana.PublicKey, uint16, error) {
	if shardCount == 0 {
		return solana.PublicKey{}, 0, fmt.Errorf("shardCount cannot be zero")
	}
	shardIndex := uint16(eventID % uint64(shardCount))
	idxBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(idxBytes, shardIndex)
	pda, _, err := solana.FindProgramAddress([][]byte{[]byte(seed), idxBytes}, programID)
	return pda, shardIndex, err
}

// DeriveShardPDA is a generic alias to DeriveWrapShardPDA for semantic clarity when the context
// is not strictly a wrap operation (e.g., unwrap shards). Maintained for readability.
func DeriveShardPDA(programID solana.PublicKey, seed string, eventID uint64, shardCount uint16) (solana.PublicKey, uint16, error) {
	return DeriveWrapShardPDA(programID, seed, eventID, shardCount)
}

// Convenience wrappers for common shard derivations -------------------------

func DeriveZenbtcWrapShardPDA(programID solana.PublicKey, eventID uint64) (solana.PublicKey, uint16, error) {
	return DeriveWrapShardPDA(programID, ZENBTC_WRAP_SHARD_SEED, eventID, ZENBTC_WRAP_SHARD_COUNT)
}

func DeriveZenbtcUnwrapShardPDA(programID solana.PublicKey, eventID uint64) (solana.PublicKey, uint16, error) {
	return DeriveWrapShardPDA(programID, ZENBTC_UNWRAP_SHARD_SEED, eventID, ZENBTC_UNWRAP_SHARD_COUNT)
}

func DeriveRockWrapShardPDA(programID solana.PublicKey, eventID uint64) (solana.PublicKey, uint16, error) {
	return DeriveWrapShardPDA(programID, ROCK_WRAP_SHARD_SEED, eventID, ROCK_WRAP_SHARD_COUNT)
}

func DeriveRockUnwrapShardPDA(programID solana.PublicKey, eventID uint64) (solana.PublicKey, uint16, error) {
	return DeriveWrapShardPDA(programID, ROCK_UNWRAP_SHARD_SEED, eventID, ROCK_UNWRAP_SHARD_COUNT)
}

// getShardPDA generates a shard PDA address
func (c *Client) getShardPDA(seed string, shardIndex uint16) (solana.PublicKey, error) {
	indexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(indexBytes, shardIndex)

	pda, _, err := solana.FindProgramAddress(
		[][]byte{[]byte(seed), indexBytes},
		c.programID,
	)
	return pda, err
}

// getAllShardAddresses returns all shard addresses for a given shard type
func (c *Client) getAllShardAddresses(seed string, shardCount uint16) ([]solana.PublicKey, error) {
	addresses := make([]solana.PublicKey, shardCount)

	for i := uint16(0); i < shardCount; i++ {
		addr, err := c.getShardPDA(seed, i)
		if err != nil {
			return nil, fmt.Errorf("failed to derive shard %d address: %w", i, err)
		}
		addresses[i] = addr
	}

	return addresses, nil
}

// GetAllEvents fetches all events from all shards in a single RPC call
func (c *Client) GetAllEvents(ctx context.Context) (*AllEvents, error) {
	// Get all shard addresses
	zenbtcWrapAddresses, err := c.getAllShardAddresses(ZENBTC_WRAP_SHARD_SEED, ZENBTC_WRAP_SHARD_COUNT)
	if err != nil {
		return nil, fmt.Errorf("failed to get zenbtc wrap addresses: %w", err)
	}

	zenbtcUnwrapAddresses, err := c.getAllShardAddresses(ZENBTC_UNWRAP_SHARD_SEED, ZENBTC_UNWRAP_SHARD_COUNT)
	if err != nil {
		return nil, fmt.Errorf("failed to get zenbtc unwrap addresses: %w", err)
	}

	rockWrapAddresses, err := c.getAllShardAddresses(ROCK_WRAP_SHARD_SEED, ROCK_WRAP_SHARD_COUNT)
	if err != nil {
		return nil, fmt.Errorf("failed to get rock wrap addresses: %w", err)
	}

	rockUnwrapAddresses, err := c.getAllShardAddresses(ROCK_UNWRAP_SHARD_SEED, ROCK_UNWRAP_SHARD_COUNT)
	if err != nil {
		return nil, fmt.Errorf("failed to get rock unwrap addresses: %w", err)
	}

	// Combine all addresses for batch fetch
	allAddresses := make([]solana.PublicKey, 0, len(zenbtcWrapAddresses)+len(zenbtcUnwrapAddresses)+len(rockWrapAddresses)+len(rockUnwrapAddresses))
	allAddresses = append(allAddresses, zenbtcWrapAddresses...)
	allAddresses = append(allAddresses, zenbtcUnwrapAddresses...)
	allAddresses = append(allAddresses, rockWrapAddresses...)
	allAddresses = append(allAddresses, rockUnwrapAddresses...)

	// Fetch all accounts in a single RPC call
	accounts, err := c.rpcClient.GetMultipleAccounts(ctx, allAddresses...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	result := &AllEvents{
		ZenbtcWrapEvents:   make([]TokensMintedWithFee, 0),
		ZenbtcUnwrapEvents: make([]ZenbtcTokenRedemption, 0),
		RockWrapEvents:     make([]TokensMintedWithFee, 0),
		RockUnwrapEvents:   make([]RockTokenRedemption, 0),
	}

	accountIndex := 0

	// Process zenbtc wrap shards
	for i := 0; i < len(zenbtcWrapAddresses); i++ {
		if accountIndex >= len(accounts.Value) || accounts.Value[accountIndex] == nil {
			accountIndex++
			continue // Skip non-existent shards
		}

		account := accounts.Value[accountIndex]
		accountIndex++

		// Skip account discriminator (8 bytes)
		if len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard ZenbtcWrapShard
		err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:])
		if err != nil {
			continue // Skip failed deserializations
		}

		result.ZenbtcWrapEvents = append(result.ZenbtcWrapEvents, shard.Events...)
	}

	// Process zenbtc unwrap shards
	for i := 0; i < len(zenbtcUnwrapAddresses); i++ {
		if accountIndex >= len(accounts.Value) || accounts.Value[accountIndex] == nil {
			accountIndex++
			continue
		}

		account := accounts.Value[accountIndex]
		accountIndex++

		if len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard ZenbtcUnwrapShard
		err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:])
		if err != nil {
			continue
		}

		result.ZenbtcUnwrapEvents = append(result.ZenbtcUnwrapEvents, shard.Events...)
	}

	// Process rock wrap shards
	for i := 0; i < len(rockWrapAddresses); i++ {
		if accountIndex >= len(accounts.Value) || accounts.Value[accountIndex] == nil {
			accountIndex++
			continue
		}

		account := accounts.Value[accountIndex]
		accountIndex++

		if len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard RockWrapShard
		err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:])
		if err != nil {
			continue
		}

		result.RockWrapEvents = append(result.RockWrapEvents, shard.Events...)
	}

	// Process rock unwrap shards
	for i := 0; i < len(rockUnwrapAddresses); i++ {
		if accountIndex >= len(accounts.Value) || accounts.Value[accountIndex] == nil {
			accountIndex++
			continue
		}

		account := accounts.Value[accountIndex]
		accountIndex++

		if len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard RockUnwrapShard
		err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:])
		if err != nil {
			continue
		}

		result.RockUnwrapEvents = append(result.RockUnwrapEvents, shard.Events...)
	}

	return result, nil
}

// GetZenbtcWrapEvents fetches only zenbtc wrap events
func (c *Client) GetZenbtcWrapEvents(ctx context.Context) ([]TokensMintedWithFee, error) {
	addresses, err := c.getAllShardAddresses(ZENBTC_WRAP_SHARD_SEED, ZENBTC_WRAP_SHARD_COUNT)
	if err != nil {
		return nil, err
	}

	accounts, err := c.rpcClient.GetMultipleAccounts(ctx, addresses...)
	if err != nil {
		return nil, err
	}

	var events []TokensMintedWithFee
	for _, account := range accounts.Value {
		if account == nil || len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard ZenbtcWrapShard
		if err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:]); err == nil {
			events = append(events, shard.Events...)
		}
	}

	return events, nil
}

// GetZenbtcUnwrapEvents fetches only zenbtc unwrap events with decoded Bitcoin addresses
func (c *Client) GetZenbtcUnwrapEvents(ctx context.Context) ([]ZenbtcTokenRedemption, error) {
	addresses, err := c.getAllShardAddresses(ZENBTC_UNWRAP_SHARD_SEED, ZENBTC_UNWRAP_SHARD_COUNT)
	if err != nil {
		return nil, err
	}

	accounts, err := c.rpcClient.GetMultipleAccounts(ctx, addresses...)
	if err != nil {
		return nil, err
	}

	var events []ZenbtcTokenRedemption
	for _, account := range accounts.Value {
		if account == nil || len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard ZenbtcUnwrapShard
		if err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:]); err == nil {
			events = append(events, shard.Events...)
		}
	}

	return events, nil
}

// GetRockWrapEvents fetches only rock wrap events
func (c *Client) GetRockWrapEvents(ctx context.Context) ([]TokensMintedWithFee, error) {
	addresses, err := c.getAllShardAddresses(ROCK_WRAP_SHARD_SEED, ROCK_WRAP_SHARD_COUNT)
	if err != nil {
		return nil, err
	}

	accounts, err := c.rpcClient.GetMultipleAccounts(ctx, addresses...)
	if err != nil {
		return nil, err
	}

	var events []TokensMintedWithFee
	for _, account := range accounts.Value {
		if account == nil || len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard RockWrapShard
		if err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:]); err == nil {
			events = append(events, shard.Events...)
		}
	}

	return events, nil
}

// GetRockUnwrapEvents fetches only rock unwrap events with decoded Bitcoin addresses
func (c *Client) GetRockUnwrapEvents(ctx context.Context) ([]RockTokenRedemption, error) {
	addresses, err := c.getAllShardAddresses(ROCK_UNWRAP_SHARD_SEED, ROCK_UNWRAP_SHARD_COUNT)
	if err != nil {
		return nil, err
	}

	accounts, err := c.rpcClient.GetMultipleAccounts(ctx, addresses...)
	if err != nil {
		return nil, err
	}

	var events []RockTokenRedemption
	for _, account := range accounts.Value {
		if account == nil || len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard RockUnwrapShard
		if err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:]); err == nil {
			events = append(events, shard.Events...)
		}
	}

	return events, nil
}

// GetEventCounts returns the current counts of events stored
func (c *Client) GetEventCounts(ctx context.Context) (map[string]int, error) {
	allEvents, err := c.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"zenbtc_wrap":   len(allEvents.ZenbtcWrapEvents),
		"zenbtc_unwrap": len(allEvents.ZenbtcUnwrapEvents),
		"rock_wrap":     len(allEvents.RockWrapEvents),
		"rock_unwrap":   len(allEvents.RockUnwrapEvents),
		"total":         len(allEvents.ZenbtcWrapEvents) + len(allEvents.ZenbtcUnwrapEvents) + len(allEvents.RockWrapEvents) + len(allEvents.RockUnwrapEvents),
	}, nil
}

// GetBitcoinAddressType returns a human-readable Bitcoin address type
func GetBitcoinAddressType(address string) string {
	if len(address) == 0 {
		return "empty"
	}

	// P2PKH (Legacy) - starts with 1
	if strings.HasPrefix(address, "1") && len(address) >= 26 && len(address) <= 35 {
		return "P2PKH (Legacy)"
	}

	// P2SH - starts with 3
	if strings.HasPrefix(address, "3") && len(address) >= 26 && len(address) <= 35 {
		return "P2SH"
	}

	// Bech32 (P2WPKH) - starts with bc1q and length ~42
	if strings.HasPrefix(address, "bc1q") && len(address) == 42 {
		return "P2WPKH (Bech32)"
	}

	// Bech32 (P2WSH) - starts with bc1q and length ~62
	if strings.HasPrefix(address, "bc1q") && len(address) == 62 {
		return "P2WSH (Bech32)"
	}

	// Bech32m (P2TR) - starts with bc1p
	// P2TR (Taproot/Bech32m) - starts with bc1p
	if strings.HasPrefix(address, "bc1p") && len(address) == 62 {
		return "P2TR (Taproot/Bech32m)"
	}

	// Testnet addresses
	if strings.HasPrefix(address, "m") || strings.HasPrefix(address, "n") {
		return "Testnet P2PKH"
	}
	if strings.HasPrefix(address, "2") {
		return "Testnet P2SH"
	}
	if strings.HasPrefix(address, "tb1q") && len(address) == 42 {
		return "Testnet P2WPKH (Bech32)"
	}
	if strings.HasPrefix(address, "tb1q") && len(address) == 62 {
		return "Testnet P2WSH (Bech32)"
	}
	if strings.HasPrefix(address, "tb1p") && len(address) == 62 {
		return "Testnet P2TR (Taproot)"
	}

	// Regtest addresses (bcrt1)
	if strings.HasPrefix(address, "bcrt1q") && len(address) == 44 {
		return "Regtest P2WPKH (Bech32)"
	}
	if strings.HasPrefix(address, "bcrt1q") && len(address) == 64 {
		return "Regtest P2WSH (Bech32)"
	}
	if strings.HasPrefix(address, "bcrt1p") && len(address) == 64 {
		return "Regtest P2TR (Taproot)"
	}

	return fmt.Sprintf("Unknown (%d chars)", len(address))
}
