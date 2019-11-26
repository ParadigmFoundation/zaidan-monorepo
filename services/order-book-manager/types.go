package obm

import (
	"fmt"
	"strconv"
	"time"
)

type Entry struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
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
	Symbol string
	Bids   Entries
	Asks   Entries
}

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("0"), nil
	}
	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}

type OrderBook struct {
	LastUpdate Time               `json:"last_update"`
	Exchange   string             `json:"exchange"`
	Symbol     string             `json:"symbol"`
	Asks       EntriesByPriceAsc  `json:"asks"`
	Bids       EntriesByPriceDesc `json:"bids"`
}
