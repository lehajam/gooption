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
	"github.com/machinebox/graphql"
	"github.com/sacOO7/gowebsocket"
	"github.com/spf13/cobra"
	"strings"
	"time"
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
		graphqlAddr, err := cmd.Flags().GetString("graphql")
		if err != nil {
			fmt.Println(err)
			return
		}
		client := graphql.NewClient(graphqlAddr)

		stocks := map[string]string{}
		var response QueryStockResponse
		err = RunGraphQLQuery(client, "{queryStock { id ticker }}", nil, &response)
		for _, stock := range response.QueryStock {
			stocks[stock.Ticker] = stock.ID
		}

		socket := gowebsocket.New("wss://socket.polygon.io/stocks") //{host}:{port})
		socket.Connect()

		socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
			if strings.Contains(string(message), "\"ev\":\"AM\"") {
				var quotes []PolygonQuote
				if err := json.Unmarshal([]byte(message), &quotes); err != nil {
					panic(err)
				}

				fmt.Printf("Received %d quotes\n", len(quotes))
				for _, quote := range quotes {
					if _, exist := stocks[quote.Sym]; !exist {
						var res AddStockResponse
						vars := map[string]interface{}{"ticker": quote.Sym}
						err = RunGraphQLQuery(client,
							`mutation ($ticker: String!) {
								addStock(input: [{ ticker: $ticker }]) {
									stock {
										id
									}
								}
							}`, vars, &res)
						if err != nil {
							fmt.Println(err.Error())
							continue
						}

						stocks[quote.Sym] = res.AddStock.Stock[0].ID
						fmt.Printf("Added stock %s\n", quote.Sym)
					}

					var res AddStockQuoteResponse
					vars := map[string]interface{}{"timestamp": time.Now(), "indexID": stocks[quote.Sym], "close": quote.C}
					err = RunGraphQLQuery(client,
						`mutation ($timestamp: DateTime!, $indexID: ID!, $close: Float!) {
							addStockQuote(input: [{ index: { id: $indexID }, datePublished: $timestamp, close: $close}]){
								stockquote {
								  id
								}
							  }
						}`,
						vars, &res)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					fmt.Printf("Added quote for stock %s\n", quote.Sym)
				}
			}
		}

		socket.SendText(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKEY))
		socket.SendText(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", CHANNELS))

		waitc := make(chan struct{})
		<-waitc
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

type AddStockResponse struct {
	AddStock struct {
		Stock []struct {
			ID string `json:"id"`
		} `json:"stock"`
	} `json:"addStock"`
}

type QueryStockResponse struct {
	QueryStock []struct {
		ID     string `json:"id"`
		Ticker string `json:"ticker"`
	} `json:"queryStock"`
}

type PolygonQuote struct {
	Ev  string  `json:"ev"`
	Sym string  `json:"sym"`
	V   int     `json:"v"`
	Av  int     `json:"av"`
	Op  float64 `json:"op"`
	Vw  float64 `json:"vw"`
	O   float64 `json:"o"`
	C   float64 `json:"c"`
	H   float64 `json:"h"`
	L   float64 `json:"l"`
	A   float64 `json:"a"`
	S   int64   `json:"s"`
	E   int64   `json:"e"`
}

type AddStockQuoteResponse struct {
	AddStockQuote struct {
		Stockquote []struct {
			ID string `json:"id"`
		} `json:"stockquote"`
	} `json:"addStockQuote"`
}
