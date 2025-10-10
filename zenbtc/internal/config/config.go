package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	EthConfig  EthConfig  `yaml:"eth_config"`
	SwapConfig SwapConfig `yaml:"swap_config"`
}

type EthConfig struct {
	Operators      []string `yaml:"operators"`
	Mnemonic       string   `yaml:"mnemonic"`
	DerivationPath string   `yaml:"derivation_path"`
	RPCUrl         string   `yaml:"rpc_url"`
	WETHAddress    string   `yaml:"weth_address"`
}

type SwapConfig struct {
	ThorNodeUrl           string `yaml:"thor_node_url"`
	BTCDestinationAddress string `yaml:"btc_destination_address"`
	ToleranceBPS          int    `yaml:"tolerance_bps"`
}

func ReadConfig(configFile string) (*Config, error) {
	cfgFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "read service config")
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(cfgFile, cfg); err != nil {
		return nil, errors.Wrap(err, "unmarshal service config")
	}
	return cfg, nil
}
