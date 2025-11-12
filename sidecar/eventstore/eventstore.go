package eventstore

import (
	"context"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

const (
	// EventStore program constants (must match on-chain program)
	TARGET_EVENTS_PER_TYPE     = 1000
	SHARD_SIZE_WRAP            = 100
	SHARD_SIZE_UNWRAP          = 60
	ZENZEC_WRAP_SHARD_COUNT    = 10 // ZenZEC uses "zenbtc_wrap" seed
	ZENZEC_UNWRAP_SHARD_COUNT  = 17 // ZenZEC uses "zenbtc_unwrap" seed
	ZENBTC2_WRAP_SHARD_COUNT   = 10 // ZenBTC uses "zenbtc2_wrap" seed
	ZENBTC2_UNWRAP_SHARD_COUNT = 17 // ZenBTC uses "zenbtc2_unwrap" seed
	ROCK_WRAP_SHARD_COUNT      = 10
	ROCK_UNWRAP_SHARD_COUNT    = 17

	// PDA seeds - IMPORTANT: These map to the Rust on-chain program seeds
	GLOBAL_CONFIG_SEED        = "global_config"
	ZENZEC_WRAP_SHARD_SEED    = "zenbtc_wrap"    // ZenZEC wrap events
	ZENZEC_UNWRAP_SHARD_SEED  = "zenbtc_unwrap"  // ZenZEC unwrap events
	ZENBTC2_WRAP_SHARD_SEED   = "zenbtc2_wrap"   // ZenBTC wrap events
	ZENBTC2_UNWRAP_SHARD_SEED = "zenbtc2_unwrap" // ZenBTC unwrap events
	ROCK_WRAP_SHARD_SEED      = "rock_wrap"
	ROCK_UNWRAP_SHARD_SEED    = "rock_unwrap"

	// Default EventStore program ID
	DEFAULT_PROGRAM_ID = "2BZ3Vi9BurkVJv5wX8H9QSxQasDJ42FVFRNS4vXSYf22"
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

// Slot wrappers retain occupancy metadata for circular buffers
type TokensMintedWithFeeSlot struct {
	Initialized uint8               `borsh_struct:"true"`
	Event       TokensMintedWithFee `borsh_struct:"true"`
}

func (s TokensMintedWithFeeSlot) IsInitialized() bool {
	return s.Initialized != 0
}

type ZenbtcTokenRedemptionSlot struct {
	Initialized uint8                 `borsh_struct:"true"`
	Event       ZenbtcTokenRedemption `borsh_struct:"true"`
}

func (s ZenbtcTokenRedemptionSlot) IsInitialized() bool {
	return s.Initialized != 0
}

type RockTokenRedemptionSlot struct {
	Initialized uint8               `borsh_struct:"true"`
	Event       RockTokenRedemption `borsh_struct:"true"`
}

func (s RockTokenRedemptionSlot) IsInitialized() bool {
	return s.Initialized != 0
}

// Shard account structures
type ZenzecWrapShard struct {
	ShardIndex   uint16                    `borsh_struct:"true"`
	EventsStored uint64                    `borsh_struct:"true"`
	Events       []TokensMintedWithFeeSlot `borsh_struct:"true"`
}

type ZenzecUnwrapShard struct {
	ShardIndex   uint16                      `borsh_struct:"true"`
	EventsStored uint64                      `borsh_struct:"true"`
	Events       []ZenbtcTokenRedemptionSlot `borsh_struct:"true"` // Same structure for ZenZEC
}

type Zenbtc2WrapShard struct {
	ShardIndex   uint16                    `borsh_struct:"true"`
	EventsStored uint64                    `borsh_struct:"true"`
	Events       []TokensMintedWithFeeSlot `borsh_struct:"true"`
}

type Zenbtc2UnwrapShard struct {
	ShardIndex   uint16                      `borsh_struct:"true"`
	EventsStored uint64                      `borsh_struct:"true"`
	Events       []ZenbtcTokenRedemptionSlot `borsh_struct:"true"`
}

type RockWrapShard struct {
	ShardIndex   uint16                    `borsh_struct:"true"`
	EventsStored uint64                    `borsh_struct:"true"`
	Events       []TokensMintedWithFeeSlot `borsh_struct:"true"`
}

type RockUnwrapShard struct {
	ShardIndex   uint16                    `borsh_struct:"true"`
	EventsStored uint64                    `borsh_struct:"true"`
	Events       []RockTokenRedemptionSlot `borsh_struct:"true"`
}

// AllEvents represents all events from all shards
type AllEvents struct {
	ZenzecWrapEvents   []TokensMintedWithFee   `json:"zenzec_wrap_events"`
	ZenzecUnwrapEvents []ZenbtcTokenRedemption `json:"zenzec_unwrap_events"`
	Zenbtc2WrapEvents  []TokensMintedWithFee   `json:"zenbtc2_wrap_events"`
	Zenbtc2UnwrapEvents []ZenbtcTokenRedemption `json:"zenbtc2_unwrap_events"`
	RockWrapEvents     []TokensMintedWithFee   `json:"rock_wrap_events"`
	RockUnwrapEvents   []RockTokenRedemption   `json:"rock_unwrap_events"`
}

func appendWrapEvents(dst []TokensMintedWithFee, slots []TokensMintedWithFeeSlot) []TokensMintedWithFee {
	for _, slot := range slots {
		if slot.IsInitialized() {
			dst = append(dst, slot.Event)
		}
	}
	return dst
}

func appendZenbtcUnwrapEvents(dst []ZenbtcTokenRedemption, slots []ZenbtcTokenRedemptionSlot) []ZenbtcTokenRedemption {
	for _, slot := range slots {
		if slot.IsInitialized() {
			dst = append(dst, slot.Event)
		}
	}
	return dst
}

func appendRockUnwrapEvents(dst []RockTokenRedemption, slots []RockTokenRedemptionSlot) []RockTokenRedemption {
	for _, slot := range slots {
		if slot.IsInitialized() {
			dst = append(dst, slot.Event)
		}
	}
	return dst
}

func compareUint128Bytes(a [16]uint8, b [16]uint8) int {
	hiA := binary.LittleEndian.Uint64(a[8:])
	hiB := binary.LittleEndian.Uint64(b[8:])
	if hiA < hiB {
		return -1
	}
	if hiA > hiB {
		return 1
	}
	loA := binary.LittleEndian.Uint64(a[:8])
	loB := binary.LittleEndian.Uint64(b[:8])
	switch {
	case loA < loB:
		return -1
	case loA > loB:
		return 1
	default:
		return 0
	}
}

func sortWrapEventSlice(events []TokensMintedWithFee) {
	sort.Slice(events, func(i, j int) bool {
		return compareUint128Bytes(events[i].ID, events[j].ID) < 0
	})
}

func sortZenbtcUnwrapEventSlice(events []ZenbtcTokenRedemption) {
	sort.Slice(events, func(i, j int) bool {
		return compareUint128Bytes(events[i].ID, events[j].ID) < 0
	})
}

func sortRockUnwrapEventSlice(events []RockTokenRedemption) {
	sort.Slice(events, func(i, j int) bool {
		return compareUint128Bytes(events[i].ID, events[j].ID) < 0
	})
}

func dedupWrapEventSlice(events []TokensMintedWithFee) []TokensMintedWithFee {
	if len(events) < 2 {
		return events
	}

	out := events[:1]
	for _, ev := range events[1:] {
		if compareUint128Bytes(ev.ID, out[len(out)-1].ID) != 0 {
			out = append(out, ev)
		}
	}
	return out
}

func dedupZenbtcUnwrapEventSlice(events []ZenbtcTokenRedemption) []ZenbtcTokenRedemption {
	if len(events) < 2 {
		return events
	}

	out := events[:1]
	for _, ev := range events[1:] {
		if compareUint128Bytes(ev.ID, out[len(out)-1].ID) != 0 {
			out = append(out, ev)
		}
	}
	return out
}

func dedupRockUnwrapEventSlice(events []RockTokenRedemption) []RockTokenRedemption {
	if len(events) < 2 {
		return events
	}

	out := events[:1]
	for _, ev := range events[1:] {
		if compareUint128Bytes(ev.ID, out[len(out)-1].ID) != 0 {
			out = append(out, ev)
		}
	}
	return out
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
	// Get all shard addresses for all asset types
	zenzecWrapAddresses, err := c.getAllShardAddresses(ZENZEC_WRAP_SHARD_SEED, ZENZEC_WRAP_SHARD_COUNT)
	if err != nil {
		return nil, fmt.Errorf("failed to get zenzec wrap addresses: %w", err)
	}

	zenzecUnwrapAddresses, err := c.getAllShardAddresses(ZENZEC_UNWRAP_SHARD_SEED, ZENZEC_UNWRAP_SHARD_COUNT)
	if err != nil {
		return nil, fmt.Errorf("failed to get zenzec unwrap addresses: %w", err)
	}

	zenbtc2WrapAddresses, err := c.getAllShardAddresses(ZENBTC2_WRAP_SHARD_SEED, ZENBTC2_WRAP_SHARD_COUNT)
	if err != nil {
		return nil, fmt.Errorf("failed to get zenbtc2 wrap addresses: %w", err)
	}

	zenbtc2UnwrapAddresses, err := c.getAllShardAddresses(ZENBTC2_UNWRAP_SHARD_SEED, ZENBTC2_UNWRAP_SHARD_COUNT)
	if err != nil {
		return nil, fmt.Errorf("failed to get zenbtc2 unwrap addresses: %w", err)
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
	allAddresses := make([]solana.PublicKey, 0,
		len(zenzecWrapAddresses)+len(zenzecUnwrapAddresses)+
			len(zenbtc2WrapAddresses)+len(zenbtc2UnwrapAddresses)+
			len(rockWrapAddresses)+len(rockUnwrapAddresses))
	allAddresses = append(allAddresses, zenzecWrapAddresses...)
	allAddresses = append(allAddresses, zenzecUnwrapAddresses...)
	allAddresses = append(allAddresses, zenbtc2WrapAddresses...)
	allAddresses = append(allAddresses, zenbtc2UnwrapAddresses...)
	allAddresses = append(allAddresses, rockWrapAddresses...)
	allAddresses = append(allAddresses, rockUnwrapAddresses...)

	// Fetch all accounts in a single RPC call
	accounts, err := c.rpcClient.GetMultipleAccounts(ctx, allAddresses...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	result := &AllEvents{
		ZenzecWrapEvents:    make([]TokensMintedWithFee, 0),
		ZenzecUnwrapEvents:  make([]ZenbtcTokenRedemption, 0),
		Zenbtc2WrapEvents:   make([]TokensMintedWithFee, 0),
		Zenbtc2UnwrapEvents: make([]ZenbtcTokenRedemption, 0),
		RockWrapEvents:      make([]TokensMintedWithFee, 0),
		RockUnwrapEvents:    make([]RockTokenRedemption, 0),
	}

	accountIndex := 0

	// Process zenzec wrap shards
	for i := 0; i < len(zenzecWrapAddresses); i++ {
		if accountIndex >= len(accounts.Value) || accounts.Value[accountIndex] == nil {
			accountIndex++
			continue
		}

		account := accounts.Value[accountIndex]
		accountIndex++

		if len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard ZenzecWrapShard
		err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:])
		if err != nil {
			continue
		}

		result.ZenzecWrapEvents = appendWrapEvents(result.ZenzecWrapEvents, shard.Events)
	}

	// Process zenzec unwrap shards
	for i := 0; i < len(zenzecUnwrapAddresses); i++ {
		if accountIndex >= len(accounts.Value) || accounts.Value[accountIndex] == nil {
			accountIndex++
			continue
		}

		account := accounts.Value[accountIndex]
		accountIndex++

		if len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard ZenzecUnwrapShard
		err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:])
		if err != nil {
			continue
		}

		result.ZenzecUnwrapEvents = appendZenbtcUnwrapEvents(result.ZenzecUnwrapEvents, shard.Events)
	}

	// Process zenbtc2 wrap shards
	for i := 0; i < len(zenbtc2WrapAddresses); i++ {
		if accountIndex >= len(accounts.Value) || accounts.Value[accountIndex] == nil {
			accountIndex++
			continue
		}

		account := accounts.Value[accountIndex]
		accountIndex++

		if len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard Zenbtc2WrapShard
		err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:])
		if err != nil {
			continue
		}

		result.Zenbtc2WrapEvents = appendWrapEvents(result.Zenbtc2WrapEvents, shard.Events)
	}

	// Process zenbtc2 unwrap shards
	for i := 0; i < len(zenbtc2UnwrapAddresses); i++ {
		if accountIndex >= len(accounts.Value) || accounts.Value[accountIndex] == nil {
			accountIndex++
			continue
		}

		account := accounts.Value[accountIndex]
		accountIndex++

		if len(account.Data.GetBinary()) < 8 {
			continue
		}

		var shard Zenbtc2UnwrapShard
		err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:])
		if err != nil {
			continue
		}

		result.Zenbtc2UnwrapEvents = appendZenbtcUnwrapEvents(result.Zenbtc2UnwrapEvents, shard.Events)
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

		result.RockWrapEvents = appendWrapEvents(result.RockWrapEvents, shard.Events)
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

		result.RockUnwrapEvents = appendRockUnwrapEvents(result.RockUnwrapEvents, shard.Events)
	}

	// Sort and dedup all event types
	sortWrapEventSlice(result.ZenzecWrapEvents)
	result.ZenzecWrapEvents = dedupWrapEventSlice(result.ZenzecWrapEvents)

	sortZenbtcUnwrapEventSlice(result.ZenzecUnwrapEvents)
	result.ZenzecUnwrapEvents = dedupZenbtcUnwrapEventSlice(result.ZenzecUnwrapEvents)

	sortWrapEventSlice(result.Zenbtc2WrapEvents)
	result.Zenbtc2WrapEvents = dedupWrapEventSlice(result.Zenbtc2WrapEvents)

	sortZenbtcUnwrapEventSlice(result.Zenbtc2UnwrapEvents)
	result.Zenbtc2UnwrapEvents = dedupZenbtcUnwrapEventSlice(result.Zenbtc2UnwrapEvents)

	sortWrapEventSlice(result.RockWrapEvents)
	result.RockWrapEvents = dedupWrapEventSlice(result.RockWrapEvents)

	sortRockUnwrapEventSlice(result.RockUnwrapEvents)
	result.RockUnwrapEvents = dedupRockUnwrapEventSlice(result.RockUnwrapEvents)

	return result, nil
}

