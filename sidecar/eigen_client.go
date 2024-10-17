package main

import (
	"context"
	"encoding/json"
	"log"

	sdkutils "github.com/Layr-Labs/eigensdk-go/utils"
	"github.com/zenrocklabs/zenrock-avs/operator"
	"github.com/zenrocklabs/zenrock-avs/types"
)

func (o *Oracle) runEigenOperator() error {
	log.Println("Initializing Operator")
	operatorConfig := types.NodeConfig{}
	if err := sdkutils.ReadYamlConfig(o.Config.OperatorConfig, &operatorConfig); err != nil {
		return err
	}
	configJson, err := json.MarshalIndent(operatorConfig, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Config:", string(configJson))

	log.Println("initializing operator")
	operator, err := operator.NewOperatorFromConfig(operatorConfig)
	if err != nil {
		return err
	}
	log.Println("initialized operator")

	log.Println("starting operator")
	if err = operator.Start(context.Background()); err != nil {
		return err
	}
	log.Println("started operator")

	return nil
}
