package solzenbtclegacy

import (
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

var wrapDiscriminator = [8]byte{0xb2, 0x28, 0x0a, 0xbd, 0xe4, 0x81, 0xba, 0x8c}

// WrapArgs captures the legacy wrap instruction arguments.
type WrapArgs struct {
	Value uint64
	Fee   uint64
}

// Wrap builds the legacy zenBTC wrap instruction that predates the event-store migration.
func Wrap(
	programID solana.PublicKey,
	args WrapArgs,
	signer solana.PublicKey,
	mint solana.PublicKey,
	multisigKey solana.PublicKey,
	feeWallet solana.PublicKey,
	feeWalletAta solana.PublicKey,
	receiver solana.PublicKey,
	receiverAta solana.PublicKey,
) (solana.Instruction, error) {
	globalConfigPDA, err := GetGlobalConfigPDA(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to derive global config PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		solana.Meta(signer).WRITE().SIGNER(),
		solana.Meta(globalConfigPDA).WRITE(),
		solana.Meta(multisigKey).WRITE(),
		solana.Meta(mint).WRITE(),
		solana.Meta(feeWallet).WRITE(),
		solana.Meta(feeWalletAta).WRITE(),
		solana.Meta(receiver),
		solana.Meta(receiverAta).WRITE(),
		solana.Meta(solana.SystemProgramID),
		solana.Meta(solana.TokenProgramID),
		solana.Meta(solana.SPLAssociatedTokenAccountProgramID),
	}

	data := make([]byte, 8+8+8)
	copy(data[:8], wrapDiscriminator[:])
	binary.LittleEndian.PutUint64(data[8:16], args.Value)
	binary.LittleEndian.PutUint64(data[16:], args.Fee)

	return solana.NewInstruction(
		programID,
		accounts,
		data,
	), nil
}
