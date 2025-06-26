package integration_test

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	ginkgo "github.com/onsi/ginkgo/v2"
)

type TestEnv struct {
	Ctx   context.Context
	Ident client.Identity
	Query *client.QueryClient
	Tx    *client.TxClient
}

func setupTestEnv(t ginkgo.FullGinkgoTInterface) *TestEnv {
	t.Helper()

	identity, err := client.NewIdentityFromSeed(
		"m/44'/118'/0'/0/0",
		"strategy social surge orange pioneer tiger skill endless lend slide one jazz pipe expose detect soup fork cube trigger frown stick wreck ring tissue",
	)
	if err != nil {
		t.Fatalf("Failed to create identity: %v", err)
	}

	queryClient, err := client.NewQueryClient("127.0.0.1:9790", true)
	if err != nil {
		t.Fatalf("Failed to create query client: %v", err)
	}

	conn, err := client.NewClientConn("127.0.0.1:9790", true)
	if err != nil {
		t.Fatalf("Failed to create GRPC conn: %v", err)
	}

	txClient, err := client.NewTxClient(identity, "docker", conn, queryClient)
	if err != nil {
		t.Fatalf("Failed to create tx client: %v", err)
	}

	return &TestEnv{
		Ctx:   context.Background(),
		Ident: identity,
		Query: queryClient,
		Tx:    txClient,
	}
}
