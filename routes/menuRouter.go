package routes

import (
	"github.com/gin-gonic/gin"
	"restaurant-management-API/controllers"
)

func MenuRoutes(r *gin.Engine) {
	r.GET("/menus", controller.GetMenus())
	r.GET("/menus/:menu_id", controller.GetMenu())
	r.POST("/menus", controller.CreateMenu())
	r.PATCH("/menus/:menu_id", controller.UpdateMenu())
}
