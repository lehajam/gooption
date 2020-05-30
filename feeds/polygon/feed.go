package polygon

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v2"
	"github.com/lehajam/gooption/core"
	"github.com/sacOO7/gowebsocket"
	"net/http"
	"strings"
	"time"
)

const APIKEY = "tib3RQoc2VasTZrVojO_yZ4_VNsm__n8gIeSzI"

type PolygonStockFeed struct {
	client  *dgo.Dgraph
	Tickers []PolygonStock `json:"tickers"`
}

func NewPolygonStockFeed(addr string) PolygonStockFeed {
	return PolygonStockFeed{client: core.NewDgraphClient(addr)}
}

func (f PolygonStockFeed) Fetch() error {
	response, err := http.Get(fmt.Sprintf("https://api.polygon.io/v2/reference/tickers?apiKey=%s", APIKEY))
	if err != nil {
		return err
	}

	err = core.ParseFromResponse(response, &f)
	if err != nil {
		return err
	}

	err = queryStocks(f.client, f)
	if err != nil {
		return err
	}

	return core.Save(f.client, f.Tickers)
}

func queryStocks(client *dgo.Dgraph, feed PolygonStockFeed) error {
	query := `
{
  tickers(func: type(Stock)) @filter(eq(ReferenceIndex.symbol, %s)) {
    uid
    ReferenceIndex.symbol
	dgraph.type
  }
}`

	symbols := make([]string, len(feed.Tickers))
	for i, s := range feed.Tickers {
		symbols[i] = s.Ticker
	}

	query = fmt.Sprintf(query, symbols)
	count, err := core.Query(client, query, "uid", &feed)
	if err != nil {
		return err
	}

	if count != uint64(len(feed.Tickers)) {
		types := []string{"ReferenceIndex", "Stock", "PolygonStock"}
		for i := range feed.Tickers {
			feed.Tickers[i].Types = types
		}
	}

	return nil
}

type PolygonStockQuoteStream struct {
	client *dgo.Dgraph
}

func (s PolygonStockQuoteStream) Stream() error {
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
			err = core.Save(s.client, quotes)
			if err != nil {
				// LOG ERROR
			}
		}
	}

	// 1 - Get Stocks from DB
	var symbols []string
	err := querySymbols(s.client, symbols)
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

func querySymbols(client *dgo.Dgraph, symbols []string) error {
	query := `
{
  tickers(func: type(Stock)) {
    ReferenceIndex.symbol
  }
}`

	_, err := core.Query(client, query, "ReferenceIndex.symbol", &symbols)
	if err != nil {
		return err
	}

	return nil
}
