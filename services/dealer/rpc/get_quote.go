package rpc

import (
	"context"
	"encoding/json"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
	"github.com/ethereum/go-ethereum/common"
)

type getQuoteResponse struct {
	quote
}

func (gqr *getQuoteResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{gqr.quote})
}

type quote struct {
	QuoteId               string                `json:"quoteId"`
	MakerAssetAddress     string                `json:"makerAssetAddress"`
	TakerAssetAddress     string                `json:"takerAssetAddress"`
	MakerAssetSize        string                `json:"makerAssetSize"`
	TakerAssetSize        string                `json:"takerAssetSize"`
	Expiration            int64                 `json:"expiration"`
	ServerTime            int64                 `json:"serverTime"`
	ZeroExTransactionHash string                `json:"zeroExTransactionHash"`
	ZeroExTransactionInfo zeroExTransactionInfo `json:"zeroExTransactionInfo"`
}

type zeroExTransactionInfo struct {
	Order       *zeroex.SignedOrder `json:"order"`
	Transaction *zrx.Transaction    `json:"transaction"`
}

func (svc *Service) GetQuote(makerAsset string, takerAsset string, makerSize *string, takerSize *string, takerAddress *string, includeOrder *bool) (*getQuoteResponse, error) {
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
		return nil, ErrInvalidParameters
	} else if makerSize != nil && takerSize != nil {
		return nil, ErrInvalidParameters
	} else if makerSize != nil {
		makerAmount = *makerSize
	} else if takerSize != nil {
		takerAmount = *takerSize
	}

	quoteReq := &grpc.GetQuoteRequest{
		TakerAsset:   takerAsset,
		MakerAsset:   makerAsset,
		TakerSize:    takerAmount,
		MakerSize:    makerAmount,
		TakerAddress: taker.Hex(),
	}

	quoteRes, err := svc.dealer.FetchQuote(context.Background(), quoteReq)
	if err != nil {
		svc.log.Errorf("failed to fetch quote: %v", err)
		return nil, err
	}

	order, err := quoteRes.ZeroExTransactionInfo.Order.ToZeroExSignedOrder()
	if err != nil {
		svc.log.Errorf("failed to convert signed order: %v", err)
		return nil, err
	}

	tx, err := quoteRes.ZeroExTransactionInfo.Transaction.ToZeroExTransaction()
	if err != nil {
		svc.log.Errorf("failed to convert 0x transaction: %v", err)
		return nil, err
	}

	res := quote{
		QuoteId:               quoteRes.QuoteId,
		MakerAssetAddress:     quoteRes.MakerAssetAddress,
		TakerAssetAddress:     quoteRes.TakerAssetAddress,
		MakerAssetSize:        quoteRes.MakerAssetSize,
		TakerAssetSize:        quoteRes.TakerAssetSize,
		Expiration:            quoteRes.Expiration,
		ServerTime:            quoteRes.ServerTime,
		ZeroExTransactionHash: quoteRes.ZeroExTransactionHash,
		ZeroExTransactionInfo: zeroExTransactionInfo{
			Order:       order,
			Transaction: tx,
		},
	}
	return &getQuoteResponse{res}, nil
}
