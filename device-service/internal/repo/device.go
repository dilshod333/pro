// package repo

// import (
// 	"context"
// 	"device-service/internal/handle"
// 	"device-service/models"
// 	"device-service/pkg"
// 	"device-service/protos/gendevice"
// 	"fmt"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// )

// type NewServ struct {
// 	serv *pkg.Server
// }

// func Conn() *NewServ {
// 	server, err := pkg.NewServer()
// 	if err != nil {
// 		log.Fatal("Error on newserver.. on repooo", err)

// 	}

// 	return &NewServ{serv: server}
// }

// // func (s *NewServ) DeviceCreate(req *gendevice.CreateDeviceRequest) (*gendevice.CreateDeviceResp, error) {
// // 	ctx := context.Background()
// // 	userid, name, err := handle.HandleUser(req.UserId)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("xatolik userhandleda yoki id topilmadi %+v", err)

// // 	}
// // 	device := bson.D{
// // 		{Key: "user_id", Value: userid},
// // 		{Key: "name", Value: name},
// // 		{Key: "device_type", Value: req.DeviceType},
// // 		{Key: "device_name", Value: req.DeviceName},
// // 		{Key: "device_status", Value: req.DeviceStatus},
// // 		{Key: "configuration_settings", Value: req.ConfigurationSettings},
// // 		{Key: "last_updated", Value: req.LastUpdated},
// // 		{Key: "location", Value: req.Location},
// // 	}

// // 	deviceID, err := s.serv.CreateDevicee(ctx, device)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	resp := &gendevice.CreateDeviceResp{
// // 		DeviceId: deviceID.Hex(),
// // 		Device:   req,
// // 	}

// // 	return resp, nil
// // }

// func (s *NewServ) DeviceCreate(req *gendevice.CreateDeviceRequest) (*gendevice.CreateDeviceResp, error) {
// 	ctx := context.Background()
// 	userid, name, err := handle.HandleUser(req.UserId)
// 	if err != nil {
// 		return nil, fmt.Errorf("error handling user or user ID not found: %+v", err)
// 	}

// 	device := models.Device{
// 		UserID:                userid,
// 		Name:                  name,
// 		DeviceType:            req.DeviceType,
// 		DeviceName:            req.DeviceName,
// 		DeviceStatus:          req.DeviceStatus,
// 		ConfigurationSettings: req.ConfigurationSettings,
// 		LastUpdated:           time.Now().Format(time.RFC3339),
// 		Location:              req.Location,
// 	}
// 	log.Println("it is req.name", name)
// 	log.Println("It is updated time", device.LastUpdated)

// 	deviceBSON := bson.D{
// 		{Key: "user_id", Value: device.UserID},
// 		{Key: "name", Value: device.Name},
// 		{Key: "device_type", Value: device.DeviceType},
// 		{Key: "device_name", Value: device.DeviceName},
// 		{Key: "device_status", Value: device.DeviceStatus},
// 		{Key: "configuration_settings", Value: device.ConfigurationSettings},
// 		{Key: "last_updated", Value: device.LastUpdated},
// 		{Key: "location", Value: device.Location},
// 	}

// 	deviceID, err := s.serv.CreateDevicee(ctx, deviceBSON)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp := &gendevice.CreateDeviceResp{
// 		DeviceId: deviceID.Hex(),
// 		Device:   req,
// 	}

// 	return resp, nil
// }

package repo

