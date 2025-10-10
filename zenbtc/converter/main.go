package main

import (
	"context"
	"flag"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/zenrocklabs/zenbtc/internal/chain"
	"github.com/zenrocklabs/zenbtc/internal/config"
	"github.com/zenrocklabs/zenbtc/internal/eigenlayer"
	"github.com/zenrocklabs/zenbtc/internal/thorchain"
	"github.com/zenrocklabs/zenbtc/internal/weth"

	eigensdkLogger "github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	flags      = flag.NewFlagSet("rewards_swap", flag.ContinueOnError)
	configFile = flags.String("config", "./config.yaml", "Path to config.yaml")
)

func init() {
	if len(os.Args) > 1 {
		flags.Parse(os.Args[1:])
	}

	envCfgFile := os.Getenv("CONFIG_FILE")
	if envCfgFile != "" {
		configFile = &envCfgFile
	}
}

func main() {
	cfg, err := config.ReadConfig(*configFile)
	if err != nil {
		panic(err)
	}

	logger := eigenlayer.GetLogger(true)

	ethAccount, err := chain.NewEthAccount(cfg.EthConfig.Mnemonic, cfg.EthConfig.DerivationPath)
	if err != nil {
		panic(err)
	}

	ethClient, err := ethclient.Dial(cfg.EthConfig.RPCUrl)
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:      "create-avs-rewards",
				Usage:     "Create AVS rewards submission",
				UsageText: "create-avs-rewards [amount] [startTimestamp] [duration] [broadcast]",
				Action:    createAVSRewardsCommand(logger, ethClient, ethAccount),
			},
			{
				Name:      "claim",
				Usage:     "claim eigenlayer rewards",
				UsageText: "claim [earner] [broadcast]",
				Action:    claimCommand(logger, ethClient, ethAccount),
			},
			{
				Name:      "balance",
				Usage:     "get eth and weth balance",
				UsageText: "balance",
				Action:    balanceCommand(cfg, logger, ethClient, ethAccount),
			},
			{
				Name:      "unwrap",
				Usage:     "unwrap weth to eth",
				UsageText: "unwrap [amount] [broadcast]",
				Action:    unwrapCommand(cfg, logger, ethClient, ethAccount),
			},
			{
				Name:      "swap",
				Usage:     "swap funds to btc",
				UsageText: "swap [destination] [amount] [broadcast]",
				Action:    swapCommand(cfg, logger, ethClient, ethAccount),
			},
			{
				Name:      "full",
				UsageText: "full [broadcast]",
				Usage:     "claim rewards, unwrap weth and swap to btc",
				Action:    fullCycleCommand(cfg, logger, ethClient, ethAccount),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func createAVSRewardsCommand(logger eigensdkLogger.Logger, ethClient *ethclient.Client, ethAccount *chain.EthAccount) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		amountParam := cCtx.Args().Get(0)
		amount, ok := new(big.Int).SetString(amountParam, 10)
		if !ok {
			return errors.New("invalid amount")
		}

		startTimestamp, err := strconv.ParseUint(cCtx.Args().Get(1), 10, 32)
		if err != nil {
			return errors.Wrap(err, "invalid start timestamp")
		}

		duration, err := strconv.ParseUint(cCtx.Args().Get(2), 10, 32)
		if err != nil {
			return errors.Wrap(err, "invalid duration")
		}

		broadcastParam := cCtx.Args().Get(3)
		broadcast := broadcastParam == "true"

		elClient, err := eigenlayer.NewEigenlayerClient(logger, ethClient, ethAccount)
		if err != nil {
			return err
		}

		_, err = elClient.CreateAVSRewardsSubmission(amount, uint32(startTimestamp), uint32(duration), broadcast)
		if err != nil {
			logger.Errorf("Error creating AVS rewards submission: %v", err)
			return errors.Wrap(err, "failed to create AVS rewards submission")
		}

		logger.Infof("AVS rewards submission successfully initiated")
		return nil
	}
}

func claimCommand(logger eigensdkLogger.Logger, ethClient *ethclient.Client, ethAccount *chain.EthAccount) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		earner := cCtx.Args().Get(0)
		if earner == "" {
			return errors.New("no earner specified")
		}
		broadcastParam := cCtx.Args().Get(1)
		broadcast := broadcastParam == "true"

		elClient, err := eigenlayer.NewEigenlayerClient(logger, ethClient, ethAccount)
		if err != nil {
			return err
		}

		_, err = elClient.ClaimRewards(earner, broadcast)
		if err != nil {
			logger.Errorf("Error while claiming rewards for operator %s: %v", earner, err)
			return errors.Wrapf(err, "error while claiming rewards for operator %s", earner)
		}
		return nil
	}
}

func balanceCommand(cfg *config.Config, logger eigensdkLogger.Logger, ethClient *ethclient.Client, ethAccount *chain.EthAccount) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		wc, err := weth.NewWETHClient(cfg.EthConfig.WETHAddress, logger, ethAccount, ethClient)
		if err != nil {
			return err
		}

		ethBalance, err := ethClient.BalanceAt(cCtx.Context, ethAccount.GetAddress(), nil)
		if err != nil {
			logger.Errorf("unable to get ETH balance for account %s, err: %w", ethAccount.GetAddress().String(), err)
			return err
		}

		wethBalance, err := wc.Balance(ethAccount.GetAddress())
		if err != nil {
			logger.Errorf("unable to get WETH balance for account %s, err: %w", ethAccount.GetAddress().String(), err)
			return err
		}
		logger.Infof("Balances for %s: %d ETH, %d WETH", ethAccount.GetAddress().String(), ethBalance.Int64(), wethBalance.Int64())
		return nil
	}
}

