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
	"encoding/json"
	"fmt"
	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dgraphClient := NewDgraphClient(":9082")
		err := importStocks(dgraphClient)
		if err != nil {
			fmt.Println(err)
		}
		err = importQuotes(dgraphClient)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func importQuotes(client *dgo.Dgraph) error {
	var dump QuoteFile
	jsonFile, _ := ioutil.ReadFile("quotes.json")
	err := json.Unmarshal(jsonFile, &dump)
	if err != nil {
		return err
	}

	txn := client.NewTxn()
	defer txn.Discard(context.Background())
	for _, quote := range dump.Quotes {
		pb, err := json.Marshal(quote)
		if err != nil {
			return err
		}

		resp, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: false, SetJson: pb})
		if err != nil {
			return err
		}

		fmt.Println(string(resp.GetJson()))
	}

	return txn.Commit(context.Background())
}

func importStocks(client *dgo.Dgraph) error {
	var dump QuoteFile
	jsonFile, _ := ioutil.ReadFile("quotes.json")
	err := json.Unmarshal(jsonFile, &dump)
	if err != nil {
		return err
	}

	txn := client.NewTxn()
	defer txn.Discard(context.Background())
	stocks := map[string]bool{}
	for _, quote := range dump.Quotes {
		if _, exist := stocks[quote.Index.Ticker]; !exist {

			fmt.Printf("adding stock %s\n", quote.Index.Ticker)
			pb, err := json.Marshal(quote.Index)
			if err != nil {
				return err
			}

			resp, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: false, SetJson: pb})
			if err != nil {
				return err
			}
			fmt.Println(string(resp.GetJson()))
			stocks[quote.Index.Ticker] = true
		}
	}

	return txn.Commit(context.Background())
}

type Stock struct {
	ID     string `json:"uid"`
	Ticker string `json:"Stock.ticker"`
}

type Quote struct {
	Index         Stock     `json:"Quote.index"`
	Close         float64   `json:"StockQuote.close"`
	DatePublished time.Time `json:"Quote.datePublished"`
}

type QuoteFile struct {
	Quotes []Quote `json:"Quotes"`
}
