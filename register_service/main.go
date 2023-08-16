package main

import (
	"log"
	"net"
	"register_service/db"
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
	err := db.UserRegisterInfoInstance.Init()
	if err != nil {
		log.Fatalf("connect database error: %s", err)
		panic(err)
	}
	localutil.LoggerInit()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		localutil.UserRegisterLog.Errorf("create network listener error: %s", err)
		panic(err)
	}
	defer lis.Close()
	s.Serve(lis)
}
