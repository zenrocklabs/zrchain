package types_test

import (
	"encoding/json"
	"os"
	"testing"

	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"

	"github.com/Zenrock-Foundation/zrchain/v6/x/genutil/types"
)

func TestAppGenesis_Marshal(t *testing.T) {
	genesis := types.AppGenesis{
		AppName:    "simapp",
		AppVersion: "0.1.0",
		ChainID:    "test",
	}

	out, err := json.Marshal(&genesis)
	require.NoError(t, err)
	require.Equal(t, string(out), `{"app_name":"simapp","app_version":"0.1.0","genesis_time":"0001-01-01T00:00:00Z","chain_id":"test","initial_height":0,"app_hash":null}`)
}

func TestAppGenesis_Unmarshal(t *testing.T) {
	jsonBlob, err := os.ReadFile("testdata/app_genesis.json")
	require.NoError(t, err)

	var genesis types.AppGenesis
	err = json.Unmarshal(jsonBlob, &genesis)
	require.NoError(t, err)

	require.Equal(t, genesis.ChainID, "demo")
	require.Equal(t, genesis.Consensus.Params.Block.MaxBytes, int64(22020096))
}

func TestAppGenesis_ValidGenesis(t *testing.T) {
	// validate can read cometbft genesis file
	genesis, err := types.AppGenesisFromFile("testdata/cmt_genesis.json")
	require.NoError(t, err)

	require.Equal(t, genesis.ChainID, "demo")
	require.Equal(t, genesis.Consensus.Validators[0].Name, "test")

	// validate the app genesis can be translated properly to cometbft genesis
	cmtGenesis, err := genesis.ToGenesisDoc()
	require.NoError(t, err)
	rawCmtGenesis, err := cmttypes.GenesisDocFromFile("testdata/cmt_genesis.json")
	require.NoError(t, err)
	require.Equal(t, cmtGenesis, rawCmtGenesis)

	// validate can properly marshal to app genesis file
	rawAppGenesis, err := json.Marshal(&genesis)
	require.NoError(t, err)
	golden.Assert(t, string(rawAppGenesis), "app_genesis.json")

	// validate the app genesis can be unmarshalled properly
	var appGenesis types.AppGenesis
	err = json.Unmarshal(rawAppGenesis, &appGenesis)
	require.NoError(t, err)
	require.Equal(t, appGenesis.Consensus.Params, genesis.Consensus.Params)

	// validate marshaling of app genesis
	rawAppGenesis, err = json.Marshal(&appGenesis)
	require.NoError(t, err)
	golden.Assert(t, string(rawAppGenesis), "app_genesis.json")
}
