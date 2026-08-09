package routes

import (
	"backend-wifi/controllers"
	"backend-wifi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAttendanceRoutes(r *gin.Engine) {
	// Employee Routes
	employeeRoutes := r.Group("/api/employee/attendance")
	employeeRoutes.Use(middlewares.RequireAuth)
	employeeRoutes.Use(middlewares.RequireRole("employee"))
	{
		employeeRoutes.POST("/clock-in", controllers.ClockIn)
		employeeRoutes.POST("/clock-out", controllers.ClockOut)
		employeeRoutes.POST("/izin", controllers.RequestIzin)
		employeeRoutes.GET("/today", controllers.GetTodayAttendance)
		employeeRoutes.GET("/history", controllers.GetAttendanceHistory)
	}

	// Admin Routes
	adminRoutes := r.Group("/api/admin/attendance")
	adminRoutes.Use(middlewares.RequireAuth)
	adminRoutes.Use(middlewares.RequireRole("admin"))
	{
		adminRoutes.GET("/", controllers.GetAllAttendance)
	}
}
