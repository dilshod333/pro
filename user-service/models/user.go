package models

import (
	"context"
	"fmt"
	"log"
	"time"

	email "user-service/internal/gmail"
	"user-service/internal/handlers"
	"user-service/internal/jwt"
	"user-service/pkg"
	"user-service/protos/genuser"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Code string
var Email string

type NewServer struct {
	Server *pkg.Server
	
}

func Connect() *NewServer {
	server, err := pkg.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	return &NewServer{
		Server: server,
	}
}

func (s *NewServer) Registerr(req *genuser.UserRegisterRequest) (*genuser.Resp, error) {
	user := pkg.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: req.GetPasswordHash(),
		Profile: pkg.Profile{
			Name:    req.GetProfile().GetName(),
			Address: req.GetProfile().GetAddress(),
		},
	}
	
	redisKey := "user:" + user.Email
	userData := map[string]interface{}{
		"Username":        user.Username,
		"Email":           user.Email,
		"PasswordHash":    user.PasswordHash,
		"Profile.Name":    user.Profile.Name,
		"Profile.Address": user.Profile.Address,
	}

	err := s.Server.RedisClient.HSet(context.Background(), redisKey, userData).Err()
	if err != nil {
		log.Println("111   ", err)
		return nil, err
	}

	Code := email.SendEmail(req.Email)
	Email = req.Email

	err = s.Server.RedisClient.Set(context.Background(), "verification:"+req.Email, Code, 120*time.Second).Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &genuser.Resp{Message: "Check your email for the verification code."}, nil
}

func (s *NewServer) Verifyy(req *genuser.CodeReqest) (*genuser.CodeResponse, error) {
	if req.Email != Email {
		log.Println("It is requested email to redis cache...", req.Email)
		return &genuser.CodeResponse{
			Message: "Email Is Not Same Bruhh",
		}, nil
	}

	storedCode, err := s.Server.RedisClient.Get(context.Background(), "verification:"+req.Email).Result()
	if err == redis.Nil {
		return &genuser.CodeResponse{
			Message: "Verification code expired or Your time is up",
		}, nil
	} else if err != nil {
		return nil, err
	}

	if storedCode != req.Code {
		return &genuser.CodeResponse{
			Message: "Wrong code bruhh",
		}, nil
	}

	filter := bson.D{{Key: "email", Value: req.GetEmail()}}
	var existingUser bson.M
	err = s.Server.Client.Database("Users").Collection("user").FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		return &genuser.CodeResponse{
			Message: "Email already registered",
		}, nil
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	redisKey := "user:" + req.Email
	userData, err := s.Server.RedisClient.HGetAll(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}

	if len(userData) == 0 {
		return &genuser.CodeResponse{
			Message: "User not found",
		}, nil
	}

	token, err := jwt.CreateToken(req.Email)
	if err != nil {
		log.Println("Xatolik models createtoken...", err)
		return nil, err
	}

	passwordHash := userData["PasswordHash"]
	pass, err := handlers.HashPassword(passwordHash)
	if err != nil {
		log.Println("error on hashpassword...", err)
		return nil, err
	}

	user := bson.D{
		{Key: "username", Value: userData["Username"]},
		{Key: "email", Value: userData["Email"]},
		{Key: "passwordHash", Value: pass},
		{Key: "profile", Value: bson.D{
			{Key: "name", Value: userData["Profile.Name"]},
			{Key: "address", Value: userData["Profile.Address"]},
		}},
	}

	userID, err := s.Server.CreateUser(context.Background(), user)
	log.Println("User id keldimi userda,,,,", userID)
	log.Println("Xatolik...", err)
	if err != nil {
		return nil, err
	}

	return &genuser.CodeResponse{
		Message: "Verification successful",
		UserId:  userID.Hex(),
		Token:   token,
	}, nil
}

