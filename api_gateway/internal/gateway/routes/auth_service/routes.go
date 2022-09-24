package auth_service

const (
	CheckTokenURL        = "/api/v1/auth_service/token/check"
	CreateTokenURL       = "/api/v1/auth_service/token/create/:id/"
	CreateTokenFormatURL = "/api/v1/auth_service/token/create/%s/"
	DockerName           = "auth_service"
	Port                 = "10001"
)
