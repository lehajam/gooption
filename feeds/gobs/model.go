package gobs

import (
	"github.com/lehajam/gooption/core"
	"time"
)

type Trade struct {
	core.GraphQL
	Quantity int `json:"Trade.quantity"`
	Contract struct {
		RefIndex []struct {
			Quotes []struct {
				DatePublished time.Time `json:"Quote.datePublished"`
				Last          float64   `json:"Quote.last"`
			} `json:"ReferenceIndex.quotes"`
		} `json:"Contract.refIndex"`
		Strike  float64   `json:"Option.strike"`
		Expiry  time.Time `json:"Option.expiry"`
		Putcall string    `json:"Option.putcall"`
	} `json:"Trade.contract"`
}

type ValuationResult struct {
	core.GraphQL
	DatePublished time.Time `dgraph:"ValuationResult.datePublished"`
	Value         float64   `dgraph:"ValuationResult.value"`
	Source        string    `dgraph:"ValuationResult.source"`
}

type TradeWithResults struct {
	core.GraphQL
	Valuations []ValuationResult `json:"Trade.valuations,omitempty"`
}
