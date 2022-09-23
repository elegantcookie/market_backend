package jwt_setup

import (
	"auth_service/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateToken(cfg *config.Config, userId string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        userId,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(8760 * time.Hour)},
	})
	return token.SignedString([]byte(cfg.Keys.JWTSignKey))
}
