package solrock

import (
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go"
)

type Token struct {
	Name     string
	Symbol   string
	Decimals uint8
	Uri      string
}

type TokenRedemptionEvent struct {
	Signature string
	Slot      uint64
	Date      time.Time
	Redeemer  solana.PublicKey
	Value     uint64
	DestAddr  [25]uint8
	Fee       uint64
	Mint      solana.PublicKey
	Id        *big.Int
}
