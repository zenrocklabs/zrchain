package thorchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/zenrocklabs/zenbtc/internal/chain"
	"github.com/zenrocklabs/zenbtc/internal/contracts"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	eigensdkLogger "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
)

type thorchainInboundAddress struct {
	Chain   string `json:"chain"`
	Address string `json:"address"`
	Router  string `json:"router"`
	Halted  bool   `json:"halted"`
}

type thorchainFees struct {
	Asset       string `json:"asset"`
	Affiliate   string `json:"affiliate"`
	Outbond     string `json:"outbound"`
	Liquidity   string `json:"liquidity"`
	Total       string `json:"total"`
	SlippageBPS int    `json:"slippage_bps"`
	TotalBPS    int    `json:"total_bps"`
}

type ThorchainQuote struct {
	Address string `json:"inbound_address"`
	Router  string `json:"router"`
	Memo    string `json:"memo"`
	Expiry  uint64 `json:"expiry"`
	Error   string `json:"error"`

	Fees thorchainFees `json:"fees"`

	InboundConfirmationBlocks  int    `json:"inbound_confirmation_blocks"`
	InboundConfirmationSeconds int    `json:"inbound_confirmation_seconds"`
	OutboundDelayBlocks        int    `json:"outbound_delay_blocks"`
	OutboundDelaySeconds       int    `json:"outbound_delay_seconds"`
	RecommendedMinAmountIn     string `json:"recommended_min_amount_in"`
	RecommendedGasRate         string `json:"recommended_gas_rate"`
	ExpectedAmountOut          string `json:"expected_amount_out"`
	MaxStreamingQuantity       int    `json:"max_streaming_quantity"`
	StreamingSwapBlocks        int    `json:"streaming_swap_blocks"`
	TotalSwapSeconds           int    `json:"total_swap_seconds"`
}

func executeRequest[T any](url string) (*T, error) {
	res := new(T)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request failed")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type ThorChainClient interface {
	GetSwapQuote(destination string, amount *big.Int, toleranceBPS int) (*ThorchainQuote, error)
	Swap(destination string, amount *big.Int, toleranceBPS int, broadcast bool) (*types.Receipt, error)
}

type thorChainClient struct {
	thorNodeUrl string
	ethClient   *ethclient.Client
	ethAccount  *chain.EthAccount
	logger      eigensdkLogger.Logger
}

var _ ThorChainClient = (*thorChainClient)(nil)

func NewThorchainClient(logger eigensdkLogger.Logger, ethClient *ethclient.Client, ethAccount *chain.EthAccount, thorNodeUrl string) ThorChainClient {
	return &thorChainClient{
		ethClient:   ethClient,
		ethAccount:  ethAccount,
		thorNodeUrl: thorNodeUrl,
		logger:      logger,
	}
}

func (c *thorChainClient) GetSwapQuote(destination string, amount *big.Int, toleranceBPS int) (*ThorchainQuote, error) {
	inboundAddresses, err := executeRequest[[]thorchainInboundAddress](fmt.Sprintf("%s/thorchain/inbound_addresses", c.thorNodeUrl))
	if err != nil {
		return nil, err
	}

	var ethAddr thorchainInboundAddress
	for _, addr := range *inboundAddresses {
		if addr.Chain == "ETH" {
			ethAddr = addr
			break
		}
	}

	if ethAddr.Halted {
		return nil, errors.New("eth destination is halted")
	}

	fromAsset := "ETH.ETH"
	toAsset := "BTC.BTC"

	quoteUrl := fmt.Sprintf(
		"%s/thorchain/quote/swap?from_asset=%s&to_asset=%s&amount=%d&destination=%s&tolerance_bps=%d",
		c.thorNodeUrl,
		fromAsset,
		toAsset,
		amount,
		destination,
		toleranceBPS,
	)
	c.logger.Debugf("requesting thorchain quote: %s", quoteUrl)

	quote, err := executeRequest[ThorchainQuote](quoteUrl)

	if quote.Error != "" {
		return nil, fmt.Errorf("error while getting quote: %s", quote.Error)
	}

	return quote, err
}

