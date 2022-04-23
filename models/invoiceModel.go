package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id"`
	InvoiceId      string             `json:"invoice_id" bson:"invoice_id,omitempty"`
	OrderId        string             `json:"order_id" bson:"order_id,omitempty"`
	PaymentMethod  string             `json:"payment_method" bson:"payment_method,omitempty" validate:"eq=CARD|eq=CASH|eq="`
	PaymentStatus  string             `json:"payment_status" bson:"payment_status,omitempty" validate:"required,eq=PENDING|eq=PAID"`
	PaymentDueDate time.Time          `json:"payment_due_date" bson:"payment_due_date,omitempty"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

type InvoiceView struct {
	InvoiceId      string      `json:"invoice_id,omitempty"`
	PaymentMethod  string      `json:"payment_method,omitempty"`
	OrderId        string      `json:"order_id,omitempty"`
	PaymentStatus  string      `json:"payment_status,omitempty"`
	PaymentDue     interface{} `json:"payment_due,omitempty"`
	TableNumber    interface{} `json:"table_number,omitempty"`
	PaymentDueDate time.Time   `json:"payment_due_date,omitempty"`
	OrderDetails   interface{} `json:"order_details,omitempty"`
}
