package helpers

import (
	"context"
	"errors"
	"log"
	"restaurant-management-API/database"
	"restaurant-management-API/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

type SignedDetails struct {
	Email  string
	Fname  string
	Lname  string
	UserId string
	jwt.StandardClaims
}

var (
	userCollection = database.GetCollection("users")
	secretKey      = "RMS-SECRET-KEY"
)

func GenerateAllTokens(email, fname, lname, id string) (string, string, error) {
	claims := &SignedDetails{
		Email:  email,
		Fname:  fname,
		Lname:  lname,
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // expired after 1 day
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(), // expired after 1 week
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))

	return token, refreshToken, err
}

func UpdateUserTokens(token, refreshToken, userId string) error {
	userUpdate := models.User{
		Token:        token,
		RefreshToken: refreshToken,
	}

	userUpdate.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

	userCollection := database.GetCollection("user")
	_, err := userCollection.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.M{"$set": userUpdate})
	return err
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("Parsewithclaims:", err)
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
