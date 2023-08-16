package main

import (
	"context"
	"log"
	"os"
	"register_client/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getFacialImages() [][]byte {
	imageData1, err := os.ReadFile("./source/images/Clara_Harris_0001.jpg")
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}
	imageData2, err := os.ReadFile("./source/images/Clara_Harris_0002.jpg")
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}
	imageData3, err := os.ReadFile("./source/images/Clara_Harris_0003.jpg")
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}
	imageData4, err := os.ReadFile("./source/images/Clara_Harris_0004.jpg")
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}
	imageData5, err := os.ReadFile("./source/images/Clara_Harris_0005.jpg")
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}
	return [][]byte{
		[]byte(imageData1),
		[]byte(imageData2),
		[]byte(imageData3),
		[]byte(imageData4),
		[]byte(imageData5),
	}
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := protos.NewRegisterServiceClient(conn)
	request := &protos.BioRegisterRequset{
		// 设置BioRegisterRequset请求参数
		AccountPublicKey: "0x00001234567890",
		FacialImages:     getFacialImages(),
	}
	response, err := client.RegisterWithBioKey(context.Background(), request)
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}
	log.Print(response.String())
	request = &protos.BioRegisterRequset{
		// 设置BioRegisterRequset请求参数
		AccountPublicKey: "0x9876543210000",
		FacialImages:     getFacialImages(),
	}
	response, err = client.RegisterWithBioKey(context.Background(), request)
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}
	log.Print(response.String())
}
