package routes

import (
	"github.com/gin-gonic/gin"
	"restaurant-management-API/controllers"
)

func FoodRoutes(r *gin.Engine) {
	r.GET("/foods", controller.GetFoods())
	r.GET("/foods/:food_id", controller.GetFood())
	r.POST("/foods", controller.CreateFood())
	r.PATCH("/foods/:food_id", controller.UpdateFood())
}
