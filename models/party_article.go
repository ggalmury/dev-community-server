package models

import (
	"gorm.io/gorm"
	"time"
)

type PartyArticleEntity struct {
	gorm.Model
	Poster      UserEntity `gorm:"foreignKey:PosterUuid;references:Uuid"`
	PosterUuid  string     `gorm:"type:varchar(36)"`
	Title       string     `gorm:"type:varchar(40)"`
	Description *string    `gorm:"type:longText"`
	TechSkill   []byte     `gorm:"type:json"`
	Position    []byte     `gorm:"type:json"`
	Process     string     `gorm:"type:varchar(6)"`
	Category    string     `gorm:"type:varchar(4)"`
	Deadline    time.Time
	StartDate   time.Time
	Span        string  `gorm:"type:varchar(6)"`
	Location    *string `gorm:"type:varchar(7)"`
}

func (PartyArticleEntity) TableName() string {
	return "party_article"
}

type PartyCommentEntity struct {
	gorm.Model
	PostId  uint   `json:"post_id"`
	Comment string `json:"comment"`
}

func (PartyCommentEntity) TableName() string {
	return "party_comment"
}
