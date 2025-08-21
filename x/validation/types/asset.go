package types

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/collections/codec"
)

// AssetKey represents a collections.KeyCodec for Asset enum
type AssetKey struct{}

var _ codec.KeyCodec[Asset] = AssetKey{}

func (k AssetKey) Encode(buffer []byte, value Asset) (int, error) {
	binary.BigEndian.PutUint32(buffer, uint32(value))
	return 4, nil
}

func (k AssetKey) Decode(buffer []byte) (int, Asset, error) {
	if len(buffer) < 4 {
		return 0, Asset_UNSPECIFIED, fmt.Errorf("invalid buffer length for AssetKey: expected at least 4, got %d", len(buffer))
	}

	return 4, Asset(binary.BigEndian.Uint32(buffer)), nil
}

func (k AssetKey) Size(_ Asset) int {
	return 4
}

func (k AssetKey) EncodeJSON(value Asset) ([]byte, error) {
	return []byte(fmt.Sprintf("%q", value.String())), nil
}

func (k AssetKey) DecodeJSON(b []byte) (Asset, error) {
	// Remove quotes if present
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if val, ok := Asset_value[s]; ok {
		return Asset(val), nil
	}
	return Asset_UNSPECIFIED, fmt.Errorf("invalid asset value: %s", s)
}

func (k AssetKey) Stringify(value Asset) string {
	return value.String()
}

func (k AssetKey) KeyType() string {
	return "asset_key"
}

// Implement non-terminal encoding methods required by KeyCodec
func (k AssetKey) EncodeNonTerminal(buffer []byte, value Asset) (int, error) {
	return k.Encode(buffer, value)
}

func (k AssetKey) DecodeNonTerminal(buffer []byte) (int, Asset, error) {
	return k.Decode(buffer)
}

func (k AssetKey) SizeNonTerminal(value Asset) int {
	return k.Size(value)
}
