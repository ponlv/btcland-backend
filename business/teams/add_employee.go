package teams

import (
	"errors"
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teammembercol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AddEmployeeRequest struct {
	Email string `json:"email" binding:"required"`
}

func AddEmployee() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][teams][add_employee]")

	return func(c *gin.Context) {
		var req AddEmployeeRequest
		if err := c.BindJSON(&req); err != nil {
			code := response.ErrorResponse(err.Error())
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		// Lấy user từ context
		userInterface, exists := c.Get("current_user")
		if !exists {
			code := response.ErrorResponse("Unauthorized")
			c.JSON(http.StatusUnauthorized, code)
			c.Abort()
			return
		}

		user, ok := userInterface.(*usercol.User)
		if !ok {
			code := response.ErrorResponse("Invalid user")
			c.JSON(http.StatusUnauthorized, code)
			c.Abort()
			return
		}

		// Kiểm tra role - chỉ manager mới có quyền
		if user.Role != usercol.RoleManager {
			code := response.ErrorResponse("Only managers can add employees")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Tìm user theo email
		employee, err := usercol.FindWithEmail(c.Request.Context(), req.Email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				code := response.ErrorResponse("User with this email not found")
				c.JSON(http.StatusNotFound, code)
				c.Abort()
				return
			}
			logger.Err(err).Msg("failed to find user by email")
			code := response.ErrorResponse("Failed to find user")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Kiểm tra employee không phải là chính manager
		if employee.GetIDString() == user.GetIDString() {
			code := response.ErrorResponse("Cannot add yourself as employee")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Kiểm tra employee đã thuộc team này chưa
		belongs, err := teammembercol.CheckEmployeeBelongsToManager(c.Request.Context(), user.GetIDString(), employee.GetIDString())
		if err != nil {
			logger.Err(err).Msg("failed to check employee relationship")
			code := response.ErrorResponse("Failed to check employee relationship")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		if belongs {
			code := response.ErrorResponse("Employee already in your team")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Tạo team member relationship
		teamMember := &teammembercol.TeamMember{
			ManagerID:  user.GetIDString(),
			EmployeeID: employee.GetIDString(),
		}

		_, err = teammembercol.Create(c.Request.Context(), teamMember)
		if err != nil {
			logger.Err(err).Msg("failed to add employee to team")
			code := response.ErrorResponse("Failed to add employee to team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy lại team member vừa tạo
		created, err := teammembercol.FindByManagerAndEmployee(c.Request.Context(), user.GetIDString(), employee.GetIDString())
		if err != nil {
			logger.Err(err).Msg("failed to get created team member")
			code := response.ErrorResponse("Employee added but failed to retrieve")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		responseData := map[string]interface{}{
			"id":             employee.GetIDString(),
			"full_name":      employee.FullName,
			"email":          employee.Email,
			"joined_at":      created.JoinedAt,
			"team_member_id": created.GetIDString(),
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}
