package core

import (
	"context"
	"math/big"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

func (hw *HotWallet) ValidateOrder(ctx context.Context, req *grpc.ValidateOrderRequest) (*grpc.ValidateOrderResponse, error) {
	order, err := req.Order.ToZeroExSignedOrder()
	if err != nil {
		return nil, err
	}

	takerAmount, ok := new(big.Int).SetString(req.TakerAssetAmount, 10)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "unable to parse value for 'takerAmount'")
	}

	if err := hw.zrxHelper.ValidateFill(ctx, order, takerAmount); err != nil {
		return &grpc.ValidateOrderResponse{
			Valid: false,
			Info:  err.Error(),
		}, nil
	} else {
		return &grpc.ValidateOrderResponse{Valid: true}, nil
	}
}
