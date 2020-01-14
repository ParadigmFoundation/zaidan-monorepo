package core

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
)

// HotWalletConfig stores configuration for a HotWallet
type HotWalletConfig struct {
	// Maximum request length used for the order validator
	OrderValidatorMaxReqLength int

	// Address (0x-prefixed) to use as the maker address for orders
	MakerAddress common.Address

	// Address (0x-prefixed) to use as the sender for order
	SenderAddress common.Address
}

// HotWallet represents a live hot wallet that can interact with the 0x contract system
type HotWallet struct {
	provider     *eth.Provider
	zrxHelper    *zrx.ZeroExHelper
	tokenManager *eth.ERC20TokenManager

	makerAddress  common.Address
	senderAddress common.Address

	senderTransactor *bind.TransactOpts

	logger log.Logger
}

// NewHotWallet creates a new hot wallet with the supplied provider and configuration
func NewHotWallet(provider *eth.Provider, cfg HotWalletConfig) (*HotWallet, error) {
	zrxHelper, err := zrx.NewZeroExHelper(provider.Client(), cfg.OrderValidatorMaxReqLength)
	if err != nil {
		return nil, err
	}

	if !provider.CanSignWithAddress(cfg.MakerAddress) || !provider.CanSignWithAddress(cfg.SenderAddress) {
		return nil, fmt.Errorf("unable to sign with maker or sender address")
	}

	acct, err := provider.GetAccount(cfg.MakerAddress)
	if err != nil {
		return nil, err
	}

	mgr, err := eth.NewERC20TokenManager(provider, acct.Address, nil)
	if err != nil {
		return nil, err
	}

	senderAcct, err := provider.GetAccount(cfg.SenderAddress)
	if err != nil {
		return nil, err
	}

	senderKey, err := provider.Wallet().PrivateKey(senderAcct)
	if err != nil {
		return nil, err
	}

	hw := &HotWallet{
		provider:         provider,
		zrxHelper:        zrxHelper,
		tokenManager:     mgr,
		makerAddress:     cfg.MakerAddress,
		senderAddress:    cfg.SenderAddress,
		senderTransactor: bind.NewKeyedTransactor(senderKey),
		logger:           log.New(context.Background()),
	}

	return hw, nil
}
