package gobs

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/dgraph-io/dgo/v2"
	"github.com/lehajam/gooption/core"
	api_pb "github.com/lehajam/gooption/gobs/api"
	"google.golang.org/grpc"
)

type GobsFeed struct {
	tick   time.Duration
	client *dgo.Dgraph
}

func (f GobsFeed) Fetch() error {
	var trades []Trade
	err := queryTrades(f.client, trades)
	if err != nil {
		return err
	}

	stream, conn, err := newGobsStream("")
	if err != nil {
		return err
	}
	defer conn.Close()

	waitc := make(chan struct{})
	go gobsHandler(f.client, stream, waitc)

	pricingDate := float64(time.Now().Unix())
	for _, trade := range trades {
		err = stream.Send(&api_pb.PriceRequest{
			ClientId:    trade.ID,
			Pricingdate: pricingDate,
			Strike:      trade.Contract.Strike,
			PutCall:     trade.Contract.Putcall,
			Expiry:      float64(trade.Contract.Expiry.Unix()),
			Spot:        trade.Contract.RefIndex[0].Quotes[0].Last,
			Vol:         0.1,
			Rate:        0.01,
		})
		if err != nil {
			fmt.Printf("Failed to send a price request: %s", err.Error())
		}
	}

	stream.CloseSend()
	<-waitc

	return nil
}

func (f GobsFeed) Stream() error {
	for range time.Tick(f.tick * time.Second) {
		err := f.Fetch()
		if err != nil {
			fmt.Printf("Failed to send a price request: %s", err.Error())
		}
	}
	return nil
}

func gobsHandler(client *dgo.Dgraph, stream api_pb.PricerService_PriceClient, waitc chan struct{}) {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			close(waitc) // read done
			return
		}

		if err != nil {
			fmt.Printf("Failed to receive a result : %v\n", err)
			continue
		}

		trade := TradeWithResults{
			GraphQL: core.GraphQL{
				ID: res.ClientId,
			},
			Valuations: []ValuationResult{
				{
					GraphQL: core.GraphQL{
						Types: []string{"ValuationResult", res.ValueType},
					},
					DatePublished: time.Now(), // replace with pricing date
					Value:         res.Value,
					Source:        "gobs",
				},
			},
		}

		err = core.Save(client, &trade)
		if err != nil {
			fmt.Printf("Failed to save result : %v\n", err)
		}
	}
}

func queryTrades(client *dgo.Dgraph, trades []Trade) error {
	query := `
{
  trades(func: type(Trade)) {
    uid
    Trade.quantity
    Trade.contract @filter(type(Option) and eq(Option.optionType, "european")) {
      Contract.refIndex {
        ReferenceIndex.quotes (orderdesc: Quote.datePublished, first:1) {
          Quote.datePublished
          Quote.last
        }
      }
      Option.strike
      Option.expiry
      Option.putcall
    }
  }
}`

	_, err := core.Query(client, query, "uid", &trades)
	if err != nil {
		return err
	}

	return nil
}

func newGobsStream(grpcAddress string) (api_pb.PricerService_PriceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	grpcClient := api_pb.NewPricerServiceClient(conn)
	stream, err := grpcClient.Price(context.Background())
	if err != nil {
		return nil, conn, err
	}

	return stream, conn, nil

}
