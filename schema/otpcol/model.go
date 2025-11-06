package otpcol

import (
	"time"

	"api/internal/mongodb"
)

type OTP struct {
	mongodb.DefaultModel `json:",inline" bson:",inline,omitnested"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	Code       string    `json:"code" bson:"code,omitempty"`
	Phone      string    `json:"phone" bson:"phone,omitempty"`
	Name       string    `json:"name" bson:"name,omitempty"`
	Email      string    `json:"email" bson:"email,omitempty"`
	UserId     string    `json:"user_id" bson:"user_id,omitempty"`
	Attempts   int       `json:"attempts" bson:"attempts,omitempty"`
	Type       string    `json:"type" bson:"type,omitempty"`
	IsVerified bool      `json:"is_verified" bson:"is_verified,omitempty"`
	ExpireAt   time.Time `json:"expire_at" bson:"expire_at"`
}

func (OTP) CollectionName() string {
	return "otp_code"
}
