package gateway

import (
	"api_gateway/pkg/logging"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var _ Service = &service{}

type service struct {
	logger logging.Logger
}

type Service interface {
	CopyRequest(r *http.Request, dockerServiceName, servicePort string) (*http.Request, error)
	DoRequest(r *http.Request) (*ResponseData, error)
}

func NewService(logger logging.Logger) Service {
	return &service{
		logger: logger,
	}
}

// Returns URI with service port e.g. http://host:port/path/to/endpoint
func (s service) getURIWithPort(dockerServiceName, servicePort, requestURI string) string {
	url := fmt.Sprintf("http://%s:%s%s", dockerServiceName, servicePort, requestURI)
	s.logger.Logger.Printf("url: %s", url)
	return url
}

func (s service) CopyRequest(r *http.Request, dockerServiceName, servicePort string) (*http.Request, error) {
	url := s.getURIWithPort(dockerServiceName, servicePort, r.RequestURI)
	if r.Method == http.MethodGet {
		request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		request.Header.Add("Content-Type", "application/json")
		return request, nil
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	newBody := io.NopCloser(strings.NewReader(string(bytes)))

	request, err := http.NewRequestWithContext(r.Context(), r.Method, url, newBody)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", r.Header.Get("Content-Type"))
	return request, nil
}

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
	log.Printf("Status code: %d", response.StatusCode)
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	data = NewResponseData(bytes, response.StatusCode)
	return data, nil
}
