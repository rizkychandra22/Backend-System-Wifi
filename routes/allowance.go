package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAllowanceRoutes(r *gin.Engine) {
	allowanceGroup := r.Group("/api/allowances")
	allowanceGroup.Use(middlewares.RequireAuth)
	{
		allowanceGroup.POST("", controllers.CreateAllowance)
		allowanceGroup.GET("", controllers.GetAllowances)
		allowanceGroup.GET("/:id", controllers.GetAllowance)
		allowanceGroup.PUT("/:id", controllers.UpdateAllowance)
		allowanceGroup.DELETE("/:id", controllers.DeleteAllowance)
	}
}
