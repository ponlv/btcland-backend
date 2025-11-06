package departments

import (
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	// Tất cả routes đều yêu cầu authentication
	r.Use(middleware.AuthMiddleware())

	// Routes cho Lãnh đạo (Leader) - quản lý tất cả teams/departments
	r.GET("", List())                                              // GET /teams - Lấy danh sách tất cả team (chỉ lãnh đạo)
	r.GET("users", ListUsers())                                    // GET /teams/users - Lấy danh sách tất cả users (chỉ lãnh đạo)
	r.PUT("users/:user_id/role", UpdateUserRole())                 // PUT /teams/users/:user_id/role - Cập nhật vai trò của user (chỉ lãnh đạo)
	r.GET(":id", GetByID())                                        // GET /teams/:id - Lấy chi tiết team
	r.POST("", Create())                                           // POST /teams - Tạo team mới
	r.PUT(":id", Update())                                         // PUT /teams/:id - Cập nhật team
	r.DELETE(":id", Delete())                                      // DELETE /teams/:id - Xóa team
	r.POST(":id/assign-manager", AssignManager())                  // POST /teams/:id/assign-manager - Gán quản lý
	r.GET(":id/employees", ListEmployees())                        // GET /teams/:id/employees - Lấy danh sách nhân viên
	r.POST(":id/add-employee", AddEmployee())                      // POST /teams/:id/add-employee - Thêm nhân viên
	r.DELETE(":id/remove-employee/:employee_id", RemoveEmployee()) // DELETE /teams/:id/remove-employee/:employee_id - Xóa nhân viên
}
