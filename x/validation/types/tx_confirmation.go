package types

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/gogoproto/proto"
)

var (
	RequestedTxConfirmationKey   = collections.NewPrefix(1000)
	RequestedTxConfirmationIndex = "requested_tx_confirmation"
)

// TxConfirmationData holds the transaction hashes that need confirmation for a key ID
type TxConfirmationData struct {
	TxHashes []string `protobuf:"bytes,1,rep,name=tx_hashes,json=txHashes,proto3" json:"tx_hashes,omitempty"`
}

// Reset implements proto.Message
func (t *TxConfirmationData) Reset() {
	*t = TxConfirmationData{}
}

// String implements proto.Message
func (t *TxConfirmationData) String() string {
	return proto.CompactTextString(t)
}

// ProtoMessage implements proto.Message
func (*TxConfirmationData) ProtoMessage() {}

// AddTxHash adds a transaction hash to the list if it doesn't already exist
func (t *TxConfirmationData) AddTxHash(txHash string) {
	for _, hash := range t.TxHashes {
		if hash == txHash {
			return
		}
	}
	t.TxHashes = append(t.TxHashes, txHash)
}

// RemoveFirstTxHash removes and returns the first transaction hash from the list
func (t *TxConfirmationData) RemoveFirstTxHash() string {
	if len(t.TxHashes) == 0 {
		return ""
	}
	firstHash := t.TxHashes[0]
	t.TxHashes = t.TxHashes[1:]
	return firstHash
}

// GetFirstTxHash returns the first transaction hash from the list without removing it
func (t *TxConfirmationData) GetFirstTxHash() string {
	if len(t.TxHashes) == 0 {
		return ""
	}
	return t.TxHashes[0]
}

// HasTxHashes returns true if there are any transaction hashes in the list
func (t *TxConfirmationData) HasTxHashes() bool {
	return len(t.TxHashes) > 0
}
