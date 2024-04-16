package db

import (
	"ScrabShortener/cache"
	"ScrabShortener/helpers"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func Connect() *mongo.Client {

	MONGO_URI := helpers.GetEnv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	return client
}

func TestDb() {

	MONGO_URI := helpers.GetEnv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Ping to your MongoDB Database. Connection established.")

	rdb := cache.RedisConnection()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ping to your Redis Database. Connection established")

}

func ConnectionString() *mongo.Collection {
	client := Connect()
	return client.Database("SCShort").Collection("SCShort")

}
