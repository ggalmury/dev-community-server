package entity

import "gorm.io/gorm"

type UserEntity struct {
	gorm.Model
	Uuid                 string  `gorm:"type:varchar(36);unique"`
	Email                *string `gorm:"type:varchar(50)"`
	Password             *string `gorm:"type:varchar(60)"`
	Nickname             *string `gorm:"type:varchar(12)"`
	ProfileImgUrl        *string
	KakaoId              *int
	KakaoEmail           *string
	KakaoNickname        *string
	KakaoProfileImgUrl   *string
	KakaoThumbnailImgUrl *string
	Platform             string
}

func (UserEntity) TableName() string {
	return "user"
}
