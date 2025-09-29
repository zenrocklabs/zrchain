package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	signingv1beta1 "cosmossdk.io/api/cosmos/tx/signing/v1beta1"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"google.golang.org/grpc"

	appparams "github.com/Zenrock-Foundation/zrchain/v6/app/params"
)

var (
	DefaultGasLimit      = uint64(300000)
	ZenBTCGasLimit       = uint64(8000000)
	InjectHeaderGasLimit = uint64(20000000)
	DefaultFees          = types.NewCoins(types.NewCoin("urock", math.NewInt(750000)))
	ZenBTCDefaultFees    = types.NewCoins(types.NewCoin("urock", math.NewInt(20000000)))

	queryTimeout = 250 * time.Millisecond
)

type AccountFetcher interface {
	Account(ctx context.Context, addr string) (types.AccountI, error)
}

var _ AccountFetcher = (*QueryClient)(nil)

// RawTxClient is the client used for sending new transactions to the chain.
type RawTxClient struct {
	Identity Identity

	chainID        string
	client         txtypes.ServiceClient
	accountFetcher AccountFetcher
	app            *app.ZenrockApp
}

func NewRawTxClient(id Identity, chainID string, c *grpc.ClientConn, accountFetcher AccountFetcher) (*RawTxClient, error) {
	client := &RawTxClient{
		Identity:       id,
		chainID:        chainID,
		client:         txtypes.NewServiceClient(c),
		accountFetcher: accountFetcher,
	}

	os.RemoveAll("./tmp")
	os.MkdirAll("./tmp", 0777)
	app.DefaultNodeHome = "./tmp"

	zrConfig := &appparams.ZRConfig{
		IsValidator: false,
		SidecarAddr: "",
	}
	app := app.NewZenrockApp(
		log.NewNopLogger(),
		db.NewMemDB(),
		nil,
		false,
		simtestutil.NewAppOptionsWithFlagHome(tempDir()),
		nil,
		zrConfig,
	)
	client.app = app

	return client, nil
}

func (c *RawTxClient) TxConfig() client.TxConfig {
	return c.app.TxConfig()
}

// Send a transaction and wait for it to be included in a block.
func (c *RawTxClient) SendWaitTx(ctx context.Context, txBytes []byte) (string, error) {
	hash, err := c.SendTx(ctx, txBytes)
	if err != nil {
		return "", err
	}

	if err = c.WaitForTx(ctx, hash); err != nil {
		return "", err
	}

	return hash, nil
}

func (c *RawTxClient) BuildTxForIdentity(ctx context.Context, identity Identity, gasLimit uint64, fees types.Coins, msgs ...types.Msg) (client.TxBuilder, xauthsigning.SignerData, error) {
	account, err := c.accountFetcher.Account(ctx, identity.Address.String())
	if err != nil {
		return nil, xauthsigning.SignerData{}, fmt.Errorf("failed to fetch account in RawTxClient.BuildTx: %w", err)
	}
	accSeq := account.GetSequence()
	accNum := account.GetAccountNumber()

	txBuilder := c.app.TxConfig().NewTxBuilder()
	signMode := c.app.TxConfig().SignModeHandler().DefaultMode()

	// build unsigned tx
	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(fees)

	if err = txBuilder.SetMsgs(msgs...); err != nil {
		return nil, xauthsigning.SignerData{}, fmt.Errorf("failed to set msgs in RawTxClient.BuildTx: %w", err)
	}

	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	sigV2 := signing.SignatureV2{
		PubKey: identity.PrivKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode(signMode),
			Signature: nil,
		},
		Sequence: accSeq,
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return nil, xauthsigning.SignerData{}, fmt.Errorf("set empty signature in RawTxClient.BuildTx: %w", err)
	}

	signerData := xauthsigning.SignerData{
		ChainID:       c.chainID,
		AccountNumber: accNum,
		Sequence:      accSeq,
		PubKey:        identity.PrivKey.PubKey(),
	}

	return txBuilder, signerData, nil
}

