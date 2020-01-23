package rpc

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type submitFillResponse struct {
	quoteId         string
	transactionHash string
	submittedAt     int64
}

func (sfr *submitFillResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		sfr.quoteId,
		sfr.transactionHash,
		sfr.submittedAt,
	})
}

// SubmitFill implements the dealer_submitFill method
func (svc *Service) SubmitFill(quoteId string, salt string, signature string, signer string, data string, gasPrice string, expirationTimeSeconds int64) (*submitFillResponse, error) {
	order, err := svc.dealer.GetOrder(quoteId)
	if err != nil {
		return nil, err
	}

	validateReq := &grpc.ValidateOrderRequest{Order: grpc.SignedOrderToProto(order)}
	if err := svc.dealer.ValidateOrder(context.TODO(), validateReq); err != nil {
		return nil, ErrFillValidationFailed
	}

	bSignature, err := hexutil.Decode(signature)
	if err != nil {
		return nil, err
	}

	fillReq := &grpc.ExecuteZeroExTransactionRequest{
		Transaction: &grpc.ZeroExTransaction{
			Salt:                  salt,
			ExpirationTimeSeconds: expirationTimeSeconds,
			GasPrice:              gasPrice,
			SignerAddress:         signer,
			Data:                  data,
		},
		Signature: bSignature,
	}

	fillRes, err := svc.dealer.ExecuteZeroExTransaction(context.TODO(), fillReq)
	if err != nil {
		return nil, err
	}

	// todo (@hrharder): do we need to do anything with the response value here?
	if _, err := svc.dealer.WatchTransaction(context.TODO(), quoteId, fillRes.TransactionHash); err != nil {
		return nil, err
	}

	return &submitFillResponse{
		quoteId:         quoteId,
		transactionHash: fillRes.TransactionHash,
		submittedAt:     fillRes.SubmittedAt,
	}, nil
}
