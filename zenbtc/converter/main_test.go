package main

import (
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/zenrocklabs/zenbtc/internal/chain"
	"github.com/zenrocklabs/zenbtc/internal/config"

	eigensdkLogger "github.com/Layr-Labs/eigensdk-go/logging"
)

var (
	cfg        *config.Config
	logger     eigensdkLogger.Logger
	ethClient  *ethclient.Client
	ethAccount *chain.EthAccount
)

func setup() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "./config.yaml"
	}

	var err error
	cfg, err = config.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if cfg.EthConfig.RPCUrl == "" {
		log.Fatalf("RPC URL is missing in config")
	}
	if len(cfg.EthConfig.Operators) == 0 {
		log.Fatalf("No operators defined in config")
	}

	logger = eigensdkLogger.NewNoopLogger()

	ethAccount, err = chain.NewEthAccount(cfg.EthConfig.Mnemonic, cfg.EthConfig.DerivationPath)
	if err != nil {
		log.Fatalf("Failed to create Ethereum account: %v", err)
	}

	ethClient, err = ethclient.Dial(cfg.EthConfig.RPCUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
}

func TestCreateAVSRewards(t *testing.T) {
	setup()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "create-avs-rewards",
				Action: createAVSRewardsCommand(logger, ethClient, ethAccount),
			},
		},
	}

	args := []string{"cmd", "create-avs-rewards", "1000000000000000000", "1710000000", "86400", "true"}
	err := app.Run(args)
	assert.Nil(t, err)
}

func TestClaimRewards(t *testing.T) {
	setup()
	if len(cfg.EthConfig.Operators) == 0 {
		t.Fatal("No operators found in config")
	}
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "claim",
				Action: claimCommand(logger, ethClient, ethAccount),
			},
		},
	}

	args := []string{"cmd", "claim", cfg.EthConfig.Operators[0], "true"}
	err := app.Run(args)
	assert.Nil(t, err)
}

func TestCheckBalance(t *testing.T) {
	setup()
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "balance",
				Action: balanceCommand(cfg, logger, ethClient, ethAccount),
			},
		},
	}

	args := []string{"cmd", "balance"}
	err := app.Run(args)
	assert.Nil(t, err)
}

func TestUnwrapWETH(t *testing.T) {
	setup()
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "unwrap",
				Action: unwrapCommand(cfg, logger, ethClient, ethAccount),
			},
		},
	}

	args := []string{"cmd", "unwrap", "1000000000000000000", "true"}
	err := app.Run(args)
	assert.Nil(t, err)
}

func TestSwapToBTC(t *testing.T) {
	setup()
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "swap",
				Action: swapCommand(cfg, logger, ethClient, ethAccount),
			},
		},
	}

	args := []string{"cmd", "swap", cfg.SwapConfig.BTCDestinationAddress, "1000000000000000000", "true"}
	err := app.Run(args)
	assert.Nil(t, err)
}

func TestFullCycle(t *testing.T) {
	setup()
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "full",
				Action: fullCycleCommand(cfg, logger, ethClient, ethAccount),
			},
		},
	}

	args := []string{"cmd", "full", "true"}
	err := app.Run(args)
	assert.Nil(t, err)
}