// GetZenbtcWrapEvents fetches only zenbtc2 wrap events (for ZenBTC)
func (c *Client) GetZenbtcWrapEvents(ctx context.Context) ([]TokensMintedWithFee, error) {
	addresses, err := c.getAllShardAddresses(ZENBTC2_WRAP_SHARD_SEED, ZENBTC2_WRAP_SHARD_COUNT)
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

		var shard Zenbtc2WrapShard
		if err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:]); err == nil {
			events = appendWrapEvents(events, shard.Events)
		}
	}

	sortWrapEventSlice(events)
	events = dedupWrapEventSlice(events)
	return events, nil
}

// GetZenzecWrapEvents fetches only zenzec wrap events (for ZenZEC)
func (c *Client) GetZenzecWrapEvents(ctx context.Context) ([]TokensMintedWithFee, error) {
	addresses, err := c.getAllShardAddresses(ZENZEC_WRAP_SHARD_SEED, ZENZEC_WRAP_SHARD_COUNT)
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

		var shard ZenzecWrapShard
		if err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:]); err == nil {
			events = appendWrapEvents(events, shard.Events)
		}
	}

	sortWrapEventSlice(events)
	events = dedupWrapEventSlice(events)
	return events, nil
}

// GetZenbtcUnwrapEvents fetches only zenbtc2 unwrap events with decoded Bitcoin addresses (for ZenBTC)
func (c *Client) GetZenbtcUnwrapEvents(ctx context.Context) ([]ZenbtcTokenRedemption, error) {
	addresses, err := c.getAllShardAddresses(ZENBTC2_UNWRAP_SHARD_SEED, ZENBTC2_UNWRAP_SHARD_COUNT)
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

		var shard Zenbtc2UnwrapShard
		if err := borsh.Deserialize(&shard, account.Data.GetBinary()[8:]); err == nil {
			events = appendZenbtcUnwrapEvents(events, shard.Events)
		}
	}

	sortZenbtcUnwrapEventSlice(events)
	events = dedupZenbtcUnwrapEventSlice(events)
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
			events = appendWrapEvents(events, shard.Events)
		}
	}

	sortWrapEventSlice(events)
	events = dedupWrapEventSlice(events)
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
			events = appendRockUnwrapEvents(events, shard.Events)
		}
	}

	sortRockUnwrapEventSlice(events)
	events = dedupRockUnwrapEventSlice(events)
	return events, nil
}

// GetEventCounts returns the current counts of events stored
func (c *Client) GetEventCounts(ctx context.Context) (map[string]int, error) {
	allEvents, err := c.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	total := len(allEvents.ZenzecWrapEvents) + len(allEvents.ZenzecUnwrapEvents) +
		len(allEvents.Zenbtc2WrapEvents) + len(allEvents.Zenbtc2UnwrapEvents) +
		len(allEvents.RockWrapEvents) + len(allEvents.RockUnwrapEvents)

	return map[string]int{
		"zenzec_wrap":    len(allEvents.ZenzecWrapEvents),
		"zenzec_unwrap":  len(allEvents.ZenzecUnwrapEvents),
		"zenbtc2_wrap":   len(allEvents.Zenbtc2WrapEvents),
		"zenbtc2_unwrap": len(allEvents.Zenbtc2UnwrapEvents),
		"rock_wrap":      len(allEvents.RockWrapEvents),
		"rock_unwrap":    len(allEvents.RockUnwrapEvents),
		"total":          total,
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
