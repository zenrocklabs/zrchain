package v2

import (
	"fmt"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
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

		fmt.Println("policy", policy)

		rawPolicy, err := types.UnpackPolicy(codec, &policy)
		if err != nil {
			return err
		}

		bpPolicy := rawPolicy.(*types.BoolparserPolicy)

		approverNumber := bpPolicy.Definition[len(bpPolicy.Definition)-1:]
		approvers := bpPolicy.GetParticipantAddresses()

		var newDefinition string
		for i, approver := range approvers {
			if i == len(approvers)-1 {
				newDefinition += approver + " > " + approverNumber
			} else {
				newDefinition += approver + " + "
			}
		}

		bpPolicy.Definition = newDefinition

		policyCol.Set(ctx, policy.Id, types.Policy{
			Creator: policy.Creator,
			Name:    policy.Name,
			Policy:  policy.Policy,
			Btl:     policy.Btl,
		})

		fmt.Println("Updated Policy", policy)
	}

	return nil
}
