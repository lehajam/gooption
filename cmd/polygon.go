// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v2"
	"github.com/lehajam/gooption/core"
	"github.com/sacOO7/gowebsocket"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
	"time"
)

const APIKEY = "K105Z7BErGCa0XmV_QWUyy88PCgVp__4NM7hG_"

// const CHANNELS = "AM.AAPL,AM.CSCO,AM.MSFT,AM.FB,AM.G,AM.C"

// polygonCmd represents the polygon command
var polygonCmd = &cobra.Command{
	Use:   "polygon",
	Short: "",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dgraphAddr, err := cmd.Flags().GetString("dgraph")
		if err != nil {
			fmt.Println(err)
			return
		}

		stock, err := cmd.Flags().GetBool("stock")
		if err != nil {
			fmt.Println(err)
			return
		}

		quotes, err := cmd.Flags().GetBool("stockquotes")
		if err != nil {
			fmt.Println(err)
			return
		}

		dgraph := NewDgraphClient(dgraphAddr)

		if stock {
			err = FeedStocks(dgraph)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if quotes {
			err = FeedStockQuotes(dgraph)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(polygonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// polygonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	polygonCmd.Flags().String("dgraph", ":9080", "dgraph grpc address")
	polygonCmd.Flags().Bool("stock", false, "feed stocks")
	polygonCmd.Flags().Bool("stockquotes", false, "feed stock quotes")
}

type PolygonStock struct {
	core.GraphQL

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

type StockQuote struct {
	RefIndex      string    `json:"sym" dgraph:"Quote.ReferenceIndex,omitempty"`
	Last          float64   `json:"c" dgraph:"Quote.last,omitempty"`
	Source        string    `dgraph:"Quote.source,omitempty"`
	DatePublished time.Time `dgraph:"Quote.datePublished,omitempty"`
}

type PolygonStockQuote struct {
	StockQuote
	core.GraphQL

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
}

func FeedStocks(client *dgo.Dgraph) error {
	response, err := http.Get(fmt.Sprintf("https://api.polygon.io/v2/reference/tickers?apiKey=%s", APIKEY))
	if err != nil {
		return err
	}

	var stocks []PolygonStock
	err = core.ParseFromResponse(response, &stocks)
	if err != nil {
		return err
	}

	stocks, err = withGraphQLStockInfo(client, stocks)
	if err != nil {
		return err
	}

	return core.Save(client, stocks)
}

func FeedStockQuotes(client *dgo.Dgraph) error {
	socket := gowebsocket.New("wss://socket.polygon.io/stocks") //{host}:{port})

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		if strings.Contains(string(message), "\"ev\":\"AM\"") {
			// 1 - Get quotes from polygon tick
			var quotes []PolygonStockQuote
			err := json.Unmarshal([]byte(message), &quotes)
			if err != nil {
				// LOG ERROR
			}

			// 2 - Add extra info
			timestamp := time.Now().UTC()
			for i := 0; i < len(quotes); i++ {
				quotes[i].Source = "polygon"
				quotes[i].DatePublished = timestamp
				quotes[i].Types = []string{"Quote"}
			}

			// 2 - Save to DB
			err = core.Save(client, quotes)
			if err != nil {
				// LOG ERROR
			}
		}
	}

	// 1 - Get Stocks from DB
	symbols, err := getAllStockSymbols(client)
	if err != nil {
		return err
	}

	// 2 - Register
	channels := "AM." + strings.Join(symbols, ",AM.") //AM.AAPL,AM.MSFT ...
	socket.SendText(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKEY))
	socket.SendText(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", channels))
	socket.Connect()

	waitc := make(chan struct{})
	<-waitc

	return nil
}

func withGraphQLStockInfo(client *dgo.Dgraph, stocks []PolygonStock) ([]PolygonStock, error) {
	query := `
{
  tickers(func: type(Stock)) @filter(eq(ReferenceIndex.symbol, %s)) {
    uid
    ReferenceIndex.symbol
	dgraph.type
  }
}`

	symbols := make([]string, len(stocks))
	for i, s := range stocks {
		symbols[i] = s.Ticker
	}

	query = fmt.Sprintf(query, symbols)
	count, err := core.Query(client, query, "uid", &stocks)
	if err != nil {
		return nil, err
	}

	if count != uint64(len(stocks)) {
		types := []string{"ReferenceIndex", "Stock"}
		for i := range stocks {
			stocks[i].Types = types
		}
	}

	return stocks, nil
}

func getAllStockSymbols(client *dgo.Dgraph) ([]string, error) {
	query := `
{
  tickers(func: type(Stock)) {
    ReferenceIndex.symbol
  }
}`

	var symbols []string
	_, err := core.Query(client, query, "ReferenceIndex.symbol", &symbols)
	if err != nil {
		return nil, err
	}

	return symbols, nil
}
