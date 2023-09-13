package dto

import (
	"dev_community_server/models"
	"time"
)

type KakaoLoginDto struct {
	AccessToken string `json:"accessToken"`
}

type UserDto struct {
	Uuid                 string    `json:"uuid"`
	Email                *string   `json:"email"`
	Nickname             *string   `json:"nickname"`
	ProfileImgUrl        *string   `json:"profileImgUrl"`
	KakaoId              *int      `json:"kakaoId"`
	KakaoEmail           *string   `json:"kakaoEmail"`
	KakaoNickname        *string   `json:"kakaoNickname"`
	KakaoProfileImgUrl   *string   `json:"kakaoProfileImgUrl"`
	KakaoThumbnailImgUrl *string   `json:"kakaoThumbnailImgUrl"`
	CreatedAt            time.Time `json:"createdAt"`
}

func UserDtoFromEntity(entity models.UserEntity) *UserDto {
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
	}
}

type KakaoResponse struct {
	ID           int             `json:"id"`
	ConnectedAt  string          `json:"connected_at"`
	Properties   KakaoProperties `json:"properties"`
	KakaoAccount KakaoAccount    `json:"kakao_account"`
}

type KakaoProperties struct {
	Nickname       string `json:"nickname"`
	ProfileImage   string `json:"profile_image"`
	ThumbnailImage string `json:"thumbnail_image"`
}

type KakaoAccount struct {
	Email string `json:"email"`
}
