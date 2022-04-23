package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID             primitive.ObjectID `bson:"_id"`
	NumberOfGuests int                `json:"number_of_guests" bson:"number_of_guests,omitempty" validate:"required"`
	TableNumber    int                `json:"table_number" bson:"table_number,omitempty" validate:"required"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	TableId        string             `json:"table_id" bson:"table_id,omitempty"`
}
