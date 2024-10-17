package cmd

import (
	"fmt"
	"os"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/cobra"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	"github.com/Zenrock-Foundation/zrchain/v4/app"
	"github.com/Zenrock-Foundation/zrchain/v4/app/params"
)

const DefaultSidecarAddr = "localhost:9191"

// NewRootCmd creates a new root command for zenrockd. It is called once in the
// main function.
func NewRootCmd() *cobra.Command {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	cfg.SetAddressVerifier(params.AddressVerifier)
	cfg.Seal()
	// we "pre"-instantiate the application for getting the injected/configured encoding configuration
	zrConfig := &params.ZRConfig{
		IsValidator: true,
		SidecarAddr: DefaultSidecarAddr,
	}
	tempApp := app.NewZenrockApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		false,
		simtestutil.NewAppOptionsWithFlagHome(tempDir()),
		[]wasmkeeper.Option{},
		zrConfig,
	)
	encodingConfig := EncodingConfig{
		InterfaceRegistry: tempApp.InterfaceRegistry(),
		Codec:             tempApp.AppCodec(),
		TxConfig:          tempApp.TxConfig(),
		Amino:             tempApp.LegacyAmino(),
	}

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:           version.AppName,
		Short:         "Zenrock Daemon",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL. This sign mode
			// is only available if the client is online.
			if !initClientCtx.Offline {
				enabledSignModes := append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
				txConfigOpts := tx.ConfigOptions{
					EnabledSignModes:           enabledSignModes,
					TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
				}
				txConfig, err := tx.NewTxConfigWithOptions(
					initClientCtx.Codec,
					txConfigOpts,
				)
				if err != nil {
					return err
				}

				initClientCtx = initClientCtx.WithTxConfig(txConfig)
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := initCometBFTConfig()

			if cmd.Name() == "start" {
				zrConfig.IsValidator = !cmd.Flags().Changed("non-validator")
				sidecarAddr, err := cmd.Flags().GetString("sidecar-addr")
				if err != nil {
					return err
				}
				if sidecarAddr == "" {
					sidecarAddr = DefaultSidecarAddr
				}
				zrConfig.SidecarAddr = sidecarAddr
			}

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
		PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
			// checks for the init subcommand
			if cmd.Name() == "init" {
				snapshotUrl, err := cmd.Flags().GetString("download-snapshot")
				if err != nil {
					return err
				}
				logger := log.NewLogger(os.Stderr)
				if snapshotUrl != "" {
					homePath, err := cmd.Flags().GetString("home")
					if err != nil {
						return fmt.Errorf("error getting --home flag: %v", err)
					}
					// Use DefaultNodeHome in case --home flag is not provided
					if homePath == "" {
						homePath = app.DefaultNodeHome
					}
					if err = downloadFromUrl(snapshotUrl, homePath, logger); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	initRootCmd(rootCmd, encodingConfig.TxConfig, tempApp.BasicModuleManager, zrConfig)

	autoCliOpts, err := enrichAutoCliOpts(tempApp.AutoCliOpts(), initClientCtx)
	if err != nil {
		panic(err)
	}

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}
