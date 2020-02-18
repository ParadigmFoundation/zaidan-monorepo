package rpc

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/peterbourgon/ff"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql"
)

type Service struct {
	dealer     *core.Dealer
	server     *rpc.Server
	policyMode PolicyMode
	policy     store.Policy
	log        *logger.Logger
}

// NewService creates a new Dealer JSONRPC service
func NewService(dealer *core.Dealer) (*Service, error) {
	srv := &Service{
		dealer: dealer,
		server: rpc.NewServer(),
		log:    logger.New("rpc", logger.HandleEthLog()),
	}

	if err := srv.server.RegisterName("dealer", srv); err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/ws" {
		srv.log.Infof("New WebSocket connection")
		srv.server.WebsocketHandler([]string{"*"}).ServeHTTP(w, r)
	} else {
		srv.log.Infof("%s %s", r.Method, path)
		srv.server.ServeHTTP(w, r)
	}
}

func (srv *Service) WithPolicy(mode PolicyMode, policy store.Policy) *Service {
	srv.policyMode = mode
	srv.policy = policy
	return srv
}

func unQuote(s string) string {
	isQuote := func(c byte) bool {
		return c == '\'' || c == '"'
	}

	n := len(s)
	if isQuote(s[0]) && s[0] == s[n-1] {
		return s[1 : n-1]
	}
	return s
}

type ServerCfg struct {
	DB            string
	DSN           string
	Bind          string
	Grpc		  string
	MakerAddr     string
	HwAddr        string
	WatcherAddr   string
	OrderDuration int64
	PolicyBlack   bool
	PolicyWhite   bool
}

// ParseServerCfg builds a new ServerCfg from user flags
func ParseServerCfg(prefix string) (*ServerCfg, error) {
	cfg := ServerCfg{}

	fs := flag.NewFlagSet("dealer", flag.ExitOnError)
	fs.StringVar(&cfg.DB, "db", "sqlite3", "Database driver [sqlite3|postgres]")
	fs.StringVar(&cfg.DSN, "dsn", ":memory:", "Database's Data Source Name (see driver's doc for details)")
	fs.StringVar(&cfg.Bind, "bind", ":8000", "Server binding address")
	fs.StringVar(&cfg.Grpc, "grpc", "localhost:44001", "Grpc binding port")
	fs.StringVar(&cfg.MakerAddr, "maker", "localhost:50051", "The address to connect to a maker server over gRPC on")
	fs.StringVar(&cfg.HwAddr, "hw", "localhost:42001", "The address to connect to a hot-wallet server over gRPC on")
	fs.StringVar(&cfg.WatcherAddr, "watcher", "localhost:43001", "The address to connect to a watcher server over gRPC on")
	fs.Int64Var(&cfg.OrderDuration, "order-duration", 600, "The duration in seconds that signed order should be valid for")
	fs.BoolVar(&cfg.PolicyBlack, "policy.blacklist", false, "Enable BlackList policy mode")
	fs.BoolVar(&cfg.PolicyWhite, "policy.whitelist", false, "Enable WhiteList policy mode")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix(prefix)); err != nil {
		return nil, err
	}

	cfg.DSN = unQuote(cfg.DSN)

	return &cfg, nil
}

func StartServer() error {
	cfg, err := ParseServerCfg("DEALER")
	if err != nil {
		return err
	}

	log := logger.New("app")
	log.WithFields(logger.Fields{"db": cfg.DB, "dsn": cfg.DSN}).Info("Initializing database")
	store, err := sql.New(cfg.DB, cfg.DSN)
	if err != nil {
		return err
	}

	dealerCfg := core.DealerConfig{
		MakerBindAddress:     cfg.MakerAddr,
		DealerBindAddress:    cfg.Bind,
		HotWalletBindAddress: cfg.HwAddr,
		WatcherBindAddress:   cfg.WatcherAddr,
		OrderDuration:        cfg.OrderDuration,
	}
	log.WithField("cfg", dealerCfg).Info("Dealer")
	dealer, err := core.NewDealer(context.Background(), store, dealerCfg)
	if err != nil {
		return err
	}

	service, err := NewService(dealer)
	if err != nil {
		return err
	}

	if cfg.PolicyBlack && cfg.PolicyWhite {
		return errors.New("can't use both -policy.blacklist and -policy.whitelist")
	}

	if cfg.PolicyBlack || cfg.PolicyWhite {
		var mode PolicyMode
		if cfg.PolicyBlack {
			log.Print("Using Blacklist mode")
			mode = PolicyBlackList
		}
		if cfg.PolicyWhite {
			log.Print("Using Whitelist mode")
			mode = PolicyWhiteList
		}

		service.WithPolicy(mode, store)
	}

	ln, err := net.Listen("tcp", cfg.Bind)
	if err != nil {
		return err
	}
	log.WithField("bind", cfg.Bind).Info("Server started")
	server := &http.Server{Handler: service}
	go grpc.CreateAndListen(store, cfg.Grpc)
	return server.Serve(ln)
}
