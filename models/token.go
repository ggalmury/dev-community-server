package models

import "gorm.io/gorm"

type TokenEntity struct {
	gorm.Model
	Uuid         string `json:"uuid"`
	RefreshToken string `json:"refreshToken"`
}
