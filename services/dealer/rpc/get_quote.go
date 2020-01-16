package rpc

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gogo/protobuf/jsonpb"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ethereum/go-ethereum/common"
)

type getQuoteResponse struct {
	*grpc.Quote
}

func (gcr *getQuoteResponse) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := new(jsonpb.Marshaler).Marshal(buf, gcr.Quote); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

	req := &grpc.GetQuoteRequest{
		TakerAsset:   takerAsset,
		MakerAsset:   makerAsset,
		TakerSize:    takerAmount,
		MakerSize:    makerAmount,
		TakerAddress: taker.Hex(),
	}

	quote, err := svc.dealer.FetchQuote(context.Background(), req)
	if err != nil {
		return nil, err
	}

	res := &getQuoteResponse{quote}

	bts, err := res.MarshalJSON()
	if err != nil {
		return nil, err
	}

	fmt.Println(string(bts))
	return res, nil
}
