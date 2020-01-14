package eth

import (
	"math/big"

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
func NewERC20TokenManager(provider *Provider, addr common.Address, tokens []common.Address) (*ERC20TokenManager, error) {
	account, err := provider.GetAccount(addr)
	if err != nil {
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

// Name calls Name on token contract
func (tm *ERC20TokenManager) Name(token common.Address) (string, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return "", err
	}

	return session.Name()
}

// Symbol calls Symbol on token contract
func (tm *ERC20TokenManager) Symbol(token common.Address) (string, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return "", err
	}

	return session.Symbol()
}

// Decimals calls Decimals on token contract
func (tm *ERC20TokenManager) Decimals(token common.Address) (uint8, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return 0, err
	}

	return session.Decimals()
}

// TotalSupply calls TotalSupply on token contract
func (tm *ERC20TokenManager) TotalSupply(token common.Address) (*big.Int, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return nil, err
	}

	return session.TotalSupply()
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

// Transfer calls Transfer on token, transfering value to the to address
func (tm *ERC20TokenManager) Transfer(token common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return nil, err
	}

	return session.Transfer(to, value)
}

// TransferFrom calls TransferFrom on token, transfering value to the to address, from the from address (allowance must be set)
func (tm *ERC20TokenManager) TransferFrom(token common.Address, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return nil, err
	}

	return session.TransferFrom(from, to, value)
}

// Approve calls Approve on token contract, and allows spender to withdraw from your account multiple times, up to the value amount
func (tm *ERC20TokenManager) Approve(token common.Address, spender common.Address, value *big.Int) (*types.Transaction, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return nil, err
	}

	return session.Approve(spender, value)
}

// Allowance calls Allowance on token contract, and returns the amount which spender is still allowed to withdraw from owner
func (tm *ERC20TokenManager) Allowance(token common.Address, owner common.Address, spender common.Address) (*big.Int, error) {
	session, err := tm.tokenSession(token)
	if err != nil {
		return nil, err
	}

	return session.Allowance(owner, spender)
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

			// setting gas limit to 0 will enable auto gas limit estimation
			GasLimit: 0,
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
