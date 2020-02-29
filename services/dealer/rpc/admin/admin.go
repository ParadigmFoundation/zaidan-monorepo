package admin

import (
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ParadigmFoundation/go-logrus-bugsnag/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/rpc/policy"
)

type Service struct {
	rpc    *rpc.Server
	log    *logger.Logger
	dealer *core.Dealer
	policy *policy.Policy
}

func NewService(dealer *core.Dealer, policy *policy.Policy) (*Service, error) {
	srv := &Service{
		rpc:    rpc.NewServer(),
		log:    logger.New("admin", logger.HandleEthLog()),
		dealer: dealer,
		policy: policy,
	}

	if err := srv.rpc.RegisterName("admin", srv); err != nil {
		return nil, err
	}
	return srv, nil
}

func (srv *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/ws" {
		srv.log.Infof("New WebSocket connection")
		srv.rpc.WebsocketHandler([]string{"*"}).ServeHTTP(w, r)
	} else {
		srv.log.Infof("%s %s", r.Method, path)
		srv.rpc.ServeHTTP(w, r)
	}
}

// GetExchangeBalances fetches all balances from all currently supported centralized exchanges.
func (srv *Service) GetExchangeBalances() error { panic("not implemented") }

// SendTokens performs an on-chain ERC-20 asset transfer to a specific Ethereum address
func (srv *Service) SendTokens() error { panic("not implemented") }

// SendEther performs an on-chain Ether transfer to a specific Ethereum address
func (srv *Service) SendEther() error { panic("not implemented") }

// WithdrawAssets initiates an asset withdrawal from a centralized exchange to the hot-wallet
func (srv *Service) WithdrawAssets() error { panic("not implemented") }

// WrapEther deposits Ether into the wrapped ETH contract to mint WETH
func (srv *Service) WrapEther() error { panic("not implemented") }

// UnwrapEther withdraws Ether from the wrapped ETH contract by redeeming WETH
func (srv *Service) UnwrapEther() error { panic("not implemented") }

// DisableMaker deactivates the Maker service until the next restart (emergency shutdown)
func (srv *Service) DisableMaker() error { panic("not implemented") }

// GetTradeHistory fetches a table of previously settled 0x trades
func (srv *Service) GetTradeHistory() error { panic("not implemented") }

// GetExchangeHistory fetches a table of past orders from integrated centralized exchanges
func (srv *Service) GetExchangeHistory() error { panic("not implemented") }

// GetOrderBook fetches a centralized exchange order book snapshot
func (srv *Service) GetOrderBook() error { panic("not implemented") }

// GetUnconfirmedOrders fetches a table of 0x fills that have been submitted but not mined
func (srv *Service) GetUnconfirmedOrders() error { panic("not implemented") }

// GetOutstandingQuotes: Fetch a table of quotes that have been sent out that are not expired
func (srv *Service) GetOutstandingQuotes() error { panic("not implemented") }
