package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-management-API/database"
	"restaurant-management-API/helpers"
	"restaurant-management-API/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		mathchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}

		userCollections := database.GetCollection("users")
		cursor, err := userCollections.Aggregate(context.TODO(), mongo.Pipeline{mathchStage, projectStage})
		if err != nil {
			log.Println("Error getting users", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
		defer cursor.Close(context.TODO())

		var users []bson.M
		if err := cursor.All(context.TODO(), &users); err != nil {
			log.Println("Error decoding results from users cursor", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get users",
			})
			return
		}

		ctx.JSON(http.StatusOK, users[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		var user models.User

		usersCollection := database.GetCollection("users")
		if err := usersCollection.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&user); err != nil {
			log.Println("Error getting user", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u models.User

		if err := ctx.BindJSON(&u); err != nil {
			log.Println("Error getting user:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		if err := validate.Struct(u); err != nil {
			log.Println("Error validating user struct:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		userCollection := database.GetCollection("users")

		if emailCount, _ := userCollection.CountDocuments(context.TODO(), bson.M{"email": u.Email}); emailCount != 0 {
			log.Println("email taken by another user", emailCount)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "email taken by another user",
			})
			return
		}

		if phoneCount, _ := userCollection.CountDocuments(context.TODO(), bson.M{"phone": u.Phone}); phoneCount != 0 {
			log.Println("phone number taken by another user", phoneCount)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "phone number taken by another user",
			})
			return
		}

		hashedPassword := EncryptPassword(u.Password)
		u.Password = hashedPassword

		u.CreatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
		u.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		u.ID = primitive.NewObjectID()
		u.UserId = u.ID.Hex()

		token, refreshToken, _ := helpers.GenerateAllTokens(u.Email, u.FirstName, u.LastName, u.UserId)
		u.Token = token
		u.RefreshToken = refreshToken

		insertResult, err := userCollection.InsertOne(context.TODO(), u)
		if err != nil {
			log.Println("could not insert user into database:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, insertResult)
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u models.User
		if err := ctx.BindJSON(&u); err != nil {
			log.Println("Error getting user:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		userCollection := database.GetCollection("users")
		var found models.User
		if err := userCollection.FindOne(context.TODO(), bson.M{"email": u.Email}).Decode(&found); err != nil {
			log.Println("user not found:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "firstname or password mismatch",
			})
			return
		}

		if err := VeriifyPasswordFromHash(found.Password, u.Password); err != nil {
			log.Println("password mismatch:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "firstname or password mismatch",
			})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(found.Email, found.FirstName, found.LastName, found.UserId)
		if err := helpers.UpdateUserTokens(token, refreshToken, found.UserId); err != nil {
			log.Println("could not update user tokens:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, found)
	}
}

func EncryptPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println("Erro hashing password")
		return ""
	}

	return string(hashed)
}

func VeriifyPasswordFromHash(wantedPassword, givenPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(wantedPassword), []byte(givenPassword))
}
