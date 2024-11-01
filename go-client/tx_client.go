package client

import (
	"fmt"

	"google.golang.org/grpc"
)

// TxClient can read/write transactions to zenrockd and endpoints provided by the treasury module.
type TxClient struct {
	*RawTxClient
	*TreasuryTxClient
	*ZenBTCTxClient
}

// NewTxClient returns a TxClient.
func NewTxClient(id Identity, chainID string, c *grpc.ClientConn, accountFetcher AccountFetcher) (*TxClient, error) {
	raw, err := NewRawTxClient(id, chainID, c, accountFetcher)
	if err != nil {
		return nil, fmt.Errorf("can't create raw tx client: %w", err)
	}
	return &TxClient{
		RawTxClient:      raw,
		TreasuryTxClient: NewTreasuryTxClient(raw),
		ZenBTCTxClient:   NewZenBTCTxClient(raw),
	}, nil
}
