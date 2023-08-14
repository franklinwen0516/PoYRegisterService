package main

import (
	"net"
	"register_service/localutil"
	"register_service/protos"

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
		localutil.UserRegisterLog.Errorf("create network listener error: %s", err)
		panic(err)
	}
	defer lis.Close()
	s.Serve(lis)
}
