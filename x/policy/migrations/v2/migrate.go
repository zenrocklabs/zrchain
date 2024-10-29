package v2

import (
	"strconv"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
)

func UpdatePolicies(ctx sdk.Context, policyCol collections.Map[uint64, types.Policy], codec codec.BinaryCodec) error {
	ctx.Logger().Info("1")
	policyStore, err := policyCol.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	ctx.Logger().Info("2")
	policies, err := policyStore.Values()
	if err != nil {
		return err
	}
	ctx.Logger().Info("3")
	for _, policy := range policies {
		ctx.Logger().Info("4")
		rawPolicy, err := types.UnpackPolicy(codec, &policy)
		if err != nil {
			return err
		}
		ctx.Logger().Info("5")
		bpPolicy := rawPolicy.(*types.BoolparserPolicy)

		approverNumber, err := bpPolicy.GetApproverNumber()
		if err != nil {
			return err
		}
		ctx.Logger().Info("6")
		approvers := bpPolicy.GetParticipantAddresses()
		ctx.Logger().Info("7")
		var newDefinition string
		for i, approver := range approvers {
			ctx.Logger().Info("8")
			if i == len(approvers)-1 {
				newDefinition += approver + " > " + strconv.Itoa(approverNumber)
			} else {
				newDefinition += approver + " + "
			}
		}
		ctx.Logger().Info("9")
		bpPolicy.Definition = newDefinition

		policyCol.Set(ctx, policy.Id, types.Policy{
			Creator: policy.Creator,
			Name:    policy.Name,
			Policy:  policy.Policy,
			Btl:     policy.Btl,
		})

		ctx.Logger().Info("10")

	}

	return nil
}
