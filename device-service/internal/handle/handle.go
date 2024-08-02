package handle

import (
	"context"
	"device-service/protos/genuser"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



func HandleUser(id string) (string, string, error) {
	conn, err := grpc.NewClient("user-service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return "s", "s", fmt.Errorf("smth wron with connecting handleuser")

	}
	client := genuser.NewUserServiceClient(conn)

	res, err := client.GetById(context.Background(), &genuser.GetByIdReq{UserId: id})

	if err != nil {
		log.Println("Xatolik userdan user_id ob kelishda....", err)
		return "n", "n", err
	}
	log.Println("Name ushalvoldi... Handler userdamnnn", res.User.Profile.Name)
	return res.User.UserId, res.User.Profile.Name, nil 
}