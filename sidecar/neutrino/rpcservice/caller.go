package rpcservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonRpcRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     int           `json:"id"`
}

type JsonRpcResponse struct {
	Result json.RawMessage
}

type RpcCaller struct {
	url      string
	username string
	password string
}

func NewRpcCaller(url, username, password string) *RpcCaller {
	return &RpcCaller{url: url, username: username, password: password}
}

func (rpc *RpcCaller) CallRpcMethod(method string, requestParams []interface{}) (*json.RawMessage, error) {
	rpcRequest := &JsonRpcRequest{
		Method: method,
		Params: requestParams,
		ID:     1,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(rpcRequest); err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	// Create a new request
	req, err := http.NewRequest(http.MethodPost, rpc.url, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set BasicAuth and header
	req.SetBasicAuth(rpc.username, rpc.password)
	req.Header.Set("Content-Type", "application/json")

	// Get a http client and make the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Response not OK: %v", response.Status)
	}

	var rpcResponse JsonRpcResponse
	if err := json.NewDecoder(response.Body).Decode(&rpcResponse); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}
	response.Body.Close()
	return &rpcResponse.Result, nil
}
