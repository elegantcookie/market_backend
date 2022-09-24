package gateway

// CreateUserResponse is user_service's create user response
type CreateUserResponse struct {
	ID string `json:"id"`
}

// CreateTokenResponse is auth_service's create token response
// Includes two json web tokens
type CreateTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ResponseData is a model that stores response body and status code
type ResponseData struct {
	Body       []byte
	StatusCode int
}

// NewResponseData ResponseData constructor
func NewResponseData(body []byte, statusCode int) *ResponseData {
	return &ResponseData{
		Body:       body,
		StatusCode: statusCode,
	}
}
