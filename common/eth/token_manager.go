package eth

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20TokenManager struct {
	provider *Provider
	opts     *bind.TransactOpts

	erc20s map[common.Address]*ERC20TokenSession
	mu     *sync.Mutex
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

	if tokens != nil {
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

	return session.BalanceOf(owner)
}

// Approve calls Approve on token for owner, approving spender to spend value
func (tm *ERC20TokenManager) Approve(token common.Address, spender common.Address, value *big.Int) (*types.Transaction, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return nil, err
	}

	stx, err := session.Approve(spender, value)
	if err := tm.provider.SendTx(stx); err != nil {
		return stx, err
	}

	return stx, err
}

// adds a new tracked ERC-20 token session
func (tm *ERC20TokenManager) addToken(address common.Address) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

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
			From:     tm.opts.From,
			Signer:   tm.opts.Signer,
			GasLimit: uint64(500000),
		},
	}

	tm.erc20s[address] = session
	return nil
}

// gets a token session for a given address, adds it if it doesn't exist
func (tm *ERC20TokenManager) tokenSession(token common.Address) (*ERC20TokenSession, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	session, has := tm.erc20s[token]
	if has {
		return session, nil
	}

	if err := tm.addToken(token); err != nil {
		return nil, err
	}
	return tm.erc20s[token], nil
}
