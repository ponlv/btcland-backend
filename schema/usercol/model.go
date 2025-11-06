package usercol

import (
	"api/services/oauth2"
	"time"

	"api/internal/mongodb"
)

type User struct {
	mongodb.DefaultModel `json:",inline" bson:",inline,omitnested"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	// Basic information
	FullName    string `json:"full_name" bson:"full_name"`       // Họ và tên
	Avatar      string `json:"avatar" bson:"avatar"`             // Ảnh đại diện
	Email       string `json:"email" bson:"email"`               // Email
	PhoneNumber string `json:"phone_number" bson:"phone_number"` // Số điện thoại

	// Work confirmation system role
	Role Role `json:"role,omitempty" bson:"role,omitempty"` // employee, manager, leader, assistant_director

	// Authentication
	Password      string `json:"password" bson:"password"`
	IsSetPassword bool   `json:"is_set_password" bson:"is_set_password"`

	// Verification
	IsVerifyPhone bool `json:"is_verify_phone" bson:"is_verify_phone"`
	IsVerifyEmail bool `json:"is_verify_email" bson:"is_verify_email"`

	// OAuth
	OAuthProvider *OAuthProvider `json:"oauth_provider,omitempty" bson:"oauth_provider,omitempty"` // Thông tin nhà cung cấp OAuth (Google, etc.)

	// Soft delete
	IsDelete  bool      `json:"is_delete,omitempty" bson:"is_delete"`
	DeletedAt time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type OAuthProvider struct {
	ProviderName oauth2.Provider `json:"provider_name" bson:"provider_name"` // Tên nhà cung cấp OAuth (ví dụ: Google)
	ProviderID   string          `json:"provider_id" bson:"provider_id"`     // ID người dùng từ nhà cung cấp OAuth
}

func (User) CollectionName() string {
	return "user"
}
