package api

import (
	"api-gateway/jwt"
	"api-gateway/methods"
	"api-gateway/protos/gendevice"
	"api-gateway/protos/genuser"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Serverr struct {
	userr   genuser.UserServiceClient
	devicee gendevice.DeviceServerClient
}

func Conn() *Serverr {
	user := methods.ConnectUser()
	device := methods.ConnectDevice()
	return &Serverr{userr: user, devicee: device}
}

func (u *Serverr) GetUserByName(c *gin.Context, name string) {
	token := c.GetHeader("Authorization")
	valid, email, err := jwt.VerifyToken(token)
	log.Println("Jwt token parse qilingandan keyin olinyapti.... oldin", email)

	if err != nil {
		c.IndentedJSON(401, gin.H{"error": "First Register or Login"})
		return
	}
	log.Println("Jwt token parse qilingandan keyin olinyapti.... keyin", email)
	if !valid {
		c.JSON(401, gin.H{"error": "Unauthorized: Invalid or expired token"})
		return
	}

	c.Header("Content-Type", "application/json")

	req := &genuser.ProfileReq{Name: name}
	resp, err := u.userr.Profile(context.Background(), req)
	if err != nil {
		log.Println("Error getting user by name:", err)
		c.JSON(500, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(200, resp)
}

func (u *Serverr) GetUserId(c *gin.Context, id string) {
	token := c.GetHeader("Authorization")
	valid, email, err := jwt.VerifyToken(token)
	if err != nil {
		c.IndentedJSON(401, gin.H{"error": "First Register or Login"})
		return
	}
	if !valid {
		c.JSON(401, gin.H{"error": "Unauthorized: Invalid or expired token"})
		return
	}
	_ = email
	c.Header("Content-Type", "application/json")
	req := &genuser.GetByIdReq{UserId: id}
	resp, err := u.userr.GetById(context.Background(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": "wronggg to get user"})
		return
	}

	c.JSON(200, resp)
}

func (u *Serverr) RegisterUser(c *gin.Context) {
	var req genuser.UserRegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		log.Println("error bindddd request:", err)
		c.JSON(400, gin.H{"error": "wrong request bruhhh"})
		return
	}

	resp, err := u.userr.Register(context.Background(), &req)
	if err != nil {
		log.Println("error registering user:", err)
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(200, resp)
}

func (s *Serverr) VeriffyUser(c *gin.Context) {

	var req genuser.CodeReqest
	if err := c.ShouldBind(&req); err != nil {
		log.Println("error bindddd request:", err)
		c.JSON(400, gin.H{"error": "wrong request bruhhh"})
		return
	}
	resp, err := s.userr.Verify(context.Background(), &req)
	if err != nil {
		log.Println("error verifyyy user:", err)
		c.JSON(500, gin.H{"error": "error  veriffyyy user"})
		return
	}
	c.IndentedJSON(200, resp)
}

func (s *Serverr) Loginn(c *gin.Context) {
	var req genuser.UserLoginRequest

	if err := c.ShouldBind(&req); err != nil {
		log.Println("errror on bindd request")
		return
	}

	resp, err := s.userr.Login(context.Background(), &req)

	if err != nil {
		log.Println("error login user:", err)
		c.JSON(500, gin.H{"error": "error  login user"})
		return
	}
	c.IndentedJSON(200, resp)
}

func (s *Serverr) CreateDevicee(c *gin.Context) {
	var req gendevice.CreateDeviceRequest

	if err := c.ShouldBind(&req); err != nil {
		c.IndentedJSON(405, gin.H{"error": "should bind error bruhh"})
		return
	}

	resp, err := s.devicee.CreateDevice(context.Background(), &req)
	log.Println("Response keldi.... ", resp)
	log.Printf("Xatolik keldi>>>>>>>> %+v", err)
	if err != nil {
		c.IndentedJSON(500, gin.H{"Error": "Internal server error"})
		return
	}

	c.IndentedJSON(200, resp)
}

func (s *Serverr) DeleteByid(c *gin.Context, id string) {
	c.Header("Content-Type", "application/json")

	req := &gendevice.DeleteDeviceReq{DeviceId: id}

	resp, err := s.devicee.DeleteById(context.Background(), req)

	if err != nil {
		c.IndentedJSON(500, gin.H{"Error": "smth wrong with delete id  error"})
		return
	}

	c.IndentedJSON(200, resp)
}

func (s *Serverr) Updateee(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var req gendevice.UpdateDeviceRequest

	if err := c.ShouldBind(&req); err != nil {
		c.IndentedJSON(405, gin.H{"error": "should bind error bruhh"})
		return
	}

	resp, err := s.devicee.UpdateDevice(context.Background(), &req)

	if err != nil {
		c.IndentedJSON(500, gin.H{"Error": "Internal server error"})
		return
	}

	c.IndentedJSON(200, resp)
}

func (s *Serverr) CreateCommand(c *gin.Context) {
	var device gendevice.DeviceControlReq
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := s.devicee.Create(c, &device)

	if err != nil {
		c.IndentedJSON(405, gin.H{"error": "error on device"})
		log.Println("xatolikkk >>>>", err)
		return 
	}

	c.JSON(http.StatusOK, gin.H{"status": resp})
}