func (s *NewServer) Loginn(req *genuser.UserLoginRequest) (*genuser.UserLoginResponse, error) {
	filter := bson.D{{Key: "email", Value: req.GetEmail()}}

	var userDec bson.M
	var ctx context.Context
	err := s.Server.Client.Database("Users").Collection("user").FindOne(ctx, filter).Decode(&userDec)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &genuser.UserLoginResponse{
				Message: "User Not Found bruhh",
			}, nil
		}
		return nil, err
	}

	StorePass, ok := userDec["passwordHash"].(string)
	if !ok {
		return nil, fmt.Errorf("error on getting password from mongodb bruhh")
	}
	var ch bool
	ch = false
	ch = handlers.ComparePassword(StorePass, req.GetPassword())

	if !ch {
		return &genuser.UserLoginResponse{
			Message: "Wrong Password bruhh",
		}, nil
	}

	token, err := jwt.CreateToken(req.GetEmail())
	if err != nil {
		log.Println("Error creating token...", err)
		return nil, err
	}

	return &genuser.UserLoginResponse{
		Message: "Login Successful bruh",
		Token:   token,
	}, nil
}

func (s *NewServer) Profilee(req *genuser.ProfileReq) (*genuser.ProfileResp, error) {
	log.Println("........Profile ichiga keldiiii......")
	ctx := context.Background()

	filter := bson.D{{Key: "profile.name", Value: req.GetName()}}

	var userDec bson.M
	err := s.Server.Client.Database("Users").Collection("user").FindOne(ctx, filter).Decode(&userDec)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &genuser.ProfileResp{
				User: nil,
			}, nil
		}
		return nil, err
	}

	profileData, ok := userDec["profile"].(bson.M)
	if !ok {
		return nil, fmt.Errorf("profile data type assertion failed")
	}

	profileName, ok := profileData["name"].(string)
	if !ok {
		profileName = "Unknown"
	}

	profileAddress, ok := profileData["address"].(string)
	if !ok {
		profileAddress = "Unknown"
	}

	email, ok := userDec["email"].(string)
	if !ok {
		email = "Unknown"
	}

	username, ok := userDec["username"].(string)
	if !ok {
		username = "Unknown"
	}
	var idHex string
	if id, ok := userDec["_id"].(primitive.ObjectID); ok {
		idHex = id.Hex()
	} else {
		idHex = "Unknown"
	}

	userResponse := &genuser.UserResponse{
		UserId:       idHex,
		Username:     username,
		Email:        email,
		PasswordHash: userDec["passwordHash"].(string),
		Profile: &genuser.UserResponse_Profile{
			Name:    profileName,
			Address: profileAddress,
		},
	}

	return &genuser.ProfileResp{
		User: userResponse,
	}, nil
}


func (s *NewServer) GetByIdd(req *genuser.GetByIdReq) (*genuser.GetResponse, error) {
    ctx := context.Background()

    
    userID, err := primitive.ObjectIDFromHex(req.GetUserId())
    if err != nil {
        return nil, fmt.Errorf("wrongggg user_id format: %v", err)
    }

   
    filter := bson.D{{Key: "_id", Value: userID}}

    var userDec bson.M
    err = s.Server.Client.Database("Users").Collection("user").FindOne(ctx, filter).Decode(&userDec)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return &genuser.GetResponse{
                User: nil,
            }, nil
        }
        return nil, err
    }

    profileData, ok := userDec["profile"].(bson.M)
    if !ok {
        profileData = bson.M{
            "name":    "Unknown",
            "address": "Unknown",
        }
    }

    userResponse := &genuser.UserResponse{
        UserId:       userID.Hex(),
        Username:     userDec["username"].(string),
        Email:        userDec["email"].(string),
        PasswordHash: userDec["passwordHash"].(string),
        Profile: &genuser.UserResponse_Profile{
            Name:    profileData["name"].(string),
            Address: profileData["address"].(string),
        },
    }

    return &genuser.GetResponse{
        User: userResponse,
    }, nil
}
