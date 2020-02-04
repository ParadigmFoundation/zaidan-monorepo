package main

import (
	"github.com/spf13/pflag"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet/grpc"
)

type config struct {
	ethurl    string
	mnemonic  string
	bind      string
	maxReqLen int

	makerIndex  uint32
	senderIndex uint32
}

func main() {
	var cfg config
	pflag.StringVarP(&cfg.ethurl, "eth", "e", "http://localhost:8545", "URL to an Ethereum JSONRPC server")
	pflag.StringVarP(&cfg.mnemonic, "mnemonic", "m", "concert load couple harbor equip island argue ramp clarify fence smart topic", "a 12 word seed phrase")
	pflag.StringVarP(&cfg.bind, "bind", "b", "0.0.0.0:42001", "address to bind the hot-wallet gRPC server to")
	pflag.IntVarP(&cfg.maxReqLen, "max-request-length", "l", 16384, "set the max request length (in bytes) used by the order validator")
	pflag.Uint32VarP(&cfg.makerIndex, "maker-index", "M", 0, "set the index of the maker account for 0x orders")
	pflag.Uint32VarP(&cfg.senderIndex, "sender-index", "S", 1, "set the index of the maker account for 0x orders")
	pflag.Parse()

	log := logger.New("app")

	log.WithFields(logger.Fields{"ethurl": cfg.ethurl, "mnemonic": cfg.mnemonic}).Info("Initializing Ethereum Driver")
	deriver := eth.NewBaseDeriver()
	provider, err := eth.NewProvider(cfg.ethurl, cfg.mnemonic, deriver.Base())
	if err != nil {
		log.Fatal(err)
	}

	makerPath := deriver.DeriveAt(cfg.makerIndex)
	senderPath := deriver.DeriveAt(cfg.senderIndex)
	if err := provider.Derive(makerPath); err != nil {
		log.Fatal(err)
	}
	if err := provider.Derive(senderPath); err != nil {
		log.Fatal(err)
	}

	makerAcct, _ := provider.AccountAt(makerPath)
	senderAcct, _ := provider.AccountAt(senderPath)
	log.WithFields(logger.Fields{
		"maker":  makerAcct.Address.Hex(),
		"sender": senderAcct.Address.Hex(),
	}).Info("HotWalletConfig")
	hwcfg := core.HotWalletConfig{
		OrderValidatorMaxReqLength: cfg.maxReqLen,
		MakerAddress:               makerAcct.Address,
		SenderAddress:              senderAcct.Address,
	}

	hw, err := core.NewHotWallet(provider, hwcfg)
	if err != nil {
		log.Fatal(err)
	}

	errChan := make(chan error)
	svr := grpc.NewServer(hw)
	go func() {
		errChan <- svr.Listen(cfg.bind)
	}()

	log.WithField("bind", cfg.bind).Info("hot wallet started")
	log.Fatal(<-errChan)
}
