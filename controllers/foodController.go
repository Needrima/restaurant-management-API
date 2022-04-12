package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFoods() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GetFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func CreateFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func UpdateFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// func round(number float64) int {
// 	return int(math.Round(number))
// }

func toFixed(number float64, precision string) float64 {
	num := fmt.Sprintf("%."+precision+"f", number)

	f, _ := strconv.ParseFloat(num, 64)

	return f
}
