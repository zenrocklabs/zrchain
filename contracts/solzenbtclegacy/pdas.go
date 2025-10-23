package solzenbtclegacy

import (
	"github.com/gagliardetto/solana-go"
)

const (
	globalConfigSeed = "global_config"
)

// GetGlobalConfigPDA derives the global config PDA for the legacy zenBTC program.
func GetGlobalConfigPDA(programID solana.PublicKey) (solana.PublicKey, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(globalConfigSeed),
		},
		programID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return addr, nil
}
