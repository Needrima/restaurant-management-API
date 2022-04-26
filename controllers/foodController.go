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
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	validate = validator.New()
)

func GetFoods() gin.HandlerFunc {
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
		// startIndex, err = strconv.Atoi(ctx.Query("startIndex"))

		matchStage := bson.M{"$match": bson.M{}} // match all documents

		// group by
		groupStage := bson.M{"$group": bson.M{
			"_id":         bson.M{"_id": "null"},     // _id field where _id is null
			"total_count": bson.M{"$sum": 1},         // new total_count field which is total document count
			"data":        bson.M{"$push": "$$ROOT"}, // new data field which is a slice of documents for each distinct _id
		}}

		// project by
		projectStage := bson.M{"$project": bson.M{
			"_id":         0, // ignoring _id field
			"total_count": 1, // including total count field
			"food_items": bson.M{ // new food_items field which is a slice of data from data field containing the required number record based on the page( 10 records if recordPerPage is not specified)
				"$slice": []interface{}{"$data", startIndex, recordPerPage},
			},
		}}

		foodCollection := database.GetCollection("food")

		cursor, err := foodCollection.Aggregate(context.TODO(), []bson.M{matchStage, groupStage, projectStage})
		if err != nil {
			log.Println("getfoods aggregation err:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "soemthing went wrong",
			})
			return
		}
		defer cursor.Close(context.TODO())

		var foods []bson.M
		if err := cursor.All(context.TODO(), &foods); err != nil {
			log.Println("Error decoding cursor:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get foods",
			})
			return
		}

		ctx.JSON(http.StatusOK, foods)
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

		food.CreatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
		food.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

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
		var food models.Food
		var menu models.Menu

		foodId := ctx.Param("food_id")

		err := ctx.BindJSON(&food)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		menuCollection := database.GetCollection("menu")
		if err := menuCollection.FindOne(context.TODO(), bson.M{"menu_id": food.MenuId}).Decode(&menu); err != nil {
			log.Println("invalid menu id:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid menu id",
			})
			return
		}

		food.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		foodCollection := database.GetCollection("food")
		updateResult, err := foodCollection.UpdateOne(context.TODO(), bson.M{"food_id": foodId}, bson.M{"$set": food})
		if err != nil {
			log.Println("error updating food:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not update food",
			})
			return
		}

		ctx.JSON(http.StatusOK, updateResult)
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
