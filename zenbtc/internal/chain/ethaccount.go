package chain

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type EthAccount struct {
	privkey *ecdsa.PrivateKey
	address common.Address
}

func NewEthAccount(mnemonic string, derivationPath string) (*EthAccount, error) {
	if mnemonic == "" || derivationPath == "" {
		return nil, errors.New("mnemonic and/or derivation path not set in config")
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	path := hdwallet.MustParseDerivationPath(derivationPath)

	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	privkey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	return &EthAccount{
		privkey: privkey,
		address: account.Address,
	}, nil
}

func (e *EthAccount) GetPrivKey() *ecdsa.PrivateKey {
	return e.privkey
}

func (e *EthAccount) GetAddress() common.Address {
	return e.address
}
