package usercol

import (
	"time"
)

// IntroductionData represents user introduction/profile information
type IntroductionData struct {
	// Tổng quan (Overview)
	Workplace       string `json:"workplace" bson:"workplace"`               // Nơi làm việc
	HighSchool      string `json:"high_school" bson:"high_school"`           // Trường trung học
	University      string `json:"university" bson:"university"`             // Trường đại học
	CurrentLocation string `json:"current_location" bson:"current_location"` // Sống tại
	FromLocation    string `json:"from_location" bson:"from_location"`       // Đến từ

	// Công việc và học vấn (Work and Education)
	WorkExperience []WorkExperience `json:"work_experience" bson:"work_experience"`
	Education      []Education      `json:"education" bson:"education"`
	Certifications []Certification  `json:"certifications" bson:"certifications"`

	// Nơi từng sống (Places Lived)
	PlacesLived []PlaceLived `json:"places_lived" bson:"places_lived"`

	// Thông tin liên hệ cơ bản (Basic Contact Info)
	Phone   string `json:"phone" bson:"phone"`     // Số điện thoại
	Website string `json:"website" bson:"website"` // Trang website
	Email   string `json:"email" bson:"email"`     // Email

	// Metadata
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// WorkExperience represents work experience entry
type WorkExperience struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	ID        string    `json:"id" bson:"id"`
	Company   string    `json:"company" bson:"company"`
	Position  string    `json:"position" bson:"position"`
	StartDate string    `json:"start_date" bson:"start_date"`
	EndDate   string    `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Current   bool      `json:"current" bson:"current"`
	IsPublic  bool      `json:"is_public" bson:"is_public"`
}

// Education represents education entry
type Education struct {
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	ID          string    `json:"id" bson:"id"`
	Institution string    `json:"institution" bson:"institution"`
	Degree      string    `json:"degree" bson:"degree"`
	Field       string    `json:"field" bson:"field"`
	StartDate   string    `json:"start_date" bson:"start_date"`
	EndDate     string    `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Current     bool      `json:"current" bson:"current"`
	IsPublic    bool      `json:"is_public" bson:"is_public"`
}

// Certification represents certification entry
type Certification struct {
	ID        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Issuer    string    `json:"issuer" bson:"issuer"`
	Date      string    `json:"date" bson:"date"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Link      string    `json:"link" bson:"link"`
	IsPublic  bool      `json:"is_public" bson:"is_public"`
}

// PlaceLived represents places lived entry
type PlaceLived struct {
	ID        string `json:"id" bson:"id"`
	Location  string `json:"location" bson:"location"`
	StartDate string `json:"start_date" bson:"start_date"`
	EndDate   string `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Current   bool   `json:"current" bson:"current"`
}

// IntroductionRequest represents the request structure for introduction data
type IntroductionRequest struct {
	Phone   *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	Website *string `json:"website,omitempty" binding:"omitempty,url"`
}

// IntroductionResponse represents the response structure for introduction data
type IntroductionResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *IntroductionData `json:"data,omitempty"`
}
