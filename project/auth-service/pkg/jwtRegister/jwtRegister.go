package jwtRegister

import (
	"auth-service/internal/model"
	"auth-service/pkg/custom_errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GenerateToken(claims model.JWTCustomClaims) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(viper.GetString("jwt.key")))
	if err != nil {
		return "", fmt.Errorf(fmt.Errorf("sign string: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return token, nil
}
