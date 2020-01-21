package rpc

import "encoding/json"

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
func (svc *Service) SubmitFill(quoteId string, salt string, signature string, signer string, data string, hash string, gasPrice string) (*submitFillResponse, error) {
	return nil, nil
}
