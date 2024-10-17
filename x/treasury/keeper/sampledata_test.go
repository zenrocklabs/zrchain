package keeper_test

import (
	policytypes "github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

var boolPol1, _ = codectypes.NewAnyWithValue(&policytypes.BoolparserPolicy{
	Definition: "u1 > 0",
	Participants: []*policytypes.PolicyParticipant{
		{
			Abbreviation: "u1",
			Address:      "testOwner",
		},
	},
})

var boolPol2, _ = codectypes.NewAnyWithValue(&policytypes.BoolparserPolicy{
	Definition: "u1 > 0",
	Participants: []*policytypes.PolicyParticipant{
		{
			Abbreviation: "u1",
			Address:      "testOwner2",
		},
	},
})

var policy1 = policytypes.Policy{
	Id:     1,
	Name:   "Policy1",
	Policy: boolPol1,
}

var policy2 = policytypes.Policy{
	Id:     2,
	Name:   "Policy2",
	Policy: boolPol2,
}
