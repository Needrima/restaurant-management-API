package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id"`
	OrderDate time.Time          `json:"order_date" validate:"required" bson:"order_date"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	OrderId   string             `json:"order_id" bson:"order_id"`
	TableId   string             `json:"table_id" validate:"required" bson:"table_id"`
}
