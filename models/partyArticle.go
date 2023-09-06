package models

import "gorm.io/gorm"

type PartyArticle struct {
	gorm.Model
	Id          int            `gorm:"primaryKey" json:"id"`
	CreatedDt   string         `json:"createdDt"`
	Poster      string         `json:"poster"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	TechSkill   []string       `json:"techSkill"`
	Position    map[string]int `json:"position"`
	Process     string         `json:"process"`
	Category    string         `json:"category"`
	Deadline    string         `json:"deadline"`
	StartDate   string         `json:"startDate"`
	Span        string         `json:"span"`
	Location    *string        `json:"location"`
}
