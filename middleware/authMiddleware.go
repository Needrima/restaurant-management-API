package middleware

import (
	"log"
	"net/http"
	"restaurant-management-API/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			log.Println("No authorization token")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "no authorization token in header",
			})
			ctx.Abort()
		}

		claims, err := helpers.ValidateToken(token)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
		}

		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.Fname)
		ctx.Set("last_name", claims.Lname)
		ctx.Set("uid", claims.UserId)

		ctx.Next()
	}
}
