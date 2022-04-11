package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(name string) *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	context, cancle := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancle()

	client, err := mongo.Connect(context, clientOptions)
	if err != nil {
		panic(fmt.Sprintf("Error creating database client: %v", err))
	}

	return client.Database("RMS").Collection(name)
}
