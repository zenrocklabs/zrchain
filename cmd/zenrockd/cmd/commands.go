package cmd

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/crisis"

	"github.com/Zenrock-Foundation/zrchain/v6/app"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	genutilcli "github.com/Zenrock-Foundation/zrchain/v6/x/genutil/client/cli"
)

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.

	type CustomAppConfig struct {
		serverconfig.Config
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()
	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In simapp, we set the min gas prices to 0.
	// srvCfg.MinGasPrices = "0.0001urock"
	// srvCfg.BaseConfig.IAVLDisableFastNode = true // disable fastnode by default

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate

	return customAppTemplate, customAppConfig
}

func initRootCmd(
	rootCmd *cobra.Command,
	txConfig client.TxConfig,
	basicManager module.BasicManager,
	zrConfig *params.ZRConfig,
) {
	cfg := sdk.GetConfig()
	cfg.Seal()

	// Add download-snapshot flag to init cmd
	genUtilInitCmd := genutilcli.InitCmd(basicManager, app.DefaultNodeHome)
	genUtilInitCmd.Flags().String("download-snapshot", "", "Initialize from a snapshot")

	rootCmd.AddCommand(
		genUtilInitCmd,
		debug.Cmd(),
		confixcmd.ConfigCommand(),
		pruning.Cmd(newAppWithConfig(zrConfig), app.DefaultNodeHome),
		snapshot.Cmd(newAppWithConfig(zrConfig)),
	)

	// Modify the AddCommands call to use a custom server.AppCreator
	server.AddCommands(rootCmd, app.DefaultNodeHome,
		func(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
			return newAppWithConfig(zrConfig)(logger, db, traceStore, appOpts)
		},
		appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, genesis, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(txConfig, basicManager),
		queryCommand(),
		txCommand(),
		keys.Commands(),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

// genesisCommand builds genesis-related `simd genesis` command. Users may provide application specific commands as a parameter
func genesisCommand(txConfig client.TxConfig, basicManager module.BasicManager, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.Commands(txConfig, basicManager, app.DefaultNodeHome)

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
	)

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)

	return cmd
}

// newAppWithConfig returns a servertypes.AppCreator function that includes a config parameter
func newAppWithConfig(zrConfig *params.ZRConfig) servertypes.AppCreator {
	return func(
		logger log.Logger,
		db dbm.DB,
		traceStore io.Writer,
		appOpts servertypes.AppOptions,
	) servertypes.Application {
		baseappOptions := server.DefaultBaseappOptions(appOpts)

		skipUpgradeHeights := make(map[int64]bool)
		for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
			skipUpgradeHeights[int64(h)] = true
		}

		return app.NewZenrockApp(
			logger,
			db,
			traceStore,
			true,
			appOpts,
			zrConfig,
			baseappOptions...,
		)
	}
}

// appExport creates a new wasm app (optionally at a given height) and exports state.
func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var ZenrockApp *app.ZenrockApp
	// this check is necessary as we use the flag in x/upgrade.
	// we can exit more gracefully by checking the flag here.
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home is not set")
	}

	viperAppOpts, ok := appOpts.(*viper.Viper)
	if !ok {
		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	// overwrite the FlagInvCheckPeriod
	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
	appOpts = viperAppOpts

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(sdkserver.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	ZenrockApp = app.NewZenrockApp(
		logger,
		db,
		traceStore,
		height == -1,
		appOpts,
		nil,
	)

	if height != -1 {
		if err := ZenrockApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return ZenrockApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

func enrichAutoCliOpts(autoCliOpts autocli.AppOptions, clientCtx client.Context) (autocli.AppOptions, error) {
	autoCliOpts.AddressCodec = authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	autoCliOpts.ValidatorAddressCodec = authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	autoCliOpts.ConsensusAddressCodec = authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix())

	var err error
	clientCtx, err = config.ReadFromClientConfig(clientCtx)
	if err != nil {
		return autocli.AppOptions{}, err
	}

	autoCliOpts.ClientCtx = clientCtx
	// autoCliOpts.Keyring, err = keyring.NewAutoCLIKeyring(clientCtx.Keyring)
	// if err != nil {
	// 	return autocli.AppOptions{}, err
	// }

	return autoCliOpts, nil
}

var tempDir = func() string {
	dir, err := os.MkdirTemp("", "zenrockd")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(dir)

	return dir
}

// downloadFromUrl downloads a file from the specified URL and saves it to the given destination directory.
func downloadFromUrl(fileUrl string, destinationDir string, log log.Logger) (err error) {
	log.Info("Downloading the snapshot tar file", "url", fileUrl)
	client := &http.Client{}

	// Define the path for the temporary file
	tmpFilePath := filepath.Join(destinationDir, "snapshot.tar.gz")
	out, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	// Send the request to download the file
	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		return fmt.Errorf("error while creating the request url %s: %v", fileUrl, err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while requesting the file url %s: %v", fileUrl, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error while copying the response body to the file: %v", err)
	}

	log.Info("Unzipping the data from the snapshot tar file", "path", destinationDir)
	log.Info("In progress...")

	// Open the file for reading
	file, err := os.Open(tmpFilePath)
	if err != nil {
		return fmt.Errorf("error opening downloaded file: %v", err)
	}
	defer file.Close()

	// Uncompress the file
	_, err = UnTar(destinationDir, file)
	if err != nil {
		return fmt.Errorf("error while unzipping the tar file: %v", err)
	}

	log.Info("File copy done.")
	log.Info("Data unzipped.")

	// Delete temporary file
	err = os.Remove(tmpFilePath)
	if err != nil {
		return fmt.Errorf("error deleting temporary file: %v", err)
	}

	return nil
}

// UnTar extracts files from a gzip-compressed tar archive read from r and saves them to the
// specified destination directory. The file names and folder structure in the tar archive
// are preserved during extraction.
func UnTar(dst string, r io.Reader) (string, error) {
	// Create a gzip reader
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return "", fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Walk through the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading tar archive: %v", err)
		}

		// Extract files to dst folder
		targetPath := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(targetPath, 0755)
			if err != nil {
				return "", fmt.Errorf("error creating directory %s: %v", targetPath, err)
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return "", fmt.Errorf("error creating file %s: %v", targetPath, err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return "", fmt.Errorf("error writing to file %s: %v", targetPath, err)
			}

		default:
			return "", fmt.Errorf("unsupported file type %s", string(header.Typeflag))
		}
	}

	return dst, nil
}
