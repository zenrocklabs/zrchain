package policy

import (
	"github.com/cosmos/gogoproto/proto"
)

type ApproverSet map[string]bool

func BuildApproverSet(approvers []string) ApproverSet {
	approverSet := make(ApproverSet, len(approvers))
	for _, a := range approvers {
		approverSet[a] = true
	}
	return approverSet
}

type Policy interface {
	// Validate checks that the policy is valid (well formed).
	// The returned error is nil if the policy is valid.
	Validate() error

	// AddressToParticipant returns the participant shorthand for the given
	// address.
	AddressToParticipant(addr string) (string, error)

	// GetParticipants returns a list of addresses of the participants
	GetParticipantAddresses() []string

	// The Verify() method will receive the list of approvers as shorthands.
	// Verify tries to verify the current policy. The returned error is nil if
	// the policy is valid.
	Verify(approvers ApproverSet, policyData map[string][]byte) error
}

type PolicyMetadata interface {
	// Metadata returns the metadata associated with the policy. This is used
	// to return additional information about the policy in query responses.
	Metadata() (proto.Message, error)
}
