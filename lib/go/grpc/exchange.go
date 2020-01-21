package grpc

import (
	"strconv"
)

func (o *ExchangeOrder) Equal(v *ExchangeOrder) bool {
	toFloat := func(s string) float64 {
		f, _ := strconv.ParseFloat(s, 64)
		return f
	}

	if o.Id != v.Id {
		return false
	}
	if toFloat(o.Amount) != toFloat(v.Amount) {
		return false
	}
	if toFloat(o.Price) != toFloat(v.Price) {
		return false
	}
	if o.Side != v.Side {
		return false
	}
	if o.Symbol != v.Symbol {
		return false
	}

	return true
}
