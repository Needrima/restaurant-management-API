package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `json:"name" validate:"required" bson:"name,omitempty"`
	Category  string             `json:"category" validate:"required" bson:"category,omitempty"`
	StartDate *time.Time         `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate   *time.Time         `json:"end_date,omitempty" bson:"end_date,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	MenuId    string             `json:"menu_id" bson:"menu_id,omitempty"`
}
