package admin

import (
	"context"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

// GetEtherBalance fetches on-chain Ether balance for the hot-wallet
func (srv *Service) GetEtherBalance(owner, taker string) (string, error) {
	r, err := srv.dealer.HWClient().GetEtherBalance(
		context.Background(),
		&grpc.GetBalanceRequest{OwnerAddress: owner, TokenAddress: taker},
	)
	if err != nil {
		return "", err
	}

	return r.GetBalance(), nil
}

// GetTokenBalance fetches on-chain balance for an ERC-20 token by address
func (srv *Service) GetTokenBalance(owner, taker string) (string, error) {
	r, err := srv.dealer.HWClient().GetTokenBalance(
		context.Background(),
		&grpc.GetBalanceRequest{OwnerAddress: owner, TokenAddress: taker},
	)
	if err != nil {
		return "", err
	}

	return r.GetBalance(), nil
}

// GetAllowance fetches 0x ERC-20 asset proxy allowances for a token by address
func (srv *Service) GetAllowance(owner, taker string) (string, error) {
	r, err := srv.dealer.HWClient().GetAllowance(
		context.Background(),
		&grpc.GetAllowanceRequest{OwnerAddress: owner, TokenAddress: taker},
	)
	if err != nil {
		return "", err
	}

	return r.GetAllowance(), nil
}

// SetAllowance set max/specific allowance for ERC-20 asset proxy contract by address.
func (srv *Service) SetAllowance(token, allowance string) error {
	_, err := srv.dealer.HWClient().SetAllowance(
		context.Background(),
		&grpc.SetAllowanceRequest{TokenAddress: token, Allowance: allowance},
	)
	return err
}
