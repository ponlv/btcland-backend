package userdevicecol

import (
	"time"

	"api/internal/mongodb"
)

type UserDevice struct {
	mongodb.DefaultModel `json:",inline" bson:",inline,omitnested"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	IP          string    `json:"ip" bson:"ip"`
	UserId      string    `json:"user_id" bson:"user_id"`
	DeviceID    string    `json:"device_id" bson:"device_id"`
	DeviceName  string    `json:"device_name" bson:"device_name"`
	DeviceToken string    `json:"device_token" bson:"device_token"`
	Platform    string    `json:"platform" bson:"platform"`
	IsEnable    bool      `json:"is_enable" bson:"is_enable"`
	IsCurrent   bool      `json:"is_current" bson:"is_current"`
	LastLoginAt time.Time `json:"last_login_at" bson:"last_login_at,omitempty"`
}

func (UserDevice) CollectionName() string {
	return "user_device"
}
