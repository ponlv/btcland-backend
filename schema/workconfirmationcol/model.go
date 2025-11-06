package workconfirmationcol

import (
	"api/schema/usercol"
	"time"

	"api/internal/mongodb"
)

type WorkConfirmationStatus string

const (
	StatusPendingManager WorkConfirmationStatus = "pending_manager"
	StatusPendingLeader  WorkConfirmationStatus = "pending_leader"
	StatusApproved       WorkConfirmationStatus = "approved"
	StatusRejected       WorkConfirmationStatus = "rejected"
)

type Photo struct {
	URL        string    `json:"url" bson:"url"`
	Filename   string    `json:"filename" bson:"filename"`
	UploadedAt time.Time `json:"uploaded_at" bson:"uploaded_at"`
}

type ApprovalInfo struct {
	ApprovedBy string    `json:"approved_by" bson:"approved_by"` // user_id
	ApprovedAt time.Time `json:"approved_at" bson:"approved_at"`
	Comment    string    `json:"comment" bson:"comment"`
}

type RejectionInfo struct {
	RejectedBy string    `json:"rejected_by" bson:"rejected_by"` // user_id
	RejectedAt time.Time `json:"rejected_at" bson:"rejected_at"`
	Reason     string    `json:"reason" bson:"reason"`
}

type WorkConfirmation struct {
	mongodb.DefaultModel `json:",inline" bson:",inline,omitnested"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	// Người tạo đơn
	CreatedBy  string      `json:"created_by" bson:"created_by"`   // user_id
	CreatorRole usercol.Role `json:"creator_role" bson:"creator_role"` // employee, manager, leader, assistant_director

	// Thông tin đơn
	Date    string  `json:"date" bson:"date"`       // YYYY-MM-DD
	Content string  `json:"content" bson:"content"` // Nội dung công tác
	Photos  []Photo `json:"photos" bson:"photos"`   // Danh sách hình ảnh

	// Trạng thái
	Status WorkConfirmationStatus `json:"status" bson:"status"`

	// Xác nhận từ quản lý (chỉ khi đơn từ nhân viên)
	ManagerApproval *ApprovalInfo `json:"manager_approval,omitempty" bson:"manager_approval,omitempty"`

	// Xác nhận từ lãnh đạo
	LeaderApproval *ApprovalInfo `json:"leader_approval,omitempty" bson:"leader_approval,omitempty"`

	// Từ chối
	Rejection *RejectionInfo `json:"rejection,omitempty" bson:"rejection,omitempty"`

	// Soft delete
	IsDelete bool      `json:"is_delete,omitempty" bson:"is_delete"`
	DeletedAt time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func (WorkConfirmation) CollectionName() string {
	return "work_confirmation"
}

