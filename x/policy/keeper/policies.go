package keeper

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreatePolicy(ctx sdk.Context, policy *types.Policy) (uint64, error) {
	count, err := k.PolicyCount.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		policy.Id = 1
	} else {
		policy.Id = count + 1
	}

	if err := k.PolicyStore.Set(ctx, policy.Id, *policy); err != nil {
		return 0, err
	}

	if err := k.PolicyCount.Set(ctx, policy.Id); err != nil {
		return 0, err
	}

	return policy.Id, nil
}

func (k Keeper) GetPolicyParticipants(ctx context.Context, policyId uint64) (map[string]struct{}, error) {
	pol, err := k.PolicyStore.Get(ctx, policyId)
	if err != nil {
		return nil, err
	}

	p, err := types.UnpackPolicy(k.cdc, &pol)
	if err != nil {
		return nil, err
	}

	participants := map[string]struct{}{}
	for _, p := range p.GetParticipantAddresses() {
		participants[p] = struct{}{}
	}

	return participants, nil
}

func (k Keeper) PolicyMembersAreOwners(ctx context.Context, policyId uint64, wsOwners []string) error {
	if policyId > 0 {
		participants, err := k.GetPolicyParticipants(ctx, policyId)
		if err != nil {
			return err
		}

		additionalIds := map[string]struct{}{}
		for _, owner := range wsOwners {
			signMethods, err := k.SignMethodsByAddress(ctx, &types.QuerySignMethodsByAddressRequest{
				Address: owner,
			})
			if err != nil {
				return err
			}

			for _, s := range signMethods.Config {
				var signMethod types.SignMethod
				if err := k.cdc.UnpackAny(s, &signMethod); err != nil {
					return err
				}
				additionalIds[signMethod.GetParticipantId()] = struct{}{}
			}
		}

		for p := range participants {
			_, isAdditonalSignMethod := additionalIds[p]
			if !slices.Contains(wsOwners, p) && !isAdditonalSignMethod {
				return fmt.Errorf("policy participant %s is not an owner of the workspace", p)
			}
		}
	}

	return nil
}
