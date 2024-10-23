package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"

	"github.com/Zenrock-Foundation/zrchain/v4/sidecar/proto/api"
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
	currentState := s.oracle.currentState.Load().(*OracleState)

	contractState, err := json.Marshal(currentState.Delegations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state: %w", err)
	}

	return &api.SidecarStateResponse{
		Delegations:    contractState,
		ROCKUSDPrice:   fmt.Sprint(currentState.ROCKUSDPrice),
		ETHUSDPrice:    fmt.Sprint(currentState.ETHUSDPrice),
		EthBlockHeight: currentState.EthBlockHeight,
		EthBlockHash:   currentState.EthBlockHash,
		EthGasLimit:    currentState.EthGasLimit,
		EthBaseFee:     currentState.EthBaseFee,
		EthTipCap:      currentState.EthTipCap,
	}, nil
}

func (s *oracleService) GetSidecarStateByEthHeight(ctx context.Context, req *api.SidecarStateByEthHeightRequest) (*api.SidecarStateResponse, error) {
	state, err := s.oracle.getStateByEthHeight(req.EthBlockHeight)
	if err != nil {
		return nil, err
	}

	contractState, err := json.Marshal(state.Delegations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state: %w", err)
	}

	return &api.SidecarStateResponse{
		Delegations:    contractState,
		ROCKUSDPrice:   fmt.Sprint(state.ROCKUSDPrice),
		ETHUSDPrice:    fmt.Sprint(state.ETHUSDPrice),
		EthBlockHeight: state.EthBlockHeight,
		EthBlockHash:   state.EthBlockHash,
		EthGasLimit:    state.EthGasLimit,
		EthBaseFee:     state.EthBaseFee,
		EthTipCap:      state.EthTipCap,
	}, nil
}

func (s *oracleService) GetBitcoinBlockHeaderByHeight(ctx context.Context, req *api.BitcoinBlockHeaderByHeightRequest) (*api.BitcoinBlockHeaderResponse, error) {
	blockheader, blockhash, tipHeight, err := s.oracle.neutrinoServer.GetBlockHeaderByHeight(req.ChainName, req.BlockHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to GetBlockHeaderByHeight: %w", err)
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
		Version:    int64(blockheader.Version),
		PrevBlock:  blockheader.PrevBlock.String(),
		MerkleRoot: blockheader.MerkleRoot.String(),
		TimeStamp:  blockheader.Timestamp.Unix(),
		Bits:       int64(blockheader.Bits),
		Nonce:      int64(blockheader.Nonce),
		BlockHash:  blockhash.String(),
	}

	return &api.BitcoinBlockHeaderResponse{
		BlockHeader: bh,
		BlockHeight: int64(tipHeight),
		TipHeight:   int64(tipHeight),
	}, nil
}

