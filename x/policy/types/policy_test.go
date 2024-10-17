package types

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v4/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestBoolparserPolicy_Validate(t *testing.T) {
	address1 := sample.AccAddress()
	address2 := sample.AccAddress()

	tests := []struct {
		name    string
		bp      BoolparserPolicy
		wantErr bool
	}{
		{
			name: "pass: valid policy",
			bp: BoolparserPolicy{
				Definition: "u1 + u2 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "u1",
						Address:      address1,
					},
					{
						Abbreviation: "u2",
						Address:      address2,
					},
				},
			},
		},
		{
			name: "fail: missing participant",
			bp: BoolparserPolicy{
				Definition: "u1 + u2 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "u1",
						Address:      address1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: unused participant",
			bp: BoolparserPolicy{
				Definition: "u1 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "u1",
						Address:      address1,
					},
					{
						Abbreviation: "u2",
						Address:      address2,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: duplicate abbrev",
			bp: BoolparserPolicy{
				Definition: "u1 + u1 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "u1",
						Address:      address1,
					},
					{
						Abbreviation: "u1",
						Address:      address2,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: duplicate address",
			bp: BoolparserPolicy{
				Definition: "u1 + u2 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "u1",
						Address:      address1,
					},
					{
						Abbreviation: "u2",
						Address:      address1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: invalid address",
			bp: BoolparserPolicy{
				Definition: "u1 + u1 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "u1",
						Address:      "some-invalid-address",
					},
					{
						Abbreviation: "u1",
						Address:      address1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "pass: passkey address",
			bp: BoolparserPolicy{
				Definition: "u1 + u2 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "u1",
						Address:      "passkey{passkey_id}",
					},
					{
						Abbreviation: "u2",
						Address:      address2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail: invalid abbreviation",
			bp: BoolparserPolicy{
				Definition: "1 + 2 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "1",
						Address:      address1,
					},
					{
						Abbreviation: "2",
						Address:      address2,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: invalid abbreviation",
			bp: BoolparserPolicy{
				Definition: "1 + 2 > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "",
						Address:      address1,
					},
					{
						Abbreviation: "2",
						Address:      address2,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "pass: valid abbreviation",
			bp: BoolparserPolicy{
				Definition: "a_b + a_c > 1",
				Participants: []*PolicyParticipant{
					{
						Abbreviation: "a_c",
						Address:      address1,
					},
					{
						Abbreviation: "a_b",
						Address:      address2,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bp.Validate()
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
