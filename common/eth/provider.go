package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// Provider enables interaction with the Ethereum blockchain through a mnemonic hot-wallet and an ETH client.
type Provider struct {
	hw  *hdwallet.Wallet
	eth *ethclient.Client
}

// NewProvider creates a new signing provider with a mnemonic and derivation path.
func NewProvider(ethurl string, mnemonic string, path accounts.DerivationPath) (*Provider, error) {
	client, err := ethclient.Dial(ethurl)
	if err != nil {
		return nil, err
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	if path == nil {
		path = accounts.DefaultBaseDerivationPath
	}

	if _, err := wallet.Derive(path, true); err != nil {
		return nil, err
	}

	return &Provider{hw: wallet, eth: client}, nil
}

// SendTx sends a signed transaction. It makes no modification (gas price, etc.).
func (pvr *Provider) SendTx(tx *types.Transaction) error {
	return pvr.eth.SendTransaction(context.Background(), tx)
}

// SignText signs a personal message with account if available.
func (pvr *Provider) SignText(account accounts.Account, text []byte) ([]byte, error) {
	if err := pvr.hasAccountOrErr(account); err != nil {
		return nil, err
	}

	return pvr.hw.SignText(account, text)
}

// SignData signs arbitrary typed data of type mimeType with account.
func (pvr *Provider) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	if err := pvr.hasAccountOrErr(account); err != nil {
		return nil, err
	}

	return pvr.hw.SignData(account, mimeType, data)
}

// SignTx signs a transaction with the specified account.
func (pvr *Provider) SignTx(account accounts.Account, tx *types.Transaction) (*types.Transaction, error) {
	if err := pvr.hasAccountOrErr(account); err != nil {
		return nil, err
	}

	id, err := pvr.networkID()
	if err != nil {
		return tx, err
	}

	return pvr.hw.SignTx(account, tx, id)
}

// Accounts gets currently supported accounts.
func (pvr *Provider) Accounts() []accounts.Account { return pvr.hw.Accounts() }

// Derive derives a new account based on path and adds it to the hot wallet.
func (pvr *Provider) Derive(path accounts.DerivationPath) error {
	if _, err := pvr.hw.Derive(path, true); err != nil {
		return err
	}
	return nil
}

// returns false if account not supported by provider
func (pvr *Provider) hasAccount(acct accounts.Account) bool {
	for _, account := range pvr.Accounts() {
		if account.Address == acct.Address {
			return true
		}
	}
	return false
}

// returns error if account doesn't exist
func (pvr *Provider) hasAccountOrErr(acct accounts.Account) error {
	if !pvr.hasAccount(acct) {
		return fmt.Errorf("unsupported account: %s", acct.Address)
	}
	return nil
}

// returns the current network ID (tries network ID, if that fails, gets chain ID)
func (pvr *Provider) networkID() (id *big.Int, err error) {
	id, err = pvr.eth.NetworkID(context.TODO())
	if err != nil {
		return pvr.eth.ChainID(context.TODO())
	}
	return id, nil
}

// returns the current nonce for account
func (pvr *Provider) nonce(acct accounts.Account) (uint64, error) {
	return pvr.eth.NonceAt(context.TODO(), acct.Address, nil)
}

// returns error if TX nonce does not match what it should
func (pvr *Provider) ensureNonce(acct accounts.Account, tx *types.Transaction) error {
	nonce, err := pvr.nonce(acct)
	if err != nil {
		return err
	}

	if nonce != tx.Nonce() {
		return fmt.Errorf("invalid nonce: expected (%d) got (%d)", nonce, tx.Nonce())
	}

	return nil
}
