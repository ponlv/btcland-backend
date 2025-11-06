package workconfirmations

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/usercol"
	"api/schema/workconfirmationcol"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Download() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][work-confirmations][download]")

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			code := response.ErrorResponse("ID is required")
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

		// Tìm đơn
		workConfirmation, err := workconfirmationcol.FindByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				code := response.ErrorResponse("Work confirmation not found")
				c.JSON(http.StatusNotFound, code)
				c.Abort()
				return
			}
			logger.Err(err).Msg("failed to get work confirmation")
			code := response.ErrorResponse("Failed to get work confirmation")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy thông tin người tạo
		creator, err := usercol.FindWithUserID(c.Request.Context(), workConfirmation.CreatedBy)
		if err != nil {
			logger.Err(err).Msg("failed to get creator info")
		}

		// Tạo file Excel
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				logger.Err(err).Msg("failed to close excel file")
			}
		}()

		sheetName := "Đơn xác nhận công tác"
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
		headers := []string{"ID", "Ngày công tác", "Nội dung", "Người tạo", "Email", "Vai trò", "Trạng thái", "Ngày tạo", "Số lượng ảnh"}
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

		// Điền dữ liệu
		creatorName := "N/A"
		creatorEmail := "N/A"
		if creator != nil {
			creatorName = creator.FullName
			creatorEmail = creator.Email
		}

		statusText := string(workConfirmation.Status)
		switch workConfirmation.Status {
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
			workConfirmation.GetIDString(),
			workConfirmation.Date,
			workConfirmation.Content,
			creatorName,
			creatorEmail,
			workConfirmation.CreatorRole.Text(),
			statusText,
			workConfirmation.CreatedAt.Format("2006-01-02 15:04:05"),
			len(workConfirmation.Photos),
		}

		for i, value := range rowData {
			cell := fmt.Sprintf("%c2", 'A'+i)
			f.SetCellValue(sheetName, cell, value)
		}

		// Thêm sheet cho danh sách ảnh
		photosSheetName := "Danh sách ảnh"
		photosIndex, err := f.NewSheet(photosSheetName)
		if err == nil {
			f.SetActiveSheet(photosIndex)
			f.SetCellValue(photosSheetName, "A1", "STT")
			f.SetCellValue(photosSheetName, "B1", "Tên file")
			f.SetCellValue(photosSheetName, "C1", "URL")
			f.SetCellValue(photosSheetName, "D1", "Ngày upload")

			for i, photo := range workConfirmation.Photos {
				row := i + 2
				f.SetCellValue(photosSheetName, fmt.Sprintf("A%d", row), i+1)
				f.SetCellValue(photosSheetName, fmt.Sprintf("B%d", row), photo.Filename)
				f.SetCellValue(photosSheetName, fmt.Sprintf("C%d", row), photo.URL)
				f.SetCellValue(photosSheetName, fmt.Sprintf("D%d", row), photo.UploadedAt.Format("2006-01-02 15:04:05"))
			}
		}

		// Lưu file tạm thời
		tmpFile, err := os.CreateTemp("", "work-confirmation-*.xlsx")
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
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=work-confirmation-%s.xlsx", id))
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
	}
}
