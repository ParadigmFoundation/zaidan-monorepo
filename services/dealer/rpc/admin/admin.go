package admin

import (
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
)

type Service struct {
	rpc *rpc.Server
	log *logger.Logger
}

func NewService() (*Service, error) {
	srv := &Service{
		rpc: rpc.NewServer(),
		log: logger.New("admin", logger.HandleEthLog()),
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

// GetEtherBalance fetches on-chain Ether balance for the hot-wallet
func (srv *Service) GetEtherBalance() error { panic("not implemented") }

// GetTokenBalance fetches on-chain balance for an ERC-20 token by address
func (srv *Service) GetTokenBalance() error { panic("not implemented") }

// GetAllowance fetches 0x ERC-20 asset proxy allowances for a token by address
func (srv *Service) GetAllowance() error { panic("not implemented") }

// SetAllowance set max/specific allowance for ERC-20 asset proxy contract by address.
func (srv *Service) SetAllowance() error { panic("not implemented") }

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

// AddTaker adds a taker to the registry to whitelist/blacklist them (depending on mode)
func (srv *Service) AddTaker() error { panic("not implemented") }

// RemoveTaker removes a taker from the registry (see _addTaker)
func (srv *Service) RemoveTaker() error { panic("not implemented") }

// GetTakers fetches a table of all takers in the registry
func (srv *Service) GetTakers() error { panic("not implemented") }

// GetTakerAuthorization fetches authorization status for a taker by their address
func (srv *Service) GetTakerAuthorization() error { panic("not implemented") }
