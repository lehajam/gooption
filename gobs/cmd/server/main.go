package main

import (
	"os"

	"google.golang.org/grpc/grpclog"
)

func main() {
	err := Run()
	if err != nil {
		grpclog.Errorf("server was shutdown with errors: %v", err)
		os.Exit(1)
	}
}
