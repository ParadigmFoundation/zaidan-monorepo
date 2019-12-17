package main

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/caarlos0/env/v6"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet/grpc"
)

// test server config
type config struct {
	Ethurl   string `env:"ETHEREUM_JSONRPC_URL" envDefault:"http://localhost:8545"`
	Mnemonic string `env:"MNEMONIC" envDefault:"concert load couple harbor equip island argue ramp clarify fence smart topic"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	path1 := accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}
	path2 := accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 1}
	provider, err := eth.NewProvider(cfg.Ethurl, cfg.Mnemonic, path1)
	if err != nil {
		log.Fatal(err)
	}

	if err := provider.Derive(path2); err != nil {
		log.Fatal(err)
	}

	accounts := provider.Accounts()
	hwcfg := core.HotWalletConfig{
		OrderValidatorMaxReqLength: 16384,
		MakerAddress:               accounts[0].Address.Hex(),
		SenderAddress:              accounts[1].Address.Hex(),
	}

	hw, err := core.NewHotWallet(provider, hwcfg)
	if err != nil {
		log.Fatal(err)
	}

	svr := grpc.NewServer(hw)

	if err := svr.Listen("0.0.0.0:8000"); err != nil {
		log.Fatal(err)
	}
}
