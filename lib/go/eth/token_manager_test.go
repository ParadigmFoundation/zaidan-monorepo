package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/ethereum/go-ethereum/common"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/require"
)

const ZERO_EX_TEST_TOKEN_ADDRESS_STR = "0x871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"

var ZERO_EX_TEST_TOKEN_ADDRESS = common.HexToAddress(ZERO_EX_TEST_TOKEN_ADDRESS_STR)

func TestTokenManager(t *testing.T) {
	cfg := TestConfig{}
	if err := env.Parse(&cfg); err != nil {
		t.Fatal(err)
	}

	// setup provider test instance
	provider, err := NewProvider(cfg.Ethurl, cfg.Mnemonic, hdwallet.DefaultBaseDerivationPath)
	require.NoError(t, err, "should have no error creating new provider")

	// test token instance
	tkn, err := NewERC20Token(ZERO_EX_TEST_TOKEN_ADDRESS, provider.eth)
	require.NoError(t, err, "should be no error creating test token instance")

	acct := provider.Accounts()[0]
	manager, err := NewERC20TokenManager(provider, acct.Address, nil)
	require.NoError(t, err, "should not have error creating token manager")

	bal, err := manager.BalanceOf(ZERO_EX_TEST_TOKEN_ADDRESS, acct.Address)
	require.NoError(t, err, "should not have error getting token balance from manager")
	expectedBal, err := tkn.BalanceOf(nil, acct.Address)
	require.NoError(t, err, "should not have error getting token balance")
	assert.Equal(t, expectedBal, bal, "should have a non-0 balance")
}
