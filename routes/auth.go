package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/login", controllers.Login)
		authRoutes.GET("/admin-contact", controllers.GetAdminContact)

		// Protected routes for authenticated users
		protectedAuth := authRoutes.Group("/")
		protectedAuth.Use(middlewares.RequireAuth)
		{
			protectedAuth.PUT("/profile", controllers.UpdateProfile)
			protectedAuth.PUT("/profile/password", controllers.UpdatePassword)
		}
	}
}
