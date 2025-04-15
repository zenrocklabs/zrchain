// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package rock_spl_token

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Add a new authority that can manage fees
// Only callable by global authority
type AddFeeAuthority struct {
	Args *AddFeeAuthorityArgs

	// [0] = [WRITE, SIGNER] signer
	// ··········· Signer must be the global authority
	//
	// [1] = [WRITE] global_config
	// ··········· Global configuration account storing authorities
	//
	// [2] = [WRITE] mint
	// ··········· Token mint controlled by this program
	//
	// [3] = [] system_program
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewAddFeeAuthorityInstructionBuilder creates a new `AddFeeAuthority` instruction builder.
func NewAddFeeAuthorityInstructionBuilder() *AddFeeAuthority {
	nd := &AddFeeAuthority{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	nd.AccountMetaSlice[3] = ag_solanago.Meta(Addresses["11111111111111111111111111111111"])
	return nd
}

// SetArgs sets the "args" parameter.
func (inst *AddFeeAuthority) SetArgs(args AddFeeAuthorityArgs) *AddFeeAuthority {
	inst.Args = &args
	return inst
}

// SetSignerAccount sets the "signer" account.
// Signer must be the global authority
func (inst *AddFeeAuthority) SetSignerAccount(signer ag_solanago.PublicKey) *AddFeeAuthority {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(signer).WRITE().SIGNER()
	return inst
}

// GetSignerAccount gets the "signer" account.
// Signer must be the global authority
func (inst *AddFeeAuthority) GetSignerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetGlobalConfigAccount sets the "global_config" account.
// Global configuration account storing authorities
func (inst *AddFeeAuthority) SetGlobalConfigAccount(globalConfig ag_solanago.PublicKey) *AddFeeAuthority {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(globalConfig).WRITE()
	return inst
}

func (inst *AddFeeAuthority) findFindGlobalConfigAddress(knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
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
func (inst *AddFeeAuthority) FindGlobalConfigAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindGlobalConfigAddress(bumpSeed)
	return
}

func (inst *AddFeeAuthority) MustFindGlobalConfigAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindGlobalConfigAddress(bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindGlobalConfigAddress finds GlobalConfig account address with given seeds.
func (inst *AddFeeAuthority) FindGlobalConfigAddress() (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindGlobalConfigAddress(0)
	return
}

func (inst *AddFeeAuthority) MustFindGlobalConfigAddress() (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindGlobalConfigAddress(0)
	if err != nil {
		panic(err)
	}
	return
}

// GetGlobalConfigAccount gets the "global_config" account.
// Global configuration account storing authorities
func (inst *AddFeeAuthority) GetGlobalConfigAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetMintAccount sets the "mint" account.
// Token mint controlled by this program
func (inst *AddFeeAuthority) SetMintAccount(mint ag_solanago.PublicKey) *AddFeeAuthority {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(mint).WRITE()
	return inst
}

// GetMintAccount gets the "mint" account.
// Token mint controlled by this program
func (inst *AddFeeAuthority) GetMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSystemProgramAccount sets the "system_program" account.
func (inst *AddFeeAuthority) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *AddFeeAuthority {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "system_program" account.
func (inst *AddFeeAuthority) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst AddFeeAuthority) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_AddFeeAuthority,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst AddFeeAuthority) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AddFeeAuthority) Validate() error {
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
			return errors.New("accounts.Mint is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *AddFeeAuthority) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("AddFeeAuthority")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Args", *inst.Args))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("        signer", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta(" global_config", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("          mint", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("system_program", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj AddFeeAuthority) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Args` param:
	err = encoder.Encode(obj.Args)
	if err != nil {
		return err
	}
	return nil
}
func (obj *AddFeeAuthority) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Args`:
	err = decoder.Decode(&obj.Args)
	if err != nil {
		return err
	}
	return nil
}

// NewAddFeeAuthorityInstruction declares a new AddFeeAuthority instruction with the provided parameters and accounts.
func NewAddFeeAuthorityInstruction(
	// Parameters:
	args AddFeeAuthorityArgs,
	// Accounts:
	signer ag_solanago.PublicKey,
	globalConfig ag_solanago.PublicKey,
	mint ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *AddFeeAuthority {
	return NewAddFeeAuthorityInstructionBuilder().
		SetArgs(args).
		SetSignerAccount(signer).
		SetGlobalConfigAccount(globalConfig).
		SetMintAccount(mint).
		SetSystemProgramAccount(systemProgram)
}
