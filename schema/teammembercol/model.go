package teammembercol

import (
	"time"

	"api/internal/mongodb"
)

type TeamMember struct {
	mongodb.DefaultModel `json:",inline" bson:",inline,omitnested"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	// Quan hệ manager - employee
	ManagerID  string    `json:"manager_id" bson:"manager_id"`   // user_id của quản lý
	EmployeeID string    `json:"employee_id" bson:"employee_id"` // user_id của nhân viên
	JoinedAt   time.Time `json:"joined_at" bson:"joined_at"`     // Ngày tham gia team

	// Soft delete
	IsDelete  bool      `json:"is_delete,omitempty" bson:"is_delete"`
	DeletedAt time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func (TeamMember) CollectionName() string {
	return "team_member"
}
