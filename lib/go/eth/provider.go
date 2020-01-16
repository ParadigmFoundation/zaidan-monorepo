package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/0xProject/0x-mesh/ethereum/signer"

	"golang.org/x/crypto/sha3"

	"github.com/btcsuite/btcd/btcec"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

// NewProviderWithDefaultPath creates a new provider with the default base derivation path
func NewProviderWithDefaultPath(ethurl string, mnemonic string) (*Provider, error) {
	return NewProvider(ethurl, mnemonic, accounts.DefaultBaseDerivationPath)
}

// Client returns a pointer to the underlying Ethereum client
func (pvr *Provider) Client() *ethclient.Client { return pvr.eth }

// Wallet returns a pointer to the underlying wallet
func (pvr *Provider) Wallet() *hdwallet.Wallet { return pvr.hw }

// AccountAt returns the account at the specified path (if applicable). Use Derive to add an address.
func (pvr *Provider) AccountAt(path accounts.DerivationPath) (accounts.Account, error) {
	return pvr.hw.Derive(path, false)
}

// SignText signs a personal message with account if available.
func (pvr *Provider) SignText(account accounts.Account, text []byte) ([]byte, error) {
	if _, err := pvr.GetAccount(account.Address); err != nil {
		return nil, err
	}

	return pvr.hw.SignText(account, text)
}

// SignData signs arbitrary typed data of type mimeType with account.
func (pvr *Provider) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	if _, err := pvr.GetAccount(account.Address); err != nil {
		return nil, err
	}

	return pvr.hw.SignData(account, mimeType, data)
}

// SignTx signs a transaction with the specified account.
func (pvr *Provider) SignTx(ctx context.Context, account accounts.Account, tx *types.Transaction) (*types.Transaction, error) {
	if _, err := pvr.GetAccount(account.Address); err != nil {
		return nil, err
	}

	id, err := pvr.NetworkID(ctx)
	if err != nil {
		return tx, err
	}

	return pvr.hw.SignTx(account, tx, id)
}

// Primarily taken from: https://github.com/0xProject/0x-mesh/blob/bd3060d3efaab759913c4de2152c2ef4e5940301/ethereum/signer/sign.go
// EthSign implements zeroex.Signer
func (pvr *Provider) EthSign(message []byte, signerAddress common.Address) (*signer.ECSignature, error) {
	acct, err := pvr.GetAccount(signerAddress)
	if err != nil {
		return nil, err
	}

	if acct.Address != signerAddress {
		return nil, fmt.Errorf("invalid signer: unsupported account")
	}

	// Add message prefix: "\x19Ethereum Signed Message:\n"${message length}
	messageWithPrefix, _ := textAndHash(message)

	privateKey, err := pvr.hw.PrivateKey(acct)
	if err != nil {
		return nil, err
	}

	signatureBytes, err := sign(messageWithPrefix, privateKey)
	if err != nil {
		return nil, err
	}

	vParam := signatureBytes[64]
	if vParam == byte(0) {
		vParam = byte(27)
	} else if vParam == byte(1) {
		vParam = byte(28)
	}

	ecSignature := &signer.ECSignature{
		V: vParam,
		R: common.BytesToHash(signatureBytes[0:32]),
		S: common.BytesToHash(signatureBytes[32:64]),
	}

	return ecSignature, nil
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

// returns the current nonce for account
func (pvr *Provider) Nonce(ctx context.Context, addr common.Address) (uint64, error) {
	return pvr.eth.NonceAt(ctx, addr, nil)
}

// NetworkID gets the current network ID (tries network ID, if that fails, gets chain ID)
func (pvr *Provider) NetworkID(ctx context.Context) (id *big.Int, err error) {
	id, err = pvr.eth.NetworkID(ctx)
	if err != nil {
		return pvr.eth.ChainID(ctx)
	}
	return id, nil
}

// NetworkID gets the current chain ID (tries chain ID, if that fails, gets network ID)
func (pvr *Provider) ChainID(ctx context.Context) (id *big.Int, err error) {
	id, err = pvr.eth.ChainID(ctx)
	if err != nil {
		return pvr.eth.NetworkID(ctx)
	}
	return id, nil
}

// CanSignWithAddress returns true if the provider can sign with the given address
func (pvr *Provider) CanSignWithAddress(addr common.Address) bool {
	if _, err := pvr.GetAccount(addr); err != nil {
		return false
	}
	return true
}

// GetAccount fetches the account object for a given address, errors if not present
func (pvr *Provider) GetAccount(addr common.Address) (accounts.Account, error) {
	for _, acct := range pvr.Accounts() {
		if acct.Address == addr {
			return acct, nil
		}
	}
	return accounts.Account{}, fmt.Errorf("unsupported address; no pinned account for: %s", addr.Hex())
}

// FROM: https://github.com/ethereum/go-ethereum/blob/master/crypto/signature_nocgo.go
//
// Original comment:
// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given hash cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func sign(hash []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	if prv.Curve != btcec.S256() {
		return nil, fmt.Errorf("private key curve is not secp256k1")
	}
	sig, err := btcec.SignCompact(btcec.S256(), (*btcec.PrivateKey)(prv), hash, false)
	if err != nil {
		return nil, err
	}
	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, nil
}

// FROM: https://github.com/0xProject/0x-mesh/blob/bd3060d3efaab759913c4de2152c2ef4e5940301/ethereum/signer/sign.go#L180-L194
//
// Original comment:
// textAndHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calulcated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func textAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	// Note: Write will never return an error here. We added placeholders in order
	// to satisfy the linter.
	_, _ = hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}
