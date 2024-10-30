package v3

import (
	"strconv"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
)

func UpdatePolicies(ctx sdk.Context, policyCol collections.Map[uint64, types.Policy], codec codec.BinaryCodec) error {
	policyStore, err := policyCol.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	policies, err := policyStore.Values()
	if err != nil {
		return err
	}
	for _, policy := range policies {
		rawPolicy, err := types.UnpackPolicy(codec, &policy)
		if err != nil {
			return err
		}
		bpPolicy := rawPolicy.(*types.BoolparserPolicy)

		approverNumber, err := bpPolicy.GetApproverNumber()
		if err != nil {
			return err
		}
		approvers := bpPolicy.GetParticipantAddresses()
		var newDefinition string
		for i, approver := range approvers {
			if i == len(approvers)-1 {
				newDefinition += approver + " > " + strconv.Itoa(approverNumber)
			} else {
				newDefinition += approver + " + "
			}
		}
		bpPolicy.Definition = newDefinition

		participants := make([]*types.PolicyParticipant, len(approvers))
		for i, addr := range approvers {
			participants[i] = &types.PolicyParticipant{
				Address: addr,
			}
		}

		policyData := types.BoolparserPolicy{
			Definition:   newDefinition,
			Participants: participants,
		}

		policyDataAny, err := codectypes.NewAnyWithValue(&policyData)
		if err != nil {
			return err
		}

		var newPolicy = types.Policy{
			Id:     policy.Id,
			Name:   policy.Name,
			Policy: policyDataAny,
			Btl:    policy.Btl,
		}

		err = policyCol.Set(ctx, policy.Id, newPolicy)
		if err != nil {
			return err
		}
	}

	return nil
}
