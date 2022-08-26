package main

import (
	"google.golang.org/grpc"
	"hezzl/grpcserver"
	"hezzl/protogrpc"
	"log"
	"net"
)

func main() {

	s := grpc.NewServer()
	srv := &grpcserver.GRPCServer{}

	protogrpc.RegisterUsersAdminServer(s, srv)

	grpcserver.InitRedisConnection()

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
