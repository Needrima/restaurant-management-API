package main

import (
	"github.com/gin-gonic/gin"
	"restaurant-management-API/middleware"
	"restaurant-management-API/routes"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router) // called before authentication to check if user is authenticated before having access to other routes
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":1000")
}
