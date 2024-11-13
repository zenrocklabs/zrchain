package neutrino

import (
	"encoding/json"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"strconv"
	"time"
)

func (ns *NeutrinoServer) ProxyGetLatestBlockHeader(chainName string) (*wire.BlockHeader, *chainhash.Hash, int32, error) {
	if ns.Proxy == nil {
		return nil, nil, 0, nil
	}
	height, err := ns.ProxyBlockHeight(chainName)
	if err != nil {
		return nil, nil, 0, err
	}
	blockHash, err := ns.ProxyGetBlockhash(chainName, height)
	if err != nil {
		return nil, nil, 0, err
	}

	blockHeader, err := ns.ProxyGetBlockHeader(chainName, blockHash.String())
	if err != nil {
		return nil, nil, 0, err
	}
	return blockHeader, blockHash, int32(height), nil
}

func (ns *NeutrinoServer) ProxyGetBlockHeaderByHeight(chainName string, height int64) (*wire.BlockHeader, *chainhash.Hash, int32, error) {
	if ns.Proxy == nil {
		return nil, nil, 0, nil
	}

	blockHash, err := ns.ProxyGetBlockhash(chainName, uint64(height))
	if err != nil {
		return nil, nil, 0, err
	}

	blockHeader, err := ns.ProxyGetBlockHeader(chainName, blockHash.String())
	if err != nil {
		return nil, nil, 0, err
	}

	tipHeight, err := ns.ProxyBlockHeight(chainName)
	if err != nil {
		return nil, nil, 0, err
	}

	return blockHeader, blockHash, int32(tipHeight), nil
}

func (ns *NeutrinoServer) ProxyBlockHeight(chainName string) (uint64, error) {
	request := BlockCountRequest{
		ChainName: chainName,
		Quiet:     true,
	}
	//Get Tip Height
	response, err := ns.Proxy.CallRpcMethod("BitcoinServer.GetBlockCount", []interface{}{request})
	if err != nil {
		return 0, err
	}
	var blockResponse BlockCountResponse
	err = json.Unmarshal(*response, &blockResponse)
	if err != nil {
		return 0, err
	}
	return blockResponse.BlockHeight, nil
}

func (ns *NeutrinoServer) ProxyGetBlockhash(chainName string, blockHeight uint64) (*chainhash.Hash, error) {
	request := BlockHashRequest{
		ChainName:   chainName,
		BlockHeight: blockHeight,
		Quiet:       true,
	}
	response, err := ns.Proxy.CallRpcMethod("BitcoinServer.GetHashCount", []interface{}{request})
	if err != nil {
		return nil, err
	}
	var blockHashResponse BlockHashResponse
	err = json.Unmarshal(*response, &blockHashResponse)
	if err != nil {
		return nil, err
	}
	hash, err := chainhash.NewHashFromStr(blockHashResponse.BlockHash)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (ns *NeutrinoServer) ProxyGetBlockHeader(chainName string, blockhash string) (*wire.BlockHeader, error) {
	request := BlockHeaderRequest{
		ChainName: chainName,
		BlockHash: blockhash,
		Verbose:   true,
		Quiet:     true,
	}
	response2, err := ns.Proxy.CallRpcMethod("BitcoinServer.GetBlockHeader", []interface{}{request})
	if err != nil {
		return nil, err
	}

	var blockHeaderResponse BlockHeaderResponse
	err = json.Unmarshal(*response2, &blockHeaderResponse)
	if err != nil {
		return nil, err
	}

	bb := blockHeaderResponse.BlockHeader

	prevHash, err := chainhash.NewHashFromStr(bb.PreviousBlockHash)
	if err != nil {
		return nil, err
	}

	merkleRootHash, err := chainhash.NewHashFromStr(bb.MerkleRoot)
	if err != nil {
		return nil, err
	}

	bits, err := strconv.ParseInt(bb.Bits, 16, 64)
	if err != nil {
		return nil, err
	}
	bh := wire.NewBlockHeader(int32(bb.Version), prevHash, merkleRootHash, uint32(bits), uint32(bb.Nonce))
	bh.Timestamp = time.Unix(bb.Time, 0)

	return bh, nil

}
