package sample

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/rand/v2"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

func WorkspaceAddress() string {
	num := rand.Uint64N(math.MaxUint64)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, num)
	addrHash := sha256.Sum256(buf)
	return sdk.MustBech32ifyAddressBytes("workspace", sdk.AccAddress(addrHash[:10]))
}

func KeyringAddress() string {
	num := rand.Uint64N(math.MaxUint64)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, num)
	addrHash := sha256.Sum256(buf)
	return sdk.MustBech32ifyAddressBytes("keyring", sdk.AccAddress(addrHash[:11]))
}

func StringLen(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = 65
	}

	return string(b)
}

func GetAuthority() string {
	return authtypes.NewModuleAddress(govtypes.ModuleName).String()
}
