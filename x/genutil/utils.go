package genutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	cfg "github.com/cometbft/cometbft/config"
	tmed25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/go-bip39"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/genutil/types"
)

// ExportGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
func ExportGenesisFile(genesis *types.AppGenesis, genFile string) error {
	if err := genesis.ValidateAndComplete(); err != nil {
		return err
	}

	return genesis.SaveAs(genFile)
}

// ExportGenesisFileWithTime creates and writes the genesis configuration to disk.
// An error is returned if building or writing the configuration to file fails.
func ExportGenesisFileWithTime(genFile, chainID string, validators []cmttypes.GenesisValidator, appState json.RawMessage, genTime time.Time) error {
	appGenesis := types.NewAppGenesisWithVersion(chainID, appState)
	appGenesis.GenesisTime = genTime
	appGenesis.Consensus.Validators = validators

	if err := appGenesis.ValidateAndComplete(); err != nil {
		return err
	}

	return appGenesis.SaveAs(genFile)
}

// InitializeNodeValidatorFiles creates private validator and p2p configuration files.
func InitializeNodeValidatorFiles(config *cfg.Config) (nodeID string, valPubKey cryptotypes.PubKey, err error) {
	return InitializeNodeValidatorFilesFromMnemonic(config, "")
}

// InitializeNodeValidatorFilesFromMnemonic creates private validator and p2p configuration files using the given mnemonic.
// If no valid mnemonic is given, a random one will be used instead.
func InitializeNodeValidatorFilesFromMnemonic(config *cfg.Config, mnemonic string) (nodeID string, valPubKey cryptotypes.PubKey, err error) {
	if len(mnemonic) > 0 && !bip39.IsMnemonicValid(mnemonic) {
		return "", nil, fmt.Errorf("invalid mnemonic")
	}
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return "", nil, err
	}

	nodeID = string(nodeKey.ID())

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := os.MkdirAll(filepath.Dir(pvKeyFile), 0o777); err != nil {
		return "", nil, fmt.Errorf("could not create directory %q: %w", filepath.Dir(pvKeyFile), err)
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := os.MkdirAll(filepath.Dir(pvStateFile), 0o777); err != nil {
		return "", nil, fmt.Errorf("could not create directory %q: %w", filepath.Dir(pvStateFile), err)
	}

	var filePV *privval.FilePV
	if len(mnemonic) == 0 {
		filePV = privval.LoadOrGenFilePV(pvKeyFile, pvStateFile)
	} else {
		privKey := tmed25519.GenPrivKeyFromSecret([]byte(mnemonic))
		filePV = privval.NewFilePV(privKey, pvKeyFile, pvStateFile)
		filePV.Save()
	}

	tmValPubKey, err := filePV.GetPubKey()
	if err != nil {
		return "", nil, err
	}

	valPubKey, err = cryptocodec.FromCmtPubKeyInterface(tmValPubKey)
	if err != nil {
		return "", nil, err
	}

	return nodeID, valPubKey, nil
}
