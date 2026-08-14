package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupWifiPackageRoutes(r *gin.Engine) {
	wsRoutes := r.Group("/api/admin/wifi-packages")
	wsRoutes.Use(middlewares.RequireAuth)
	wsRoutes.Use(middlewares.RequireRole("admin"))
	{
		wsRoutes.POST("", controllers.CreateWifiPackage)
		wsRoutes.GET("", controllers.GetWifiPackages)
		wsRoutes.PUT("/:id", controllers.UpdateWifiPackage)
		wsRoutes.DELETE("/:id", controllers.DeleteWifiPackage)
	}
}
