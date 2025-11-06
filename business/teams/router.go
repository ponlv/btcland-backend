package teams

import (
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	// Tất cả routes đều yêu cầu authentication
	r.Use(middleware.AuthMiddleware())

	// Routes cho Quản lý (Manager) - quản lý team của mình
	r.GET("my-team", GetMyTeam())                              // Get my team
	r.GET("my-employees", ListMyEmployees())                   // List my employees
	r.GET("my-employees/:employee_id", GetEmployeeByID())      // Get employee by ID
	r.POST("add-employee", AddEmployee())                      // Add employee
	r.DELETE("remove-employee/:employee_id", RemoveEmployee()) // Remove employee
	r.PUT("my-employees/:employee_id", UpdateEmployee())       // Update employee
}
