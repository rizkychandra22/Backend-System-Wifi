package routes

import (
	"backend-wifi/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/controllers.Login", controllers.Login)
	}
}
