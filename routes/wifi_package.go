package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupWifiPackageRoutes(r *gin.Engine) {
	// Public GET route
	r.GET("/api/admin/wifi-packages", controllers.GetWifiPackages)

	// Protected write routes
	wsRoutes := r.Group("/api/admin/wifi-packages")
	wsRoutes.Use(middlewares.RequireAuth)
	{
		adminRoutes := wsRoutes.Group("")
		adminRoutes.Use(middlewares.RequireRole("admin"))
		{
			adminRoutes.POST("", controllers.CreateWifiPackage)
			adminRoutes.PUT("/:id", controllers.UpdateWifiPackage)
			adminRoutes.DELETE("/:id", controllers.DeleteWifiPackage)
		}
	}
}
