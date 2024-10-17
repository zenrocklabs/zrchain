package neutrino

//
//import (
//	"fmt"
//	"gopkg.in/yaml.v3"
//	"log"
//	"os"
//)
//
//func LoadConfig() Config {
//	configFile := getConfigFile()
//	cfg, err := readConfig(configFile)
//	if err != nil {
//		log.Fatalf("Failed to read config: %v", err)
//	}
//	return cfg
//}
//
//func getConfigFile() string {
//	configFile := os.Getenv("SIDECAR_CONFIG_FILE")
//	if configFile == "" {
//		configFile = "config.yaml"
//	}
//	return configFile
//}
//
//func readConfig(configFile string) (Config, error) {
//	yamlFile, err := os.ReadFile(configFile)
//	if err != nil {
//		return Config{}, fmt.Errorf("unable to read config from %s: %v", configFile, err)
//	}
//
//	rootConfig := Config{}
//	if err = yaml.Unmarshal(yamlFile, &rootConfig); err != nil {
//		return Config{}, fmt.Errorf("error unmarshalling config from %s: %v", configFile, err)
//	}
//
//	return rootConfig, nil
//}
