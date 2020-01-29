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

type jsonMarket struct {
	MakerAssetAddress   string        `json:"makerAssetAddress"`
	TakerAssetAddresses []string      `json:"takerAssetAddresses"`
	TradeInfo           jsonTradeInfo `json:"tradeInfo"`
	QuoteInfo           jsonQuoteInfo `json:"quoteInfo"`
}

type jsonTradeInfo struct {
	ChainID  uint32 `json:"chainId"`
	GasPrice string `json:"gasPrice"`
	GasLimit string `json:"gasLimit"`
}

type jsonQuoteInfo struct {
	MinSize string `json:"minSize"`
	MaxSize string `json:"maxSize"`
}

func (r *getMarketsResponse) MarshalJSON() ([]byte, error) {
	jsonMarkets := make([]jsonMarket, r.total)
	for i, mkt := range r.record {
		jsonMarkets[i] = jsonMarket{
			MakerAssetAddress:   mkt.MakerAssetAddress,
			TakerAssetAddresses: mkt.TakerAssetAddresses,
			TradeInfo: jsonTradeInfo{
				ChainID:  mkt.TradeInfo.ChainId,
				GasPrice: mkt.TradeInfo.GasPrice,
				GasLimit: mkt.TradeInfo.GasLimit,
			},
			QuoteInfo: jsonQuoteInfo{
				MinSize: mkt.QuoteInfo.MinSize,
				MaxSize: mkt.QuoteInfo.MaxSize,
			},
		}
	}
	return json.Marshal([]interface{}{
		jsonMarkets, r.total, r.page, r.perPage,
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
