package main

import (
	"log"
	"net"
	"user-service/logs"
	"user-service/protos/genuser"
	"user-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lg := logs.GetLogger("logs/logger.log")
	log.Println("Loggeeer bu", lg)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	server := service.Conn(lg)
	if server == nil {
		log.Fatalf("failed to initialize server")
	}

	grpcServer := grpc.NewServer()

	genuser.RegisterUserServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	log.Println("Server running on :8080")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
