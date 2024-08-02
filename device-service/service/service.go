package service

import (
	"context"
	"device-service/internal/repo"
	"device-service/protos/gendevice"
	"device-service/rabbitmq"
	"encoding/json"
	"fmt"
	"log"
)

type Serv struct {
	gendevice.UnimplementedDeviceServerServer
	n *repo.NewServ
}

func Connect() *Serv {
	s := repo.Conn()
	return &Serv{n: s}
}

func (s *Serv) CreateDevice(ctx context.Context, req *gendevice.CreateDeviceRequest) (*gendevice.CreateDeviceResp, error) {
	log.Println("CreateDevice ichida hozi bizning request>>>>>", req)
	resp, err := s.n.DeviceCreate(req)
	log.Println("response>>>>>>>>> Createdevice Servicedan", resp)

	if err != nil {
		log.Println("Devicecreate databasaga borib xatolik keldi>>>>>", err)
		return nil, fmt.Errorf("error on devicecreate")

	}

	return resp, nil
}

func (s *Serv) DeleteById(ctx context.Context, req *gendevice.DeleteDeviceReq) (*gendevice.DeleteResp, error) {
	resp, err := s.n.DeleteByIdd(req)
	if err != nil {
		return nil, fmt.Errorf("not found that id")
	}

	return resp, nil
}

func (s *Serv) UpdateDevice(ctx context.Context, req *gendevice.UpdateDeviceRequest) (*gendevice.UpdateDeviceResp, error) {
	resp, err := s.n.UpdateByIdd(req)

	if err != nil {
		return nil, fmt.Errorf("wrong on update")

	}

	return resp, nil
}

func (s *Serv) Create(ctx context.Context, req *gendevice.DeviceControlReq) (*gendevice.DeviceControlResp, error) {
	log.Println("CReate command ichiga kirdiiii")
	byte, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("error on marshaling")
	}
	err = rabbitmq.PublishOrder(byte)

	if err != nil {
		return nil, fmt.Errorf("publishda xatolik")
	}

	return &gendevice.DeviceControlResp{
		Message: "Updated Successfullyy",
	}, nil

}
