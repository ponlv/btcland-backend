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

	// Standard information
	FullName                string   `json:"full_name" bson:"full_name"`                                   // Họ và tên đầy đủ của công dân
	ArtistTitle             []string `json:"artist_title" bson:"artist_title"`                             // Danh hiệu nghệ sĩ (vd Nghệ sĩ ưu tú, Nghệ sĩ nhân dân, Nhà giáo ưu tú, Nhà giáo nhân dân...)
	Avatar                  string   `json:"avatar" bson:"avatar"`                                         // Ảnh đại diện
	Email                   string   `json:"email" bson:"email"`                                           // Họ và tên đầy đủ của công dân
	PhoneNumber             string   `json:"phone_number" bson:"phone_number"`                             // Thông tin số điện thoại
	BirthDate               string   `json:"birth_date" bson:"birth_date"`                                 // Ngày tháng năm sinh của công dân, định dạng dd-MM-yyyy
	IdCardExpireDate        string   `json:"id_card_expire_date" bson:"id_card_expire_date"`               // Ngày hết hạn của CCCD (mới nhất đến thời điểm chia sẻ)
	PreviousCardNumber      string   `json:"previous_card_number" bson:"previous_card_number"`             // Số CCCD cũ (nếu có)
	CitizenPid              string   `json:"citizen_id" bson:"citizen_id"`                                 // Số định danh/ CCCD của công dân
	DateOfIssue             string   `json:"date_of_issue" bson:"date_of_issue"`                           // Ngày cấp số định danh/ CCCD của công dân
	IssuingAuthority        string   `json:"issuing_authority" bson:"issuing_authority"`                   // Nơi cấp định danh/ CCCD của công dân
	NationalityCode         string   `json:"nationality_code" bson:"nationality_code"`                     // Mã quốc tịch của công dân
	PermanentAddress        string   `json:"permanent_address" bson:"permanent_address"`                   // Địa chỉ thường trú (chi tiết) của công dân
	PermanentVillageCode    string   `json:"permanent_village_code" bson:"permanent_village_code"`         // Mã xã thường trú (Ref danh mục xã)
	PermanentVillageText    string   `json:"permanent_village_text" bson:"permanent_village_text"`         // Tên xã thường trú
	PermanentDistrictCode   string   `json:"permanent_district_code" bson:"permanent_district_code"`       // Mã huyện thường trú (Ref danh mục huyện)
	PermanentDistrictText   string   `json:"permanent_district_text" bson:"permanent_district_text"`       // Tên huyện thường trú
	PermanentCityCode       string   `json:"permanent_city_code" bson:"permanent_city_code"`               // Mã tỉnh thường trú (Ref danh mục tỉnh)
	PermanentCityText       string   `json:"permanent_city_text" bson:"permanent_city_text"`               // Tên tỉnh thường trú
	LivingPlaceAddress      string   `json:"living_place_address" bson:"living_place_address"`             // Địa chỉ nơi ở hiện tại (chi tiết) của công dân
	LivingPlaceVillageCode  string   `json:"living_place_village_code" bson:"living_place_village_code"`   // Mã xã nơi ở hiện tại (Ref danh mục xã)
	LivingPlaceVillageText  string   `json:"living_place_village_text" bson:"living_place_village_text"`   // Tên xã nơi ở hiện tại
	LivingPlaceDistrictCode string   `json:"living_place_district_code" bson:"living_place_district_code"` // Mã huyện nơi ở hiện tại (Ref danh mục huyện)
	LivingPlaceDistrictText string   `json:"living_place_district_text" bson:"living_place_district_text"` // Tên huyện nơi ở hiện tại
	LivingPlaceCityCode     string   `json:"living_place_city_code" bson:"living_place_city_code"`         // Mã tỉnh nơi ở hiện tại (Ref danh mục tỉnh)
	LivingPlaceCityText     string   `json:"living_place_city_text" bson:"living_place_city_text"`         // Tên tỉnh nơi ở hiện tại
	Religion                string   `json:"religion" bson:"religion"`                                     // Tôn giáo của công dân

	FatherName         string `json:"father_name" bson:"father_name"`                 // Họ và tên cha của công dân
	MotherName         string `json:"mother_name" bson:"mother_name"`                 // Họ và tên mẹ của công dân
	PartnerName        string `json:"partner_name" bson:"partner_name"`               // Họ và tên vợ/chồng của công dân
	RepresentativeName string `json:"representative_name" bson:"representative_name"` // Họ và tên người đại diện (nếu có)
	RepresentativeJob  string `json:"representative_job" bson:"representative_job"`   // Nghề nghiệp người đại diện (nếu có)

	IdentifyLevel int      `json:"identify_level,omitempty" bson:"identify_level,omitempty"`
	UserType      UserType `json:"user_type,omitempty" bson:"user_type,omitempty"`

	IsKYC         bool   `json:"is_kyc,omitempty" bson:"is_kyc"`
	IsShareInfo   bool   `json:"is_share_info" bson:"is_share_info"`
	IsVerifyPhone bool   `json:"is_verify_phone" bson:"is_verify_phone"`
	IsVerifyEmail bool   `json:"is_verify_email" bson:"is_verify_email"`
	IsSetPassword bool   `json:"is_set_password" bson:"is_set_password"`
	ShareInfoId   string `json:"share_info_id" bson:"share_info_id"`
	KYCSignature  string `json:"kyc_signature" bson:"kyc_signature"`
	KYCResponseId string `json:"kyc_response_id" bson:"kyc_response_id"`

	// Veriff KYC fields
	VeriffSessionId  string    `json:"veriff_session_id" bson:"veriff_session_id"`
	VeriffAttemptId  string    `json:"veriff_attempt_id" bson:"veriff_attempt_id"`
	VeriffStatus     string    `json:"veriff_status" bson:"veriff_status"` // submitted, approved, declined, etc.
	VeriffCode       int       `json:"veriff_code" bson:"veriff_code"`
	VeriffAction     string    `json:"veriff_action" bson:"veriff_action"`
	VeriffFeature    string    `json:"veriff_feature" bson:"veriff_feature"`
	VeriffEndUserId  string    `json:"veriff_end_user_id" bson:"veriff_end_user_id"`
	VeriffVendorData string    `json:"veriff_vendor_data" bson:"veriff_vendor_data"`
	VeriffUpdatedAt  time.Time `json:"veriff_updated_at" bson:"veriff_updated_at"`

	IsDelete         bool      `json:"is_delete,omitempty" bson:"is_delete"`
	DeletedAt        time.Time `json:"deleted_at" bson:"deleted_at,omitempty"`
	DeleteReasonCode string    `json:"delete_reason_code" bson:"delete_reason_code"`
	DeleteReasonNote string    `json:"delete_reason_note" bson:"delete_reason_note"`
	GenderCode       string    `json:"gender_code" bson:"gender_code"`

	Password string `json:"password" bson:"password"`

	// OCR
	ScanOCRDate time.Time `json:"scan_ocr_date" bson:"scan_ocr_date,omitempty"`
	MRZData     MRZData   `json:"mrz_data" bson:"mrz_data,omitempty"` // Dữ liệu MRZ được quét từ CCCD

	IsVerifyFace         bool    `json:"is_verify_face" bson:"is_verify_face"`                   // Trạng thái xác minh khuôn mặt
	FaceAccuracy         float64 `json:"face_accuracy" bson:"face_accuracy"`                     // Độ chính xác của khuôn mặt khi so sánh với ảnh trên CCCD
	VerifyFaceSignature  string  `json:"verify_face_signature" bson:"verify_face_signature"`     // Chữ ký xác minh khuôn mặt
	VerifyFaceResponseId string  `json:"verify_face_response_id" bson:"verify_face_response_id"` // ID phản hồi xác minh khuôn mặt

	OAuthProvider *OAuthProvider `json:"oauth_provider,omitempty" bson:"oauth_provider,omitempty"` // Thông tin nhà cung cấp OAuth (nếu có)
	IsAuthor      bool           `json:"is_author" bson:"is_author"`

	WalletAddress string `json:"wallet_address" bson:"wallet_address"`

	// User Settings
	TwoFactorAuthEnabled      bool `json:"two_factor_auth_enabled" bson:"two_factor_auth_enabled"`         // 2FA enabled
	EmailNotificationsEnabled bool `json:"email_notifications_enabled" bson:"email_notifications_enabled"` // Email notifications enabled
	ProfileProtectionEnabled  bool `json:"profile_protection_enabled" bson:"profile_protection_enabled"`   // Profile protection enabled
	BlockConnectionsEnabled   bool `json:"block_connections_enabled" bson:"block_connections_enabled"`     // Block connections enabled

	// 2FA specific fields
	TwoFactorSecret      string   `json:"two_factor_secret" bson:"two_factor_secret"`             // 2FA secret key
	TwoFactorBackupCodes []string `json:"two_factor_backup_codes" bson:"two_factor_backup_codes"` // 2FA backup codes

	// Introduction/Profile data
	IntroductionData *IntroductionData `json:"introduction_data,omitempty" bson:"introduction_data,omitempty"`

	// Package subscription fields
	AccountUpgraded  bool      `json:"account_upgraded,omitempty" bson:"account_upgraded,omitempty"`     // Account upgraded status
	PackageID        string    `json:"package_id,omitempty" bson:"package_id,omitempty"`                 // Current package ID (free, basic, advanced, premium)
	PackageExpiresAt time.Time `json:"package_expires_at,omitempty" bson:"package_expires_at,omitempty"` // Package expiration date
}

type OAuthProvider struct {
	ProviderName oauth2.Provider `json:"provider_name" bson:"provider_name"` // Tên nhà cung cấp OAuth (ví dụ: Google, Facebook)
	ProviderID   string          `json:"provider_id" bson:"provider_id"`     // ID người dùng từ nhà cung cấp OAuth
}

type MRZData struct {
	MRZ    []string            `json:"mrz" bson:"mrz"`
	Fields map[string]MRZField `json:"fields" bson:"fields"`
}

type MRZField struct {
	Value      string `json:"value" bson:"value"`
	RawValue   string `json:"raw_value" bson:"raw_value"`
	CheckDigit string `json:"check_digit" bson:"check_digit"`
	IsValid    bool   `json:"is_valid" bson:"is_valid"`
}

type Blockchain struct {
	ChainId   string `json:"chain_id" bson:"chain_id,omitempty"`
	ChainName string `json:"chain_name" bson:"chain_name,omitempty"`
	ChainType string `json:"chain_type" bson:"chain_type,omitempty"`
	PublicKey string `json:"public_key" bson:"public_key,omitempty"`
}

func (User) CollectionName() string {
	return "user"
}
