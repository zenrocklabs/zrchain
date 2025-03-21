package keeper

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/Zenrock-Foundation/zrchain/v6/policy"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var policyDataProviders = map[string]PolicyDataProvider{}

func init() {
	policyDataProviders["passkey{"] = &PasskeyPolicyDataProvider{}
}

// RegisterActionHandler registers a handler for a specific action type.
// The handler function is called when the action is executed.
func RegisterActionHandler[ResT any](k types.PolicyKeeper, actionType string, handlerFn func(sdk.Context, *types.Action) (ResT, error)) {
	if _, ok := k.ActionHandler(actionType); ok {
		// To be safe and prevent mistakes we shouldn't allow to register
		// multiple handlers for the same action type.
		// However, in the current implementation of Cosmos SDK, this is called
		// twice so we'll ignore the second call.
		return
	}
	k.RegisterActionHandler(actionType, func(ctx sdk.Context, a *types.Action) (any, error) {
		return handlerFn(ctx, a)
	})
}

func RegisterPolicyGeneratorHandler[ReqT any](k types.PolicyKeeper, reqType string, handlerFn func(sdk.Context, ReqT) (policy.Policy, error)) {
	if _, ok := k.GeneratorHandler(reqType); ok {
		// To be safe and prevent mistakes we shouldn't allow to register
		// multiple handlers for the same action type.
		// However, in the current implementation of Cosmos SDK, this is called
		// twice so we'll ignore the second call.
		return
	}

	k.RegisterPolicyGeneratorHandler(reqType, func(ctx sdk.Context, a *codectypes.Any) (policy.Policy, error) {
		var m sdk.Msg
		if err := k.Codec().UnpackAny(a, &m); err != nil {
			return nil, err
		}
		return handlerFn(ctx, m.(ReqT))
	})
}

// TryExecuteAction checks if the policy attached to the action is satisfied
// and executes it.
//
// If the policy is satisfied, the provided handler function is executed and
// its response returned. If the policy is still not satisfied, nil is returned.
//
// This function should be called:
// - after an action is created
// - every time there is a change in the approvers set
func TryExecuteAction[ReqT sdk.Msg, ResT any](
	k types.PolicyKeeper,
	cdc codec.BinaryCodec,
	ctx sdk.Context,
	act *types.Action,
	handlerFn func(sdk.Context, ReqT) (*ResT, error),
) (*ResT, error) {
	var m sdk.Msg
	err := k.Codec().UnpackAny(act.Msg, &m)
	if err != nil {
		return nil, err
	}

	msg, ok := m.(ReqT)
	if !ok {
		return nil, fmt.Errorf("invalid message type, expected %T", new(ReqT))
	}

	pol, err := k.PolicyForAction(ctx, act)
	if err != nil {
		return nil, err
	}

	signersSet := policy.BuildApproverSet(act.Approvers)

	if err := pol.Verify(signersSet, act.GetPolicyDataMap()); err == nil {
		act.Status = types.ActionStatus_ACTION_STATUS_COMPLETED

		if err := k.SetAction(ctx, act); err != nil {
			return nil, err
		}

		return handlerFn(ctx, msg)
	}

	var res ResT
	return &res, nil
}

func (k Keeper) PolicyForAction(ctx sdk.Context, act *types.Action) (policy.Policy, error) {
	var (
		pol policy.Policy
		err error
	)

	if act.PolicyId == 0 {
		// if no explicit policy ID specified, try to generate one
		polGen, found := k.GeneratorHandler(act.Msg.TypeUrl)
		if !found {
			return nil, fmt.Errorf("no policy ID specied for action and no policy generator registered for %s", act.Msg.TypeUrl)
		}

		pol, err = polGen(ctx, act.Msg)
		if err != nil {
			return nil, err
		}
	} else {
		p, err := k.PolicyStore.Get(ctx, act.PolicyId)
		if err != nil {
			return nil, fmt.Errorf("policy not found: %d", act.PolicyId)
		}

		pol, err = k.Unpack(&p)
		if err != nil {
			return nil, err
		}
	}

	return pol, nil
}

