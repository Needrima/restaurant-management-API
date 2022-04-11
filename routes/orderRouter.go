package routes

import (
	"github.com/gin-gonic/gin"
	"restaurant-management-API/controllers"
)

func OrderRoutes(r *gin.Engine) {
	r.GET("/orders", controller.GetOrders())
	r.GET("/orders/:order_id", controller.GetOrder())
	r.POST("/orders", controller.CreateOrder())
	r.PATCH("/orders/:order_id", controller.UpdateOrder())
}
