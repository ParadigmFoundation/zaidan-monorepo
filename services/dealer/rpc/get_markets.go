package rpc

type getMarketsResponse []interface{}

func (svc *Service) GetMarkets(mAddr, tAddr string, page, perPage int) (getMarketsResponse, error) {
	mkts, err := svc.dealer.GetMarkets(mAddr, tAddr, page, perPage)
	if err != nil {
		return nil, err
	}

	return getMarketsResponse([]interface{}{mkts}), nil
}
