package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guluzadehh/go_chat/internal/config"
)

func AccessToken(username string, config *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(config.JWT.Access.Expire).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func RefreshToken(username string, config *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(config.JWT.Refresh.Expire).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func Verify(tokenStr string, config *config.Config) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
