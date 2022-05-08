package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID
	FirstName    string    `json:"firstname" bson:"firstname,omitempty" validate:"min=3,max=20"`
	LastName     string    `json:"lastname" bson:"lastname,omitempty" validate:"min=3,max=20"`
	Password     string    `json:"password" bson:"password,omitempty" validate:"required,min=5"`
	Email        string    `json:"email" bson:"email,omitempty" validate:"email,required"`
	Avatar       string    `json:"avatar" bson:"avatar,omitempty"`
	Phone        string    `json:"phone" bson:"phone,omitempty" validate:"required"`
	Token        string    `json:"token" bson:"token,omitempty"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token,omitempty"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at,omitempty"`
	UserId       string    `json:"user_id" bson:"user_id,omitempty"`
}
