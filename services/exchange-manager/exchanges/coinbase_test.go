package exchanges

import (
	"os"
	"testing"

	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
	"github.com/stretchr/testify/suite"
)

func TestCoinbase(t *testing.T) {
	cfg := &coinbasepro.ClientConfig{
		BaseURL:    "https://api-public.sandbox.pro.coinbase.com",
		Key:        os.Getenv("EM_COINBASE_TEST_KEY"),
		Passphrase: os.Getenv("EM_COINBASE_TEST_PASSPHRASE"),
		Secret:     os.Getenv("EM_COINBASE_TEST_SECRET"),
	}
	if cfg.Key == "" && cfg.Passphrase == "" && cfg.Secret == "" {
		t.Skip("EM_COINBASE_ env not defined")
	}

	exchange := NewCoinbase(cfg)
	suite.Run(t, &ExchangesSuite{exchange: exchange})
}
