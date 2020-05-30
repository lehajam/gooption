package polygon

import (
	"github.com/lehajam/gooption/core"
	"time"
)

type PolygonStock struct {
	core.GraphQL

	// from polygon feed
	Ticker      string `json:"ticker" dgraph:"ReferenceIndex.symbol,omitempty"`
	Name        string `json:"name" dgraph:"PolygonStock.name,omitempty"`
	Market      string `json:"market" dgraph:"PolygonStock.market,omitempty"`
	Locale      string `json:"locale" dgraph:"PolygonStock.locale,omitempty"`
	Type        string `json:"type" dgraph:"PolygonStock.type,omitempty"`
	Currency    string `json:"currency" dgraph:"Stock.currency,omitempty"`
	Active      bool   `json:"active" dgraph:"PolygonStock.active,omitempty"`
	PrimaryExch string `json:"primaryExch" dgraph:"PolygonStock.primaryexchange,omitempty"`
	Updated     string `json:"updated" dgraph:"PolygonStock.updated,omitempty"`
	URL         string `json:"url" dgraph:"PolygonStock.url,omitempty"`
}

type PolygonStockQuote struct {
	core.GraphQL

	// from polygon feed
	Ev  string  `json:"ev"`
	Sym string  `json:"sym" dgraph:"Quote.symbol,omitempty"`
	V   int     `json:"v" dgraph:"Quote.volume,omitempty"`
	Av  int     `json:"av"`
	Op  float64 `json:"op"`
	Vw  float64 `json:"vw"`
	O   float64 `json:"o" dgraph:"Quote.open,omitempty"`
	C   float64 `json:"c" dgraph:"Quote.close,omitempty"`
	H   float64 `json:"h" dgraph:"Quote.high,omitempty"`
	L   float64 `json:"l" dgraph:"Quote.low,omitempty"`
	A   float64 `json:"a"`
	S   int64   `json:"s"`
	E   int64   `json:"e"`

	// From graphql schema
	RefIndex      string    `json:"sym" dgraph:"Quote.ReferenceIndex,omitempty"`
	Last          float64   `json:"c" dgraph:"Quote.last,omitempty"`
	Source        string    `dgraph:"Quote.source,omitempty"`
	DatePublished time.Time `dgraph:"Quote.datePublished,omitempty"`
}
