package models

import (
	"gorm.io/gorm"
	"time"
)

type PartyArticleEntity struct {
	gorm.Model
	Poster      string  `gorm:"type:varchar(10)"`
	Title       string  `gorm:"type:varchar(40)"`
	Description *string `gorm:"type:longText"`
	TechSkill   []byte  `gorm:"type:json"`
	Position    []byte  `gorm:"type:json"`
	Process     string  `gorm:"type:varchar(6)"`
	Category    string  `gorm:"type:varchar(4)"`
	Deadline    time.Time
	StartDate   time.Time
	Span        string  `gorm:"type:varchar(6)"`
	Location    *string `gorm:"type:varchar(7)"`
}

func (PartyArticleEntity) TableName() string {
	return "party_article"
}
