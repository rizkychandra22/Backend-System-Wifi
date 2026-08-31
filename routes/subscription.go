package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupSubscriptionRoutes(r *gin.Engine) {
	subRoutes := r.Group("/api/subscriptions")
	subRoutes.Use(middlewares.RequireAuth)
	{
		// Admin & Employee can manage subscriptions
		subRoutes.GET("", middlewares.RequireRole("admin", "employee"), controllers.GetAllSubscriptions)
		subRoutes.POST("", middlewares.RequireRole("admin", "employee"), controllers.CreateOrUpdateSubscription)
		subRoutes.DELETE("/:id", middlewares.RequireRole("admin"), controllers.DeleteSubscription)

		// Customers can get their own, Admin & Employee can get any
		subRoutes.GET("/customer/:id", controllers.GetSubscriptionByCustomerID)
	}
}
