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
		customerRoutes.GET("/:id/subscription", controllers.GetCustomerSubscription)
	}
}
