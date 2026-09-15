package db

import (
	"Orbit/configs"
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Client   *mongo.Client
	instance *mongo.Database
	once     sync.Once
)

func NewMongoDb() {
	cfg := configs.LoadConfig()
	clientOptions := options.Client().ApplyURI(cfg.Mongo.Url)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Some Error Occured Connecting with Database..!!")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil { // checking with pinging if successfull connection happened or not
		log.Fatal("Mongo Ping Error:", err)
	}

	log.Println("Mongodb Database Connected Successfully..!!")

	instance = client.Database("Orbit")
	Client = client
}

func GetInstance() *mongo.Database {
	once.Do(func() {
		NewMongoDb()
	})

	return instance
}

func GetClient() *mongo.Client {
	once.Do(func() {
		NewMongoDb()
	})
	return Client
}
