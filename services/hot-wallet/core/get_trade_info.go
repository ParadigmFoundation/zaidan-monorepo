package core

import (
	"context"
	"strconv"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

func (hw *HotWallet) GetTradeInfo(ctx context.Context, req *empty.Empty) (*grpc.TradeInfo, error) {
	gasPrice, err := hw.provider.Client().SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return &grpc.TradeInfo{
		ChainId:  uint32(hw.chainId),
		GasLimit: strconv.FormatUint(zrx.EXECUTE_FILL_TX_GAS_LIMIT, 10),
		GasPrice: gasPrice.String(),
	}, nil
}
