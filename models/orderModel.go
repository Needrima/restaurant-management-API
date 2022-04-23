package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id"`
	OrderDate time.Time          `json:"order_date" validate:"required" bson:"order_date,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	OrderId   string             `json:"order_id" bson:"order_id,omitempty"`
	TableId   string             `json:"table_id" validate:"required" bson:"table_id,omitempty"`
}
