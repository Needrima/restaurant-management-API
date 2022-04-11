package routes

import (
	"github.com/gin-gonic/gin"
	"restaurant-management-API/controllers"
)

func InvoiceRoutes(r *gin.Engine) {
	r.GET("/invoices", controller.GetInvoices())
	r.GET("/invoices/:invoice_id", controller.GetInvoice())
	r.POST("/invoices", controller.CreateInvoice())
	r.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())
}
