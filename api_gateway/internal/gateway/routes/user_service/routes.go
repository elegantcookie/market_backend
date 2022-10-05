package user_service

const (
	CreateByPhoneNumberURL = "/api/v1/user_service/phone"
	CreateByVkURL          = "/api/v1/user_service/vk/:vk_token/"
	SendVerificationURL    = "/api/v1/user_service/send_verification/:phone_number/"
	GetUsersURL            = "/api/v1/user_service/get/all"
	GetByIdURL             = "/api/v1/user_service/get/id/:id/"
	GetByPhoneNumberURL    = "/api/v1/user_service/get/num/:phone_number/"
	UpdateByIdURL          = "/api/v1/user_service/upd/"
	DeleteByIdURL          = "/api/v1/user_service/del/id/:id/"

	//DocsURL                = "/api/v1/user_service/docs/"

	Port = "10002"
)
