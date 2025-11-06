package workconfirmations

import (
	"fmt"
	"net/http"
	"os"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/usercol"
	"api/schema/workconfirmationcol"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type DownloadMultipleRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

func DownloadMultiple() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][work-confirmations][download_multiple]")

	return func(c *gin.Context) {
		var req DownloadMultipleRequest
		if err := c.BindJSON(&req); err != nil {
			code := response.ErrorResponse(err.Error())
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		if len(req.IDs) == 0 {
			code := response.ErrorResponse("IDs are required")
			c.JSON(http.StatusBadRequest, code)
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

		// Kiểm tra quyền - Assistant Director và Leader có quyền tải
		userRole := user.Role
		if userRole != usercol.RoleAssistantDirector && userRole != usercol.RoleLeader {
			code := response.ErrorResponse("Only Assistant Director and Leader can download work confirmations")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Tìm tất cả các đơn
		workConfirmations, err := workconfirmationcol.FindByIDs(c.Request.Context(), req.IDs)
		if err != nil {
			logger.Err(err).Msg("failed to get work confirmations")
			code := response.ErrorResponse("Failed to get work confirmations")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		if len(workConfirmations) == 0 {
			code := response.ErrorResponse("No work confirmations found")
			c.JSON(http.StatusNotFound, code)
			c.Abort()
			return
		}

		// Tạo file Excel
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				logger.Err(err).Msg("failed to close excel file")
			}
		}()

		sheetName := "Danh sách đơn xác nhận công tác"
		index, err := f.NewSheet(sheetName)
		if err != nil {
			logger.Err(err).Msg("failed to create sheet")
			code := response.ErrorResponse("Failed to create Excel file")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}
		f.SetActiveSheet(index)

		// Xóa sheet mặc định
		f.DeleteSheet("Sheet1")

		// Đặt header
		headers := []string{"STT", "ID", "Ngày công tác", "Nội dung", "Người tạo", "Email", "Vai trò", "Trạng thái", "Ngày tạo", "Số lượng ảnh"}
		for i, header := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue(sheetName, cell, header)
		}

		// Style cho header
		headerStyle, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		})
		if err == nil {
			f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), headerStyle)
		}

		// Điền dữ liệu cho từng đơn
		for idx, wc := range workConfirmations {
			row := idx + 2

			// Lấy thông tin người tạo
			creator, err := usercol.FindWithUserID(c.Request.Context(), wc.CreatedBy)
			creatorName := "N/A"
			creatorEmail := "N/A"
			if err == nil && creator != nil {
				creatorName = creator.FullName
				creatorEmail = creator.Email
			}

			statusText := string(wc.Status)
			switch wc.Status {
			case workconfirmationcol.StatusPendingManager:
				statusText = "Chờ quản lý xác nhận"
			case workconfirmationcol.StatusPendingLeader:
				statusText = "Chờ lãnh đạo xác nhận"
			case workconfirmationcol.StatusApproved:
				statusText = "Đã duyệt"
			case workconfirmationcol.StatusRejected:
				statusText = "Đã từ chối"
			}

			rowData := []interface{}{
				idx + 1,
				wc.GetIDString(),
				wc.Date,
				wc.Content,
				creatorName,
				creatorEmail,
				wc.CreatorRole.Text(),
				statusText,
				wc.CreatedAt.Format("2006-01-02 15:04:05"),
				len(wc.Photos),
			}

			for i, value := range rowData {
				cell := fmt.Sprintf("%c%d", 'A'+i, row)
				f.SetCellValue(sheetName, cell, value)
			}
		}

		// Set column width
		for i := 0; i < len(headers); i++ {
			col := string(rune('A' + i))
			f.SetColWidth(sheetName, col, col, 15)
		}

		// Lưu file tạm thời
		tmpFile, err := os.CreateTemp("", "work-confirmations-*.xlsx")
		if err != nil {
			logger.Err(err).Msg("failed to create temp file")
			code := response.ErrorResponse("Failed to create Excel file")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		// Lưu Excel file
		if err := f.SaveAs(tmpFile.Name()); err != nil {
			logger.Err(err).Msg("failed to save Excel file")
			code := response.ErrorResponse("Failed to save Excel file")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Đọc file Excel
		excelData, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			logger.Err(err).Msg("failed to read Excel file")
			code := response.ErrorResponse("Failed to read Excel file")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Set headers để download
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=work-confirmations.xlsx")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
	}
}
