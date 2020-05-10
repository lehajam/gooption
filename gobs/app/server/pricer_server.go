package server

import (
	"io"
	"math"
	"strings"
	"time"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"gonum.org/v1/gonum/stat/distuv"

	api_pb "gobs/api"
)

const (
	ivSeed    = 0.1 // solver starting point
	maxIter   = 1000
	putLBound = 0.20
)

var (
	phi             = distuv.Normal{Mu: 0, Sigma: 1}.CDF
	dphi            = distuv.Normal{Mu: 0, Sigma: 1}.Prob
	mapToMultiplier = map[string]float64{"call": 1.0, "put": -1.0}
	allGreeks       = []string{"delta", "gamma", "vega", "theta", "rho"}
)

// PricerServiceServer is a composite interface of api_pb.PricerServiceServer and grapiserver.Server.
type PricerServiceServer interface {
	api_pb.PricerServiceServer
	grapiserver.Server
}

// NewPricerServiceServer creates a new PricerServiceServer instance.
func NewPricerServiceServer() PricerServiceServer {
	return &pricerServiceServerImpl{}
}

type pricerServiceServerImpl struct {
}

/*
Price computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func (srv *pricerServiceServerImpl) Price(stream api_pb.PricerService_PriceServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// get inputs
		s, v, k, r := req.Spot, req.Vol, req.Strike, req.Rate
		t := time.Unix(int64(req.Expiry), 0).Sub(time.Unix(int64(req.Pricingdate), 0)).Hours() / 24.0 / 365.250
		mult := mapToMultiplier[strings.ToLower(req.PutCall)]

		d1 := d1(s, k, t, v, r)
		d2 := d2(d1, v, t)

		results := map[string]float64{
			"price": bs(s, v, r, k, t, mult),
			"delta": delta(d1, mult),
			"vega":  vega(s, t, d1),
			"gamma": gamma(s, t, v, d1),
			"rho":   rho(k, t, r, d2, mult),
		}

		for valueType, value := range results {
			if err := stream.Send(&api_pb.PriceResponse{Value: value, ValueType: valueType}); err != nil {
				return err
			}
		}
	}
}

/*
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func bs(s, v, r, k, t, mult float64) float64 {
	d1 := d1(s, k, t, v, r)
	d2 := d2(d1, v, t)

	return mult * (s*phi(mult*d1) - k*phi(mult*d2)*math.Exp(-r*t))
}

func d1(S, K, T, Sigma, R float64) float64 {
	return (1.0 / Sigma * math.Sqrt(T)) * (math.Log(S/K) + (R+Sigma*Sigma*0.5)*T)
}

func d2(d1, Sigma, T float64) float64 {
	return d1 - Sigma*math.Sqrt(T)
}

func delta(d1, mult float64) float64 {
	return mult * phi(mult*d1)
}

func gamma(s, t, sigma, d1 float64) float64 {
	return dphi(d1) / (s * sigma * math.Sqrt(t))
}

func vega(s, t, d1 float64) float64 {
	return s * dphi(d1) * math.Sqrt(t)
}

func theta(s, k, t, sigma, r, d1, d2, mult float64) float64 {
	return -0.5*(s*dphi(d1)*sigma/math.Sqrt(t)) - (mult * r * k * math.Exp(-r*t) * phi(mult*d2))
}

func rho(k, t, r, d2, mult float64) float64 {
	return mult * k * t * math.Exp(-r*t) * phi(mult*d2)
}
