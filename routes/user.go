package routes

import (
	"backend-wifi/middlewares"
	"backend-wifi/controllers"
"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	// Group di /api/admin dengan perlindungan JWT dan Role Admin
	adminRoutes := r.Group("/api/admin")
	adminRoutes.Use(middlewares.RequireAuth)
	adminRoutes.Use(middlewares.RequireRole("admin"))
	{
		adminRoutes.POST("/users", controllers.CreateUser)
		adminRoutes.GET("/users", controllers.GetUsers)
		adminRoutes.GET("/users/:id", controllers.GetUserByID)
		adminRoutes.PUT("/users/:id", controllers.UpdateUser)
		adminRoutes.DELETE("/users/:id", controllers.DeleteUser)
		adminRoutes.PUT("/users/:id/reset-ip", controllers.ResetUserIP)
	}
}