// AddAction creates a new action for the provided message with initial approvers.
// Who calls this function should also immediately check if the action can be
// executed with the provided initialApprovers, by calling TryExecuteAction.
func (k Keeper) AddAction(ctx sdk.Context, creator string, msg sdk.Msg, policyID, reqBtl uint64, policyData map[string][]byte, wsOwners []string) (*types.Action, error) {
	wrappedMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	policyBtl := uint64(0)
	if policyID != 0 {
		p, err := k.PolicyStore.Get(ctx, policyID)
		if err == nil {
			policyBtl = p.Btl
		}
	}
	btl := k.getCheckedBTL(ctx, reqBtl, policyBtl)

	policyDataKv := mapToDeterministicSlice(policyData)
	// create action object
	act := types.Action{
		Status:     types.ActionStatus_ACTION_STATUS_PENDING,
		Approvers:  []string{},
		PolicyId:   policyID,
		Msg:        wrappedMsg,
		Creator:    creator,
		Btl:        uint64(ctx.BlockHeight()) + btl,
		PolicyData: policyDataKv,
	}

	// add initial approver
	pol, err := k.PolicyForAction(ctx, &act)
	if err != nil {
		return nil, err
	}

	if !addressInWorkspace(wsOwners, creator) {
		return nil, fmt.Errorf("creator %s is not an owner of this workspace", creator)
	}

	if addressInPolicy(pol.GetParticipantAddresses(), creator) {
		if err := act.AddApprover(creator); err != nil {
			return nil, err
		}
	}

	// get policy participants
	participants := []string{}
	switch p := pol.(type) {
	case *types.BoolparserPolicy:
		for _, participant := range p.Participants {
			participants = append(participants, participant.Address)
		}
	}

	if policyData == nil {
		policyData = make(map[string][]byte)
	}

	for _, participant := range participants {
		for prefix, prov := range policyDataProviders {
			if strings.HasPrefix(participant, prefix) {
				dataK, dataV, err := prov.GetData(participant, &act, k.cdc)
				if err != nil {
					return nil, err
				}
				policyData[dataK] = dataV[:]
			}
		}
	}

	if len(policyData) > 0 {
		act.PolicyData = mapToDeterministicSlice(policyData)
	}

	// store and return generated action
	k.CreateAction(ctx, &act)

	// emit new-action events for ws rpc clients
	actionId := strconv.FormatUint(act.GetId(), 10)
	for _, participant := range participants {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventNewAction,
				sdk.NewAttribute(types.AttributeActionId, actionId),
				sdk.NewAttribute(types.AttributeParticipantAddr, participant),
			),
		})
	}

	return &act, nil
}

// getCheckedBTL verifies that the specified btl is above the minimum, if not the minimum is returned
// if the btl is specified as zero, the default will be returned
func (k Keeper) getCheckedBTL(goCtx sdk.Context, reqBtl uint64, polBtl uint64) (res uint64) {
	res = reqBtl
	params, err := k.ParamStore.Get(goCtx)
	if err != nil {
		return res
	}

	if reqBtl == 0 && polBtl > 0 {
		res = polBtl
	} else if reqBtl == 0 {
		res = params.DefaultBtl
	} else if reqBtl < params.MinimumBtl {
		res = params.MinimumBtl
	}

	return res
}

func mapToDeterministicSlice(policyData map[string][]byte) []*types.KeyValue {
	policyDataKv := make([]*types.KeyValue, 0, len(policyData))
	for k, v := range policyData { // Iterate over map (non-deterministic outcome)
		policyDataKv = append(policyDataKv, &types.KeyValue{Key: k, Value: v})
	}
	// Rank slice by Key.
	sort.Slice(policyDataKv, func(i, j int) bool {
		return strings.Compare(policyDataKv[i].Key, policyDataKv[j].Key) < 0
	})
	return policyDataKv
}

func addressInWorkspace(wsOwners []string, creator string) bool {
	return slices.Contains(wsOwners, creator)
}

func addressInPolicy(policyParticipants []string, creator string) bool {
	return slices.Contains(policyParticipants, creator)
}
