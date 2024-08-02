package methods

import (
	"api-gateway/protos/gendevice"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectDevice() gendevice.DeviceServerClient {
	conn, err := grpc.NewClient("device-service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connect user micro...", err)
	}

	client := gendevice.NewDeviceServerClient(conn)
	return client
}
