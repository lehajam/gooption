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
	"fmt"
	"io"
	"log"
	"time"

	api_pb "github.com/lehajam/gooption/gobs/api"
	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
)

// gobsCmd represents the gobs command
var gobsCmd = &cobra.Command{
	Use:   "gobs",
	Short: "Price trades using gobs and save results every tick",
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

		gobsAddr, err := cmd.Flags().GetString("gobs")
		if err != nil {
			fmt.Println(err)
			return
		}

		tick, err := cmd.Flags().GetInt64("tick")
		if err != nil {
			fmt.Println(err)
			return
		}

		client := graphql.NewClient(graphqlAddr)
		for range time.Tick(time.Duration(tick) * time.Second) {
			var queryTradeResp QueryTradeResponse
			err := RunGraphQLQuery(client, `{
				queryTrade {
					contract {
						ticker
						index {
							quotes (order: { desc: datePublished }, first: 1) {
								... on StockQuote {
									datePublished
									close
								}
							}
						}
						... on EuropeanContract {
							strike
							expiry
							putcall
						}
					}
				}
			}`, nil, &queryTradeResp)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			stream, conn, err := NewGobsStream(gobsAddr)
			if conn != nil {
				defer conn.Close()
			}
			if err != nil {
				fmt.Println(err)
				return
			}

			waitc := make(chan struct{})
			go gobsResponseHandler(client, waitc, stream)

			for _, trade := range queryTradeResp.QueryTrade {
				err = stream.Send(&api_pb.PriceRequest{
					ClientId:    trade.Contract.Ticker,
					Pricingdate: float64(time.Now().Unix()),
					Strike:      trade.Contract.Strike,
					PutCall:     trade.Contract.Putcall,
					Expiry:      float64(trade.Contract.Expiry.Unix()),
					Spot:        trade.Contract.Index[0].Quotes[0].Close,
					Vol:         0.1,
					Rate:        0.01,
				})
				if err != nil {
					log.Fatalf("Failed to send a price request: %v", err)
				}
			}

			stream.CloseSend()
			<-waitc
		}
	},
}

func init() {
	rootCmd.AddCommand(gobsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gobsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	gobsCmd.Flags().String("graphql", "http://localhost:8080/graphql", "graphql server URL")
	gobsCmd.Flags().String("gobs", ":5050", "gobs URL")
	gobsCmd.Flags().Int64("tick", 10, "the time interval in seconds between")
}

type QueryTradeResponse struct {
	QueryTrade []struct {
		Contract struct {
			Ticker string `json:"ticker"`
			Index  []struct {
				Quotes []struct {
					DatePublished time.Time `json:"datePublished"`
					Close         float64   `json:"close"`
				} `json:"quotes"`
			} `json:"index"`
			Strike  float64   `json:"strike"`
			Expiry  time.Time `json:"expiry"`
			Putcall string    `json:"putcall"`
		} `json:"contract"`
	} `json:"queryTrade"`
}

type PriceResultMutaionResponse struct {
	AddPriceResult struct {
		Priceresult []struct {
			ID string `json:"id"`
		} `json:"priceresult"`
	} `json:"addPriceResult"`
}

func gobsResponseHandler(client *graphql.Client, waitc chan struct{}, stream api_pb.PricerService_PriceClient) {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(waitc)
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a result : %v", err)
		}

		var mutationResp PriceResultMutaionResponse
		vars := map[string]interface{}{"timestamp": time.Now(), "contractID": resp.ClientId, "value": resp.Value, "type": resp.ValueType, "source": "gobs"}
		err = RunGraphQLQuery(client, `mutation ($timestamp: DateTime!, $contractID: String!, $value: Float!, $type: String!, $source: String!) {
			addPriceResult(input: [{
				datePublished: $timestamp,
				contract: { ticker: $contractID },
				value: $value,
				resultType: $type,
				source: $source
			}]) {
				priceresult {
					id
					value
				}
			}
		}`, vars, &mutationResp)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("Added %s for contract %s\n", resp.ValueType, resp.ClientId)
	}
}
