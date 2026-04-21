package db

import "go.mongodb.org/mongo-driver/v2/mongo"

func NewMongoDatabase(client *mongo.Client) *mongo.Database {
	return client.Database("go_test")
}
