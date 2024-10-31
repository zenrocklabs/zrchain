package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Zenrock-Foundation/zrchain/v5/boolparser"
	"github.com/Zenrock-Foundation/zrchain/v5/policy"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint:stylecheck,st1003
// revive:disable-next-line var-naming
func (a *Policy) SetId(id uint64) {
	a.Id = id
}

func UnpackPolicy(cdc codec.BinaryCodec, policyPb *Policy) (policy.Policy, error) {
	var p policy.Policy
	err := cdc.UnpackAny(policyPb.Policy, &p)
	if err != nil {
		return nil, fmt.Errorf("unpacking Any: %w", err)
	}

	return p, nil
}

var _ (policy.Policy) = (*BoolparserPolicy)(nil)

func (p *BoolparserPolicy) Validate() error {
	if len(p.Participants) == 0 {
		return fmt.Errorf("no participants")
	}
	if len(p.Definition) == 0 {
		return fmt.Errorf("no definition")
	}

	existingAddresses := map[string]struct{}{}

	for _, participant := range p.Participants {
		if len(participant.Address) == 0 {
			return fmt.Errorf("no address for %s", participant.Address)
		}

		// TODO: address verification

		if !strings.Contains(p.Definition, participant.Address) {
			return fmt.Errorf("participant %s not found in expression", participant.Address)
		}

		if _, ok := existingAddresses[participant.Address]; ok {
			return fmt.Errorf("duplicate address for %s", participant.Address)
		}
		existingAddresses[participant.Address] = struct{}{}

		if !strings.HasPrefix(participant.Address, "passkey{") {
			_, err := sdk.AccAddressFromBech32(participant.Address)
			if err != nil {
				return fmt.Errorf("invalid address %s", err)
			}
		}

		approverNumber, err := p.GetApproverNumber()

		if err != nil {
			return fmt.Errorf("error getting approver number: %w", err)
		}

		if len(p.Participants) <= approverNumber {
			return fmt.Errorf("number of participants is less than the approver number")
		}
	}

	parser := boolparser.NewParser(strings.NewReader(p.Definition))
	stack, _ := parser.Parse()
outer:
	for _, part := range stack.Values {
		if part.Type != boolparser.CONSTANT {
			continue
		}
		for _, part2 := range p.Participants {
			if strings.EqualFold(part.Value, part2.Address) {
				continue outer
			}
		}
		return fmt.Errorf("participant %s not provided", part.Value)
	}
	return nil
}

func (p *BoolparserPolicy) AddressToParticipant(addr string) (string, error) {
	for _, participant := range p.Participants {
		if participant.Address == addr {
			return participant.Address, nil
		}
	}
	return "", fmt.Errorf("address not a participant of this policy")
}

func (p *BoolparserPolicy) GetParticipantAddresses() []string {
	addresses := []string{}
	for _, part := range p.Participants {
		addresses = append(addresses, part.Address)
	}
	return addresses
}

func (p *BoolparserPolicy) GetApproverNumber() (int, error) {

	// Split the string into parts
	parts := strings.Fields(p.Definition)

	// Check if there are any parts
	if len(parts) == 0 {
		return 0, fmt.Errorf("no values found in the input string")
	}

	// Extract the last part with the approver number
	lastValue := parts[len(parts)-1]

	// Convert the last part to an integer
	approverNumber, err := strconv.Atoi(lastValue)
	if err != nil {
		return 0, fmt.Errorf("error converting '%s' to int: %v", lastValue, err)
	}

	return approverNumber, nil
}

func (p *BoolparserPolicy) Verify(approvers policy.ApproverSet, policyData map[string][]byte) error {
	expression := p.Definition
	for addr := range approvers {
		expression = strings.ReplaceAll(expression, addr, "1")
	}

	for valueName, value := range policyData {
		expression = strings.ReplaceAll(expression, valueName, string(value))
	}

	if boolparser.BoolSolve(expression) {
		return nil
	}
	return fmt.Errorf("expression not satisfied")
}
