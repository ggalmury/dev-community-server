package dto

import (
	"dev_community_server/entity"
	"dev_community_server/model"
	"time"
)

type UserDto struct {
	Uuid                 string       `json:"uuid"`
	Email                *string      `json:"email"`
	Nickname             *string      `json:"nickname"`
	ProfileImgUrl        *string      `json:"profileImgUrl"`
	KakaoId              *int         `json:"kakaoId"`
	KakaoEmail           *string      `json:"kakaoEmail"`
	KakaoNickname        *string      `json:"kakaoNickname"`
	KakaoProfileImgUrl   *string      `json:"kakaoProfileImgUrl"`
	KakaoThumbnailImgUrl *string      `json:"kakaoThumbnailImgUrl"`
	CreatedAt            time.Time    `json:"createdAt"`
	Platform             string       `json:"platform"`
	Token                model.Tokens `json:"token"`
}

func UserDtoFromEntity(entity entity.UserEntity, tokens model.Tokens) *UserDto {
	return &UserDto{
		Uuid:                 entity.Uuid,
		Email:                entity.Email,
		Nickname:             entity.Nickname,
		ProfileImgUrl:        entity.ProfileImgUrl,
		KakaoId:              entity.KakaoId,
		KakaoEmail:           entity.KakaoEmail,
		KakaoNickname:        entity.KakaoNickname,
		KakaoProfileImgUrl:   entity.KakaoProfileImgUrl,
		KakaoThumbnailImgUrl: entity.KakaoThumbnailImgUrl,
		CreatedAt:            entity.CreatedAt,
		Platform:             entity.Platform,
		Token:                tokens,
	}
}

type LogoutDto struct {
	Uuid string `json:"uuid" bind:"required"`
}
