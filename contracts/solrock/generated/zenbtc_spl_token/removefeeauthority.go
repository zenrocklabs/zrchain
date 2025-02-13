// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package zenbtc_spl_token

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// RemoveFeeAuthority is the `remove_fee_authority` instruction.
type RemoveFeeAuthority struct {
	Args *RemoveFeeAuthorityArgs

	// [0] = [WRITE, SIGNER] signer
	//
	// [1] = [WRITE] global_config
	//
	// [2] = [WRITE] token_config
	//
	// [3] = [WRITE] mint
	//
	// [4] = [] system_program
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewRemoveFeeAuthorityInstructionBuilder creates a new `RemoveFeeAuthority` instruction builder.
func NewRemoveFeeAuthorityInstructionBuilder() *RemoveFeeAuthority {
	nd := &RemoveFeeAuthority{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 5),
	}
	nd.AccountMetaSlice[4] = ag_solanago.Meta(Addresses["11111111111111111111111111111111"])
	return nd
}

// SetArgs sets the "args" parameter.
func (inst *RemoveFeeAuthority) SetArgs(args RemoveFeeAuthorityArgs) *RemoveFeeAuthority {
	inst.Args = &args
	return inst
}

// SetSignerAccount sets the "signer" account.
func (inst *RemoveFeeAuthority) SetSignerAccount(signer ag_solanago.PublicKey) *RemoveFeeAuthority {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(signer).WRITE().SIGNER()
	return inst
}

// GetSignerAccount gets the "signer" account.
func (inst *RemoveFeeAuthority) GetSignerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetGlobalConfigAccount sets the "global_config" account.
func (inst *RemoveFeeAuthority) SetGlobalConfigAccount(globalConfig ag_solanago.PublicKey) *RemoveFeeAuthority {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(globalConfig).WRITE()
	return inst
}

func (inst *RemoveFeeAuthority) findFindGlobalConfigAddress(knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
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
func (inst *RemoveFeeAuthority) FindGlobalConfigAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindGlobalConfigAddress(bumpSeed)
	return
}

func (inst *RemoveFeeAuthority) MustFindGlobalConfigAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindGlobalConfigAddress(bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindGlobalConfigAddress finds GlobalConfig account address with given seeds.
func (inst *RemoveFeeAuthority) FindGlobalConfigAddress() (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindGlobalConfigAddress(0)
	return
}

func (inst *RemoveFeeAuthority) MustFindGlobalConfigAddress() (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindGlobalConfigAddress(0)
	if err != nil {
		panic(err)
	}
	return
}

// GetGlobalConfigAccount gets the "global_config" account.
func (inst *RemoveFeeAuthority) GetGlobalConfigAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetTokenConfigAccount sets the "token_config" account.
func (inst *RemoveFeeAuthority) SetTokenConfigAccount(tokenConfig ag_solanago.PublicKey) *RemoveFeeAuthority {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(tokenConfig).WRITE()
	return inst
}

func (inst *RemoveFeeAuthority) findFindTokenConfigAddress(mint ag_solanago.PublicKey, knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// const: token_config
	seeds = append(seeds, []byte{byte(0x74), byte(0x6f), byte(0x6b), byte(0x65), byte(0x6e), byte(0x5f), byte(0x63), byte(0x6f), byte(0x6e), byte(0x66), byte(0x69), byte(0x67)})
	// path: mint
	seeds = append(seeds, mint.Bytes())

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, ProgramID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, ProgramID)
	}
	return
}

// FindTokenConfigAddressWithBumpSeed calculates TokenConfig account address with given seeds and a known bump seed.
func (inst *RemoveFeeAuthority) FindTokenConfigAddressWithBumpSeed(mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindTokenConfigAddress(mint, bumpSeed)
	return
}

func (inst *RemoveFeeAuthority) MustFindTokenConfigAddressWithBumpSeed(mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindTokenConfigAddress(mint, bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindTokenConfigAddress finds TokenConfig account address with given seeds.
func (inst *RemoveFeeAuthority) FindTokenConfigAddress(mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindTokenConfigAddress(mint, 0)
	return
}

func (inst *RemoveFeeAuthority) MustFindTokenConfigAddress(mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindTokenConfigAddress(mint, 0)
	if err != nil {
		panic(err)
	}
	return
}

// GetTokenConfigAccount gets the "token_config" account.
func (inst *RemoveFeeAuthority) GetTokenConfigAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetMintAccount sets the "mint" account.
func (inst *RemoveFeeAuthority) SetMintAccount(mint ag_solanago.PublicKey) *RemoveFeeAuthority {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(mint).WRITE()
	return inst
}

// GetMintAccount gets the "mint" account.
func (inst *RemoveFeeAuthority) GetMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetSystemProgramAccount sets the "system_program" account.
func (inst *RemoveFeeAuthority) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *RemoveFeeAuthority {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "system_program" account.
func (inst *RemoveFeeAuthority) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

func (inst RemoveFeeAuthority) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_RemoveFeeAuthority,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst RemoveFeeAuthority) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RemoveFeeAuthority) Validate() error {
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
			return errors.New("accounts.TokenConfig is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *RemoveFeeAuthority) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("RemoveFeeAuthority")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Args", *inst.Args))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=5]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("        signer", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta(" global_config", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("  token_config", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("          mint", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("system_program", inst.AccountMetaSlice.Get(4)))
					})
				})
		})
}

func (obj RemoveFeeAuthority) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Args` param:
	err = encoder.Encode(obj.Args)
	if err != nil {
		return err
	}
	return nil
}
func (obj *RemoveFeeAuthority) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Args`:
	err = decoder.Decode(&obj.Args)
	if err != nil {
		return err
	}
	return nil
}

// NewRemoveFeeAuthorityInstruction declares a new RemoveFeeAuthority instruction with the provided parameters and accounts.
func NewRemoveFeeAuthorityInstruction(
	// Parameters:
	args RemoveFeeAuthorityArgs,
	// Accounts:
	signer ag_solanago.PublicKey,
	globalConfig ag_solanago.PublicKey,
	tokenConfig ag_solanago.PublicKey,
	mint ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *RemoveFeeAuthority {
	return NewRemoveFeeAuthorityInstructionBuilder().
		SetArgs(args).
		SetSignerAccount(signer).
		SetGlobalConfigAccount(globalConfig).
		SetTokenConfigAccount(tokenConfig).
		SetMintAccount(mint).
		SetSystemProgramAccount(systemProgram)
}
