// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package zenbtc_spl_token

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Unwrap is the `unwrap` instruction.
type Unwrap struct {
	Args *UnwrapArgs

	// [0] = [WRITE, SIGNER] signer
	//
	// [1] = [WRITE] global_config
	//
	// [2] = [WRITE] spl_multisig
	//
	// [3] = [WRITE] mint
	//
	// [4] = [WRITE] signer_ata
	//
	// [5] = [WRITE] fee_wallet
	//
	// [6] = [WRITE] fee_wallet_ata
	//
	// [7] = [] system_program
	//
	// [8] = [] token_program
	//
	// [9] = [] associated_token_program
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewUnwrapInstructionBuilder creates a new `Unwrap` instruction builder.
func NewUnwrapInstructionBuilder() *Unwrap {
	nd := &Unwrap{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 10),
	}
	nd.AccountMetaSlice[7] = ag_solanago.Meta(Addresses["11111111111111111111111111111111"])
	nd.AccountMetaSlice[8] = ag_solanago.Meta(Addresses["TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"])
	nd.AccountMetaSlice[9] = ag_solanago.Meta(Addresses["ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"])
	return nd
}

// SetArgs sets the "args" parameter.
func (inst *Unwrap) SetArgs(args UnwrapArgs) *Unwrap {
	inst.Args = &args
	return inst
}

// SetSignerAccount sets the "signer" account.
func (inst *Unwrap) SetSignerAccount(signer ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(signer).WRITE().SIGNER()
	return inst
}

// GetSignerAccount gets the "signer" account.
func (inst *Unwrap) GetSignerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetGlobalConfigAccount sets the "global_config" account.
func (inst *Unwrap) SetGlobalConfigAccount(globalConfig ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(globalConfig).WRITE()
	return inst
}

func (inst *Unwrap) findFindGlobalConfigAddress(knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// const: global_config
	seeds = append(seeds, []byte{byte(0x67), byte(0x6c), byte(0x6f), byte(0x62), byte(0x61), byte(0x6c), byte(0x5f), byte(0x63), byte(0x6f), byte(0x6e), byte(0x66), byte(0x69), byte(0x67)})

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, ProgramID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, ProgramID)
	}
	return
}

// FindGlobalConfigAddressWithBumpSeed calculates GlobalConfig account address with given seeds and a known bump seed.
func (inst *Unwrap) FindGlobalConfigAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindGlobalConfigAddress(bumpSeed)
	return
}

func (inst *Unwrap) MustFindGlobalConfigAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindGlobalConfigAddress(bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindGlobalConfigAddress finds GlobalConfig account address with given seeds.
func (inst *Unwrap) FindGlobalConfigAddress() (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindGlobalConfigAddress(0)
	return
}

func (inst *Unwrap) MustFindGlobalConfigAddress() (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindGlobalConfigAddress(0)
	if err != nil {
		panic(err)
	}
	return
}

// GetGlobalConfigAccount gets the "global_config" account.
func (inst *Unwrap) GetGlobalConfigAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSplMultisigAccount sets the "spl_multisig" account.
func (inst *Unwrap) SetSplMultisigAccount(splMultisig ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(splMultisig).WRITE()
	return inst
}

// GetSplMultisigAccount gets the "spl_multisig" account.
func (inst *Unwrap) GetSplMultisigAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetMintAccount sets the "mint" account.
func (inst *Unwrap) SetMintAccount(mint ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(mint).WRITE()
	return inst
}

// GetMintAccount gets the "mint" account.
func (inst *Unwrap) GetMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetSignerAtaAccount sets the "signer_ata" account.
func (inst *Unwrap) SetSignerAtaAccount(signerAta ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(signerAta).WRITE()
	return inst
}

func (inst *Unwrap) findFindSignerAtaAddress(signer ag_solanago.PublicKey, mint ag_solanago.PublicKey, knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// path: signer
	seeds = append(seeds, signer.Bytes())
	// const (raw): [6 221 246 225 215 101 161 147 217 203 225 70 206 235 121 172 28 180 133 237 95 91 55 145 58 140 245 133 126 255 0 169]
	seeds = append(seeds, []byte{byte(0x6), byte(0xdd), byte(0xf6), byte(0xe1), byte(0xd7), byte(0x65), byte(0xa1), byte(0x93), byte(0xd9), byte(0xcb), byte(0xe1), byte(0x46), byte(0xce), byte(0xeb), byte(0x79), byte(0xac), byte(0x1c), byte(0xb4), byte(0x85), byte(0xed), byte(0x5f), byte(0x5b), byte(0x37), byte(0x91), byte(0x3a), byte(0x8c), byte(0xf5), byte(0x85), byte(0x7e), byte(0xff), byte(0x0), byte(0xa9)})
	// path: mint
	seeds = append(seeds, mint.Bytes())

	programID := Addresses["ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"]

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, programID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, programID)
	}
	return
}

