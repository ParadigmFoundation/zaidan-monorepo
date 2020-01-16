package exchanges

import types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"

func SideFromString(side string) types.ExchangeOrder_Side {
	switch side {
	case "buy":
		return types.ExchangeOrder_BUY
	case "sell":
		return types.ExchangeOrder_SELL
	}
	return 0
}