func unwrapCommand(cfg *config.Config, logger eigensdkLogger.Logger, ethClient *ethclient.Client, ethAccount *chain.EthAccount) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		amountParam := cCtx.Args().Get(0)
		amount, err := strconv.ParseInt(amountParam, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error while parsing amount %s", amountParam)
		}
		broadcastParam := cCtx.Args().Get(1)
		broadcast := broadcastParam == "true"

		wc, err := weth.NewWETHClient(cfg.EthConfig.WETHAddress, logger, ethAccount, ethClient)
		if err != nil {
			return err
		}

		_, err = wc.Unwrap(big.NewInt(amount), broadcast)
		if err != nil {
			return errors.Wrap(err, "error while unwrapping weth")
		}
		return nil
	}
}

func swapCommand(cfg *config.Config, logger eigensdkLogger.Logger, ethClient *ethclient.Client, ethAccount *chain.EthAccount) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		dest := cCtx.Args().Get(0)
		if dest == "" {
			return errors.New("no destination specified")
		}
		if dest == "" {
			return errors.New("no destination specified")
		}
		amountParam := cCtx.Args().Get(1)
		amount, err := strconv.ParseInt(amountParam, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error while parsing amount %s", amountParam)
		}

		broadcastParam := cCtx.Args().Get(2)
		broadcast := broadcastParam == "true"

		tc := thorchain.NewThorchainClient(logger, ethClient, ethAccount, cfg.SwapConfig.ThorNodeUrl)

		rcpt, err := tc.Swap(dest, big.NewInt(amount), cfg.SwapConfig.ToleranceBPS, broadcast)
		if err != nil {
			return err
		}

		logger.Debugf("receipt: %+v", rcpt)
		return nil
	}
}

func fullCycleCommand(cfg *config.Config, logger eigensdkLogger.Logger, ethClient *ethclient.Client, ethAccount *chain.EthAccount) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		broadcastParam := cCtx.Args().Get(0)
		broadcast := broadcastParam == "true"

		elClient, err := eigenlayer.NewEigenlayerClient(logger, ethClient, ethAccount)
		if err != nil {
			return err
		}

		tc := thorchain.NewThorchainClient(logger, ethClient, ethAccount, cfg.SwapConfig.ThorNodeUrl)

		wc, err := weth.NewWETHClient(cfg.EthConfig.WETHAddress, logger, ethAccount, ethClient)
		if err != nil {
			return err
		}

		// claim operator rewards
		for _, e := range cfg.EthConfig.Operators {
			logger.Infof("claiming rewards for %s", e)
			receipt, err := elClient.ClaimRewards(e, broadcast)
			if err != nil {
				logger.Errorf("error while claiming rewards for %s, skipping, err: %w", e, err)
				continue
			}

			if broadcast {
				logger.Infof("claimed rewards for %s succesfully, tx hash: %s", e, receipt.TxHash.String())
			}
		}

		// unwrap all available WETH
		balance, err := wc.Balance(ethAccount.GetAddress())
		if err != nil {
			logger.Errorf("unable to get WETH balance for account %s, err: %w", ethAccount.GetAddress().String(), err)
			return err
		}

		logger.Infof("WETH balance (%s): %d", ethAccount.GetAddress().String(), balance.Int64())

		if balance.Int64() > 0 {
			logger.Infof("unwrapping WETH")
			receipt, err := wc.Unwrap(balance, broadcast)
			if err != nil {
				logger.Errorf("unable to unwrap WETH for account %s, balance: %d, err: %w", ethAccount.GetAddress().String(), balance, err)
				return err
			}

			if broadcast {
				logger.Infof("unwrapped %d WETH for %s succesfully, tx hash: %s", balance, receipt.TxHash.String())
			}
		}

		// swap available ETH to BTC
		balance, err = ethClient.BalanceAt(context.Background(), ethAccount.GetAddress(), nil)
		if err != nil {
			logger.Errorf("error while readin ETH balance for %s, err: %w", ethAccount.GetAddress().String(), err)
			return err
		}
		logger.Infof("ETH balance (%s): %d", ethAccount.GetAddress().String(), balance.Int64())

		gasMargin := big.NewInt(0.002 * 1e18)
		amount := new(big.Int).Sub(balance, gasMargin)
		amountToSwap := new(big.Int).Div(amount, big.NewInt(1e10))

		logger.Infof("ETH amount minus gas margin: %d", amountToSwap.Int64())

		if amountToSwap.Int64() > 0 {
			quote, err := tc.GetSwapQuote(cfg.SwapConfig.BTCDestinationAddress, amountToSwap, cfg.SwapConfig.ToleranceBPS)
			if err != nil {
				logger.Errorf("unable to get swap quote, err: %w", err)
				return err
			}

			minAmount, err := strconv.ParseInt(quote.RecommendedMinAmountIn, 10, 64)
			if err != nil {
				logger.Errorf("unable to parse swap quote min amount, err: %w", err)
				return err
			}

			if amount.Int64() < minAmount {
				logger.Errorf("available amount to swap is smaller then min amount (%d < %d), aborting", amount, minAmount)
				return errors.New("available amount to swap is smaller then min amount")
			}

			receipt, err := tc.Swap(cfg.SwapConfig.BTCDestinationAddress, amountToSwap, cfg.SwapConfig.ToleranceBPS, broadcast)
			if err != nil {
				logger.Errorf("error while executing swap: %w", err)
				return err
			}

			if broadcast {
				logger.Infof("swap requested successfully, tx hash: %s", receipt.TxHash.String())
			}
		}

		return nil
	}
}
