package main_test

import (
	"context"
	"fmt"
	"testing"

	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar"
	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	"github.com/test-go/testify/require"
)

func TestGetLatestEthereumNonce(t *testing.T) {
	oracle := initTestOracle()
	oracleService := sidecar.NewOracleService(oracle)
	out, err := oracleService.GetLatestEthereumNonceForAccount(context.Background(), &api.LatestEthereumNonceForAccountRequest{
		Address: "0xF198A1Af2682538E834bBc5F1af1847Cf50603E1",
	})
	require.NoError(t, err)
	fmt.Println("nonce is", out.Nonce)
}
