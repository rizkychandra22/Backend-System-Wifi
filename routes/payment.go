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
		paymentRoutes.GET("/", middlewares.RequireRole("admin"), controllers.GetAllPayments)
		paymentRoutes.POST("/", middlewares.RequireRole("admin"), controllers.CreatePayment)
		paymentRoutes.PUT("/:id", middlewares.RequireRole("admin"), controllers.UpdatePayment)
		paymentRoutes.DELETE("/:id", middlewares.RequireRole("admin"), controllers.DeletePayment)
		
		// History & PDF (Admin & Employee? Actually Customer & Admin)
		paymentRoutes.GET("/history/:customer_id", controllers.GetCustomerPayments)
		paymentRoutes.GET("/:id/pdf", controllers.GeneratePaymentPDF)
	}
}
