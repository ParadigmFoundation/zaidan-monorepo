package rpc

import (
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

func (svc *Service) LoadQuote(quoteId string) (*grpc.Quote, error) {
	return svc.dealer.GetQuote(quoteId)
}
