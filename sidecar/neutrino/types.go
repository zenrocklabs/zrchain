package neutrino

import (
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/neutrino/rpcservice"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/lightninglabs/neutrino"
)

type Config struct {
	GRPCPort int `yaml:"grpc_port"`
}

type LiteNode struct {
	Node *neutrino.ChainService
	DB   walletdb.DB
}

type NeutrinoServer struct {
	Nodes map[string]LiteNode
	Proxy *rpcservice.RpcCaller
}

type BlockCountRequest struct {
	ChainName string
	Quiet     bool
}

type BlockCountResponse struct {
	BlockHeight uint64
	Error       string
}

type BlockHashRequest struct {
	ChainName   string
	BlockHeight uint64
	Quiet       bool
}

type BlockHashResponse struct {
	BlockHash string
	Error     string
}

type BlockHeaderRequest struct {
	ChainName string
	BlockHash string // The block hash
	Verbose   bool   // Optional, default=true. If false, returns a string, otherwise returns a json object
	Quiet     bool
}

type BlockHeaderResponse struct {
	BlockHeader BitcoinBlock
	Error       string
}

type BitcoinBlock struct {
	Hash              string   `json:"hash"`
	Confirmations     int      `json:"confirmations"`
	Size              int      `json:"size"`
	Height            int      `json:"height"`
	Version           int      `json:"version"`
	MerkleRoot        string   `json:"merkleroot"`
	Tx                []string `json:"tx"`
	Time              int64    `json:"time"`
	Nonce             int64    `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	ChainWork         string   `json:"chainwork"`
	PreviousBlockHash string   `json:"previousblockhash"`
	NextBlockHash     string   `json:"nextblockhash,omitempty"`
}
