package crypto

import (
	"dev_community_server/dto"
	"dev_community_server/entity"
	"dev_community_server/initializers"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/shyunku-libraries/go-logger"
	"os"
	"strconv"
	"time"
)

func SaveTokens(uuid string, token string) error {
	if err := initializers.InMemoryDB.Set(uuid, token); err != nil {
		return err
	}

	log.Info("Token saved in in-memory DB / [uuid]:", uuid)
	return nil
}

func DeleteTokens(uuid string) error {
	if err := initializers.InMemoryDB.Del(uuid); err != nil {
		return err
	}

	log.Info("Token deleted in in-memory DB / [uuid]:", uuid)
	return nil
}

func GenerateTokens(e *entity.UserEntity) (*dto.TokenDto, error) {
	atKey := os.Getenv("JWT_ACCESS_SECRET")
	rtKey := os.Getenv("JWT_REFRESH_SECRET")
	atExp, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXP_DATE"))
	rtExp, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXP_DATE"))

	atClaims := dto.AccessTokenClaims{
		Uuid:      e.Uuid,
		CreatedAt: e.CreatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(atExp))),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atString, err := at.SignedString([]byte(atKey))
	if err != nil {
		return nil, err
	}

	rtClaims := dto.RefreshTokenClaims{
		Id:        e.ID,
		Uuid:      e.Uuid,
		CreatedAt: e.CreatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(rtExp))),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtString, err := rt.SignedString([]byte(rtKey))
	if err != nil {
		return nil, err
	}

	token := dto.TokenDto{
		AccessToken:  atString,
		RefreshToken: rtString,
	}

	log.Info("Auth token generated / [uuid]:", e.Uuid)
	return &token, nil
}

func ValidateAccessToken(at string) (*dto.AccessTokenClaims, error) {
	claim := &dto.AccessTokenClaims{}

	_, err := jwt.ParseWithClaims(at, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err != nil {
		return claim, err
	}
	return claim, nil
}

func ValidateRefreshToken(rt string) (*dto.RefreshTokenClaims, error) {
	claim := &dto.RefreshTokenClaims{}

	_, err := jwt.ParseWithClaims(rt, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil {
		return claim, err
	}
	return claim, nil
}
