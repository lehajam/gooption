package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	jsoniter "github.com/json-iterator/go"
	api_pb "github.com/lehajam/gooption/gobs/api"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
)

func ParseFromResponse(response *http.Response, obj interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}

	return nil
}

func ParseFromFile(file string, obj interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonFile, obj)
}

func NewDgraphClient(addr string) *dgo.Dgraph {
	d, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func Query(client *dgo.Dgraph, query string, index string, obj interface{}) (count uint64, err error) {
	queryRes, err := client.NewTxn().Query(context.Background(), query)
	if err != nil {
		return 0, err
	}

	if queryRes.Metrics.NumUids[index] > 0 {
		err = json.Unmarshal(queryRes.GetJson(), obj)
		if err != nil {
			return 0, err
		}
	}

	return queryRes.Metrics.NumUids[index], nil
}

func Save(client *dgo.Dgraph, obj ...interface{}) error {
	if len(obj) == 0 {
		return fmt.Errorf("nothing to save")
	}

	txn := client.NewTxn()
	defer txn.Discard(context.Background())

	json := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "dgraph",
	}.Froze()

	commitNow := len(obj) == 1
	for _, o := range obj {
		toSave, err := json.Marshal(o)
		if err != nil {
			return err
		}

		_, err = txn.Mutate(context.Background(), &api.Mutation{CommitNow: commitNow, SetJson: toSave})
		if err != nil {
			return err
		}
	}

	if !commitNow {
		return txn.Commit(context.Background())
	}

	return nil
}

func NewGobsPriceStream(grpcAddress string) (api_pb.PricerService_PriceClient, *grpc.ClientConn, error) {
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
