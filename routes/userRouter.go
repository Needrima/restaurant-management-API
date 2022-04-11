package routes

import "github.com/gin-gonic/gin"

func UserRoutes(r *gin.Engine) {
	r.GET("/users", controller.GetUsers())
	r.GET("/users/:user_id", controller.GetUser())
	r.POST("/users/signup", controller.SignUp())
	r.POST("/users/login", controller.Login())
}