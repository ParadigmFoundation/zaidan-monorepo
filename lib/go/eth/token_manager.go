package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20TokenManager struct {
	provider *Provider
	opts     *bind.TransactOpts

	erc20s map[common.Address]*ERC20TokenSession
}

// NewERC20TokenManager creates a new manager with provider, where account is the signer. Adds tokens if provided.
func NewERC20TokenManager(provider *Provider, account accounts.Account, tokens []common.Address) (*ERC20TokenManager, error) {
	if err := provider.hasAccountOrErr(account); err != nil {
		return nil, err
	}

	key, err := provider.hw.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	mgr := &ERC20TokenManager{
		provider: provider,
		opts:     bind.NewKeyedTransactor(key),
		erc20s:   make(map[common.Address]*ERC20TokenSession),
	}

	if len(tokens) != 0 {
		for _, token := range tokens {
			if err := mgr.addToken(token); err != nil {
				return nil, err
			}
		}
	}

	return mgr, nil
}

// BalanceOf calls BalanceOf on token for owner
func (tm *ERC20TokenManager) BalanceOf(token common.Address, owner common.Address) (*big.Int, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return nil, err
	}

	bal, err := session.BalanceOf(owner)
	if err != nil {
		return nil, err
	}
	return bal, nil
}

// Approve calls Approve on token for owner, approving spender to spend value
func (tm *ERC20TokenManager) Approve(ctx context.Context, token common.Address, spender common.Address, value *big.Int) (*types.Transaction, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return nil, err
	}

	return session.Approve(spender, value)
}

// adds a new tracked ERC-20 token session
func (tm *ERC20TokenManager) addToken(address common.Address) error {
	token, err := NewERC20Token(address, tm.provider.eth)
	if err != nil {
		return err
	}

	session := &ERC20TokenSession{
		Contract: token,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:   tm.opts.From,
			Signer: tm.opts.Signer,

			// todo(@hrharder): reconsider where this value should come from
			GasLimit: 1000000,
		},
	}

	tm.erc20s[address] = session
	return nil
}

// gets a token session for a given address, adds it if it doesn't exist
func (tm *ERC20TokenManager) tokenSession(token common.Address) (*ERC20TokenSession, error) {
	session, has := tm.erc20s[token]
	if has {
		return session, nil
	}

	if err := tm.addToken(token); err != nil {
		return nil, err
	}
	return tm.erc20s[token], nil
}
