package workconfirmations

import (
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	// Tất cả routes đều yêu cầu authentication
	r.Use(middleware.AuthMiddleware())

	r.POST("", Create())                            // Create a new work confirmation
	r.GET("", List())                               // List all work confirmations
	r.GET(":id", GetByID())                         // Get a work confirmation by ID
	r.PUT(":id", Update())                          // Update a work confirmation by ID
	r.POST(":id/approve", Approve())                // Approve a work confirmation by ID
	r.POST(":id/reject", Reject())                  // Reject a work confirmation by ID
	r.GET(":id/download", Download())               // Download a work confirmation by ID
	r.POST("download-multiple", DownloadMultiple()) // Download multiple work confirmations
}
