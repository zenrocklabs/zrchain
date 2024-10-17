package policy

import (
	"fmt"
)

// AnyInGroupPolicy is a simple policy where any member of a group can verify it.
type AnyInGroupPolicy struct {
	group []string
}

var _ Policy = &AnyInGroupPolicy{}

func (*AnyInGroupPolicy) Validate() error { return nil }

func (p *AnyInGroupPolicy) AddressToParticipant(addr string) (string, error) {
	for _, s := range p.group {
		if s == addr {
			return addr, nil
		}
	}
	return "", fmt.Errorf("address not a participant of this policy")
}

func (p *AnyInGroupPolicy) GetParticipantAddresses() []string {
	return p.group
}

func (p *AnyInGroupPolicy) Verify(approvers ApproverSet, _ map[string][]byte) error {
	if len(approvers) == 0 {
		return fmt.Errorf("no approvers")
	}

	for _, s := range p.group {
		if _, found := approvers[s]; found {
			return nil
		}
	}

	return fmt.Errorf("approvers are not in the group")
}

func NewAnyInGroupPolicy(group []string) *AnyInGroupPolicy {
	return &AnyInGroupPolicy{
		group: group,
	}
}
