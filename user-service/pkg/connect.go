package pkg

import (
	"context"
	"log"

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

func (s *Server) CreateUser(ctx context.Context, user bson.D) (primitive.ObjectID, error) {
	collection := s.Client.Database("Users").Collection("user")

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	
	userID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, err
	}

	return userID, nil
}
