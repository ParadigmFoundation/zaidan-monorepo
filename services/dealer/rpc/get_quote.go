package rpc

import (
	"context"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ethereum/go-ethereum/common"
)

func (svc *Service) GetQuote(makerAsset string, takerAsset string, makerSize *string, takerSize *string, takerAddress *string, includeOrder *bool) error {
	// var inclOrder bool
	// if includeOrder == nil {
	// 	inclOrder = true
	// } else {
	// 	inclOrder = *includeOrder
	// }

	var taker common.Address
	if takerAddress == nil {
		taker = eth.NULL_ADDRESS
	} else {
		taker = common.HexToAddress(*takerAddress)
	}

	var makerAmount, takerAmount string
	if makerSize == nil && takerSize == nil {
		return ErrInvalidParameters
	} else if makerSize != nil && takerSize != nil {
		return ErrInvalidParameters
	} else if makerSize != nil {
		makerAmount = *makerSize
	} else if takerSize != nil {
		takerAmount = *takerSize
	}

	req := &grpc.GetQuoteRequest{
		TakerAsset:   takerAsset,
		MakerAsset:   makerAsset,
		TakerSize:    takerAmount,
		MakerSize:    makerAmount,
		TakerAddress: taker.Hex(),
	}

	if err := svc.dealer.FetchQuote(context.Background(), req); err != nil {
		return err
	}

	return ErrTemporaryRestriction
}
