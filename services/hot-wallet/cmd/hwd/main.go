package main

import (
	"log"

	"github.com/spf13/pflag"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"
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

var cfg = config{}

func init() {
	pflag.StringVar(&cfg.ethurl, "eth", "http://localhost:8545", "URL to an Ethereum JSONRPC server")
	pflag.StringVar(&cfg.mnemonic, "mnemonic", "concert load couple harbor equip island argue ramp clarify fence smart topic", "a 12 word seed phrase")
	pflag.StringVar(&cfg.bind, "bind", "0.0.0.0:42001", "address to bind the hot-wallet gRPC server to")
	pflag.IntVar(&cfg.maxReqLen, "max-request-length", 16384, "set the max request length (in bytes) used by the order validator")
	pflag.Uint32Var(&cfg.makerIndex, "maker-index", 0, "set the index of the maker account for 0x orders")
	pflag.Uint32Var(&cfg.senderIndex, "sender-index", 1, "set the index of the maker account for 0x orders")
}

func main() {
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
	log.Println("maker address", makerAcct.Address.Hex())
	log.Println("sender address", senderAcct.Address.Hex())

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

	log.Println("hot wallet started on", cfg.bind)
	log.Fatal(<-errChan)
}
