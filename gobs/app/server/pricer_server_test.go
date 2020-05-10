package server

import (
	"context"
	"encoding/json"
	"fmt"
	api_pb "gobs/api"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_pricerServiceServerImpl_Price(t *testing.T) {
	svr := NewPricerServiceServer()
	ctx := context.Background()

	priceReq := &api_pb.PriceRequest{
		Pricingdate: 1564669224,
		Strike:    159.76,
		Expiry:    1582890485,
		PutCall:   "put",
		Vol:  0.1,
		Rate: 0.01,
		Spot:  264.16,
	}

	resp, err := svr.Price(ctx, priceReq)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Price)

	jsonReq, _ := json.MarshalIndent(priceReq, "", "\t")
	fmt.Printf("%s \n", jsonReq)

	jsonResp, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Printf("%s \n", jsonResp)
}
