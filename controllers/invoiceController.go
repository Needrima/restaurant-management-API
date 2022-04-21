package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-management-API/database"
	"restaurant-management-API/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

		ctx.JSON(http.StatusOK, invoice)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
