package user

import (
	"auth_service/internal/config"
	jwt_setup "auth_service/pkg/jwt-setup"
	"auth_service/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type service struct {
	logger  *logging.Logger
	storage Storage
}

func NewService(storage Storage, logger *logging.Logger) Service {
	return &service{
		logger:  logger,
		storage: storage,
	}
}

type Service interface {
	SignIn(ctx context.Context, body io.ReadCloser) (string, error)
}

// SignIn method is called when some JSON request got to /auth/ as login: login, pwd:pwd
// To find out what to do with that we:
// At first check if the given login is in the db if yes, we authorize the user
// If not then create new user
func (s service) SignIn(ctx context.Context, body io.ReadCloser) (token string, err error) {
	var dto UserDTO
	err = json.NewDecoder(body).Decode(&dto)
	if err != nil {
		return "", fmt.Errorf("unable to decode user data due to: %v", err)
	}
	uuid, err := s.storage.FindByUsername(ctx, dto.Username)
	s.logger.Printf("uuid outside: %s", uuid)
	if err != nil {
		return "", err
	}
	// if FindByUsername can't find user then it creates new user
	if uuid == "" {
		s.logger.Printf("NEW USER")
		err := s.storage.Create(ctx, dto)
		if err != nil {
			return "", fmt.Errorf("unable to create new user due to: %v", err)
		}
	} else {
		s.logger.Printf("FINDING USER IN DB BY USERNAME AND PASSWORD")
		uuid, err = s.storage.FindByUsernameAndPassword(ctx, dto)
		if err != nil {
			return "", fmt.Errorf("unable to authorize due to: %v", err)
		}
	}
	token, err = jwt_setup.CreateToken(config.GetConfig(), uuid)
	if err != nil {
		return "", fmt.Errorf("unable to create jwt token due to: %v", err)
	}

	return token, nil
}
