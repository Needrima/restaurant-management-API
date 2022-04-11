package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(name string) (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	context, cancle := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancle()

	client, err := mongo.Connect(context, clientOptions)
	if err != nil {
		return nil, err
	}

	db := client.Database("RMS")

	return db.Collection(name), nil
}
