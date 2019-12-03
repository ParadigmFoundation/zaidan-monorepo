package eth

import (
	"testing"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/ethereum/go-ethereum/common"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/require"
)

func TestTokenManager(t *testing.T) {
	cfg := TestConfig{}
	if err := env.Parse(&cfg); err != nil {
		t.Fatal(err)
	}

	// setup provider test instance
	provider, err := NewProvider(cfg.Ethurl, cfg.Mnemonic, hdwallet.DefaultBaseDerivationPath)
	require.NoError(t, err, "should have no error creating new provider")

	acct := provider.Accounts()[0]
	manager, err := NewERC20TokenManager(provider, acct, nil)
	require.NoError(t, err, "should not have error creating token manager")

	// the ZRX test token
	zrxTokenAddress := common.HexToAddress("0x871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c")

	bal, err := manager.BalanceOf(zrxTokenAddress, acct.Address)
	require.NoError(t, err, "should not have error getting token balance")
	t.Fatal(bal.String())
}