// FindSignerAtaAddressWithBumpSeed calculates SignerAta account address with given seeds and a known bump seed.
func (inst *Unwrap) FindSignerAtaAddressWithBumpSeed(signer ag_solanago.PublicKey, mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindSignerAtaAddress(signer, mint, bumpSeed)
	return
}

func (inst *Unwrap) MustFindSignerAtaAddressWithBumpSeed(signer ag_solanago.PublicKey, mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindSignerAtaAddress(signer, mint, bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindSignerAtaAddress finds SignerAta account address with given seeds.
func (inst *Unwrap) FindSignerAtaAddress(signer ag_solanago.PublicKey, mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindSignerAtaAddress(signer, mint, 0)
	return
}

func (inst *Unwrap) MustFindSignerAtaAddress(signer ag_solanago.PublicKey, mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindSignerAtaAddress(signer, mint, 0)
	if err != nil {
		panic(err)
	}
	return
}

// GetSignerAtaAccount gets the "signer_ata" account.
func (inst *Unwrap) GetSignerAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetFeeWalletAccount sets the "fee_wallet" account.
func (inst *Unwrap) SetFeeWalletAccount(feeWallet ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(feeWallet).WRITE()
	return inst
}

// GetFeeWalletAccount gets the "fee_wallet" account.
func (inst *Unwrap) GetFeeWalletAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetFeeWalletAtaAccount sets the "fee_wallet_ata" account.
func (inst *Unwrap) SetFeeWalletAtaAccount(feeWalletAta ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(feeWalletAta).WRITE()
	return inst
}

func (inst *Unwrap) findFindFeeWalletAtaAddress(feeWallet ag_solanago.PublicKey, mint ag_solanago.PublicKey, knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// path: feeWallet
	seeds = append(seeds, feeWallet.Bytes())
	// const (raw): [6 221 246 225 215 101 161 147 217 203 225 70 206 235 121 172 28 180 133 237 95 91 55 145 58 140 245 133 126 255 0 169]
	seeds = append(seeds, []byte{byte(0x6), byte(0xdd), byte(0xf6), byte(0xe1), byte(0xd7), byte(0x65), byte(0xa1), byte(0x93), byte(0xd9), byte(0xcb), byte(0xe1), byte(0x46), byte(0xce), byte(0xeb), byte(0x79), byte(0xac), byte(0x1c), byte(0xb4), byte(0x85), byte(0xed), byte(0x5f), byte(0x5b), byte(0x37), byte(0x91), byte(0x3a), byte(0x8c), byte(0xf5), byte(0x85), byte(0x7e), byte(0xff), byte(0x0), byte(0xa9)})
	// path: mint
	seeds = append(seeds, mint.Bytes())

	programID := Addresses["ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"]

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, programID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, programID)
	}
	return
}

