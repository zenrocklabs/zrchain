package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"google.golang.org/grpc"
)

func startGRPCServer(oracle *Oracle, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterSidecarServiceServer(s, &oracleService{oracle: oracle})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type oracleService struct {
	api.UnimplementedSidecarServiceServer
	oracle *Oracle
}

func NewOracleService(oracle *Oracle) *oracleService {
	return &oracleService{
		oracle: oracle,
	}
}

func (s *oracleService) GetSidecarState(ctx context.Context, req *api.SidecarStateRequest) (*api.SidecarStateResponse, error) {
	currentState := s.oracle.currentState.Load()

	b, _ := json.Marshal(struct{}{})

	return &api.SidecarStateResponse{
		EigenDelegations:   b,
		EthBlockHeight:     currentState.EthBlockHeight,
		EthGasLimit:        currentState.EthGasLimit,
		EthBaseFee:         currentState.EthBaseFee,
		EthTipCap:          currentState.EthTipCap,
		EthBurnEvents:      currentState.EthBurnEvents,
		Redemptions:        currentState.Redemptions,
		ROCKUSDPrice:       currentState.ROCKUSDPrice.String(),
		BTCUSDPrice:        currentState.BTCUSDPrice.String(),
		ETHUSDPrice:        currentState.ETHUSDPrice.String(),
		ZECUSDPrice:        currentState.ZECUSDPrice.String(),
		SolanaBurnEvents:   currentState.SolanaBurnEvents,
		SolanaMintEvents:   currentState.SolanaMintEvents,
		SidecarVersionName: sidecartypes.SidecarVersionName,
	}, nil
}

func (s *oracleService) GetSidecarStateByEthHeight(ctx context.Context, req *api.SidecarStateByEthHeightRequest) (*api.SidecarStateResponse, error) {
	state, err := s.oracle.getStateByEthHeight(req.EthBlockHeight)
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(struct{}{})

	return &api.SidecarStateResponse{
		EigenDelegations:   b,
		EthBlockHeight:     state.EthBlockHeight,
		EthGasLimit:        state.EthGasLimit,
		EthBaseFee:         state.EthBaseFee,
		EthTipCap:          state.EthTipCap,
		EthBurnEvents:      state.EthBurnEvents,
		Redemptions:        state.Redemptions,
		ROCKUSDPrice:       state.ROCKUSDPrice.String(),
		BTCUSDPrice:        state.BTCUSDPrice.String(),
		ETHUSDPrice:        state.ETHUSDPrice.String(),
		ZECUSDPrice:        state.ZECUSDPrice.String(),
		SolanaMintEvents:   state.SolanaMintEvents,
		SolanaBurnEvents:   state.SolanaBurnEvents,
		SidecarVersionName: sidecartypes.SidecarVersionName,
	}, nil
}

func (s *oracleService) GetBitcoinBlockHeaderByHeight(ctx context.Context, req *api.BitcoinBlockHeaderByHeightRequest) (*api.BitcoinBlockHeaderResponse, error) {
	blockheader, blockhash, tipHeight, err := s.oracle.neutrinoServer.GetBlockHeaderByHeight(req.ChainName, req.BlockHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to GetBlockHeaderByHeight: %w", err)
	}

	bh := &api.BTCBlockHeader{
		Version:     int64(blockheader.Version),
		PrevBlock:   blockheader.PrevBlock.String(),
		MerkleRoot:  blockheader.MerkleRoot.String(),
		TimeStamp:   blockheader.Timestamp.Unix(),
		Bits:        int64(blockheader.Bits),
		Nonce:       int64(blockheader.Nonce),
		BlockHash:   blockhash.String(),
		BlockHeight: req.BlockHeight,
	}

	return &api.BitcoinBlockHeaderResponse{
		BlockHeader: bh,
		BlockHeight: req.BlockHeight,
		TipHeight:   int64(tipHeight),
	}, nil
}

func (s *oracleService) GetLatestBitcoinBlockHeader(ctx context.Context, req *api.LatestBitcoinBlockHeaderRequest) (*api.BitcoinBlockHeaderResponse, error) {
	blockheader, blockhash, tipHeight, err := s.oracle.neutrinoServer.GetLatestBlockHeader(req.ChainName)
	if err != nil {
		return nil, fmt.Errorf("failed to GetBlockHeaderByHeight: %w", err)
	}

	bh := &api.BTCBlockHeader{
		Version:     int64(blockheader.Version),
		PrevBlock:   blockheader.PrevBlock.String(),
		MerkleRoot:  blockheader.MerkleRoot.String(),
		TimeStamp:   blockheader.Timestamp.Unix(),
		Bits:        int64(blockheader.Bits),
		Nonce:       int64(blockheader.Nonce),
		BlockHash:   blockhash.String(),
		BlockHeight: int64(tipHeight),
	}

	return &api.BitcoinBlockHeaderResponse{
		BlockHeader: bh,
		BlockHeight: int64(tipHeight),
		TipHeight:   int64(tipHeight),
	}, nil
}

func (s *oracleService) GetZcashBlockHeaderByHeight(ctx context.Context, req *api.BitcoinBlockHeaderByHeightRequest) (*api.BitcoinBlockHeaderResponse, error) {
	if s.oracle.zcashClient == nil {
		return nil, fmt.Errorf("ZCash client not initialized")
	}

	blockHeader, err := s.oracle.zcashClient.GetBlockHeaderByHeight(ctx, req.BlockHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to GetZcashBlockHeaderByHeight: %w", err)
	}

	tipHeight, err := s.oracle.zcashClient.GetBlockCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ZCash tip height: %w", err)
	}

	return &api.BitcoinBlockHeaderResponse{
		BlockHeader: blockHeader,
		BlockHeight: req.BlockHeight,
		TipHeight:   tipHeight,
	}, nil
}