func (c *thorChainClient) Swap(destination string, amount *big.Int, toleranceBPS int, broadcast bool) (*types.Receipt, error) {
	quote, err := c.GetSwapQuote(destination, amount, toleranceBPS)
	if err != nil {
		return nil, err
	}

	vaultAddress := common.HexToAddress(quote.Address)
	routerAddress := common.HexToAddress(quote.Router)
	memo := quote.Memo
	expiry := time.Now().Add(time.Hour * 2).Unix()

	// Convert amount from 1e8 (THORChain precision) to 1e18 (ETH Wei)
	// Multiple by 1e10 since 1e18/1e8 = 1e10
	nativeAmount := new(big.Int).Mul(amount, big.NewInt(1e10))

	// First send ETH to vault with memo
	memoData := common.Hex2Bytes(common.Bytes2Hex([]byte(memo)))

	txMgr, err := c.getTxMgr()
	if err != nil {
		return nil, err
	}

	noSend, err := txMgr.GetNoSendTxOpts()
	if err != nil {
		return nil, err
	}

	// Get nonce for the sender
	nonce, err := c.ethClient.PendingNonceAt(context.Background(), noSend.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price suggestions
	gasTipCap, err := c.ethClient.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas tip cap: %w", err)
	}

	baseFee, err := c.ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get base fee: %w", err)
	}

	// Calculate fee cap as baseFee*2 + tip
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee.BaseFee, big.NewInt(2)),
		gasTipCap,
	)

	// Create transaction data for sending ETH
	ethTx := &types.DynamicFeeTx{
		To:        &vaultAddress,
		Value:     nativeAmount, // Use converted amount for ETH transfer
		Data:      memoData,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
	}

	// Estimate gas for ETH transfer
	gasLimit, err := c.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  noSend.From,
		To:    ethTx.To,
		Value: ethTx.Value,
		Data:  ethTx.Data,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas for ETH transfer: %w", err)
	}
	ethTx.Gas = gasLimit

	tx := types.NewTx(ethTx)

	// Send ETH transaction and wait for receipt
	if broadcast {
		_, err := txMgr.Send(context.Background(), tx)
		if err != nil {
			return nil, fmt.Errorf("failed to send ETH: %w", err)
		}

		// Then call router
		tr, err := contracts.NewThorchainrouter(routerAddress, c.ethClient)
		if err != nil {
			return nil, err
		}

		noSend, err = txMgr.GetNoSendTxOpts()
		if err != nil {
			return nil, err
		}

		// Use same nativeAmount for router call since both need 1e18 precision
		routerTx, err := tr.DepositWithExpiry(noSend, vaultAddress, common.Address{}, nativeAmount, memo, big.NewInt(expiry))
		if err != nil {
			return nil, fmt.Errorf("failed to call router: %w", err)
		}

		// Send router transaction and wait for receipt
		receipt, err := txMgr.Send(context.Background(), routerTx)
		return receipt, err
	}

	return nil, nil
}

func (c *thorChainClient) getTxMgr() (*txmgr.SimpleTxManager, error) {
	signerCfg := signerv2.Config{
		PrivateKey: c.ethAccount.GetPrivKey(),
	}

	chainId, err := c.ethClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	if chainId.Int64() != 1 {
		return nil, fmt.Errorf("thorchain only supported on mainnet, found chainId %d", chainId.Int64())
	}

	sgn, sender, err := signerv2.SignerFromConfig(signerCfg, chainId)
	if err != nil {
		return nil, err
	}
	keyWallet, err := wallet.NewPrivateKeyWallet(c.ethClient, sgn, sender, c.logger)
	if err != nil {
		return nil, err
	}

	txMgr := txmgr.NewSimpleTxManager(keyWallet, c.ethClient, c.logger, sender)

	return txMgr, nil
}