// FindFeeWalletAtaAddressWithBumpSeed calculates FeeWalletAta account address with given seeds and a known bump seed.
func (inst *Unwrap) FindFeeWalletAtaAddressWithBumpSeed(feeWallet ag_solanago.PublicKey, mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindFeeWalletAtaAddress(feeWallet, mint, bumpSeed)
	return
}

func (inst *Unwrap) MustFindFeeWalletAtaAddressWithBumpSeed(feeWallet ag_solanago.PublicKey, mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindFeeWalletAtaAddress(feeWallet, mint, bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindFeeWalletAtaAddress finds FeeWalletAta account address with given seeds.
func (inst *Unwrap) FindFeeWalletAtaAddress(feeWallet ag_solanago.PublicKey, mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindFeeWalletAtaAddress(feeWallet, mint, 0)
	return
}

func (inst *Unwrap) MustFindFeeWalletAtaAddress(feeWallet ag_solanago.PublicKey, mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindFeeWalletAtaAddress(feeWallet, mint, 0)
	if err != nil {
		panic(err)
	}
	return
}

// GetFeeWalletAtaAccount gets the "fee_wallet_ata" account.
func (inst *Unwrap) GetFeeWalletAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetSystemProgramAccount sets the "system_program" account.
func (inst *Unwrap) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "system_program" account.
func (inst *Unwrap) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetTokenProgramAccount sets the "token_program" account.
func (inst *Unwrap) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "token_program" account.
func (inst *Unwrap) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetAssociatedTokenProgramAccount sets the "associated_token_program" account.
func (inst *Unwrap) SetAssociatedTokenProgramAccount(associatedTokenProgram ag_solanago.PublicKey) *Unwrap {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(associatedTokenProgram)
	return inst
}

// GetAssociatedTokenProgramAccount gets the "associated_token_program" account.
func (inst *Unwrap) GetAssociatedTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

func (inst Unwrap) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_Unwrap,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Unwrap) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Unwrap) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Args == nil {
			return errors.New("Args parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Signer is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.GlobalConfig is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.SplMultisig is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.SignerAta is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.FeeWallet is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.FeeWalletAta is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.AssociatedTokenProgram is not set")
		}
	}
	return nil
}

func (inst *Unwrap) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Unwrap")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Args", *inst.Args))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=10]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("                  signer", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("           global_config", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("            spl_multisig", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("                    mint", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("              signer_ata", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("              fee_wallet", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("          fee_wallet_ata", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("          system_program", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("           token_program", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("associated_token_program", inst.AccountMetaSlice.Get(9)))
					})
				})
		})
}

func (obj Unwrap) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Args` param:
	err = encoder.Encode(obj.Args)
	if err != nil {
		return err
	}
	return nil
}
func (obj *Unwrap) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Args`:
	err = decoder.Decode(&obj.Args)
	if err != nil {
		return err
	}
	return nil
}

// NewUnwrapInstruction declares a new Unwrap instruction with the provided parameters and accounts.
func NewUnwrapInstruction(
	// Parameters:
	args UnwrapArgs,
	// Accounts:
	signer ag_solanago.PublicKey,
	globalConfig ag_solanago.PublicKey,
	splMultisig ag_solanago.PublicKey,
	mint ag_solanago.PublicKey,
	signerAta ag_solanago.PublicKey,
	feeWallet ag_solanago.PublicKey,
	feeWalletAta ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	associatedTokenProgram ag_solanago.PublicKey) *Unwrap {
	return NewUnwrapInstructionBuilder().
		SetArgs(args).
		SetSignerAccount(signer).
		SetGlobalConfigAccount(globalConfig).
		SetSplMultisigAccount(splMultisig).
		SetMintAccount(mint).
		SetSignerAtaAccount(signerAta).
		SetFeeWalletAccount(feeWallet).
		SetFeeWalletAtaAccount(feeWalletAta).
		SetSystemProgramAccount(systemProgram).
		SetTokenProgramAccount(tokenProgram).
		SetAssociatedTokenProgramAccount(associatedTokenProgram)
}
