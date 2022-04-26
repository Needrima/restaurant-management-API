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
	"go.mongodb.org/mongo-driver/mongo"
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

		if err := ctx.BindJSON(&orderItemPack); err != nil {
			log.Println("wrong format from fron end:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid data",
			})
			return
		}

		var order models.Order
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
	matchStage := bson.D{
		{
			Key:   "$match",
			Value: bson.D{{Key: "order_id", Value: id}},
		},
	}

	foodLookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "food"},
				{Key: "localField", Value: "food_id"},
				{Key: "foreignField", Value: "food_id"},
				{Key: "as", Value: "food"},
			},
		},
	}

	foodUnwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.D{
				{
					Key:   "path",
					Value: "$food",
				},
				{
					Key:   "preserveNullAndEmptyArrays",
					Value: true,
				},
			},
		},
	}

	orderLookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "order"},
				{Key: "localField", Value: "order_id"},
				{Key: "foreignField", Value: "order_id"},
				{Key: "as", Value: "order"},
			},
		},
	}

	orderUnwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.D{
				{
					Key:   "path",
					Value: "$order",
				},
				{
					Key:   "preserveNullAndEmptyArrays",
					Value: true,
				},
			},
		},
	}

	tableLookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "table"},
				{Key: "localField", Value: "order.table_id"},
				{Key: "foreignField", Value: "table_id"},
				{Key: "as", Value: "table"},
			},
		},
	}

	tableUnwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.D{
				{
					Key:   "path",
					Value: "$table",
				},
				{
					Key:   "preserveNullAndEmptyArrays",
					Value: true,
				},
			},
		},
	}

	projectStage := bson.D{
		{
			Key: "$project",
			Value: bson.D{
				{"id", 0},
				{"total_count", 1},
				{"amount", "$food.price"},
				{"food_name", "$food.name"},
				{"food_image", "$food.food_image"},
				{"table_number", "$table.table_number"},
				{"table_id", "$table.table_id"},
				{"order_id", "$order.order_id"},
				{"price", "$food.price"},
				{"quantity", 1},
			},
		},
	}

	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{
			{"order_id", "$order_id"},
			{"table_id", "$table_id"},
			{"table_number", "$table_number"},
		}},
		{"payment_due", bson.D{{"$sum", "$amount"}}},
		{"total_count", bson.D{{"$sum", 1}}},
		{"order_items", bson.D{{"$push", "$$ROOT"}}},
	}}}

	nextProjectStage := bson.D{{"$project", bson.D{
		{"_id", 0},
		{"payment_due", 1},
		{"total_count", 1},
		{"table_number", "$_id.table_number"},
		{"order_items", 1},
	}}}

	orderItemCollection := database.GetCollection("orderitem")
	cursor, err := orderItemCollection.Aggregate(context.TODO(), mongo.Pipeline{
		matchStage,
		foodLookupStage,
		foodUnwindStage,
		orderLookupStage,
		orderUnwindStage,
		tableLookupStage,
		tableUnwindStage,
		projectStage,
		groupStage,
		nextProjectStage,
	})

	if err != nil {
		return []primitive.M{}, err
	}
	defer cursor.Close(context.TODO())

	var orderItems []primitive.M
	if err := cursor.All(context.TODO(), &orderItems); err != nil {
		return []primitive.M{}, err
	}

	return orderItems, nil
}
