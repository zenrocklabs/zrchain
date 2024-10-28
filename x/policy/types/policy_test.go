package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoolparserPolicy_Validate(t *testing.T) {

	tests := []struct {
		name    string
		bp      BoolparserPolicy
		wantErr bool
	}{
		{
			name: "pass: valid policy",
			bp: BoolparserPolicy{
				Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1",
				Participants: []*PolicyParticipant{
					{
						Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					},
					{
						Address: "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
					},
				},
			},
		},
		{
			name: "fail: missing participant",
			bp: BoolparserPolicy{
				Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + u2 > 1",
				Participants: []*PolicyParticipant{
					{
						Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: unused participant",
			bp: BoolparserPolicy{
				Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty > 1",
				Participants: []*PolicyParticipant{
					{
						Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					},
					{
						Address: "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: approver number can't be fulfilled",
			bp: BoolparserPolicy{
				Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty > 1",
				Participants: []*PolicyParticipant{
					{
						Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: approver number can't be fulfilled",
			bp: BoolparserPolicy{
				Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty > 2",
				Participants: []*PolicyParticipant{
					{
						Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: duplicate address",
			bp: BoolparserPolicy{
				Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty > 1",
				Participants: []*PolicyParticipant{
					{
						Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					},
					{
						Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail: invalid address",
			bp: BoolparserPolicy{
				Definition: "some-invalid-address + zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty > 1",
				Participants: []*PolicyParticipant{
					{
						Address: "some-invalid-address",
					},
					{
						Address: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "pass: passkey address",
			bp: BoolparserPolicy{
				Definition: "passkey{passkey_id} + zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1",
				Participants: []*PolicyParticipant{
					{
						Address: "passkey{passkey_id}",
					},
					{
						Address: "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
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
