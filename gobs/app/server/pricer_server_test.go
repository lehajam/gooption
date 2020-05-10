package server

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	api_pb "gobs/api"
//	"github.com/stretchr/testify/require"
//	"io"
//	"log"
//	"testing"
//)

//func Test_pricerServiceServerImpl_Price(t *testing.T) {
//	stream, err := client.RouteChat(context.Background())
//	waitc := make(chan struct{})
//	go func() {
//		for {
//			in, err := stream.Recv()
//			if err == io.EOF {
//				// read done.
//				close(waitc)
//				return
//			}
//			if err != nil {
//				log.Fatalf("Failed to receive a note : %v", err)
//			}
//			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
//		}
//	}()
//	for _, note := range notes {
//		if err := svr.Price.Send(note); err != nil {
//			log.Fatalf("Failed to send a note: %v", err)
//		}
//	}
//	stream.CloseSend()
//	<-waitc
//
//
//
//	svr := NewPricerServiceServer()
//	ctx := context.Background()
//
//	priceReq := &api_pb.PriceRequest{
//		Pricingdate: 1564669224,
//		Strike:    159.76,
//		Expiry:    1582890485,
//		PutCall:   "put",
//		Vol:  0.1,
//		Rate: 0.01,
//		Spot:  264.16,
//	}
//
//	resp, err := svr.Price(ctx, priceReq)
//
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//	require.NotEmpty(t, resp.Price)
//
//	jsonReq, _ := json.MarshalIndent(priceReq, "", "\t")
//	fmt.Printf("%s \n", jsonReq)
//
//	jsonResp, _ := json.MarshalIndent(resp, "", "\t")
//	fmt.Printf("%s \n", jsonResp)
//}
