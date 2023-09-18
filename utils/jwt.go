package utils

import (
	"dev_community_server/dto"
	"dev_community_server/models"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

func GenerateAccessToken(e *models.UserEntity) (*string, error) {
	secretKey := os.Getenv("JWT_ACCESS_SECRET")
	expDate, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXP_DATE"))

	claims := jwt.MapClaims{
		"uuid":      e.Uuid,
		"createdAt": e.CreatedAt,
		"exp":       time.Now().Add(time.Hour * 24 * time.Duration(expDate)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func GenerateRefreshToken(e *models.UserEntity) (*string, error) {
	secretKey := os.Getenv("JWT_REFRESH_SECRET")
	expDate, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXP_DATE"))

	claims := jwt.MapClaims{
		"id":        e.ID,
		"uuid":      e.Uuid,
		"createdAt": e.CreatedAt,
		"exp":       time.Now().Add(time.Hour * 24 * time.Duration(expDate)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func GenerateTokens(e *models.UserEntity) (*dto.TokenDto, error, error) {
	accessToken, accErr := GenerateAccessToken(e)
	refreshToken, refErr := GenerateRefreshToken(e)

	if accErr != nil || refErr != nil {
		return nil, accErr, refErr
	}

	return &dto.TokenDto{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil, nil
}
