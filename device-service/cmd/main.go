package main

import (
	"device-service/protos/gendevice"
	"device-service/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	les, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("smth wrong on portttt: %v", err)
	}
	defer les.Close()
	// rabbitmq.ConsumeOrders()

	server := service.Connect()
	if server == nil {
		log.Fatalf("wrong to initialize server")
	}

	grpcServer := grpc.NewServer()

	gendevice.RegisterDeviceServerServer(grpcServer, server)

	reflection.Register(grpcServer)

	log.Println("Server running on :8080")

	if err := grpcServer.Serve(les); err != nil {
		log.Fatalf("wrronngg to listennn: %v", err)
	}

}
