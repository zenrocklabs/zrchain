package solzenbtc

// import (
// 	"context"
// 	"math/big"
// 	"testing"

// 	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc/generated/zenbtc_spl_token"
// 	"github.com/gagliardetto/solana-go"
// 	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
// 	"github.com/gagliardetto/solana-go/programs/system"
// 	token "github.com/gagliardetto/solana-go/programs/token"
// 	"github.com/gagliardetto/solana-go/rpc"

// 	"github.com/stretchr/testify/require"
// )

// var endpoint = rpc.LocalNet.RPC

// // var authority, _ = LoadLocalWallet()
// var authority = solana.NewWallet().PrivateKey
// var programID = solana.MustPublicKeyFromBase58("3jo4mdc6QbGRigia2jvmKShbmz3aWq4Y8bgUXfur5StT")

// var client = rpc.New(endpoint)

// var tokenParams = Token{
// 	Name:     "Zenrock BTC",
// 	Symbol:   "zenBTC",
// 	Decimals: 8,
// 	Uri:      "https://www.zenrocklabs.io/metadata.json",
// }

// var mintAddress, _ = GetMintAddress(programID) // You can also just use the token address

// var userWallet = solana.NewWallet()
// var feeWallet = authority.PublicKey()
// var multisigKey solana.PublicKey

// func TestInitialize(t *testing.T) {
// 	signer := authority

// 	err := RequestAirdrop(context.Background(), client, userWallet.PublicKey(), 1000000000)
// 	require.NoError(t, err)

// 	err = RequestAirdrop(context.Background(), client, signer.PublicKey(), 1000000000)
// 	require.NoError(t, err)

// 	recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
// 	require.NoError(t, err)

// 	tx, err := solana.NewTransaction(
// 		[]solana.Instruction{
// 			Initialize(
// 				programID,
// 				zenbtc_spl_token.InitializeArgs{
// 					GlobalAuthority: authority.PublicKey(),
// 					MintAuthorities: []solana.PublicKey{authority.PublicKey()},
// 					FeeAuthorities:  []solana.PublicKey{authority.PublicKey()},
// 					FeeWallet:       authority.PublicKey(),
// 					BurnFeeBps:      2,
// 					TokenName:       tokenParams.Name,
// 					TokenSymbol:     tokenParams.Symbol,
// 					TokenDecimals:   tokenParams.Decimals,
// 					TokenUri:        tokenParams.Uri,
// 				},
// 				signer.PublicKey(),
// 				mintAddress,
// 			),
// 		},
// 		recent.Value.Blockhash,
// 		solana.TransactionPayer(signer.PublicKey()),
// 	)
// 	require.NoError(t, err)

// 	_, err = SignTransaction(tx, signer)
// 	require.NoError(t, err)

// 	signature, err := SendTransaction(client, context.Background(), tx)
// 	require.NoError(t, err)

// 	confirmedTx, err := WaitForTransactionConfirmation(
// 		context.Background(),
// 		client,
// 		signature,
// 		rpc.CommitmentConfirmed,
// 	)
// 	require.NoError(t, err)

// 	require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error
// }

// func TestTransferMintAuthority(t *testing.T) {
// 	signer := authority

// 	globalConfigPDA, err := GetGlobalConfigPDA(programID)
// 	require.NoError(t, err)

// 	t.Run("Creates a multisig account", func(t *testing.T) {
// 		multisigKeyPair := solana.NewWallet()
// 		multisigKey = multisigKeyPair.PublicKey()

// 		recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
// 		require.NoError(t, err)

// 		instructions := []solana.Instruction{
// 			system.NewCreateAccountInstruction(
// 				uint64(0.004*float64(solana.LAMPORTS_PER_SOL)),
// 				MULTISIG_SIZE,
// 				solana.TokenProgramID,
// 				signer.PublicKey(),
// 				multisigKey,
// 			).Build(),
// 			NewInitializeMultisigInstruction(
// 				1,
// 				multisigKey,
// 				solana.SysVarRentPubkey,
// 				[]solana.PublicKey{globalConfigPDA},
// 			).Build(),
// 		}

// 		tx, err := solana.NewTransaction(
// 			instructions,
// 			recent.Value.Blockhash,
// 			solana.TransactionPayer(signer.PublicKey()),
// 		)
// 		require.NoError(t, err)

