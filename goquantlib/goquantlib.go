//go:generate sh -c "protoc --proto_path=pb --proto_path=$GOPATH/src/github.com/lehajam/gooption/pb --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogofast_out=plugins=grpc:pb $GOPATH/src/github.com/lehajam/gooption/pb/*.proto pb/*.proto"
//go:generate sh -c "protoc --proto_path=pb --proto_path=$GOPATH/src/github.com/lehajam/gooption/pb --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:pb pb/*.proto"
// -- //go:generate sh -c "protoc --proto_path=pb --proto_path=$GOPATH/src/github.com/lehajam/gooption/pb --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogofast_out=plugins=grpc:pb --grpc-gateway_out=logtostderr=true:pb $GOPATH/src/github.com/lehajam/gooption/pb/*.proto pb/*.proto"
//go:generate gooption-cli -p gobs -r Price -r Greek -r ImpliedVol
package main

import (
	"flag"
	"net/http"
	"strings"
	"sync"

	context "golang.org/x/net/context"

	"errors"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/lehajam/gooption/gobs/pb"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

var (
	tcpPort      = ":50051"
	httpPort     = ":8081"
	gobsEndpoint = flag.String(
		"gobs_endpoint",
		"localhost:50051",
		"endpoint of YourService")
)

// server is used to implement pb.ModerlServer.
type server struct{}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func httpServer() error {
	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGobsHandlerFromEndpoint(ctx, mux, *gobsEndpoint, opts)
	if err != nil {
		return err
	}

	logrus.Infoln("http server ready on port ", httpPort)
	return http.ListenAndServe(httpPort, cors.Default().Handler(mux))
}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func tcpServer() error {
	lis, err := net.Listen("tcp", tcpPort)
	if err != nil {
		return err
	}
	defer lis.Close()

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(methodFullName string, err error) bool {
			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if err == nil && methodFullName == "main.server.healthcheck" {
				return false
			}

			// by default you will log all calls
			return true
		}),
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()), opts...)))
	pb.RegisterGobsServer(s, &server{})
	reflection.Register(s)

	logrus.Infoln("grpc server ready on port ", tcpPort)
	return s.Serve(lis)
}

func start(entrypoint func() error) {
	defer glog.Flush()
	if err := entrypoint(); err != nil {
		logrus.Fatal(err)
		panic(err)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go start(tcpServer)
	go start(httpServer)
	wg.Wait()
}

/*
Price computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func (srv *server) Price(ctx context.Context, in *pb.PriceRequest) (*pb.PriceResponse, error) {
	var (
		mult = putCallMap[strings.ToLower(in.Contract.Putcall)]

		s = in.Marketdata.Spot.Index.Value
		v = in.Marketdata.Vol.Index.Value
		r = in.Marketdata.Rate.Index.Value
		k = in.Contract.Strike
		p = EuropeanFlatVol(
			s, 
			r, 
			0, 
			v, 
			k, 
			in.Pricingdate, 
			in.Contract.Expiry, 
			mult	
	)


	return &pb.PriceResponse{
		Price: p,
	}, nil
}

/*
Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func (srv *server) Greek(ctx context.Context, in *pb.GreekRequest) (*pb.GreekResponse, error) {
	var (
		mult = putCallMap[strings.ToLower(in.Request.Contract.Putcall)]

		s = in.Request.Marketdata.Spot.Index.Value
		v = in.Request.Marketdata.Vol.Index.Value
		r = in.Request.Marketdata.Rate.Index.Value
		k = in.Request.Contract.Strike
		t = time.Unix(int64(in.Request.Contract.Expiry), 0).Sub(
			time.Unix(int64(in.Request.Pricingdate), 0)).Hours() / 24.0 / 365.250
		d1 = d1(s, k, t, v, r)
		d2 = d2(d1, v, t)
	)

	logrus.Debugf("%+v\n", proto.MarshalTextString(in))
	return &pb.GreekResponse{}, errors.New("Not Implemented")
}

/*
ImpliedVol computes volatility matching the option quote passed as Quote using Newton Raphson solver
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
The second argument returned is the number of iteration used to converge
*/
func (srv *server) ImpliedVol(ctx context.Context, in *pb.ImpliedVolRequest) (*pb.ImpliedVolResponse, error) {
	var (
		mult = func(q *pb.OptionQuote) float64 { return putCallMap[strings.ToLower(q.Putcall)] }

		s = in.Marketdata.Spot.Index.Value
		r = in.Marketdata.Rate.Index.Value
		k = func(q *pb.OptionQuote) float64 { return q.Strike }
		p = func(q *pb.OptionQuote) float64 { return (q.Ask + q.Bid) / 2.0 }
		t = func(q *pb.OptionQuoteSlice) float64 {
			return time.Unix(int64(q.Expiry), 0).Sub(time.Unix(int64(in.Pricingdate), 0)).Hours() / 24.0 / 365.250
		}
	)

	logrus.Debugf("%+v\n", proto.MarshalTextString(in))
	return &pb.ImpliedVolResponse{}, errors.New("Not Implemented")
}
