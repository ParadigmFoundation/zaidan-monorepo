package exchange

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Symbol represent a currency pair definition. i.e BTC-USD
type Symbol struct {
	Base  string
	Quote string
}

// NewSymbol returns a new symbol
func NewSymbol(b, q string) *Symbol { return &Symbol{Base: b, Quote: q} }

func NewSymbolFromString(str, sep string) (*Symbol, error) {
	split := strings.Split(str, sep)
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid symbol format: %s", str)
	}

	return NewSymbol(split[0], split[1]), nil
}

// String returns the formatted symbol
func (s Symbol) String() string { return fmt.Sprintf("%s-%s", s.Base, s.Quote) }

type Entry struct {
	Price    float64
	Quantity float64
}

func NewEntryFromStrings(p, q string) (*Entry, error) {
	price, err := strconv.ParseFloat(p, 64)
	if err != nil {
		return nil, err
	}

	quantity, err := strconv.ParseFloat(q, 64)
	if err != nil {
		return nil, err
	}

	return &Entry{Price: price, Quantity: quantity}, nil
}

type Entries []*Entry

// EntriesByPriceAsc implements the sort interface. Returns the entries sorted by ascending Price
type EntriesByPriceAsc Entries

func (e EntriesByPriceAsc) Len() int           { return len(e) }
func (e EntriesByPriceAsc) Less(i, j int) bool { return e[i].Price < e[j].Price }
func (e EntriesByPriceAsc) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// EntriesByPriceDesc implements the sort interface. Returns the entries sorted by descending Price
type EntriesByPriceDesc Entries

func (e EntriesByPriceDesc) Len() int           { return len(e) }
func (e EntriesByPriceDesc) Less(i, j int) bool { return e[i].Price > e[j].Price }
func (e EntriesByPriceDesc) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

type Update struct {
	Symbol Symbol
	Bids   Entries
	Asks   Entries
}

type Market struct {
	Symbol Symbol
	Bids   map[float64]float64
	Asks   map[float64]float64
}

func NewMarket(sym Symbol) *Market {
	return &Market{
		Symbol: sym,
		Bids:   make(map[float64]float64),
		Asks:   make(map[float64]float64),
	}
}

type Callbacks struct {
	OnSnapshot func(*Update) error
	OnChange   func(*Update) error
}

type Exchange interface {
	// Subscribe subscribes to one or more symbols with a given set of callbacks
	Subscribe(context.Context, Callbacks, Symbol, ...Symbol) error
}
