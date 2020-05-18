package main

import (
	"github.com/srvc/appctx"

	"github.com/izumin5210/grapi/pkg/grapiserver"

	"github.com/lehajam/gooption/gobs/app/server"
)

func Run() error {
	// Application context
	ctx := appctx.Global()

	s := grapiserver.New(
		grapiserver.WithGrpcAddr("tcp", ":5050"),
		grapiserver.WithGatewayAddr("tcp", ":4200"),
		grapiserver.WithDefaultLogger(),
		grapiserver.WithServers(
			server.NewPricerServiceServer(),
		),
	)
	return s.Serve(ctx)
}
