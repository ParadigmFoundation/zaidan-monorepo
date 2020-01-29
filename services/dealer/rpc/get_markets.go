package rpc

import (
	"encoding/json"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/utils/ptr"
)

const DEFAULT_PER_PAGE = 10

type getMarketsResponse struct {
	record  []*types.Market
	total   int
	page    int
	perPage int
}

func (r *getMarketsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		r.record, r.total, r.page, r.perPage,
	})
}

func (svc *Service) GetMarkets(mAddr, tAddr *string, page, perPage *int) (*getMarketsResponse, error) {

	mkts, err := svc.dealer.GetMarkets(
		ptr.String(mAddr), ptr.String(tAddr), ptr.Int(page), ptr.IntOr(perPage, DEFAULT_PER_PAGE),
	)
	if err != nil {
		return nil, err
	}

	return &getMarketsResponse{
		record:  mkts,
		total:   len(mkts),
		page:    ptr.Int(page),
		perPage: ptr.IntOr(perPage, DEFAULT_PER_PAGE),
	}, nil
}
