package teamcol

import (
	"time"

	"api/internal/mongodb"
)

type Team struct {
	mongodb.DefaultModel `json:",inline" bson:",inline,omitnested"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	// Thông tin phòng ban/team
	Name        string `json:"name" bson:"name"`                 // Tên phòng ban/team
	Description string `json:"description" bson:"description"`   // Mô tả
	ManagerID   string `json:"manager_id" bson:"manager_id"`     // ID của quản lý (user_id), có thể null nếu chưa gán

	// Soft delete
	IsDelete  bool      `json:"is_delete,omitempty" bson:"is_delete"`
	DeletedAt time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func (Team) CollectionName() string {
	return "team"
}

