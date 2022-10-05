package gateway

import (
	"api_gateway/internal/gateway/routes/auth_service"
	"api_gateway/pkg/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var _ Service = &service{}

type service struct {
	logger logging.Logger
}

type Service interface {
	CreateJWT(r *http.Request, userID string) (tr CreateTokenResponse, err error)
	CopyRequest(r *http.Request, dockerServiceName, servicePort string) (*http.Request, error)
	GetURIWithPort(dockerServiceName, servicePort, requestURI string) string
	DoRequest(r *http.Request) (*ResponseData, error)
}

func NewService(logger logging.Logger) Service {
	return &service{
		logger: logger,
	}
}

// Returns URI with service port e.g. http://host:port/path/to/endpoint
func (s service) GetURIWithPort(dockerServiceName, servicePort, requestURI string) string {
	// renders uri from docker service name, port and requested uri
	url := fmt.Sprintf("http://%s:%s%s", dockerServiceName, servicePort, requestURI)
	//s.logger.Logger.Printf("url: %s", url)
	return url
}

// CopyRequest copies the request so, that it can be redirected from api_gateway to another service
// Applies requested method, copies body and sets "Content-Type" header to "application/json"
func (s service) CopyRequest(r *http.Request, dockerServiceName, servicePort string) (*http.Request, error) {
	url := s.GetURIWithPort(dockerServiceName, servicePort, r.RequestURI)

	// if request method is GET then it doesn't add body
	if r.Method == http.MethodGet {
		request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		request.Header.Add("Content-Type", "application/json")
		return request, nil
	}

	// copies given r request body
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	newBody := io.NopCloser(strings.NewReader(string(bytes)))

	// makes new request with url from docker container
	request, err := http.NewRequestWithContext(r.Context(), r.Method, url, newBody)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", r.Header.Get("Content-Type"))
	return request, nil
}

// DoRequest does a request, handles errors and returns *ResponseData pointer
func (s service) DoRequest(r *http.Request) (data *ResponseData, err error) {
	var client http.Client
	response, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %v", err)
	}
	if response == nil {
		//log.Printf("Response is null")
		return nil, fmt.Errorf("response is null")
	}
	//log.Printf("Status code: %d", response.StatusCode)
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	data = NewResponseData(bytes, response.StatusCode)
	return data, nil
}

// CreateJWT does a request to auth_service using his docker container name
func (s service) CreateJWT(r *http.Request, userID string) (tr CreateTokenResponse, err error) {
	uri := fmt.Sprintf(auth_service.CreateTokenFormatURL, userID)
	url := fmt.Sprintf("http://%s:%s%s", auth_service.DockerName, auth_service.Port, uri)

	// Create a request to auth service to generate JWT
	request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		return tr, err
	}
	var client http.Client

	// Do the response
	response, err := client.Do(request)
	if err != nil {
		return tr, err
	}

	var dto CreateTokenResponse

	// Get JWT from response body
	err = json.NewDecoder(response.Body).Decode(&dto)
	if err != nil {
		return tr, err
	}
	return dto, nil
}
