package neutrino

import (
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
}
