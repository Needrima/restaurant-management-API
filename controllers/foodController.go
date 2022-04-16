package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-management-API/database"
	"restaurant-management-API/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/bluesuncorp/validator.v5"
)

var (
	validate = validator.New()
)

func GetFoods() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GetFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		foodId := ctx.Param("food_id")

		var food models.Food

		foodCollection := database.GetCollection("food")
		if err := foodCollection.FindOne(context.TODO(), bson.M{"food_id": foodId}).Decode(&food); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "no food with specified id",
			})
			return
		}

		ctx.JSON(http.StatusOK, food) // similar to json.NewEncoder(w).Encode(food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var menu models.Menu
		var food models.Food

		err := ctx.BindJSON(&food) // similar to json.NewDecoder(r.Body).Decode(&food)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validateErr := validate.Struct(food)
		if validateErr != nil {
			log.Println("Validating food:", validateErr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "error validating food",
			})
			return
		}

		menuCollection := database.GetCollection("menu")
		if err := menuCollection.FindOne(context.TODO(), bson.M{"menu_id": food.MenuId}).Decode(&menu); err != nil {
			log.Println("menu not found:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "menu to add food to does not exist",
			})
			return
		}

		food.CreatedAt = time.Now().Format(time.ANSIC)
		food.UpdateAt = time.Now().Format(time.ANSIC)

		food.ID = primitive.NewObjectID()
		food.FoodId = food.ID.Hex()

		var num = toFixed(food.Price, 2)
		food.Price = num

		foodCollection := database.GetCollection("food")
		insertResult, err := foodCollection.InsertOne(context.TODO(), food)
		if err != nil {
			log.Println("inserting food into database:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not insert add food to database",
			})
			return
		}

		ctx.JSON(http.StatusOK, insertResult)
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// func round(number float64) int {
// 	return int(math.Round(number))
// }

func toFixed(number float64, precision int) float64 {
	p := strconv.Itoa(precision)
	num := fmt.Sprintf("%."+p+"f", number)

	f, _ := strconv.ParseFloat(num, 64)

	return f
}
