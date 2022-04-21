package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `json:"name" validate:"required,min=2,max=100" bson:"name,omitempty"`
	Price     float64            `json:"price" validate:"required" bson:"price,omitempty"`
	FoodImage string             `json:"food_image" validate:"required" bson:"food_image,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	FoodId    string             `json:"food_id" bson:"food_id,omitempty"`
	MenuId    string             `json:"menu_id" validate:"required" bson:"menu_id,omitempty"`
}
