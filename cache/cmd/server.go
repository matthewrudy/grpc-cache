package main

import (
	"fmt"
	"net"

	"github.com/matthewrudy/grpc-cache/cache"
	proto "github.com/matthewrudy/grpc-cache/cache/proto"
	"google.golang.org/grpc"
)

func main() {
	runServer()
}

func runServer() error {
	srv := grpc.NewServer()
	service := cache.NewService()
	proto.RegisterCacheServer(srv, service)
	l, err := net.Listen("tcp", "localhost:5051")
	if err != nil {
		return err
	}

	fmt.Println("listening on localhost:5051")

	return srv.Serve(l)
}
