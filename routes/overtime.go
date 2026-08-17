package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupOvertimeRoutes(r *gin.Engine) {
	overtimeGroup := r.Group("/api/overtimes")
	overtimeGroup.Use(middlewares.RequireAuth)
	{
		overtimeGroup.POST("", controllers.CreateOvertime)
		overtimeGroup.GET("", controllers.GetOvertimes)
		overtimeGroup.GET("/:id", controllers.GetOvertime)
		overtimeGroup.PUT("/:id", controllers.UpdateOvertime)
		overtimeGroup.DELETE("/:id", controllers.DeleteOvertime)
	}
}
