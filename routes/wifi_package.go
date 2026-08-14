package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupWifiServiceRoutes(r *gin.Engine) {
	wsRoutes := r.Group("/api/admin/wifi-services")
	wsRoutes.Use(middlewares.RequireAuth)
	wsRoutes.Use(middlewares.RequireRole("admin"))
	{
		wsRoutes.POST("", controllers.CreateWifiService)
		wsRoutes.GET("", controllers.GetWifiServices)
		wsRoutes.PUT("/:id", controllers.UpdateWifiService)
		wsRoutes.DELETE("/:id", controllers.DeleteWifiService)
	}
}
