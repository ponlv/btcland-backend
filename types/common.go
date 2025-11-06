package types

import (
	"time"
)

type Asset struct {
	Owner          string     `json:"owner" bson:"owner,omitempty"`
	Address        string     `json:"collection" bson:"collection,omitempty"`
	Name           string     `json:"name" bson:"name,omitempty"`
	CollectionType string     `json:"collection_type" bson:"collection_type,omitempty"`
	Symbol         string     `json:"symbol" bson:"symbol,omitempty"`
	CollectionURI  string     `json:"collection_uri" bson:"collection_uri,omitempty"`
	BaseURI        string     `json:"base_uri" bson:"base_uri,omitempty"`
	AvatarURI      string     `json:"avatar_uri" bson:"avatar_uri,omitempty"`
	Description    string     `json:"description" bson:"description,omitempty"`
	ProductType    string     `json:"product_type" bson:"product_type,omitempty"`
	Traits         []string   `json:"traits" bson:"traits,omitempty"`
	ChainId        string     `json:"chain_id" bson:"chain_id,omitempty"`
	Activity       []Activity `json:"activity" bson:"activity,omitempty"`
}

type Activity struct {
	From         string `json:"from" bson:"from,omitempty"`
	To           string `json:"to" bson:"to,omitempty"`
	Action       string `json:"action" bson:"action,omitempty"`
	Price        string `json:"price" bson:"price,omitempty"`
	TokenId      string `json:"token_id" bson:"token_id,omitempty"`
	TokenName    string `json:"token_name" bson:"token_name,omitempty"`
	Image        string `json:"image" bson:"image,omitempty"`
	Tx           string `json:"tx" bson:"tx,omitempty"`
	Date         string `json:"date" bson:"date,omitempty"`
	PaymentToken struct {
		Address string `json:"address" bson:"address,omitempty"`
		Name    string `json:"name" bson:"name,omitempty"`
		Symbol  string `json:"symbol" bson:"symbol,omitempty"`
	} `json:"payment_token" bson:"payment_token,omitempty"`
}

type NFT struct {
	ItemId         int64  `json:"item_id" bson:"item_id,omitempty"`
	Price          string `json:"price" bson:"price,omitempty"`
	Owner          string `json:"owner" bson:"owner,omitempty"`
	CollectionAddr string `json:"collection_address" bson:"collection_address,omitempty"`
	TokenId        string `json:"token_id" bson:"token_id,omitempty"`
	TokenName      string `json:"token_name" bson:"token_name,omitempty"`
	TokenURI       string `json:"token_uri" bson:"token_uri,omitempty"`
	Description    string `json:"description" bson:"description,omitempty"`
	NFTType        string `json:"nft_type" bson:"nft_type,omitempty"`
	ExpireTime     int64  `json:"expire_time" bson:"expire_time,omitempty"`
	Value          int64  `json:"value" bson:"value,omitempty"`
	ChainId        string `json:"chain_id" bson:"chain_id,omitempty"`
}

