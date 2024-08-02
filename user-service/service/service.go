package service

import (
	"context"
	"fmt"
	"log"
	"user-service/models"
	"user-service/protos/genuser"
)

type Serverr struct {
	*models.NewServer
	genuser.UnimplementedUserServiceServer
	Lg *log.Logger
}

func Conn(lg *log.Logger) *Serverr {
	serv := models.Connect()
	return &Serverr{NewServer: serv, Lg: lg}
}

func (s *Serverr) Register(ctx context.Context, req *genuser.UserRegisterRequest) (*genuser.Resp, error) {
	resp, err := s.Registerr(req)
	s.Lg.Println("Register ishladiiiiii")
	if err != nil {
		s.Lg.Println("Error keldimi ....", err)
		log.Println("Smth wrong getting response on service .... ", err)
		return nil, err
	}

	return resp, nil
}

func (s *Serverr) Verify(ctx context.Context, req *genuser.CodeReqest) (*genuser.CodeResponse, error) {
	resp, err := s.Verifyy(req)

	if err != nil {
		log.Println("Xatolik bor veriffyyda...", err)
		return nil, err
	}

	return resp, nil
}

func (s *Serverr) Login(ctx context.Context, req *genuser.UserLoginRequest) (*genuser.UserLoginResponse, error) {
	resp, err := s.Loginn(req)

	if err != nil {
		return nil, fmt.Errorf("error on loginnn")
	}

	return resp, nil
}

func (s *Serverr) Profile(ctx context.Context, req *genuser.ProfileReq) (*genuser.ProfileResp, error) {
	resp, err := s.Profilee(req)

	if err != nil {
		return nil, fmt.Errorf("error on getting profile info")
	}

	return resp, nil 
}


func (s *Serverr) GetById(ctx context.Context, req *genuser.GetByIdReq) (*genuser.GetResponse, error) {
	resp, err := s.GetByIdd(req)

	if err != nil {
		log.Println("Error getbyidd,,,", err)
		return nil, err 
	}

	return resp, nil 
}