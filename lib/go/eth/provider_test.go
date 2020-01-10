package eth

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/caarlos0/env/v6"
)

var (
	TestAccountPath_A = accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 7}
	TestAccountPath_B = accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 8}
	TestAccountPath_C = accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 9}
)

type TestConfig struct {
	Ethurl   string `env:"ETHEREUM_JSONRPC_URL" envDefault:"http://localhost:8545"`
	Mnemonic string `env:"ETHEREUM_MNEMONIC" envDefault:"concert load couple harbor equip island argue ramp clarify fence smart topic"`
}

func TestProvider(t *testing.T) {
	cfg := TestConfig{}
	if err := env.Parse(&cfg); err != nil {
		t.Fatal(err)
	}

	// setup provider test instance
	provider, err := NewProvider(cfg.Ethurl, cfg.Mnemonic, TestAccountPath_A)
	require.NoError(t, err, "should have no error creating new provider")

	// run test runners
	t.Run("derive accounts", func(t *testing.T) { testDeriveAccounts(provider, t) })
	t.Run("transfer ether", func(t *testing.T) { testTransferEther(provider, t) })
}

func testDeriveAccounts(provider *Provider, t *testing.T) {
	assert.Equal(t, 1, len(provider.Accounts()), "should have one account available after construction")
	require.NoError(t, provider.Derive(TestAccountPath_B))
	assert.Equal(t, 2, len(provider.Accounts()), "should have two accounts available after deriving another")
}

func testTransferEther(provider *Provider, t *testing.T) {
	fromAccount := provider.Accounts()[0]
	toAccount := provider.Accounts()[1]

	// fetch before balances
	fromAccountBeforeBalance, err := provider.eth.BalanceAt(context.Background(), fromAccount.Address, nil)
	require.NoError(t, err, "should be no error getting balance")
	toAccountBeforeBalance, err := provider.eth.BalanceAt(context.Background(), toAccount.Address, nil)
	require.NoError(t, err, "should be no error getting balance")

	// prep transaction
	transferAmt, ok := new(big.Int).SetString("10000000000000000000", 10) // 1e19 wei (10 ether)
	assert.True(t, ok, "should be able to set transfer amount as 10 ether")
	nonce, err := provider.Nonce(context.Background(), fromAccount.Address)
	require.NoError(t, err, "should be no error getting nonce")
	expectedNonce, err := provider.eth.NonceAt(context.Background(), fromAccount.Address, nil)
	require.NoError(t, err, "should be no error getting nonce")
	assert.Equal(t, expectedNonce, nonce, "nonce should be 0 before any transactions are submitter")
	tx := types.NewTransaction(nonce, toAccount.Address, transferAmt, 21000, big.NewInt(1), nil)

	// sign transaction
	stx, err := provider.SignTx(context.Background(), fromAccount, tx)
	require.NoError(t, err, "should be no error signing transaction")

	// send transaction
	err = provider.eth.SendTransaction(context.Background(), stx)
	require.NoError(t, err, "should be no error sending transaction")

	// wait for a block then get post-transfer balances
	time.Sleep(1 * time.Second)
	fromAccountAfterBalance, err := provider.eth.BalanceAt(context.Background(), fromAccount.Address, nil)
	require.NoError(t, err, "should be no error getting balance")
	toAccountAfterBalance, err := provider.eth.BalanceAt(context.Background(), toAccount.Address, nil)
	require.NoError(t, err, "should be no error getting balance")

	// get balance diffs and assert logical differences
	fromAccountDiff := new(big.Int).Sub(fromAccountAfterBalance, fromAccountBeforeBalance)
	toAccountDiff := new(big.Int).Sub(toAccountAfterBalance, toAccountBeforeBalance)
	require.Equal(t, -1, fromAccountDiff.Cmp(new(big.Int).Sub(fromAccountBeforeBalance, transferAmt)), "from account balance should be < (before - amt)")
	require.Equal(t, new(big.Int).Add(toAccountBeforeBalance, transferAmt), toAccountAfterBalance, "to account should be before + amt")
	require.Equal(t, transferAmt, toAccountDiff, "to account diff should be transfer amount")
}
