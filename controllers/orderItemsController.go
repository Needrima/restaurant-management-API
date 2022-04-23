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

func GetOrderItems() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderItemCollection := database.GetCollection("orderitem")

		var orderItems []models.OrderItem

		cursor, err := orderItemCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Println("Getting all orderitems error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "soemthing went wrong",
			})
			return
		}
		defer cursor.Close(context.TODO())

		if err := cursor.All(context.TODO(), &orderItems); err != nil {
			log.Println("Decoding into orderItems error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get order items",
			})
			return
		}

		ctx.JSON(http.StatusOK, orderItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderItemId := ctx.Param("order_item_id")
		var orderItem models.OrderItem

		orderItemCollection := database.GetCollection("orderitem")

		if err := orderItemCollection.FindOne(context.TODO(), bson.M{"order_item_id": orderItemId}).Decode(&orderItem); err != nil {
			log.Println("getting order item error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get order item",
			})
			return
		}

		ctx.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var orderItemPack models.OrderItemPack
		var order models.Order

		if err := ctx.BindJSON(&orderItemPack); err != nil {
			log.Println("wrong format from fron end:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid data",
			})
			return
		}

		order.OrderDate, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		orderItemsToBeInserted := []interface{}{}
		order.TableId = orderItemPack.TableId
		orderId := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.OrderId = orderId

			if err := validate.Struct(orderItem); err != nil {
				log.Println("invalid input in orderitem field:", err)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid input in one or more field",
				})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.OrderItemId = orderItem.ID.Hex()
			orderItem.CreatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
			orderItem.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
			orderItem.UnitPrice = toFixed(orderItem.UnitPrice, 2)

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		orderItemcollection := database.GetCollection("orderitem")
		insertresult, err := orderItemcollection.InsertMany(context.TODO(), orderItemsToBeInserted)
		if err != nil {
			log.Println("could not insert order items into database:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not insert order items into database",
			})
			return
		}

		ctx.JSON(http.StatusOK, insertresult)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderItemId := ctx.Param("order_item_id")
		var updatedOrderItem models.OrderItem
		if err := ctx.BindJSON(&updatedOrderItem); err != nil {
			log.Println("getting updated order item:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "something went wrong",
			})
			return
		}

		if updatedOrderItem.FoodID != "" {
			foodCollection := database.GetCollection("food")
			var food models.Food
			if err := foodCollection.FindOne(context.TODO(), bson.M{"food_id": updatedOrderItem.FoodID}).Decode(&food); err != nil {
				log.Println("invalid food id:", err)
				ctx.JSON(http.StatusBadGateway, gin.H{
					"error": "invalid food id",
				})
				return
			}
		}

		updatedOrderItem.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		orderItemCollection := database.GetCollection("orderitem")
		updateResult, err := orderItemCollection.UpdateOne(context.TODO(), bson.M{"order_item_id": orderItemId}, bson.M{"$set": updatedOrderItem})
		if err != nil {
			log.Println("could not update order item:", err)
			ctx.JSON(http.StatusBadGateway, gin.H{
				"error": "could not update order item",
			})
			return
		}

		ctx.JSON(http.StatusOK, updateResult)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId := ctx.Param("order_id")
		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			log.Println("getting order items by order error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string) ([]primitive.M, error) {

}