// 		_, err = SignTransaction(tx, signer)
// 		require.NoError(t, err)

// 		_, err = SignTransaction(tx, multisigKeyPair.PrivateKey)
// 		require.NoError(t, err)

// 		signature, err := SendTransaction(client, context.Background(), tx)
// 		require.NoError(t, err)

// 		confirmedTx, err := WaitForTransactionConfirmation(
// 			context.Background(),
// 			client,
// 			signature,
// 			rpc.CommitmentConfirmed,
// 		)
// 		require.NoError(t, err)
// 		require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error
// 	})

// 	t.Run("Transfers the mint authority to the multisig", func(t *testing.T) {
// 		recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
// 		require.NoError(t, err)

// 		instructions := []solana.Instruction{
// 			token.NewSetAuthorityInstruction(
// 				token.AuthorityMintTokens,
// 				multisigKey,
// 				mintAddress,
// 				signer.PublicKey(),
// 				[]solana.PublicKey{},
// 			).Build(),
// 		}

// 		tx, err := solana.NewTransaction(
// 			instructions,
// 			recent.Value.Blockhash,
// 			solana.TransactionPayer(signer.PublicKey()),
// 		)
// 		require.NoError(t, err)

// 		_, err = SignTransaction(tx, signer)
// 		require.NoError(t, err)

// 		signature, err := SendTransaction(client, context.Background(), tx)
// 		require.NoError(t, err)

// 		confirmedTx, err := WaitForTransactionConfirmation(
// 			context.Background(),
// 			client,
// 			signature,
// 			rpc.CommitmentConfirmed,
// 		)
// 		require.NoError(t, err)
// 		require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error

// 	})
// }

// func TestWrap(t *testing.T) {
// 	signer := authority

// 	receiver := userWallet.PublicKey()
// 	value := uint64(10000)
// 	fee := uint64(20)

// 	recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
// 	require.NoError(t, err)

// 	instructions := []solana.Instruction{}

// 	feeWalletAta, _, err := solana.FindAssociatedTokenAddress(feeWallet, mintAddress)
// 	require.NoError(t, err)

// 	_, err = GetTokenAccount(context.Background(), client, feeWalletAta)

// 	if err != nil && err.Error() == "not found" {
// 		instructions = append(
// 			instructions,
// 			ata.NewCreateInstruction(
// 				signer.PublicKey(),
// 				feeWallet,
// 				mintAddress,
// 			).Build(),
// 		)
// 	} else {
// 		require.NoError(t, err)
// 	}

// 	receiverAta, _, err := solana.FindAssociatedTokenAddress(receiver, mintAddress)
// 	require.NoError(t, err)

// 	_, err = GetTokenAccount(context.Background(), client, receiverAta)

// 	if err != nil && err.Error() == "not found" {
// 		instructions = append(
// 			instructions,
// 			ata.NewCreateInstruction(
// 				signer.PublicKey(),
// 				receiver,
// 				mintAddress,
// 			).Build(),
// 		)
// 	} else {
// 		require.NoError(t, err)
// 	}

// 	instructions = append(instructions, Wrap(
// 		programID,
// 		zenbtc_spl_token.WrapArgs{
// 			Value: value,
// 			Fee:   fee,
// 		},
// 		signer.PublicKey(),
// 		mintAddress,
// 		multisigKey,
// 		feeWallet,
// 		feeWalletAta,
// 		receiver,
// 		receiverAta,
// 	))

// 	tx, err := solana.NewTransaction(
// 		instructions,
// 		recent.Value.Blockhash,
// 		solana.TransactionPayer(signer.PublicKey()),
// 	)
// 	require.NoError(t, err)

// 	_, err = SignTransaction(tx, signer)
// 	require.NoError(t, err)

// 	// // Pretty print the transaction:
// 	// fmt.Println(tx.String())

// 	signature, err := SendTransaction(client, context.Background(), tx)
// 	require.NoError(t, err)

// 	confirmedTx, err := WaitForTransactionConfirmation(
// 		context.Background(),
// 		client,
// 		signature,
// 		rpc.CommitmentConfirmed,
// 	)
// 	require.NoError(t, err)
// 	require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error
// }

