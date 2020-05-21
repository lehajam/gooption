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
	api_pb "github.com/lehajam/gooption/gobs/api"
	"github.com/machinebox/graphql"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// gobsPnlCmd represents the gobsPnl command
var gobsPnlCmd = &cobra.Command{
	Use:   "pnl",
	Short: "A brief description of your command",
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
		conn, gobs, err := NewGobsClient(gobsAddr)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		for range time.Tick(time.Duration(tick) * time.Second) {
			requests, err := getPnLRequests(client)
			if err != nil {
				log.Fatalf("Failed to send a price request: %s", err.Error())
				continue
			}

			var responses = make([]*api_pb.PnLResponse, len(requests))
			for i, req := range requests {
				responses[i], err = gobs.PnL(context.Background(), req)
				if err != nil {
					log.Fatalf("Failed to send a price request: %s", err.Error())
				}
			}

			savePnLResponses(client, responses)
		}
	},
}

func init() {
	gobsCmd.AddCommand(gobsPnlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gobsPnlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gobsPnlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getPnLRequests(client *graphql.Client) ([]*api_pb.PnLRequest, error) {
	query := `
{
  queryTrade {
    contract {
      ticker
      index {
        eodSpot: quotes(
          filter: { eod: true },
          order: {desc: datePublished}, first: 1) {          
          ... on StockQuote {
            close
          }
        }
        spot: quotes(order: {desc: datePublished}, first: 1) {          
          ... on StockQuote {
            close
          }
        }
      }
      eodPrice: results(
        filter: { eod: true and: { resultType: { eq: "price" }}},
        order: {desc: datePublished}, first: 1) {
        value
      }
      price: results(
        filter: { resultType: { eq: "price" } },
        order: {desc: datePublished}, first: 1) {
        value
      }
      delta: results(
        filter: { eod: true and: { resultType: { eq: "delta" }}},
        order: {desc: datePublished}, first: 1) {
        value
      }
      gamma: results(
        filter: { eod: true and: { resultType: { eq: "gamma" }}},
        order: {desc: datePublished}, first: 1) {
        value
      }
      theta: results(
        filter: { eod: true and: { resultType: { eq: "theta" }}},
        order: {desc: datePublished}, first: 1) {
        datePublished
        value
      }
    }
  }
}`
	type queryResponse struct {
		QueryTrade []struct {
			Contract struct {
				Ticker string `json:"ticker"`
				Index  []struct {
					EodSpot []struct {
						Close float64 `json:"close"`
					} `json:"eodSpot"`
					Spot []struct {
						Close float64 `json:"close"`
					} `json:"spot"`
				} `json:"index"`
				EodPrice []struct {
					Value float64 `json:"value"`
				} `json:"eodPrice"`
				Price []struct {
					Value float64 `json:"value"`
				} `json:"price"`
				Delta []struct {
					Value float64 `json:"value"`
				} `json:"delta"`
				Gamma []struct {
					Value float64 `json:"value"`
				} `json:"gamma"`
				Theta []struct {
					DatePublished time.Time `json:"datePublished"`
					Value         float64   `json:"value"`
				} `json:"theta"`
			} `json:"contract,omitempty"`
		} `json:"queryTrade"`
	}

	var response queryResponse
	err := RunGraphQLQuery(client, query, nil, &response)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var requests = make([]*api_pb.PnLRequest, len(response.QueryTrade))
	for i, trade := range response.QueryTrade {
		requests[i] = &api_pb.PnLRequest{
			ClientId: trade.Contract.Ticker,

			PricingDate:     float64(time.Now().Unix()),
			PreviousEodDate: float64(trade.Contract.Theta[0].DatePublished.Unix()),

			Spot:    trade.Contract.Index[0].Spot[0].Close,
			EodSpot: trade.Contract.Index[0].EodSpot[0].Close,

			Price:    trade.Contract.EodPrice[0].Value,
			EodPrice: trade.Contract.EodPrice[0].Value,
			EodDelta: trade.Contract.Delta[0].Value,
			EodGamma: trade.Contract.Gamma[0].Value,
			EodTheta: trade.Contract.Theta[0].Value,
		}
	}

	return requests, nil
}

func savePnLResponses(client *graphql.Client, pnlResponses []*api_pb.PnLResponse) {
	query := `
mutation ($timestamp: DateTime!, $contractID: String!, actual: Float!, expected: Float!, residual: Float!, $source: String!) {
  addPnLAttribution(input: [
    {datePublished:$timestamp, contract: { ticker: $contractID }, actual:$actual, expected:$expected, residual:$residual, source:"gobs"}
  ]) {
    pnlattribution {
      id
    }
  }
}`
	type queryResponse struct {
		AddPnLAttribution struct {
			PnLAttribution []struct {
				ID string `json:"id"`
			} `json:"pnlattribution"`
		} `json:"addPnLAttribution"`
	}

	for _, pnlResponse := range pnlResponses {
		vars := map[string]interface{}{
			"timestamp":  time.Now(),
			"contractID": pnlResponse.ClientId,
			"actual":     pnlResponse.Realized,
			"expected":   pnlResponse.Expected,
			"residual":   pnlResponse.Residual,
			"source":     "gobs",
		}

		var response queryResponse
		err := RunGraphQLQuery(client, query, vars, &response)
		if err != nil {
			log.Fatalf("Failed to insert pnl attribution for contract %s: %s", pnlResponse.ClientId, err.Error())
		}

		fmt.Printf("Added PnL attribution for contract %s\n", pnlResponse.ClientId)
	}
}
