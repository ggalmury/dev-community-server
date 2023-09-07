package models

import (
	"gorm.io/gorm"
	"time"
)

type PartyArticle struct {
	gorm.Model
	Id          int            `json:"id" gorm:"primaryKey" `
	Poster      string         `json:"poster" gorm:"type:varchar(10)"`
	Title       string         `json:"title" gorm:"type:varchar(40)"`
	Description string         `json:"description" gorm:"type:longText"`
	TechSkill   []string       `json:"techSkill" gorm:"type:json"`
	Position    map[string]int `json:"position" gorm:"type:json"`
	Process     string         `json:"process" gorm:"type:varchar(6)"`
	Category    string         `json:"category" gorm:"type:varchar(4)"`
	Deadline    time.Time      `json:"deadline"`
	StartDate   time.Time      `json:"startDate"`
	Span        string         `json:"span" gorm:"type:varchar(6)"`
	Location    *string        `json:"location" gorm:"type:varchar(7)"`
}