// BuildTx builds an unsigned transaction with the given messages.
func (c *RawTxClient) BuildTx(ctx context.Context, gasLimit uint64, fees types.Coins, msgs ...types.Msg) (client.TxBuilder, xauthsigning.SignerData, error) {
	return c.BuildTxForIdentity(ctx, c.Identity, gasLimit, fees, msgs...)
}

func (c *RawTxClient) SignTxForIdentity(ctx context.Context, identity Identity, txBuilder client.TxBuilder, signerData xauthsigning.SignerData, signMode signingv1beta1.SignMode) ([]byte, error) {
	sigV2, err := tx.SignWithPrivKey(
		ctx,
		signing.SignMode(signMode),
		signerData,
		txBuilder,
		identity.PrivKey,
		c.app.TxConfig(),
		signerData.Sequence,
	)
	if err != nil {
		return nil, fmt.Errorf("sign with priv key: %w", err)
	}

	if err = txBuilder.SetSignatures(sigV2); err != nil {
		return nil, fmt.Errorf("set signature: %w", err)
	}

	txBytes, err := c.app.TxConfig().TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, fmt.Errorf("encode tx: %w", err)
	}

	return txBytes, nil
}

// SignTx signs the transaction.
func (c *RawTxClient) SignTx(ctx context.Context, txBuilder client.TxBuilder, signerData xauthsigning.SignerData, signMode signingv1beta1.SignMode) ([]byte, error) {
	return c.SignTxForIdentity(ctx, c.Identity, txBuilder, signerData, signMode)
}

func (c *RawTxClient) BuildAndSignTx(ctx context.Context, gasLimit uint64, fees types.Coins, msgs ...types.Msg) ([]byte, error) {
	return c.BuildAndSignTxForIdentity(ctx, c.Identity, gasLimit, fees, msgs...)
}

// BuildAndSignTx builds and signs a transaction.
func (c *RawTxClient) BuildAndSignTxForIdentity(ctx context.Context, identity Identity, gasLimit uint64, fees types.Coins, msgs ...types.Msg) ([]byte, error) {
	txBuilder, signerData, err := c.BuildTxForIdentity(ctx, identity, gasLimit, fees, msgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction in BuildAndSignTx: %w", err)
	}
	txBytes, err := c.SignTxForIdentity(ctx, identity, txBuilder, signerData, c.app.TxConfig().SignModeHandler().DefaultMode())
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction in BuildAndSignTx: %w", err)
	}
	return txBytes, nil
}

// SendTx broadcasts a signed transaction and returns its hash.
// This method does not wait until the transaction is actually added to the,
// blockchain. Use SendWaitForTx for that.
func (c *RawTxClient) SendTx(ctx context.Context, txBytes []byte) (string, error) {
	grpcRes, err := c.client.BroadcastTx(
		ctx,
		&txtypes.BroadcastTxRequest{
			Mode:    txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return "", err
	}

	if grpcRes.TxResponse.Code != 0 {
		return "", fmt.Errorf("tx failed: %s", grpcRes.TxResponse.RawLog)
	}

	return grpcRes.TxResponse.TxHash, nil
}

// WaitForTx requests the tx from hash, if not found, waits for some time and
// tries again. Returns an error if ctx is canceled.
func (c *RawTxClient) WaitForTx(ctx context.Context, hash string) error {
	tick := time.NewTicker(queryTimeout)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			_, err := c.client.GetTx(ctx, &txtypes.GetTxRequest{Hash: hash})
			if err == nil {
				return nil
			}

			if !strings.Contains(err.Error(), "not found") {
				return err
			}
		}
	}
}

func (c *RawTxClient) GetTx(ctx context.Context, hash string) (*txtypes.GetTxResponse, error) {
	return c.client.GetTx(ctx, &txtypes.GetTxRequest{Hash: hash})
}

func (c *RawTxClient) GetTxConfig() client.TxConfig {
	return c.app.TxConfig()
}

var tempDir = func() string {
	dir, err := os.MkdirTemp("", "zrchain")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(dir)

	return dir
}
