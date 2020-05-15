package cmd

import (
	"context"

	api_pb "github.com/lehajam/gooption/gobs/api"
	"github.com/machinebox/graphql"
	"google.golang.org/grpc"
)

func RunGraphQLQuery(client *graphql.Client, queryString string, vars map[string]interface{}, response interface{}) error {
	// make a request
	req := graphql.NewRequest(queryString)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	if vars != nil {
		for k, v := range vars {
			req.Var(k, v)
		}
	}

	// run it and capture the response
	err := client.Run(context.Background(), req, response)
	if err != nil {
		return err
	}

	// jreq, _ := json.MarshalIndent(response, "", "\t")
	// fmt.Printf("%s \n", jreq)

	return nil
}

func NewGobsStream(grpcAddress string) (api_pb.PricerService_PriceClient, *grpc.ClientConn, error) {
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
