package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID
	FirstName    string `json:"firstname" validate:"min=3,max=20"`
	LastName     string `json:"lastname" validate:"min=3,max=20"`
	Password     string `json:"password" validate:"required,min=5"`
	Email        string `json:"email" validate:"email,required"`
	Avatar       string `json:"avatar"`
	Phone        string `json:"phone" validate:"required"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	UserId       string `json:"user_id"`
}
