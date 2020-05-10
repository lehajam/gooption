package main

import (
	"github.com/srvc/appctx"

	"github.com/izumin5210/grapi/pkg/grapiserver"

	"gobs/app/server"
)

func run() error {
	// Application context
	ctx := appctx.Global()

	s := grapiserver.New(
		grapiserver.WithDefaultLogger(),
		grapiserver.WithServers(
			server.NewPricerServiceServer(),
		),
	)
	return s.Serve(ctx)
}
