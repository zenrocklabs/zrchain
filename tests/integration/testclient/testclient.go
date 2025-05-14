package testclient

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/hd"

	client "github.com/Zenrock-Foundation/zrchain/v6/go-client"

	identitytypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	policytypes "github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"

	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type TestClient struct {
	txc *client.TxClient

	iqc identitytypes.QueryClient
	tqc treasurytypes.QueryClient
	pqc policytypes.QueryClient

	aqc authztypes.QueryClient
	bqc banktypes.QueryClient

	IdentitySigner  client.Identity
	IdentityAlice   client.Identity
	IdentityBob     client.Identity
	IdentityCharlie client.Identity

	Keyring string
}

const defaultTimeout = time.Second * 30

func GetTestClient() (context.Context, *TestClient) {
	url := os.Getenv("TEST_GRPC_URL")
	chainId := os.Getenv("TEST_CHAIN_ID")
	insecure := os.Getenv("TEST_USE_TLS") != "true"
	keyring := os.Getenv("TEST_KEYRING")
	timeoutStr := os.Getenv("TEST_TIMEOUT")
	timeoutSeconds, err := strconv.ParseInt(timeoutStr, 10, 64)
	timeout := defaultTimeout
	if err == nil {
		timeout = time.Second * time.Duration(timeoutSeconds)
	}

	env := os.Getenv("TEST_ENV")
	if env == "local" {
		url = "localhost:9090"
		chainId = "zenrock"
		insecure = true
		keyring = "keyring1pfnq7r04rept47gaf5cpdew2"
		timeout = defaultTimeout
	} else if env == "dev" {
		url = "grpc.dev.zenrock.tech"
		chainId = "amber-1"
		insecure = false
		keyring = "keyring1k6vc6vhp6e6l3rxalue9v4ux"
		timeout = time.Second * 60
	} else if env == "gardia" {
		url = "grpc.gardia.zenrocklabs.io"
		chainId = "gardia-5"
		insecure = false
		keyring = "keyring1k6vc6vhp6e6l3rxalue9v4ux"
		timeout = time.Second * 60
	}

	if url == "" {
		url = "localhost:9090"
	}

	if chainId == "" {
		chainId = "zenrock"
	}

	if keyring == "" {
		keyring = "keyring1pfnq7r04rept47gaf5cpdew2" //gitleaks:allow
	}

	c, err := NewTestClient(url, chainId, insecure, keyring)
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx, c
}

func NewTestClient(url string, chainId string, insecure bool, keyring string) (*TestClient, error) {
	c := &TestClient{
		Keyring: keyring,
	}

	var err error

	derivationPath := hd.BIP44Params{Purpose: 44, CoinType: 118, Account: 0, Change: false, AddressIndex: 0}.String()
	c.IdentityAlice, err = client.NewIdentityFromSeed(derivationPath, "strategy social surge orange pioneer tiger skill endless lend slide one jazz pipe expose detect soup fork cube trigger frown stick wreck ring tissue")
	if err != nil {
		return nil, err
	}
	c.IdentitySigner = c.IdentityAlice

	c.IdentityBob, err = client.NewIdentityFromSeed(derivationPath, "fee buzz avocado dolphin syrup rule access cave close puppy lemon round able bronze fame give spoon company since fog error trip toast unable")
	if err != nil {
		return nil, err
	}

	c.IdentityCharlie, err = client.NewIdentityFromSeed(derivationPath, "clip broken sight warfare boring borrow orchard trumpet isolate wire police behave round dream pattern")
	if err != nil {
		return nil, err
	}

	grpcConn, err := c.createGRPCClient(url, insecure)
	if err != nil {
		return nil, err
	}

	txc, err := c.createZenrockClient(grpcConn, chainId)
	if err != nil {
		return nil, err
	}
	c.txc = txc

	c.iqc = identitytypes.NewQueryClient(grpcConn)
	c.pqc = policytypes.NewQueryClient(grpcConn)
	c.tqc = treasurytypes.NewQueryClient(grpcConn)
	c.aqc = authztypes.NewQueryClient(grpcConn)
	c.bqc = banktypes.NewQueryClient(grpcConn)

	return c, nil
}

func (c *TestClient) createGRPCClient(url string, notls bool) (*grpc.ClientConn, error) {
	var grpcOpts []grpc.DialOption
	if notls {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})))
	}

	return grpc.NewClient(
		url,
		grpcOpts...,
	)
}

func (c *TestClient) createZenrockClient(grpcConn *grpc.ClientConn, chainId string) (*client.TxClient, error) {
	queryClient := client.NewQueryClientWithConn(grpcConn)
	txclient, err := client.NewTxClient(c.IdentitySigner, chainId, grpcConn, queryClient)

	return txclient, err
}

func (c *TestClient) executeTx(ctx context.Context, msgs ...cosmostypes.Msg) (*cosmostypes.TxMsgData, error) {
	res, _, err := c.executeTxWithIdentity(ctx, c.IdentitySigner, msgs...)
	return res, err
}

func (c *TestClient) executeTxWithIdentity(ctx context.Context, identity client.Identity, msgs ...cosmostypes.Msg) (*cosmostypes.TxMsgData, []abcitypes.Event, error) {
	txb, err := c.txc.BuildAndSignTxForIdentity(ctx, identity, client.DefaultGasLimit, client.DefaultFees, msgs...)
	if err != nil {
		return nil, nil, err
	}

	hash, err := c.txc.SendWaitTx(ctx, txb)
	if err != nil {
		return nil, nil, err
	}

	txres, _ := c.txc.GetTx(ctx, hash)
	if txres.TxResponse.Code != 0 {
		return nil, nil, errors.New(txres.TxResponse.RawLog)
	}

	txm := &cosmostypes.TxMsgData{}
	txresb, err := hex.DecodeString(txres.TxResponse.Data)
	if err != nil {
		return nil, nil, err
	}

	err = proto.Unmarshal(txresb, txm)
	return txm, txres.TxResponse.Events, err
}
