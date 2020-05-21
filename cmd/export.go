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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dgraph-io/dgo/v2"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dgraphClient := NewDgraphClient(":9080")
		err := exportStocks(dgraphClient)
		if err != nil {
			fmt.Println(err)
		}
		err = exportQuotes(dgraphClient)
		if err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func exportQuotes(client *dgo.Dgraph) error {
	query := `
{
	Quotes(func: has(Quote.index)) {
    uid
    Quote.index {
      uid
      Stock.ticker
    }
    StockQuote.close
    Quote.datePublished
  }
}`
	return export(client, query, "quotes.json")
}

func exportStocks(client *dgo.Dgraph) error {
	query := `
{
	Stocks(func: has(Stock.ticker)) {
		uid
		Stock.ticker
  }
}`
	return export(client, query, "stocks.json")
}

func export(client *dgo.Dgraph, query string, file string) error {
	res, err := client.NewTxn().Query(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
	dst := &bytes.Buffer{}
	json.Indent(dst, res.GetJson(), "", "\t")
	return ioutil.WriteFile(file, dst.Bytes(), 0644)
}
