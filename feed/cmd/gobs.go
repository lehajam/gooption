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
	"io"
	"log"
	"time"

	api_pb "github.com/lehajam/gooption/gobs/api"
	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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
		feed()
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
	// gobsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func feed() {

	client := graphql.NewClient("http://localhost:8080/graphql")

	// make a request
	req := graphql.NewRequest(`{
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
	}`)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var resp QueryTradeResponse
	err := client.Run(ctx, req, &resp)
	if err != nil {
		fmt.Println(err.Error())
	}

	// jreq, _ := json.MarshalIndent(resp, "", "\t")
	// fmt.Printf("%s \n", jreq)

	conn, err := grpc.Dial(":5050", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	grpcClient := api_pb.NewPricerServiceClient(conn)
	stream, err := grpcClient.Price(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}

			// jreq, _ = json.MarshalIndent(resp, "", "\t")
			// fmt.Printf("%s \n", jreq)

			// make a request
			req = graphql.NewRequest(`mutation ($timestamp: DateTime!, $contractID: String!, $value: Float!, $type: String!, $source: String!) {
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
			}`)

			// set any variables
			req.Var("timestamp", time.Now())
			req.Var("contractID", resp.ClientId)
			req.Var("value", resp.Value)
			req.Var("type", resp.ValueType)
			req.Var("source", "gobs")

			// set header fields
			req.Header.Set("Cache-Control", "no-cache")

			// run it and capture the response
			obj := &PriceResultMutaionResponse{}
			err = client.Run(ctx, req, obj)
			if err != nil {
				fmt.Println(err.Error())
			}

			// jreq, _ = json.MarshalIndent(obj, "", "\t")
			// fmt.Printf("%s \n", jreq)
		}
	}()

	for _, trade := range resp.QueryTrade {
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
