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
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/spf13/cobra"
	"time"
)

// flashCmd represents the flash command
var eodCmd = &cobra.Command{
	Use:   "eod",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dgraphClient := NewDgraphClient(":9080")
		err := flash(dgraphClient)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(eodCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flashCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	eodCmd.Flags().String("datetime", time.Now().String(), "The datetime to use to flash P&L")
}

func flash(client *dgo.Dgraph) error {
	var (
		query = `
query eod($eod: string) {
    var(func: type(Contract)) {
      EodResults as Contract.results @filter(le(PriceResult.datePublished, $eod)) (orderdesc: Quote.datePublished) {
        expand(_all_)
      }
      Contract.index {
        EodQuotes as Index.quotes @filter(le(Quote.datePublished, $eod)) (orderdesc:Quote.datePublished) (first:1) {
          expand(_all_)
        }
      }
    }

    quotes(func: type(Contract)) {
      Contract.index {
        Index.quotes @filter(uid(EodQuotes)) {
			uid
        }
      }
      price: Contract.results @filter(uid(EodResults) and eq(PriceResult.resultType, "price")) (first: 1) {
		uid
      }

      delta: Contract.results @filter(uid(EodResults) and eq(PriceResult.resultType, "delta")) (first: 1) {
		uid
      }

      gamma: Contract.results @filter(uid(EodResults) and eq(PriceResult.resultType, "gamma")) (first: 1) {
		uid
      }

      theta: Contract.results @filter(uid(EodResults) and eq(PriceResult.resultType, "theta")) (first: 1) {
		uid
      }
    }
}`
		response struct {
			EODFlash []struct {
				ContractIndex []struct {
					IndexQuotes []struct {
						UID string `json:"uid"`
						EOD bool   `json:"eod"`
					} `json:"Index.quotes"`
				} `json:"Contract.index"`
				Price []struct {
					UID string `json:"uid"`
					EOD bool   `json:"eod"`
				} `json:"price,omitempty"`
				Delta []struct {
					UID string `json:"uid"`
					EOD bool   `json:"eod"`
				} `json:"delta,omitempty"`
				Gamma []struct {
					UID string `json:"uid"`
					EOD bool   `json:"eod"`
				} `json:"gamma,omitempty"`
			} `json:"quotes"`
		}
	)

	queryResult, err := client.NewTxn().QueryWithVars(context.Background(), query, map[string]string{"$eod": "2020-05-18T19:00:00"})
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(queryResult.GetJson(), &response)
	if err != nil {
		return err
	}

	for _, obj := range response.EODFlash {
		for _, idx := range obj.ContractIndex {
			for _, q := range idx.IndexQuotes {
				q.EOD = true
			}
		}
		for _, p := range obj.Price {
			p.EOD = true
		}
		for _, d := range obj.Delta {
			d.EOD = true
		}
		for _, g := range obj.Gamma {
			g.EOD = true
		}
	}

	pb, err := json.Marshal(response.EODFlash)
	if err != nil {
		return err
	}

	fmt.Println(string(pb))

	txn := client.NewTxn()
	defer txn.Discard(context.Background())
	_, err = txn.Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: pb})
	if err != nil {
		return err
	}
	return nil
}