type User struct {
	ID            string    `json:"id,omitempty" bson:"id,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt     time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	FullName      string    `json:"full_name" bson:"full_name"`
	Email         string    `json:"email" bson:"email"`
	PhoneNumber   string    `json:"phone_number" bson:"phone_number"`
	Role          string    `json:"role" bson:"role"`
	IsVerifyEmail bool      `json:"is_verify_email" bson:"is_verify_email"`
	IsVerifyPhone bool      `json:"is_verify_phone" bson:"is_verify_phone"`
	IsSetPassword bool      `json:"is_set_password" bson:"is_set_password"`
	IsDelete      bool      `json:"is_delete" bson:"is_delete"`
	// Standard information
	BirthDate               string `json:"birth_date" bson:"birth_date"`                                 // Ngày tháng năm sinh của công dân, định dạng dd-MM-yyyy
	IdCardExpireDate        string `json:"id_card_expire_date" bson:"id_card_expire_date"`               // Ngày hết hạn của CCCD (mới nhất đến thời điểm chia sẻ)
	CitizenPid              string `json:"citizen_id" bson:"citizen_id"`                                 // Số định danh/ CCCD của công dân
	DateOfIssue             string `json:"date_of_issue" bson:"date_of_issue"`                           // Ngày cấp số định danh/ CCCD của công dân
	IssuingAuthority        string `json:"issuing_authority" bson:"issuing_authority"`                   // Nơi cấp định danh/ CCCD của công dân
	NationalityCode         string `json:"nationality_code" bson:"nationality_code"`                     // Mã quốc tịch của công dân
	PermanentAddress        string `json:"permanent_address" bson:"permanent_address"`                   // Địa chỉ thường trú (chi tiết) của công dân
	PermanentVillageCode    string `json:"permanent_village_code" bson:"permanent_village_code"`         // Mã xã thường trú (Ref danh mục xã)
	PermanentVillageText    string `json:"permanent_village_text" bson:"permanent_village_text"`         // Tên xã thường trú
	PermanentDistrictCode   string `json:"permanent_district_code" bson:"permanent_district_code"`       // Mã huyện thường trú (Ref danh mục huyện)
	PermanentDistrictText   string `json:"permanent_district_text" bson:"permanent_district_text"`       // Tên huyện thường trú
	PermanentCityCode       string `json:"permanent_city_code" bson:"permanent_city_code"`               // Mã tỉnh thường trú (Ref danh mục tỉnh)
	PermanentCityText       string `json:"permanent_city_text" bson:"permanent_city_text"`               // Tên tỉnh thường trú
	LivingPlaceAddress      string `json:"living_place_address" bson:"living_place_address"`             // Địa chỉ nơi ở hiện tại (chi tiết) của công dân
	LivingPlaceVillageCode  string `json:"living_place_village_code" bson:"living_place_village_code"`   // Mã xã nơi ở hiện tại (Ref danh mục xã)
	LivingPlaceVillageText  string `json:"living_place_village_text" bson:"living_place_village_text"`   // Tên xã nơi ở hiện tại
	LivingPlaceDistrictCode string `json:"living_place_district_code" bson:"living_place_district_code"` // Mã huyện nơi ở hiện tại (Ref danh mục huyện)
	LivingPlaceDistrictText string `json:"living_place_district_text" bson:"living_place_district_text"` // Tên huyện nơi ở hiện tại
	LivingPlaceCityCode     string `json:"living_place_city_code" bson:"living_place_city_code"`         // Mã tỉnh nơi ở hiện tại (Ref danh mục tỉnh)
	LivingPlaceCityText     string `json:"living_place_city_text" bson:"living_place_city_text"`         // Tên tỉnh nơi ở hiện tại

	IdentifyLevel int `json:"identify_level,omitempty" bson:"identify_level,omitempty"`
}

type ChipCardData struct {
}

type OCRCardData struct {
	DocumentNumber   string `json:"document_number" bson:"document_number"`
	FullName         string `json:"full_name" bson:"full_name"`
	BirthDate        string `json:"birth_date" bson:"birth_date"`
	PlaceOfOrigin    string `json:"place_of_origin" bson:"place_of_origin"`
	Gender           string `json:"gender" bson:"gender"`
	PlaceOfResidence string `json:"place_of_residence" bson:"place_of_residence"`
	Province         string `json:"province" bson:"province"`
	District         string `json:"district" bson:"district"`
	Ward             string `json:"ward" bson:"ward"`
	DistrictCode     string `json:"district_code" bson:"district_code"`
	WardCode         string `json:"ward_code" bson:"ward_code"`
	Street           string `json:"street" bson:"street"`
	Nationality      string `json:"nationality" bson:"nationality"`
	Religion         string `json:"religion" bson:"religion"`
	Ethnicity        string `json:"ethnicity" bson:"ethnicity"`
	ExpiryDate       string `json:"expiry_date" bson:"expiry_date"`
	IssuanceDate     string `json:"issuance_date" bson:"issuance_date"`
	IssuanceBy       string `json:"issuance_by" bson:"issuance_by"`
	DocumentType     string `json:"document_type" bson:"document_type"`
	Identification   string `json:"identification" bson:"identification"`
	PortraitImage    string `json:"portrait_image" bson:"portrait_image"`
	Mrz              string `json:"mrz" bson:"mrz"`
	PassportNumber   string `json:"passport_number" bson:"passport_number"`
}

type Blockchain struct {
	ChainId       string `json:"chain_id" bson:"chain_id,omitempty"`
	ChainName     string `json:"chain_name" bson:"chain_name,omitempty"`
	ChainType     string `json:"chain_type" bson:"chain_type,omitempty"`
	PublicAddress string `json:"public_address" bson:"public_address,omitempty"`
}

type Material struct {
	Title     string    `json:"title" bson:"title"`
	Code      string    `json:"code" bson:"code"`
	IsDisable bool      `json:"is_disable" bson:"is_disable"`
	Product   []Product `json:"product"`
}

type Product struct {
	Code       string       `json:"code" bson:"code"`
	Title      string       `json:"title" bson:"title"`
	IsDisable  bool         `json:"is_disable" bson:"is_disable"`
	Icon       string       `json:"icon" bson:"icon"`
	SubProduct []SubProduct `json:"sub_product" bson:"sub_product"`
}

type SubProduct struct {
	Code  string `json:"code" bson:"code"`
	Title string `json:"title" bson:"title"`
}

type RegistrationAsset struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty" `
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty" `

	CitizenId                 string    `json:"citizen_id" bson:"citizen_id,omitempty"`
	Title                     string    `json:"title" bson:"title"`
	IPCSection                string    `json:"ipc_section" bson:"ipc_section"`
	IPCClass                  string    `json:"ipc_class" bson:"ipc_class"`
	IPCSubClass               string    `json:"ipc_sub_class" bson:"ipc_sub_class"`
	IPCGroup                  string    `json:"ipc_group" bson:"ipc_group"`
	IsSubmittedPCT            bool      `json:"is_submitted_pct" bson:"is_submitted_pct"`
	IDSubmittedPCT            string    `json:"id_submitted_pct" bson:"id_submitted_pct"`
	DateSubmittedPCT          time.Time `json:"date_submitted_pct" bson:"date_submitted_pct"`
	IsSeparated               bool      `json:"is_separated" bson:"is_separated"`
	ParentSeparatedId         string    `json:"parent_separated_id" bson:"parent_separated_id"`
	ParentSeparatedDate       time.Time `json:"parent_separated_date" bson:"parent_separated_date"`
	IsConverted               bool      `json:"is_converted" bson:"is_converted"`
	ParentConvertedId         string    `json:"parent_converted_id" bson:"parent_converted_id"`
	ParentConvertedDate       time.Time `json:"parent_converted_date" bson:"parent_converted_date"`
	IsScientificMission       bool      `json:"is_scientific_mission" bson:"is_scientific_mission"`
	ScientificMissionId       string    `json:"scientific_mission_id" bson:"scientific_mission_id"`
	ScientificMissionMinistry string    `json:"scientific_mission_ministry" bson:"scientific_mission_ministry"`
	ScientificMissionName     string    `json:"scientific_mission_name" bson:"scientific_mission_name"`
	RepresentativeId          string    `json:"representative_id" bson:"representative_id"`
	RepresentativeType        string    `json:"representative_type"`
	OwnerId                   string    `json:"owner_id" bson:"owner_id"`
	OwnerType                 string    `json:"owner_type"`
	AuthorId                  string    `json:"author_id" bson:"author_id"`
	CoAuthor                  []string  `json:"co_author" bson:"co_author"`
	Image                     string    `json:"image" bson:"image"`
	DocumentList              []string  `json:"document_list" bson:"document_list"`
	Status                    string    `json:"status" bson:"status"`
}

type Pagination struct {
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
	Page       int64 `json:"page"` // current page number
	Size       int64 `json:"size"` // items per page
	HasMore    bool  `json:"has_more"`
}

type Document struct {
	DocumentID   string `json:"document_id" bson:"document_id"`
	DocumentName string `json:"document_name" bson:"document_name"`
}

type JSON map[string]interface{}