// func TestUnwrap(t *testing.T) {
// 	signer := userWallet.PrivateKey

// 	value := uint64(5000)
// 	destAddr := [25]uint8{
// 		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
// 		22, 23, 24, 25,
// 	}

// 	recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
// 	require.NoError(t, err)

// 	tx, err := solana.NewTransaction(
// 		[]solana.Instruction{
// 			Unwrap(
// 				programID,
// 				zenbtc_spl_token.UnwrapArgs{
// 					Value:    value,
// 					DestAddr: destAddr,
// 				},
// 				signer.PublicKey(),
// 				mintAddress,
// 				multisigKey,
// 				feeWallet,
// 			),
// 		},
// 		recent.Value.Blockhash,
// 		solana.TransactionPayer(signer.PublicKey()),
// 	)
// 	require.NoError(t, err)

// 	_, err = SignTransaction(tx, signer)
// 	require.NoError(t, err)

// 	signature, err := SendTransaction(client, context.Background(), tx)
// 	require.NoError(t, err)

// 	confirmedTx, err := WaitForTransactionConfirmation(
// 		context.Background(),
// 		client,
// 		signature,
// 		rpc.CommitmentConfirmed,
// 	)
// 	require.NoError(t, err)
// 	require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error
// }

// func TestGetTokenBalance(t *testing.T) {
// 	userAta, _, err := solana.FindAssociatedTokenAddress(userWallet.PublicKey(), mintAddress)
// 	require.NoError(t, err)

// 	accountAtaInfo, err := client.GetTokenAccountBalance(context.Background(), userAta, rpc.CommitmentConfirmed)
// 	require.NoError(t, err)

// 	balance := accountAtaInfo.Value.Amount
// 	require.NotEmpty(t, balance)
// }

// func TestGetMint(t *testing.T) {
// 	mint, err := GetMint(context.Background(), client, mintAddress)
// 	require.NoError(t, err)

// 	require.Equal(t, *mint.MintAuthority, multisigKey)
// 	require.Greater(t, mint.Decimals, uint8(0))
// }

// func TestGetGlobalConfig(t *testing.T) {
// 	globalConfig, err := GetGlobalConfig(context.Background(), client, programID)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, globalConfig.GlobalAuthority.String())
// }

// func TestGetTokenRedemptionEvents(t *testing.T) {
// 	limit := 1000

// 	signatures, err := client.GetSignaturesForAddressWithOpts(context.Background(), programID, &rpc.GetSignaturesForAddressOpts{
// 		Limit:      &limit,
// 		Commitment: rpc.CommitmentConfirmed,
// 	})
// 	require.NoError(t, err)
// 	require.NotEmpty(t, signatures)

// 	tokenRedemptionEvents := []TokenRedemptionEvent{}

// 	for _, signature := range signatures {
// 		tx, err := client.GetTransaction(context.Background(), signature.Signature, &rpc.GetTransactionOpts{
// 			Commitment: rpc.CommitmentConfirmed,
// 		})
// 		require.NoError(t, err)
// 		require.NotEmpty(t, tx)

// 		events, err := zenbtc_spl_token.DecodeEvents(tx, programID)
// 		require.NoError(t, err)

// 		for _, event := range events {
// 			if event.Name == "TokenRedemption" {
// 				eventData := event.Data.(*zenbtc_spl_token.TokenRedemptionEventData)

// 				tokenRedemptionEvents = append(tokenRedemptionEvents, TokenRedemptionEvent{
// 					Signature: signature.Signature.String(),
// 					Slot:      tx.Slot,
// 					Date:      tx.BlockTime.Time(),
// 					Redeemer:  eventData.Redeemer,
// 					Value:     eventData.Value,
// 					DestAddr:  eventData.DestAddr,
// 					Fee:       eventData.Fee,
// 					Mint:      eventData.Mint,
// 					Id:        eventData.Id.BigInt(),
// 				})
// 			}
// 		}
// 	}

// 	require.NotEmpty(t, tokenRedemptionEvents)
// 	require.Len(t, tokenRedemptionEvents, 1)
// 	require.Equal(t, tokenRedemptionEvents[0].Id.Cmp(big.NewInt(0)), 0)
// }

