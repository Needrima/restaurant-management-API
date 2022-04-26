package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-management-API/database"
	"restaurant-management-API/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTables() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tableCollection := database.GetCollection("table")

		var tables []models.Table

		cursor, err := tableCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Println("Getting all tables error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "soemthing went wrong",
			})
			return
		}
		defer cursor.Close(context.TODO())

		if err := cursor.All(context.TODO(), &tables); err != nil {
			log.Println("Decoding into tables error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get order items",
			})
			return
		}

		ctx.JSON(http.StatusOK, tables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tableId := ctx.Param("table_id")
		var table models.Table

		tableCollection := database.GetCollection("table")

		if err := tableCollection.FindOne(context.TODO(), bson.M{"table_id": tableId}).Decode(&table); err != nil {
			log.Println("getting table error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get table ",
			})
			return
		}

		ctx.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var table models.Table

		if err := ctx.BindJSON(&table); err != nil {
			log.Println("getting table error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		if err := validate.Struct(table); err != nil {
			log.Println("validating table error:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid input in fileds",
			})
			return
		}

		table.CreatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
		table.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		table.ID = primitive.NewObjectID()
		table.TableId = table.ID.Hex()

		tableCollection := database.GetCollection("table")
		insertResult, err := tableCollection.InsertOne(context.TODO(), table)
		if err != nil {
			log.Println("inserting table error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, insertResult)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updatedtable models.Table
		tableId := ctx.Param("table_id")

		if err := ctx.BindJSON(&updatedtable); err != nil {
			log.Println("getting table error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		updatedtable.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		tableCollection := database.GetCollection("table")
		updateResult, err := tableCollection.UpdateOne(context.TODO(), bson.M{"table_id": tableId}, bson.M{"$set": updatedtable})
		if err != nil {
			log.Println("updating table error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not update table",
			})
			return
		}

		ctx.JSON(http.StatusOK, updateResult)
	}
}
