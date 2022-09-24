package auth

import (
	"auth_service/internal/apperror"
	"auth_service/internal/config"
	"auth_service/pkg/logging"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"time"
)

type service struct {
	logger  *logging.Logger
	storage Storage
}

func NewService(logger *logging.Logger) Service {
	return &service{
		logger: logger,
	}
}

type Service interface {
	CheckToken(ctx context.Context, token string) error
	CreateToken(ctx context.Context, userID string) (JWT, error)
	SignIn(ctx context.Context, body io.ReadCloser) (string, error)
}

var signKey *rsa.PrivateKey

type RegisteredClaims struct {
	jwt.RegisteredClaims
}

func (s service) CheckToken(ctx context.Context, tokenString string) error {
	s.logger.Printf("CHECK TOKEN")
	if tokenString == "" {
		return apperror.ErrEmptyTokenString
	}

	token, err := jwt.ParseWithClaims(tokenString, &RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Keys.JWTSignKey), nil
	})
	if err != nil {
		s.logger.Printf("TOKEN PARSED WITH ERRORS: %v", err)
		return apperror.ErrInvalidTokenSignature
	}

	s.logger.Printf("TOKEN PARSED WITHOUT ERRORS")

	if claims, ok := token.Claims.(*RegisteredClaims); ok && token.Valid {
		return nil
	} else {
		s.logger.Printf("wrong token: %v", claims.RegisteredClaims.Issuer)
		return apperror.ErrInvalidTokenSignature

	}
}

func (s service) CreateToken(ctx context.Context, userID string) (tokens JWT, err error) {
	if userID == "" {
		return tokens, fmt.Errorf("userID can't be empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        userID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(48 * time.Hour)},
	})
	accessToken, err := token.SignedString([]byte(config.GetConfig().Keys.JWTSignKey))
	if err != nil {
		return tokens, err
	}
	tokens = JWT{
		AccessToken:  accessToken,
		RefreshToken: "",
	}
	return tokens, nil
}

// SignIn method is called when some JSON request got to /auth/ as login: login, pwd:pwd
// To find out what to do with that we:
// At first check if the given login is in the db if yes, we authorize the auth
// If not then create new auth
func (s service) SignIn(ctx context.Context, body io.ReadCloser) (token string, err error) {
	var dto UserDTO
	err = json.NewDecoder(body).Decode(&dto)
	if err != nil {
		return "", fmt.Errorf("unable to decode auth data due to: %v", err)
	}
	uuid, err := s.storage.FindByUsername(ctx, dto.Username)
	s.logger.Printf("uuid outside: %s", uuid)
	if err != nil {
		return "", err
	}
	// if FindByUsername can't find auth then it creates new auth
	if uuid == "" {
		s.logger.Printf("NEW USER")
		err := s.storage.Create(ctx, dto)
		if err != nil {
			return "", fmt.Errorf("unable to create new auth due to: %v", err)
		}
	} else {
		s.logger.Printf("FINDING USER IN DB BY USERNAME AND PASSWORD")
		uuid, err = s.storage.FindByUsernameAndPassword(ctx, dto)
		if err != nil {
			return "", fmt.Errorf("unable to authorize due to: %v", err)
		}
	}
	//token, err = jwt_setup.CreateToken(config.GetConfig(), uuid)
	if err != nil {
		return "", fmt.Errorf("unable to create jwt token due to: %v", err)
	}

	return token, nil
}
