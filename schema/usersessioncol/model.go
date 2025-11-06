package usersessioncol

import (
	"time"

	"api/internal/mongodb"
)

type UserSession struct {
	mongodb.DefaultModel `json:",inline" bson:",inline,omitnested"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at,omitempty"`
	ActiveAt             time.Time `json:"active_at" bson:"active_at,omitempty"`
	DeletedAt            time.Time `json:"deleted_at" bson:"deleted_at,omitempty"`

	IP          string `json:"ip" bson:"ip"`
	UserId      string `json:"user_id" bson:"user_id"`
	IsDelete    bool   `json:"is_delete" bson:"is_delete"`
	DeviceID    string `json:"device_id" bson:"device_id"`
	DeviceName  string `json:"device_name" bson:"device_name"`
	DeviceToken string `json:"device_token" bson:"device_token"`
	BrowserName string `json:"browser_name" bson:"browser_name"`
	Platform    string `json:"platform" bson:"platform"`
	IsEnable    bool   `json:"is_enable" bson:"is_enable"`
}

func (UserSession) CollectionName() string {
	return "user_session"
}
