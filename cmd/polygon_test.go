package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fetchPolygonStocks(t *testing.T) {
	var tickers PolygonTickers
	err := tickers.Fetch()
	assert.NoError(t, err)
	fmt.Println(tickers)

	dgraphClient := NewDgraphClient(":9080")
	err = tickers.Save(dgraphClient)
	assert.NoError(t, err)
}
