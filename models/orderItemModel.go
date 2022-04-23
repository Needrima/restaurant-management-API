package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	Quantity    string             `json:"quantity" bson:"quantity,omitempty" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice   float64            `json:"unit_price" bson:"unit_price,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	FoodID      string             `json:"food_id" bson:"food_id,omitempty" validate:"required"`
	OrderItemId string             `json:"order_item_id" bson:"order_item_id,omitempty"`
	OrderId     string             `json:"order_id" bson:"order_id,omitempty" validate:"required"`
}

type OrderItemPack struct {
	TableId    string
	OrderItems []OrderItem
}
