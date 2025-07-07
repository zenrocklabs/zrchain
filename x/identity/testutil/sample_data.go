package testutil

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	policytypes "github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

var DefaultKr = types.Keyring{
	Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
	Creator:     "testCreator",
	Description: "testDescription",
	Admins:      []string{"testCreator"},
	KeyReqFee:   0,
	SigReqFee:   0,
	IsActive:    true,
}

var WantKr = types.Keyring{
	Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
	Creator:     "testCreator",
	Description: "testDescription",
	Admins:      []string{"testCreator", "testAdmin"},
	KeyReqFee:   0,
	SigReqFee:   0,
	IsActive:    true,
}

var DefaultWs = types.Workspace{
	Address: "workspace14a2hpadpsy9h4auve2z8lw",
	Creator: "testOwner",
	Owners:  []string{"testOwner"},
}

var DefaultWsWithOwners = types.Workspace{
	Address: "workspace14a2hpadpsy9h4auve2z8lw",
	Creator: "testOwner",
	Owners:  []string{"testOwner", "testOwner2"},
}

var Policy, _ = codectypes.NewAnyWithValue(&policytypes.BoolparserPolicy{
	Definition: "testOwner + testOwner2 > 1",
	Participants: []*policytypes.PolicyParticipant{
		{
			Address: "testOwner",
		},
		{
			Address: "testOwner2",
		},
	},
})

var Policy1 = policytypes.Policy{
	Id:     1,
	Name:   "Policy1",
	Policy: Policy,
}

var Policy2 = policytypes.Policy{
	Id:     2,
	Name:   "Policy2",
	Policy: Policy,
}

var ChildWs = types.Workspace{
	Address: "childWs",
	Creator: "testOwner",
	Owners:  []string{"testOwner"},
}

var InvalidChildWs = types.Workspace{
	Address: "invalidChildWs",
	Creator: "testOwner2",
	Owners:  []string{"testOwner2"},
}

var WsWithChild = types.Workspace{
	Address:         "workspace14a2hpadpsy9h4auve2z8lw",
	Creator:         "testOwner",
	Owners:          []string{"testOwner"},
	ChildWorkspaces: []string{"childWs"},
}

var Keyring = types.Keyring{
	Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
	Creator:     "testCreator",
	Description: "testDescription",
	Admins:      []string{"testCreator"},
	IsActive:    true,
}

var WantKeyring = types.Keyring{
	Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
	Creator:     "testCreator",
	Description: "testDescription",
	Admins:      []string{"testCreator"},
	IsActive:    false,
}

var DefaultKrWithAdmins = types.Keyring{
	Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
	Creator:     "testCreator",
	Description: "testDescription",
	Admins:      []string{"testCreator", "admin1", "admin2", "admin3"},
	KeyReqFee:   0,
	SigReqFee:   0,
	IsActive:    true,
}

var DefaultKrWithParties = types.Keyring{
	Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
	Creator:     "testCreator",
	Description: "testDescription",
	Admins:      []string{"testCreator"},
	Parties:     []string{"party1", "party2", "party3"},
	KeyReqFee:   0,
	SigReqFee:   0,
	IsActive:    true,
}
