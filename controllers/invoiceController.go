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

func GetInvoices() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var invoices []models.Invoice

		invoiceCollection := database.GetCollection("invoice")
		cursor, err := invoiceCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Println("could not decode invoices:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
		defer cursor.Close(context.TODO())

		if err := cursor.All(context.TODO(), &invoices); err != nil {
			log.Println("could not get invoices:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get invoices",
			})
			return
		}

		ctx.JSON(http.StatusOK, invoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		invoiceId := ctx.Param("invoice_id")

		var invoice models.Invoice
		invoiceCollection := database.GetCollection("invoice")
		if err := invoiceCollection.FindOne(context.TODO(), bson.M{"invoice_id": invoiceId}).Decode(&invoice); err != nil {
			log.Println("could not decode invoice:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get invoice with specified id",
			})
			return
		}

		var invoiceView models.InvoiceView

		allOrderItems, err := ItemsByOrder(invoice.OrderId)
		if err != nil {
			log.Println("error getting all order items:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		invoiceView.OrderId = invoice.OrderId
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.PaymentMethod = invoice.PaymentMethod
		invoiceView.InvoiceId = invoice.InvoiceId
		invoiceView.PaymentStatus = invoice.PaymentStatus
		invoiceView.PaymentDue = allOrderItems[0]["payment_due"]
		invoiceView.TableNumber = allOrderItems[0]["table_number"]
		invoiceView.OrderDetails = allOrderItems[0]["order_items"]

		ctx.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var invoice models.Invoice
		if err := ctx.BindJSON(&invoice); err != nil {
			log.Println("error getting all invoice from frontend:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		if invoice.PaymentStatus == "" {
			invoice.PaymentStatus = "PENDING"
		}

		var order models.Order
		orderCollection := database.GetCollection("order")
		if err := orderCollection.FindOne(context.TODO(), bson.M{"order_id": invoice.OrderId}).Decode(&order); err != nil {
			log.Println("invalid orderId:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid order id, check order id and try again",
			})
			return
		}

		invoice.PaymentDueDate, _ = time.Parse(time.ANSIC, time.Now().AddDate(0, 0, 1).Format(time.ANSIC))
		invoice.CreatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
		invoice.UpdatedAt, _ = time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))

		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceId = invoice.ID.Hex()

		if err := validate.Struct(invoice); err != nil {
			log.Println("invalid invoice:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "something went wrong",
			})
			return
		}

		invoiceCollection := database.GetCollection("invoice")
		insertResult, err := invoiceCollection.InsertOne(context.TODO(), invoice)
		if err != nil {
			log.Println("invalid invoice:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		ctx.JSON(http.StatusOK, insertResult)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		invoiceId := ctx.Param("invoice_id")
		var invoice models.Invoice
		if err := ctx.BindJSON(&invoice); err != nil {
			log.Println("error getting all update invoice from frontend:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		if invoice.PaymentStatus == "" {
			invoice.PaymentStatus = "PENDING"
		}

		invoiceCollection := database.GetCollection("invoice")
		updateResult, err := invoiceCollection.UpdateOne(context.TODO(), bson.M{"invoice_id": invoiceId}, bson.M{"$set": invoice})
		if err != nil {
			log.Println("error updating invoice:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not update invoice",
			})
			return
		}

		ctx.JSON(http.StatusOK, updateResult)
	}
}
