package entity

import (
	"gorm.io/gorm"
	"time"
)

type PartyEntity struct {
	gorm.Model
	Poster      UserEntity `gorm:"foreignKey:PosterUuid;references:uuid"`
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

func (PartyEntity) TableName() string {
	return "party"
}

type PartyCommentEntity struct {
	gorm.Model
	Poster     UserEntity `gorm:"foreignKey:PosterUuid;references:uuid"`
	PosterUuid string     `gorm:"type:varchar(36)"`
	PostId     uint
	Comment    string `gorm:"type:varchar(255)"`
	Group      *uint
	Depth      int
}

func (PartyCommentEntity) TableName() string {
	return "party_comment"
}
