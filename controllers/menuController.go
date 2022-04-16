package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-management-API/database"
	"restaurant-management-API/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var ()

func GetMenus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		menuCollection := database.GetCollection("menu")

		menuCursor, err := menuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Println("Could not get menus cursor:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		var menus []models.Menu

		for menuCursor.Next(context.TODO()) {
			var menu models.Menu
			if err := menuCursor.Decode(&menu); err != nil {
				log.Println("error getting menus:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "something went wrong",
				})
				return
			}

			menus = append(menus, menu)
		}

		ctx.JSON(http.StatusOK, menus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func CreateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