// func TestDurableNonces(t *testing.T) {
// 	nonceAuth := solana.NewWallet()
// 	nonceWallet := solana.NewWallet()
// 	sender := solana.NewWallet()
// 	receiver := solana.NewWallet().PublicKey()

// 	err := RequestAirdrop(context.Background(), client, nonceAuth.PublicKey(), 1000*solana.LAMPORTS_PER_SOL)
// 	require.NoError(t, err)

// 	err = RequestAirdrop(context.Background(), client, sender.PublicKey(), 1000*solana.LAMPORTS_PER_SOL)
// 	require.NoError(t, err)

// 	t.Run("Creates a nonce account", func(t *testing.T) {
// 		recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
// 		require.NoError(t, err)

// 		tx, err := solana.NewTransaction(
// 			[]solana.Instruction{
// 				system.NewCreateAccountInstruction(
// 					uint64(0.0015*float64(solana.LAMPORTS_PER_SOL)),
// 					NONCE_ACCOUNT_LENGTH,
// 					solana.SystemProgramID,
// 					nonceAuth.PublicKey(),
// 					nonceWallet.PublicKey(),
// 				).Build(),

// 				system.NewInitializeNonceAccountInstruction(
// 					nonceAuth.PublicKey(),
// 					nonceWallet.PublicKey(),
// 					solana.SysVarRecentBlockHashesPubkey,
// 					solana.SysVarRentPubkey,
// 				).Build(),
// 			},
// 			recent.Value.Blockhash,
// 			solana.TransactionPayer(nonceAuth.PublicKey()),
// 		)
// 		require.NoError(t, err)

// 		_, err = SignTransaction(tx, nonceAuth.PrivateKey)
// 		require.NoError(t, err)

// 		_, err = SignTransaction(tx, nonceWallet.PrivateKey)
// 		require.NoError(t, err)

// 		signature, err := SendTransaction(client, context.Background(), tx)
// 		require.NoError(t, err)

// 		confirmedTx, err := WaitForTransactionConfirmation(
// 			context.Background(),
// 			client,
// 			signature,
// 			rpc.CommitmentConfirmed,
// 		)
// 		require.NoError(t, err)

// 		require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error

// 		nonceAccount, err := GetNonceAccount(context.Background(), client, nonceWallet.PublicKey())
// 		require.NoError(t, err)

// 		require.Equal(t, nonceAccount.AuthorizedPubkey, nonceAuth.PublicKey())
// 	})

// 	t.Run("Performs a transfer using a durable nonce", func(t *testing.T) {
// 		noncePubkey := nonceWallet.PublicKey()

// 		nonceAccountBefore, err := GetNonceAccount(context.Background(), client, noncePubkey)
// 		require.NoError(t, err)

// 		tx, err := solana.NewTransaction(
// 			[]solana.Instruction{
// 				system.NewAdvanceNonceAccountInstruction(
// 					noncePubkey,
// 					solana.SysVarRecentBlockHashesPubkey,
// 					nonceAuth.PublicKey(),
// 				).Build(),
// 				system.NewTransferInstruction(
// 					uint64(0.01*float64(solana.LAMPORTS_PER_SOL)),
// 					sender.PublicKey(),
// 					receiver,
// 				).Build(),
// 			},
// 			solana.Hash(nonceAccountBefore.Nonce),
// 			solana.TransactionPayer(sender.PublicKey()),
// 		)
// 		require.NoError(t, err)

// 		_, err = SignTransaction(tx, nonceAuth.PrivateKey)
// 		require.NoError(t, err)

// 		_, err = SignTransaction(tx, sender.PrivateKey)
// 		require.NoError(t, err)

// 		signature, err := SendTransaction(client, context.Background(), tx)
// 		require.NoError(t, err)

// 		confirmedTx, err := WaitForTransactionConfirmation(
// 			context.Background(),
// 			client,
// 			signature,
// 			rpc.CommitmentConfirmed,
// 		)
// 		require.NoError(t, err)

// 		require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error

// 		nonceAccountAfter, err := GetNonceAccount(context.Background(), client, noncePubkey)
// 		require.NoError(t, err)

// 		require.NotEqual(t, nonceAccountBefore.Nonce, nonceAccountAfter.Nonce)
// 	})
// }