import (
	"context"
	"device-service/internal/handle"
	"device-service/models"
	"device-service/pkg"
	"device-service/protos/gendevice"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewServ struct {
	serv *pkg.Server
}

func Conn() *NewServ {
	server, err := pkg.NewServer()
	if err != nil {
		log.Fatal("Error on newserver.. on repooo", err)

	}

	return &NewServ{serv: server}
}

func (s *NewServ) DeviceCreate(req *gendevice.CreateDeviceRequest) (*gendevice.CreateDeviceResp, error) {
	ctx := context.Background()
	
	obj, err := primitive.ObjectIDFromHex(req.UserId) 
	if err != nil {
        return nil, fmt.Errorf("invalid user_id format: %s", req.UserId)

	}

	log.Println("Objeccttt...it is in repo device service....", obj)
   

	userid, name, err := handle.HandleUser(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("error handling user or user ID not found: %+v", err)
	}
	log.Println("Bu userId .... bu name on devicecreateichida", userid, name)

	device := models.Device{
		UserID:                userid,
		Name:                  name,
		DeviceType:            req.DeviceType,
		DeviceName:            req.DeviceName,
		DeviceStatus:          req.DeviceStatus,
		ConfigurationSettings: req.ConfigurationSettings,
		LastUpdated:           time.Now().Format(time.RFC3339),
		Location:              req.Location,
	}

	deviceBSON := bson.D{
		{Key: "user_id", Value: device.UserID},
		{Key: "name", Value: device.Name},
		{Key: "device_type", Value: device.DeviceType},
		{Key: "device_name", Value: device.DeviceName},
		{Key: "device_status", Value: device.DeviceStatus},
		{Key: "configuration_settings", Value: device.ConfigurationSettings},
		{Key: "last_updated", Value: device.LastUpdated},
		{Key: "location", Value: device.Location},
	}

	deviceID, err := s.serv.CreateDevicee(ctx, deviceBSON)
	if err != nil {
		return nil, err
	}

	resp := &gendevice.CreateDeviceResp{
		DeviceId:    deviceID.Hex(),
		Name:        name,
		Device:      req,
		LastUpdated: device.LastUpdated,
	}

	return resp, nil
}

func (s *NewServ) DeleteByIdd(req *gendevice.DeleteDeviceReq) (*gendevice.DeleteResp, error) {
	ctx := context.Background()

	deviceID, err := primitive.ObjectIDFromHex(req.DeviceId)
	if err != nil {
		return nil, fmt.Errorf("wronggg device ID format: %v", err)
	}


	filter := bson.D{{Key: "_id", Value: deviceID}}
	log.Println("Filter used for delete:", filter)
	log.Println("DeviceId from request:", req.DeviceId)


	result, err := s.serv.DeleteOnee(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error deleting device: %v", err)
	}

	
	if result.DeletedCount == 0 {
		return &gendevice.DeleteResp{
			Message: "No ID found with that",
		}, nil
	}

	return &gendevice.DeleteResp{
		Message: " deleted successfully",
	}, nil
}



func (s *NewServ) UpdateByIdd(req *gendevice.UpdateDeviceRequest) (*gendevice.UpdateDeviceResp, error) {
    ctx := context.Background()

    deviceID, err := primitive.ObjectIDFromHex(req.DeviceId)
    if err != nil {
        return nil, fmt.Errorf("invalid device ID format: %v", err)
    }

    update := bson.D{
        {Key: "$set", Value: bson.D{
            {Key: "device_type", Value: req.DeviceType},
            {Key: "device_name", Value: req.DeviceName},
            {Key: "device_status", Value: req.DeviceStatus},
            {Key: "configuration_settings", Value: req.ConfigurationSettings},
            {Key: "location", Value: req.Location},
            {Key: "last_updated", Value: time.Now().Format(time.RFC3339)},
        }},
    }

    filter := bson.D{{Key: "_id", Value: deviceID}}

    result, err := s.serv.UpdateOnee(ctx, filter, update)
    if err != nil {
        return nil, fmt.Errorf("error updating device: %v", err)
    }

    if result.MatchedCount == 0 {
        return &gendevice.UpdateDeviceResp{
            DeviceId:      req.DeviceId,
            DeviceType:    req.DeviceType,
            DeviceName:    req.DeviceName,
            DeviceStatus:  req.DeviceStatus,
            ConfigurationSettings: req.ConfigurationSettings,
            Location:      req.Location,
            LastUpdated:   time.Now().Format(time.RFC3339),
        }, nil
    }

    return &gendevice.UpdateDeviceResp{
        DeviceId:      req.DeviceId,
        DeviceType:    req.DeviceType,
        DeviceName:    req.DeviceName,
        DeviceStatus:  req.DeviceStatus,
        ConfigurationSettings: req.ConfigurationSettings,
        Location:      req.Location,
        LastUpdated:   time.Now().Format(time.RFC3339),
    }, nil
}
