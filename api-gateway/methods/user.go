package methods

import (
	"api-gateway/protos/genuser"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectUser() genuser.UserServiceClient {
	log.Println("Connecting to user service...")
	conn, err := grpc.NewClient("user-service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error connecting to user microservice:", err)
	}

	client := genuser.NewUserServiceClient(conn)
	return client
}
