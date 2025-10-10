package weth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/zenrocklabs/zenbtc/internal/chain"
	"github.com/zenrocklabs/zenbtc/internal/contracts"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	eigensdkLogger "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
)

type WETHClient interface {
	Balance(address common.Address) (*big.Int, error)
	Unwrap(amount *big.Int, broadcast bool) (*types.Receipt, error)
	Wrap(amount *big.Int, broadcast bool) (*types.Receipt, error)
}

type wethClient struct {
	logger     eigensdkLogger.Logger
	ethAccount *chain.EthAccount
	ethClient  *ethclient.Client
	chainId    *big.Int
	weth       *contracts.Weth
}

var _ WETHClient = (*wethClient)(nil)

func NewWETHClient(address string, logger eigensdkLogger.Logger, ethAccount *chain.EthAccount, ethClient *ethclient.Client) (WETHClient, error) {
	wethAddr := common.HexToAddress(address)
	weth, err := contracts.NewWeth(wethAddr, ethClient)
	if err != nil {
		return nil, err
	}

	chainId, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return &wethClient{
		ethAccount: ethAccount,
		ethClient:  ethClient,
		logger:     logger,
		weth:       weth,
		chainId:    chainId,
	}, nil
}

func (c *wethClient) Balance(address common.Address) (*big.Int, error) {
	return c.weth.BalanceOf(&bind.CallOpts{}, address)
}

func (c *wethClient) Unwrap(amount *big.Int, broadcast bool) (*types.Receipt, error) {
	ctx := context.Background()
	txMgr, err := c.getTxMgr()
	if err != nil {
		return nil, err
	}

	noSendTxOpts, err := txMgr.GetNoSendTxOpts()
	if err != nil {
		return nil, err
	}

	tx, err := c.weth.Withdraw(noSendTxOpts, amount)
	if err != nil {
		return nil, err
	}

	if broadcast {
		c.logger.Info("Broadcasting tx")
		receipt, err := txMgr.Send(ctx, tx)
		if err != nil {
			return nil, errors.Wrap(err, "failed broadcast withdraw")
		}

		c.logger.Infof("Unwrap transaction submitted successfully, %s, %d", receipt.TxHash.String(), c.chainId)
		return receipt, nil
	}

	return nil, nil
}

func (c *wethClient) Wrap(amount *big.Int, broadcast bool) (*types.Receipt, error) {
	ctx := context.Background()
	txMgr, err := c.getTxMgr()
	if err != nil {
		return nil, err
	}

	noSendTxOpts, err := txMgr.GetNoSendTxOpts()
	if err != nil {
		return nil, err
	}

	noSendTxOpts.Value = amount
	tx, err := c.weth.Deposit(noSendTxOpts)
	if err != nil {
		return nil, err
	}

	if broadcast {
		c.logger.Info("Broadcasting tx")
		receipt, err := txMgr.Send(ctx, tx)
		if err != nil {
			return nil, errors.Wrap(err, "failed broadcast withdraw")
		}

		c.logger.Infof("Wrap transaction submitted successfully, %s, %d", receipt.TxHash.String(), c.chainId)
		return receipt, nil
	}

	return nil, nil
}

func (c *wethClient) getTxMgr() (*txmgr.SimpleTxManager, error) {
	signerCfg := signerv2.Config{
		PrivateKey: c.ethAccount.GetPrivKey(),
	}

	sgn, sender, err := signerv2.SignerFromConfig(signerCfg, c.chainId)
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