func (s *oracleService) GetLatestZcashBlockHeader(ctx context.Context, req *api.LatestBitcoinBlockHeaderRequest) (*api.BitcoinBlockHeaderResponse, error) {
	if s.oracle.zcashClient == nil {
		return nil, fmt.Errorf("ZCash client not initialized")
	}

	blockHeader, tipHeight, err := s.oracle.zcashClient.GetLatestBlockHeader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to GetLatestZcashBlockHeader: %w", err)
	}

	return &api.BitcoinBlockHeaderResponse{
		BlockHeader: blockHeader,
		BlockHeight: blockHeader.BlockHeight,
		TipHeight:   tipHeight,
	}, nil
}

func (s *oracleService) GetLatestEthereumNonceForAccount(ctx context.Context, req *api.LatestEthereumNonceForAccountRequest) (*api.LatestEthereumNonceForAccountResponse, error) {
	nonce, err := s.oracle.EthClient.NonceAt(ctx, common.HexToAddress(req.Address), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	return &api.LatestEthereumNonceForAccountResponse{
		Nonce: nonce,
	}, nil
}

func (s *oracleService) GetSolanaAccountInfo(ctx context.Context, req *api.SolanaAccountInfoRequest) (*api.SolanaAccountInfoResponse, error) {
	recipientKey, err := solana.PublicKeyFromBase58(req.PubKey)
	if err != nil {
		return nil, err
	}
	accountInfo, err := s.oracle.solanaClient.GetAccountInfoWithOpts(
		ctx,
		recipientKey,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
			DataSlice:  nil,
		},
	)
	if err != nil {
		if errors.Is(err, rpc.ErrNotFound) {
			log.Printf("Solana account %s not found (ErrNotFound)", req.PubKey)
			return &api.SolanaAccountInfoResponse{}, nil
		}
		log.Printf("Failed to get Solana account info for %s: %v", req.PubKey, err)
		return nil, fmt.Errorf("failed to get Solana account info: %w", err)
	}

	if accountInfo == nil || accountInfo.Value == nil {
		log.Printf("Solana account %s returned nil value", req.PubKey)
		return &api.SolanaAccountInfoResponse{}, nil
	}

	binary := accountInfo.GetBinary()
	log.Printf("Retrieved Solana account %s: %d bytes, owner: %s", req.PubKey, len(binary), accountInfo.Value.Owner.String())

	return &api.SolanaAccountInfoResponse{
		Account: binary,
	}, nil
}
