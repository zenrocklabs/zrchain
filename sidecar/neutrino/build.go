package neutrino

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino/rpcservice"

	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
	"github.com/lightninglabs/neutrino"

	"github.com/btcsuite/btclog"
)

func buildNeutrinoNode(chainParams chaincfg.Params, logLevel btclog.Level, nodes *map[string]LiteNode, path string) (map[string]LiteNode, error) {

	chainName := chainParams.Name
	dataDir := path + "/neutrino_" + chainName
	dbPath := dataDir + "/neutrino_" + chainName + ".db"

	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Neutrino directory: %v", dataDir)
	}

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

func (ns *NeutrinoServer) Initialize(url, user, password, path string) {
	var err error
	nodes := make(map[string]LiteNode)

	//Testnet
	nodes, err = buildNeutrinoNode(chaincfg.TestNet3Params, btclog.LevelError, &nodes, path)
	if err != nil {
		log.Printf("Failed to Start Node Testnet3")
	}
	nodes["testnet"] = nodes["testnet3"]
	ns.Nodes = nodes

	//Mainnet
	nodes, err = buildNeutrinoNode(chaincfg.MainNetParams, btclog.LevelError, &nodes, path)
	if err != nil {
		log.Printf("Failed to Start Node Mainnet")
	}
	ns.Nodes = nodes

	// Register RPC server
	rpc.Register(ns)
	rpc.HandleHTTP()
	// Listen for requests on port 12345
	l, e := net.Listen("tcp", ":12345")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	//Connect to Proxy
	var caller *rpcservice.RpcCaller
	if url != "" {
		caller = rpcservice.NewRpcCaller(url, user, password)
		ns.Proxy = caller
	}
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

	//If we can't get the blockheader and we are not on Mainnet, try from the proxy
	var returnedError error
	if chainName != "mainnet" {
		blockHeader, hash, height, err := ns.ProxyGetBlockHeaderByHeight(chainName, height)
		if err == nil {
			return blockHeader, hash, height, err
		}
		returnedError = fmt.Errorf("Failed ProxyGetBlockHeaderByHeight %d does not exist error:%w", height, returnedError)
		//ignore this error - we can't get testnet data using the proxy fallback mechanism
	}
	return nil, nil, 0, fmt.Errorf("Node %s does not exist %w", chainName, returnedError)

}

func (ns *NeutrinoServer) GetLatestBlockHeader(chainName string) (*wire.BlockHeader, *chainhash.Hash, int32, error) {
	if liteNode, exists := ns.Nodes[chainName]; exists {
		node := liteNode.Node
		blockStamp, err := node.BestBlock()
		if err != nil {
			log.Fatalf("Failed to get Tip Height %s: %v", node.ChainParams().Name, err)
		}
		blockHeader, err := node.GetBlockHeader(&blockStamp.Hash)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("Failed to get blockheader for tip %d: %v", blockStamp.Height, node.ChainParams().Name)
		}
		return blockHeader, &blockStamp.Hash, blockStamp.Height, nil
	}

	//If we can't get the blockheader and we are not on Mainnet, try from the proxy
	var returnedError error
	if chainName != "mainnet" {
		blockHeader, hash, height, err := ns.ProxyGetLatestBlockHeader(chainName)
		if err == nil {
			return blockHeader, hash, height, err
		}
		returnedError = fmt.Errorf("Failed ProxyGetLatestBlockHeader %d does not exist error:%w", height, returnedError)
		//ignore this error - we can't get testnet data using the proxy fallback mechanism
	}
	return nil, nil, 0, fmt.Errorf("Node %s does not exist %w", chainName, returnedError)
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
