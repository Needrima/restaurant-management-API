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

func GetOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var orders []models.Order

		ordersCollection := database.GetCollection("order")

		cursor, err := ordersCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Println("orders cursor error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
		defer cursor.Close(context.TODO())

		if err := cursor.All(context.TODO(), &orders); err != nil {
			log.Println("error getting orders:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get orders",
			})
			return
		}

		ctx.JSON(http.StatusOK, orders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId := ctx.Param("order_id")
		ordersCollection := database.GetCollection("order")

		var order models.Order

		if err := ordersCollection.FindOne(context.TODO(), bson.M{"order_id": orderId}).Decode(&order); err != nil {
			log.Println("error getting order:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get order",
			})
			return
		}

		ctx.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var order models.Order

		if err := ctx.BindJSON(&order); err != nil {
			log.Println("error binding order:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong. check for correct entry in fields try again",
			})
			return
		}

		if err := validate.Struct(order); err != nil {
			log.Println("invalid field entries:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong. check for correct entry in fields try again" + err.Error(),
			})
			return
		}

		var table models.Table
		tablesCollection := database.GetCollection("table")
		if err := tablesCollection.FindOne(context.TODO(), bson.M{"table_id": order.TableId}).Decode(&table); err != nil {
			log.Println("invalid table id:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid table id",
			})
			return
		}

		order.CreatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
		order.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		order.ID = primitive.NewObjectID()
		order.OrderId = order.ID.Hex()

		orderCollection := database.GetCollection("order")

		insertResult, err := orderCollection.InsertOne(context.TODO(), order)
		if err != nil {
			log.Println("inserting order into database:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not add order to database",
			})
			return
		}

		ctx.JSON(http.StatusOK, insertResult)
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId := ctx.Param("order_id")

		var order models.Order
		if err := ctx.BindJSON(&order); err != nil {
			log.Println("error getting updated order:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		if order.TableId != "" {
			var table models.Table
			tableCollection := database.GetCollection("table")
			if err := tableCollection.FindOne(context.TODO(), bson.M{"table_id": order.TableId}).Decode(&table); err != nil {
				log.Println("error finding table id to update:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "invalid table id",
				})
				return
			}
		}

		order.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		ordersCollection := database.GetCollection("order")
		updateResult, err := ordersCollection.UpdateOne(context.TODO(), bson.M{"order_id": orderId}, bson.M{"$set": order})
		if err != nil {
			log.Println("error updating order:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not update order/ may be due to invalid order id",
			})
			return
		}

		ctx.JSON(http.StatusOK, updateResult)
	}
}

func OrderItemOrderCreator(order models.Order) string {
	order.CreatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
	order.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()

	orderCollection := database.GetCollection("order")

	_, err := orderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		log.Println("order item order creator err:", err)
		return ""
	}

	return order.OrderId
}
