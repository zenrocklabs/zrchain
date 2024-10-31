package neutrino

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
	"github.com/lightninglabs/neutrino"

	"github.com/btcsuite/btclog"
)

func buildNeutrinoNode(chainParams chaincfg.Params, logLevel btclog.Level, nodes *map[string]LiteNode) (map[string]LiteNode, error) {

	chainName := chainParams.Name
	dataDir := "./neutrino/neutrino_" + chainName
	dbPath := dataDir + "/neutrino_" + chainName + ".db"

	backendLogger := btclog.NewBackend(os.Stdout)

	// Enable logging for Neutrino and its subsystems
	logger := backendLogger.Logger("Neutrino_" + chainName)
	logger.SetLevel(logLevel) // Set to debug level
	//logger.SetLevel(btclog.LevelError) // Set to debug level

	// Set Neutrino to use the logger
	neutrino.UseLogger(logger)

	// Open the BerkeleyDB database

	db, err := walletdb.Create("bdb", dbPath, true, time.Second*10)
	if err != nil {
		log.Fatalf("Failed to open DB: %s %v", chainName, err)
	}

	// Configure the Neutrino node with Testnet parameters
	spvConfig := neutrino.Config{
		DataDir:     dataDir,     // Directory to store chain data
		Database:    db,          // BerkeleyDB database
		ChainParams: chainParams, // Bitcoin parameters
	}

	// Create a new Neutrino ChainService (node)
	node, err := neutrino.NewChainService(spvConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Neutrino %s node: %v", chainName, err)
	}

	litenode := LiteNode{
		Node: node,
		DB:   db,
	}
	//shutDown(cancel, node)

	nodeMap := *nodes
	nodeMap[chainName] = litenode

	err = node.Start()
	if err != nil {
		log.Fatalf("Failed to start Neutrino node: %v", err)
	}

	log.Printf("Started Neutrino node '%s'", chainName)

	return nodeMap, err
}

func (ns *NeutrinoServer) Initialize() {
	var err error
	nodes := make(map[string]LiteNode)

	//Testnet
	nodes, err = buildNeutrinoNode(chaincfg.TestNet3Params, btclog.LevelError, &nodes)
	if err != nil {
		log.Printf("Failed to Start Node Testnet3")
	}
	nodes["testnet"] = nodes["testnet3"]
	ns.Nodes = nodes

	//Mainnet
	nodes, err = buildNeutrinoNode(chaincfg.MainNetParams, btclog.LevelError, &nodes)
	if err != nil {
		log.Printf("Failed to Start Node Mainnet")
	}
	ns.Nodes = nodes

	// Register RPC server
	rpc.Register(ns)
	rpc.HandleHTTP()
	// Listen for requests on port 1234
	l, e := net.Listen("tcp", ":12345")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (ns *NeutrinoServer) Stop() {
	for name, liteNode := range ns.Nodes {
		liteNode.Node.Stop()
		liteNode.DB.Close()
		log.Printf("Shutdown Neutrino Node %s \n" + name)
	}
}

func (ns *NeutrinoServer) GetBlockHeaderByHeight(chainName string, height int64) (*wire.BlockHeader, *chainhash.Hash, int32, error) {
	if liteNode, exists := ns.Nodes[chainName]; exists {
		node := liteNode.Node
		blockHash, err := node.GetBlockHash(height)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("Failed to get blockhash for height %d: %v", height, node.ChainParams().Name)
		}
		blockHeader, err := node.GetBlockHeader(blockHash)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("Failed to get blockheader for height %d: %v", height, node.ChainParams().Name)
		}
		blockStamp, err := node.BestBlock()
		if err != nil {
			log.Fatalf("Failed to get Tip Height %d: %v", height, node.ChainParams().Name)
		}
		return blockHeader, blockHash, blockStamp.Height, nil
	}
	return nil, nil, 0, fmt.Errorf("Node %s does not exist", chainName)

}

func (ns *NeutrinoServer) GetLatestBlockHeader(chainName string) (*wire.BlockHeader, *chainhash.Hash, int32, error) {
	if liteNode, exists := ns.Nodes[chainName]; exists {
		node := liteNode.Node
		blockStamp, err := node.BestBlock()
		if err != nil {
			log.Fatalf("Failed to get Tip Height: %v", err)
		}
		blockHeader, err := node.GetBlockHeader(&blockStamp.Hash)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("Failed to get blockheader for tip %d: %v", blockStamp.Height, node.ChainParams().Name)
		}
		return blockHeader, &blockStamp.Hash, blockStamp.Height, nil
	}
	return nil, nil, 0, fmt.Errorf("Node %s does not exist", chainName)
}

// RPC Method Arguments and Reply Types
type GetBlockHeaderByHeightArgs struct {
	ChainName string
	Height    int64
}

type GetBlockHeaderByHeightReply struct {
	BlockHeader *api.BTCBlockHeader
	BlockHash   *chainhash.Hash
	Height      int32
}

func (ns *NeutrinoServer) BlockHeaderByHeight(args *GetBlockHeaderByHeightArgs, reply *GetBlockHeaderByHeightReply) error {
	log.Printf("BlockHeaderByHeight called with args: %+v", args)
	blockheader, blockhash, height, err := ns.GetBlockHeaderByHeight(args.ChainName, args.Height)
	if err != nil {
		log.Printf("Error fetching block header by height: %v", err)
		return err
	}

	bh := &api.BTCBlockHeader{
		Version:    int64(blockheader.Version),
		PrevBlock:  blockheader.PrevBlock.String(),
		MerkleRoot: blockheader.MerkleRoot.String(),
		TimeStamp:  blockheader.Timestamp.Unix(),
		Bits:       int64(blockheader.Bits),
		Nonce:      int64(blockheader.Nonce),
		BlockHash:  blockhash.String(),
	}

	reply.BlockHeader = bh
	reply.BlockHash = blockhash
	reply.Height = height
	log.Printf("BlockHeaderByHeight successfully fetched block: %+v", reply)
	return nil
}
