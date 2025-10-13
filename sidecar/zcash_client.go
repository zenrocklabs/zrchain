package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
)

// ZcashClient handles communication with ZCash RPC node
type ZcashClient struct {
	rpcURL     string
	httpClient *http.Client
	enabled    bool
}

// ZCashRPCRequest represents a JSON-RPC request to ZCash node
type ZCashRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// ZCashRPCResponse represents a JSON-RPC response from ZCash node
type ZCashRPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID string `json:"id"`
}

// ZCashBlockHeader represents a ZCash block header from RPC
type ZCashBlockHeader struct {
	Hash              string  `json:"hash"`
	Confirmations     int64   `json:"confirmations"`
	Height            int64   `json:"height"`
	Version           int64   `json:"version"`
	MerkleRoot        string  `json:"merkleroot"`
	Time              int64   `json:"time"`
	Nonce             string  `json:"nonce"` // ZCash uses string for nonce in some versions
	Bits              string  `json:"bits"`
	Difficulty        float64 `json:"difficulty"`
	PreviousBlockHash string  `json:"previousblockhash"`
	NextBlockHash     string  `json:"nextblockhash,omitempty"`
}

// NewZcashClient creates a new ZCash RPC client
func NewZcashClient(rpcURL string, enabled bool) *ZcashClient {
	if !enabled || rpcURL == "" {
		return &ZcashClient{enabled: false}
	}

	return &ZcashClient{
		rpcURL: rpcURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		enabled: enabled,
	}
}

// IsEnabled returns whether the ZCash client is enabled
func (zc *ZcashClient) IsEnabled() bool {
	return zc.enabled
}

// call makes a JSON-RPC call to the ZCash node
func (zc *ZcashClient) call(ctx context.Context, method string, params []interface{}) (json.RawMessage, error) {
	if !zc.enabled {
		return nil, fmt.Errorf("zcash client is not enabled")
	}

	request := ZCashRPCRequest{
		JSONRPC: "1.0",
		ID:      "sidecar",
		Method:  method,
		Params:  params,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", zc.rpcURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := zc.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp ZCashRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

// GetBlockCount returns the current block height
func (zc *ZcashClient) GetBlockCount(ctx context.Context) (int64, error) {
	result, err := zc.call(ctx, "getblockcount", []interface{}{})
	if err != nil {
		return 0, err
	}

	var height int64
	if err := json.Unmarshal(result, &height); err != nil {
		return 0, fmt.Errorf("failed to parse block count: %w", err)
	}

	return height, nil
}

// GetBlockHash returns the block hash for a given height
func (zc *ZcashClient) GetBlockHash(ctx context.Context, height int64) (string, error) {
	result, err := zc.call(ctx, "getblockhash", []interface{}{height})
	if err != nil {
		return "", err
	}

	var hash string
	if err := json.Unmarshal(result, &hash); err != nil {
		return "", fmt.Errorf("failed to parse block hash: %w", err)
	}

	return hash, nil
}

// GetBlockHeader returns the block header for a given block hash
func (zc *ZcashClient) GetBlockHeader(ctx context.Context, blockHash string) (*ZCashBlockHeader, error) {
	// verbose=true returns JSON object instead of hex string
	result, err := zc.call(ctx, "getblockheader", []interface{}{blockHash, true})
	if err != nil {
		return nil, err
	}

	var header ZCashBlockHeader
	if err := json.Unmarshal(result, &header); err != nil {
		return nil, fmt.Errorf("failed to parse block header: %w", err)
	}

	return &header, nil
}

// GetBlockHeaderByHeight returns the block header for a given height
func (zc *ZcashClient) GetBlockHeaderByHeight(ctx context.Context, height int64) (*api.BTCBlockHeader, error) {
	if !zc.enabled {
		return nil, fmt.Errorf("zcash client is not enabled")
	}

	hash, err := zc.GetBlockHash(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("failed to get block hash for height %d: %w", height, err)
	}

	header, err := zc.GetBlockHeader(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get block header for hash %s: %w", hash, err)
	}

	// Convert ZCash header to BTCBlockHeader format (they're compatible)
	return &api.BTCBlockHeader{
		Version:    header.Version,
		PrevBlock:  header.PreviousBlockHash,
		MerkleRoot: header.MerkleRoot,
		TimeStamp:  header.Time,
		Bits:       parseBits(header.Bits),
		Nonce:      parseNonce(header.Nonce),
	}, nil
}

// GetLatestBlockHeader returns the latest block header
func (zc *ZcashClient) GetLatestBlockHeader(ctx context.Context) (*api.BTCBlockHeader, int64, error) {
	if !zc.enabled {
		return nil, 0, fmt.Errorf("zcash client is not enabled")
	}

	height, err := zc.GetBlockCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get block count: %w", err)
	}

	header, err := zc.GetBlockHeaderByHeight(ctx, height)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get latest block header: %w", err)
	}

	return header, height, nil
}

// parseBits converts hex bits string to int64
func parseBits(bitsHex string) int64 {
	var bits int64
	fmt.Sscanf(bitsHex, "%x", &bits)
	return bits
}

// parseNonce converts nonce string to int64
// ZCash may return nonce as hex string or decimal
func parseNonce(nonceStr string) int64 {
	var nonce int64
	// Try parsing as hex first
	if _, err := fmt.Sscanf(nonceStr, "%x", &nonce); err != nil {
		// Try parsing as decimal
		fmt.Sscanf(nonceStr, "%d", &nonce)
	}
	return nonce
}
