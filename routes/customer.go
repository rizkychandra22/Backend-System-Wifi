package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(r *gin.Engine) {
	customerRoutes := r.Group("/api/customers")
	customerRoutes.Use(middlewares.RequireAuth)
	{
		customerRoutes.POST("", controllers.CreateCustomer)
		customerRoutes.GET("", controllers.GetCustomers)
		customerRoutes.GET("/:id/subscription", controllers.GetCustomerSubscription)
	}
}