func (s *oracleService) GetSolanaTransaction(ctx context.Context, req *api.SolanaTransactionRequest) (*api.SolanaTransactionResponse, error) {
	out, err := s.oracle.solanaClient.GetTransaction(
		ctx,
		solana.MustSignatureFromBase58(req.TxSignature),
		&rpc.GetTransactionOpts{
			Encoding: solana.EncodingBase64,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &api.SolanaTransactionResponse{
		TxSlot: out.Slot,
	}, nil
}

type DebugSolanaTransactionResponse struct {
	Tx *rpc.GetTransactionResult
}

func (s *oracleService) DebugGetSolanaTransaction(ctx context.Context, req *api.SolanaTransactionRequest) (*DebugSolanaTransactionResponse, error) {
	out, err := s.oracle.solanaClient.GetTransaction(
		ctx,
		solana.MustSignatureFromBase58(req.TxSignature),
		&rpc.GetTransactionOpts{
			Encoding: solana.EncodingBase64,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &DebugSolanaTransactionResponse{
		Tx: out,
	}, nil
}

func (s *oracleService) GetEthereumNonceAtHeight(ctx context.Context, req *api.EthereumNonceAtHeightRequest) (*api.EthereumNonceAtHeightResponse, error) {
	// nonce, err := s.oracle.EthClient.NonceAt(ctx, common.HexToAddress(req.Address), big.NewInt(int64(req.Height)))
	nonce, err := s.oracle.EthClient.NonceAt(ctx, common.HexToAddress(req.Address), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	return &api.EthereumNonceAtHeightResponse{
		Nonce: nonce,
	}, nil
}

// TODO: clean up stuff we don't need
func (s *oracleService) GetEthereumTransaction(ctx context.Context, req *api.EthereumTransactionRequest) (*api.EthereumTransactionResponse, error) {
	txHash := common.HexToHash(req.GetTxHash())

	_, isPending, err := s.oracle.EthClient.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}

	chainID, err := s.oracle.EthClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	if chainID != big.NewInt(17000) {
		return nil, fmt.Errorf("unsupported chain id: %s", chainID.String())
	}

	// signer := types.LatestSignerForChainID(chainID)

	// fromAddress, err := types.Sender(signer, tx)
	// if err != nil {
	// 	return nil, err
	// }

	// fromAddressHex := fromAddress.Hex()

	var blockNumber uint64
	// var blockTime uint64

	if !isPending {
		receipt, err := s.oracle.EthClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			return nil, err
		}

		blockNumber = receipt.BlockNumber.Uint64()

		// block, err := s.oracle.EthClient.BlockByNumber(ctx, receipt.BlockNumber)
		// if err != nil {
		// 	return nil, err
		// }

		// blockTime = block.Time()
	}

	// var txResponse *api.EthereumTransactionResponse

	// var toAddress string
	// if tx.To() != nil {
	// 	toAddress = tx.To().Hex()
	// } else {
	// 	toAddress = "" // Contract creation transaction
	// }

	// switch tx.Type() {
	// case types.LegacyTxType:
	// 	legacyTx := &api.EthereumLegacyTransaction{
	// 		Nonce:       tx.Nonce(),
	// 		GasPrice:    tx.GasPrice().String(),
	// 		GasLimit:    tx.Gas(),
	// 		To:          toAddress,
	// 		From:        fromAddressHex,
	// 		Value:       tx.Value().String(),
	// 		Data:        fmt.Sprintf("0x%x", tx.Data()),
	// 		BlockHeight: blockNumber,
	// 		TxHash:      tx.Hash().Hex(),
	// 		BlockTime:   blockTime,
	// 	}
	// 	txResponse = &api.EthereumTransactionResponse{
	// 		Tx: &api.EthereumTransactionResponse_LegacyTx{
	// 			LegacyTx: legacyTx,
	// 		},
	// 		IsPending: isPending,
	// 	}
	// case types.DynamicFeeTxType:
	// 	eip1559Tx := &api.EthereumEIP1559Transaction{
	// 		Nonce:                tx.Nonce(),
	// 		MaxPriorityFeePerGas: tx.GasTipCap().String(),
	// 		MaxFeePerGas:         tx.GasFeeCap().String(),
	// 		GasLimit:             tx.Gas(),
	// 		To:                   toAddress,
	// 		From:                 fromAddressHex,
	// 		Value:                tx.Value().String(),
	// 		Data:                 fmt.Sprintf("0x%x", tx.Data()),
	// 		BlockHeight:          blockNumber,
	// 		TxHash:               tx.Hash().Hex(),
	// 		BlockTime:            blockTime,
	// 	}
	// 	txResponse = &api.EthereumTransactionResponse{
	// 		Tx: &api.EthereumTransactionResponse_Eip1559Tx{
	// 			Eip1559Tx: eip1559Tx,
	// 		},
	// 		IsPending: isPending,
	// 	}
	// default:
	// 	return nil, fmt.Errorf("unsupported transaction type: %d", tx.Type())
	// }

	return &api.EthereumTransactionResponse{TxHeight: blockNumber}, nil
}
