package main

import (
	"RegisterService/protos"
	"log"
	"net"

	"google.golang.org/grpc"
)

type RegisterServiceImpl struct {
	protos.UnsafeRegisterServiceServer
}

func main() {
	s := grpc.NewServer()
	protos.RegisterRegisterServiceServer(s, &RegisterServiceImpl{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("create network listener error: %s", err)
		panic(err)
	}
	defer lis.Close()
	s.Serve(lis)
}
