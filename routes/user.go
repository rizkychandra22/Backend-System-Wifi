package routes

import (
	"backend-wifi/middlewares"
	"backend-wifi/controllers"
"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/api/users")
	userRoutes.Use(middlewares.RequireAuth)
	{
		userRoutes.POST("", middlewares.RequireRole("admin", "employee"), controllers.CreateUser)
		userRoutes.GET("", middlewares.RequireRole("admin", "employee"), controllers.GetUsers)
		userRoutes.GET("/:id", middlewares.RequireRole("admin", "employee"), controllers.GetUserByID)
		userRoutes.PUT("/:id", middlewares.RequireRole("admin", "employee"), controllers.UpdateUser)
		userRoutes.DELETE("/:id", middlewares.RequireRole("admin"), controllers.DeleteUser)
		userRoutes.PUT("/:id/reset-ip", middlewares.RequireRole("admin"), controllers.ResetUserIP)
	}
}
