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
		// Only admin can input payment
		paymentRoutes.POST("", middlewares.RequireRole("admin"), controllers.CreatePayment)
		
		// History & PDF (Admin & Employee? Actually Customer & Admin)
		paymentRoutes.GET("/history/:customer_id", controllers.GetCustomerPayments)
		paymentRoutes.GET("/:id/pdf", controllers.GeneratePaymentPDF)
	}
}
