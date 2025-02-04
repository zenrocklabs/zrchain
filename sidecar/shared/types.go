package shared

import (
	"math/big"

	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
)

type OracleState struct {
	EigenDelegations           map[string]map[string]*big.Int `json:"eigenDelegations"`
	EthBlockHeight             uint64                         `json:"ethBlockHeight"`
	EthGasLimit                uint64                         `json:"ethGasLimit"`
	EthBaseFee                 uint64                         `json:"ethBaseFee"`
	EthTipCap                  uint64                         `json:"ethTipCap"`
	SolanaLamportsPerSignature uint64                         `json:"solanaLamportsPerSignature"`
	EthBurnEvents              []api.BurnEvent                `json:"ethBurnEvents"`
	Redemptions                []api.Redemption               `json:"redemptions"`
	ROCKUSDPrice               float64                        `json:"rockUSDPrice"`
	BTCUSDPrice                float64                        `json:"btcUSDPrice"`
	ETHUSDPrice                float64                        `json:"ethUSDPrice"` // TODO: remove field if we won't use ETH stake?
}
