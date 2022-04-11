package routes

import (
	"github.com/gin-gonic/gin"
	"restaurant-management-API/controllers"
)

func OrderItemRoutes(r *gin.Engine) {
	r.GET("/orderItems", controller.GetOrderItems())
	r.GET("/orderItems/:orderItem_id", controller.GetOrderItem())
	r.GET("/orderItems-order/:order_id", controller.GetOrderItemsByOrder())
	r.POST("/orderItems", controller.CreateOrderItem())
	r.PATCH("/orderItems/:orderItem_id", controller.UpdateOrderItem())
}
