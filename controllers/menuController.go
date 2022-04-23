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
		menuId := ctx.Param("menu_id")

		var menu models.Menu

		menuCollection := database.GetCollection("menu")
		if err := menuCollection.FindOne(context.TODO(), bson.M{"menu_id": menuId}).Decode(&menu); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "no menu with specified id",
			})
			return
		}

		ctx.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var menu models.Menu

		err := ctx.BindJSON(&menu)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validateErr := validate.Struct(menu)
		if validateErr != nil {
			log.Println("Validating menu:", validateErr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "error validating menu",
			})
			return
		}

		menu.CreatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
		menu.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		menu.ID = primitive.NewObjectID()
		menu.MenuId = menu.ID.Hex()

		menuCollection := database.GetCollection("menu")
		insertResult, err := menuCollection.InsertOne(context.TODO(), menu)
		if err != nil {
			log.Println("inserting menu into database:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not insert add menu to database",
			})
			return
		}

		ctx.JSON(http.StatusOK, insertResult)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		menuId := ctx.Param("menu_id")

		var menu models.Menu
		err := ctx.BindJSON(&menu)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if menu.StartDate != nil && menu.EndDate != nil {
			if !inTimeSpan(*menu.StartDate, *menu.EndDate) {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "re-type the time",
				})
				return
			}

			menu.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

			menuCollection := database.GetCollection("menu")
			updateResult, err := menuCollection.UpdateOne(context.TODO(), bson.M{"menu_id": menuId}, bson.M{"$set": menu})
			if err != nil {
				log.Println("Updating menu error:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "update failed",
				})
				return
			}

			ctx.JSON(http.StatusOK, updateResult)
		}
	}
}

func inTimeSpan(startTime, endTime time.Time) bool {
	return startTime.After(time.Now()) && endTime.After(startTime)
}
