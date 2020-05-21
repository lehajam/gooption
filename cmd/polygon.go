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
	"context"
	"fmt"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

const APIKEY = "K105Z7BErGCa0XmV_QWUyy88PCgVp__4NM7hG_"
const CHANNELS = "AM.AAPL,AM.CSCO,AM.MSFT,AM.FB,AM.G,AM.C"

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

		var stocks PolygonTickers
		client := NewDgraphClient(":9080")

		err := stocks.Fetch()
		if err != nil {
			panic(err)
		}

		err = stocks.Save(client)
		if err != nil {
			panic(err)
		}

		//graphqlAddr, err := cmd.Flags().GetString("graphql")
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//client := graphql.NewClient(graphqlAddr)
		//
		//stocks := map[string]string{}
		//var response QueryStockResponse
		//err = RunGraphQLQuery(client, "{queryStock { id ticker }}", nil, &response)
		//for _, stock := range response.QueryStock {
		//	stocks[stock.Ticker] = stock.ID
		//}
		//
		//socket := gowebsocket.New("wss://socket.polygon.io/stocks") //{host}:{port})
		//socket.Connect()
		//
		//socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		//	if strings.Contains(string(message), "\"ev\":\"AM\"") {
		//		var quotes []PolygonQuote
		//		if err := json.Unmarshal([]byte(message), &quotes); err != nil {
		//			panic(err)
		//		}
		//
		//		fmt.Printf("Received %d quotes\n", len(quotes))
		//		for _, quote := range quotes {
		//			if _, exist := stocks[quote.Sym]; !exist {
		//				var res AddStockResponse
		//				vars := map[string]interface{}{"ticker": quote.Sym}
		//				err = RunGraphQLQuery(client,
		//					`mutation ($ticker: String!) {
		//						addStock(input: [{ ticker: $ticker }]) {
		//							stock {
		//								id
		//							}
		//						}
		//					}`, vars, &res)
		//				if err != nil {
		//					fmt.Println(err.Error())
		//					continue
		//				}
		//
		//				stocks[quote.Sym] = res.AddStock.Stock[0].ID
		//				fmt.Printf("Added stock %s\n", quote.Sym)
		//			}
		//
		//			var res AddStockQuoteResponse
		//			vars := map[string]interface{}{"timestamp": time.Now(), "indexID": stocks[quote.Sym], "close": quote.C}
		//			err = RunGraphQLQuery(client,
		//				`mutation ($timestamp: DateTime!, $indexID: ID!, $close: Float!) {
		//					addStockQuote(input: [{ index: { id: $indexID }, datePublished: $timestamp, close: $close}]){
		//						stockquote {
		//						  id
		//						}
		//					  }
		//				}`,
		//				vars, &res)
		//			if err != nil {
		//				fmt.Println(err.Error())
		//				continue
		//			}
		//			fmt.Printf("Added quote for stock %s\n", quote.Sym)
		//		}
		//	}
		//}
		//
		//socket.SendText(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKEY))
		//socket.SendText(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", CHANNELS))
		//
		//waitc := make(chan struct{})
		//<-waitc
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
	polygonCmd.Flags().String("graphql", "http://localhost:8080/graphql", "graphql server URL")
}

type DgraphObject struct {
	ID          string   `json:"uid" dgraph:"uid,omitempty"`
	GraphqlType []string `json:"dgraph.type" dgraph:"dgraph.type,omitempty"`
}

type PolygonStock struct {
	DgraphObject

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

type PolygonTickers struct {
	Tickers []PolygonStock `json:"tickers,omitempty"`
}

func (p *PolygonTickers) Fetch() error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	response, err := http.Get(fmt.Sprintf("https://api.polygon.io/v2/reference/tickers?apiKey=%s", APIKEY))
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, p)
	if err != nil {
		return err
	}

	return nil
}

func (p *PolygonTickers) WithTypeAndID(client *dgo.Dgraph) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	symbols := make([]string, len(p.Tickers))
	for i, t := range p.Tickers {
		symbols[i] = t.Ticker
	}

	query := fmt.Sprintf(`{  
  tickers(func: type(Stock)) @filter(eq(ReferenceIndex.symbol, %s)) {
    uid
    ReferenceIndex.symbol
	dgraph.type
  }
}`, symbols)

	stocks, err := client.NewTxn().Query(context.Background(), query)
	if err != nil {
		return err
	}

	if stocks.Metrics.NumUids["uid"] > 0 {
		err = json.Unmarshal(stocks.GetJson(), p)
		if err != nil {
			return err
		}
	}

	if stocks.Metrics.NumUids["uid"] != uint64(len(p.Tickers)) {
		types := []string{"ReferenceIndex", "Stock"}
		for i := range p.Tickers {
			p.Tickers[i].GraphqlType = types
		}
	}

	return nil
}

func (p *PolygonTickers) Save(client *dgo.Dgraph) error {
	var json = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "dgraph",
	}.Froze()

	err := p.WithTypeAndID(client)
	if err != nil {
		return err
	}

	txn := client.NewTxn()
	defer txn.Discard(context.Background())

	stocks, err := json.Marshal(p)
	fmt.Println(string(stocks))
	_, err = txn.Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: stocks})
	if err != nil {
		return err
	}

	return nil
}

//type PolygonQuote struct {
//	Ev  string  `json:"ev"`
//	Sym string  `json:"sym"`
//	V   int     `json:"v"`
//	Av  int     `json:"av"`
//	Op  float64 `json:"op"`
//	Vw  float64 `json:"vw"`
//	O   float64 `json:"o"`
//	C   float64 `json:"c"`
//	H   float64 `json:"h"`
//	L   float64 `json:"l"`
//	A   float64 `json:"a"`
//	S   int64   `json:"s"`
//	E   int64   `json:"e"`
//}
