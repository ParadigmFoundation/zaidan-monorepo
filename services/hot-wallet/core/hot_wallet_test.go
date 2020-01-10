package core

import (
	"math"
	"testing"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Ethurl   string `env:"ETHEREUM_JSONRPC_URL" envDefault:"http://localhost:8545"`
	Mnemonic string `env:"ETHEREUM_MNEMONIC" envDefault:"concert load couple harbor equip island argue ramp clarify fence smart topic"`
}

func TestHotWallet(t *testing.T) {
	var cfg TestConfig
	require.NoError(t, env.Parse(&cfg))

	deriver := eth.NewBaseDeriver()
	provider, err := eth.NewProvider(cfg.Ethurl, cfg.Mnemonic, deriver.Base())
	require.NoError(t, err)

	assert.Equal(t, 1, len(provider.Accounts()))
	require.NoError(t, provider.Derive(deriver.DeriveNext()))
	assert.Equal(t, 2, len(provider.Accounts()))

	hwCfg := HotWalletConfig{
		OrderValidatorMaxReqLength: math.MaxInt16,
		MakerAddress:               provider.Accounts()[0].Address,
		SenderAddress:              provider.Accounts()[1].Address,
	}

	hw, err := NewHotWallet(provider, hwCfg)
	assert.NoError(t, err)

	t.Run("test order creation", func(t *testing.T) { testCreateOrder(hw, t) })
	t.Run("test ether transfer", func(t *testing.T) { testTransferEther(hw, t) })
	t.Run("test send transaction", func(t *testing.T) { testSendTransaction(hw, t) })
	t.Run("test get/set allowance", func(t *testing.T) { testGetSetAllowance(hw, t) })
}
