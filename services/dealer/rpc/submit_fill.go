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
func (svc *Service) SubmitFill(quoteId string, salt string, signature string, signer string, data string, gasPrice string, expirationTimeSeconds string) (*submitFillResponse, error) {
	order, err := svc.dealer.GetOrder(quoteId)
	if err != nil {
		return nil, err
	}

	validateReq := &grpc.ValidateOrderRequest{Order: grpc.SignedOrderToProto(order)}
	if err := svc.dealer.ValidateOrder(context.Background(), validateReq); err != nil {
		return nil, ErrFillValidationFailed
	}

	bData, err := hexutil.Decode(data)
	if err != nil {
		return nil, err
	}

	bSignature, err := hexutil.Decode(signature)
	if err != nil {
		return nil, err
	}

	fillReq := &grpc.ExecuteZeroExTransactionRequest{
		Salt:                  salt,
		ExpirationTimeSeconds: expirationTimeSeconds,
		GasPrice:              gasPrice,
		SignerAddress:         signer,
		Data:                  bData,
		Signature:             bSignature,
	}

	fillRes, err := svc.dealer.ExecuteZeroExTransaction(context.Background(), fillReq)
	if err != nil {
		return nil, err
	}

	return &submitFillResponse{
		quoteId:         quoteId,
		transactionHash: fillRes.TransactionHash,
		submittedAt:     fillRes.SubmittedAt,
	}, nil
}
