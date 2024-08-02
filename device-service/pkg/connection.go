package pkg

import (
	"context"
	"device-service/protos/gendevice"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	Client      *mongo.Client
	RedisClient *redis.Client
}

func ConnectMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://mongo-db:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return client, nil
}

func ConnectRedis() (*redis.Client, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-db:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis")

	return rdb, nil
}

func NewServer() (*Server, error) {
	mongoClient, err := ConnectMongo()
	if err != nil {
		return nil, err
	}

	redisClient, err := ConnectRedis()
	if err != nil {
		return nil, err
	}

	return &Server{
		Client:      mongoClient,
		RedisClient: redisClient,
	}, nil
}

func (s *Server) CreateDevicee(ctx context.Context, device bson.D) (primitive.ObjectID, error) {
	log.Println("CreateDevice database ichidaa....", device)
	collection := s.Client.Database("Devices").Collection("device")

	result, err := collection.InsertOne(ctx, device)
	if err != nil {
		log.Println("Xatolik insert qilyotganda...", err)
		return primitive.NilObjectID, err
	}

	deviceID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("Primitive object idda xatolik...", err)
		return primitive.NilObjectID, fmt.Errorf("smth wrong on nilobejectid")
	}

	return deviceID, nil
}

func (s *Server) DeleteOnee(ctx context.Context, filter bson.D) (*mongo.DeleteResult, error) {
	collection := s.Client.Database("Devices").Collection("device")

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Server) UpdateOnee(ctx context.Context, filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	collection := s.Client.Database("Devices").Collection("device")
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) UpdateControl(ctx context.Context, req *gendevice.DeviceControlReq) error {
	log.Println("DAtaabase ichiga kirdi....")
	collection := s.Client.Database("Devices").Collection("device")
	object_id, err := primitive.ObjectIDFromHex(req.DeviceId)
	if err != nil {
		return fmt.Errorf("invalid object id")
	}
	filter := bson.M{"_id": object_id}
	log.Println("Filte>>>>", filter)
	var dev gendevice.CreateDeviceResp
	db := collection.FindOne(ctx, filter).Decode(&dev)
	if db.Error() != "" {
		return fmt.Errorf("not found")
	}

	dev.Device.DeviceStatus = req.Status
	dev.Device.DeviceType = req.CommandType
	dev.LastUpdated = time.Now().Format(time.RFC3339)
	log.Println("Dev>>>>>>>>>>>>", &dev)
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": &dev})
	if err != nil {
		return fmt.Errorf("erron on updateeeeee")
	}
	
	return nil

}
