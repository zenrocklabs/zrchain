package e2e

import (
	"context"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	ginkgo "github.com/onsi/ginkgo/v2"
)

type TestEnv struct {
	Ctx    context.Context
	Ident  client.Identity
	Query  *client.QueryClient
	Tx     *client.TxClient
	Docker *Docker
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

	docker := &Docker{}

	// Create a context with timeout to prevent tests from hanging indefinitely
	testCtx, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	// Note: We don't call cancel() here because the context needs to live for the duration of the test
	// The test framework will handle cleanup when the test completes

	return &TestEnv{
		Ctx:    testCtx,
		Ident:  identity,
		Query:  queryClient,
		Tx:     txClient,
		Docker: docker,
	}
}
