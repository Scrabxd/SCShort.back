package db

import (
	"ScrabShortener/helpers"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

var rdb = RbdConn()

var Address = helpers.GetEnv("REDIS_URI")
var MONGO_URI = helpers.GetEnv("MONGO_URI")

func Connect() *mongo.Client {

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	return client
}

func RbdConn() *redis.Client {

	opt, err := redis.ParseURL(Address)
	if err != nil {
		log.Fatal("No Connection to redis Stablished.")
	}
	rdb := redis.NewClient(opt)

	return rdb
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
