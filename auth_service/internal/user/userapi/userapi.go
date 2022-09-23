package userapi

import (
	"auth_service/internal/user"
	"auth_service/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	usersURL     = "http://localhost:10002/api/users"
	usersAuthURL = "http://localhost:10002/api/users/auth"
	usernameURL  = "http://localhost:10002/api/users/username"
)

type UserAPI struct {
	logger *logging.Logger
	client http.Client
}

func NewStorage(logger *logging.Logger) *UserAPI {
	return &UserAPI{logger: logger}
}

func (ua *UserAPI) Create(ctx context.Context, dto user.UserDTO) error {
	ua.logger.Println(dto)
	userbytes, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal dto due to: %v", err)
	}
	reader := strings.NewReader(string(userbytes))
	request, err := http.NewRequest(http.MethodPost, usersURL, reader)
	if err != nil {
		return fmt.Errorf("failed to build a request due to: %v", err)
	}
	_, err = ua.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to do a request due to: %v", err)
	}
	return nil
}

func (ua *UserAPI) FindByUsername(ctx context.Context, username string) (uuid string, err error) {
	url := fmt.Sprintf("%s/%s", usernameURL, username)
	ua.logger.Printf("URL: %s", url)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}
	response, err := ua.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to do a request due to: %v", err)
	}
	var dto user.ResponseUserDTO
	err = json.NewDecoder(response.Body).Decode(&dto)
	if err != nil {
		return "", fmt.Errorf("failed to decode the response data due to: %v", err)
	}
	ua.logger.Printf("found user with id: %s", dto.ID)
	return dto.ID, nil
}

func (ua *UserAPI) FindByUsernameAndPassword(ctx context.Context, dto user.UserDTO) (uuid string, err error) {
	userbytes, err := json.Marshal(dto)
	if err != nil {
		return "", fmt.Errorf("failed to marshal dto due to: %v", err)
	}
	reader := strings.NewReader(string(userbytes))
	request, err := http.NewRequest(http.MethodPost, usersAuthURL, reader)
	ua.logger.Printf("REQUEST PAYLOAD: %v", reader)
	if err != nil {
		return "", fmt.Errorf("failed to build a request due to: %v", err)
	}

	response, err := ua.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to do a request due to: %v", err)
	}
	var rdto user.ResponseUserDTO
	err = json.NewDecoder(response.Body).Decode(&rdto)
	if err != nil {
		return "", fmt.Errorf("failed to decode the response data due to: %v", err)
	}

	if rdto.Username == "" && rdto.ID == "" {
		return "", fmt.Errorf("password doesn't match")
	}

	return rdto.ID, nil
}
