package gateway

type ResponseData struct {
	Body       []byte
	StatusCode int
}

func NewResponseData(body []byte, statusCode int) *ResponseData {
	return &ResponseData{
		Body:       body,
		StatusCode: statusCode,
	}
}
