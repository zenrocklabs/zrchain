package main

const (
	noVoteMsg       = "<no vote extension>"
	missingFieldMsg = "<missing field>"
)

// ConsensusVote represents a validator's vote in the consensus process
type ConsensusVote struct {
	Validator struct {
		Address string `json:"address"`
		Power   int    `json:"power"`
	} `json:"validator"`
	VoteExtension      string `json:"vote_extension"`
	ExtensionSignature string `json:"extension_signature"`
	BlockIDFlag        int    `json:"block_id_flag"`
}

// BlockData represents the block data structure
type BlockData struct {
	ConsensusData struct {
		Votes []ConsensusVote `json:"votes"`
	} `json:"ConsensusData"`
}

// ValidatorInfo stores processed information about a validator
type ValidatorInfo struct {
	Address     string
	Power       int
	Index       int
	HasVote     bool
	DecodedAddr string
	Moniker     string
}

// ValidatorSetEntry represents a validator in the validator set
type ValidatorSetEntry struct {
	Address          string `yaml:"address"`
	ProposerPriority string `yaml:"proposer_priority"`
	PubKey           struct {
		Type  string `yaml:"type"`
		Value string `yaml:"value"`
	} `yaml:"pub_key"`
	VotingPower string `yaml:"voting_power"`
}

// ValidatorSetResponse is the response structure for the validator set query
type ValidatorSetResponse struct {
	BlockHeight string              `yaml:"block_height"`
	Validators  []ValidatorSetEntry `yaml:"validators"`
	Pagination  struct {
		Total string `yaml:"total"`
	} `yaml:"pagination"`
}

// ValidatorEntry represents a validator in the validators query
type ValidatorEntry struct {
	ConsensusPublicKey struct {
		Type  string `yaml:"type"`
		Value string `yaml:"value"`
	} `yaml:"consensus_pubkey"`
	Description struct {
		Moniker string `yaml:"moniker"`
		Details string `yaml:"details,omitempty"`
		Website string `yaml:"website,omitempty"`
	} `yaml:"description"`
}

// ValidatorsResponse is the response structure for the validators query
type ValidatorsResponse struct {
	Validators []ValidatorEntry `yaml:"validators"`
	Pagination struct {
		Total string `yaml:"total"`
	} `yaml:"pagination"`
}

// Config holds the application configuration
type Config struct {
	UseFile             string
	RPCNode             string
	Network             string
	BlockHeight         string
	MissingOnly         bool
	ConsensusReportMode bool
}

// Stats holds the statistics from vote processing
type Stats struct {
	ValidatorsWithExtensions int
	TotalValidators          int
	TotalVotedPower          int
	TotalVotingPower         int
}

// ValueCount represents a count of validators with a specific value
type ValueCount struct {
	Value      string
	Validators []ValidatorInfo
	TotalPower int
}

// --- Structs for Consensus Report Mode ---

// RPCBlockResponse wraps the result of a /block RPC call
type RPCBlockResponse struct {
	Result struct {
		Block struct {
			Header struct {
				AppHash         string `json:"app_hash"`
				ProposerAddress string `json:"proposer_address"`
				Height          string `json:"height"`
			} `json:"header"`
			LastCommit struct {
				Height     string         `json:"height"`
				Signatures []RPCSignature `json:"signatures"`
			} `json:"last_commit"`
		} `json:"block"`
	} `json:"result"`
}

// RPCSignature represents a signature from a block's last_commit
type RPCSignature struct {
	BlockIDFlag      int    `json:"block_id_flag"`
	ValidatorAddress string `json:"validator_address"`
	Signature        string `json:"signature"`
}

// RPCValidatorsResponse wraps the result of a /validators RPC call
type RPCValidatorsResponse struct {
	Result struct {
		Count      string         `json:"count"`
		Total      string         `json:"total"`
		Validators []RPCValidator `json:"validators"`
	} `json:"result"`
}

// RPCValidator represents a validator from the /validators RPC call
type RPCValidator struct {
	Address     string `json:"address"`
	VotingPower string `json:"voting_power"`
	PubKey      struct {
		Type  string `json:"type"`  // e.g. "tendermint/PubKeyEd25519"
		Value string `json:"value"` // Base64 encoded pubkey
	} `json:"pub_key"`
}

// ConsensusReportData holds processed information for the consensus report
type ConsensusReportData struct {
	Height                     int64
	AppHash                    string
	ProposerAddress            string
	ProposerMoniker            string
	TotalValidators            int
	AgreedValidators           []ValidatorVoteInfo
	VotedNilValidators         []ValidatorVoteInfo
	AbsentValidators           []ValidatorVoteInfo
	MissingSignatureValidators []ValidatorVoteInfo
	TotalVotingPower           int64
	AgreedVotingPower          int64
}

// ValidatorVoteInfo stores validator info for the consensus report
type ValidatorVoteInfo struct {
	Address     string
	Moniker     string
	VotingPower int64
}
