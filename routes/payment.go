package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(r *gin.Engine) {
	paymentRoutes := r.Group("/api/payments")
	paymentRoutes.Use(middlewares.RequireAuth)
	{
		// Admin full CRUD
		paymentRoutes.GET("/invoice", middlewares.RequireRole("admin"), controllers.GetAllPayments)
		paymentRoutes.POST("/invoice", middlewares.RequireRole("admin"), controllers.CreatePayment)
		paymentRoutes.PUT("/invoice/:id", middlewares.RequireRole("admin"), controllers.UpdatePayment)
		paymentRoutes.DELETE("/invoice/:id", middlewares.RequireRole("admin"), controllers.DeletePayment)
		
		// History & PDF (Admin & Employee? Actually Customer & Admin)
		paymentRoutes.GET("/history/:customer_id", controllers.GetCustomerPayments)
		paymentRoutes.GET("/invoice/:id/pdf", controllers.GeneratePaymentPDF)
	}
}
