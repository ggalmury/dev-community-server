package dto

import (
	"dev_community_server/models"
	"encoding/json"
	"time"
)

type PartyArticleDto struct {
	Id          uint           `json:"id"`
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	TechSkill   []string       `json:"techSkill"`
	Position    map[string]int `json:"position"`
	Process     string         `json:"process"`
	Category    string         `json:"category"`
	Deadline    time.Time      `json:"deadline"`
	StartDate   time.Time      `json:"startDate"`
	Span        string         `json:"span"`
	Location    *string        `json:"location"`
	CreatedAt   time.Time      `json:"createdAt"`
}

func NewPartyArticleDto(entity models.PartyArticleEntity) (*PartyArticleDto, error) {
	var (
		techSkill []string
		position  map[string]int
	)

	tsErr := json.Unmarshal(entity.TechSkill, &techSkill)
	posErr := json.Unmarshal(entity.Position, &position)

	if tsErr != nil || posErr != nil {
		return nil, tsErr
	}

	return &PartyArticleDto{
		Id:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		TechSkill:   techSkill,
		Position:    position,
		Process:     entity.Process,
		Category:    entity.Category,
		Deadline:    entity.Deadline,
		StartDate:   entity.StartDate,
		Span:        entity.Span,
		Location:    entity.Location,
		CreatedAt:   entity.CreatedAt,
	}, nil
}

type PartyArticleCreateDto struct {
	Category    string         `json:"category"`
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	TechSkill   []string       `json:"techSkill"`
	Position    map[string]int `json:"position"`
	Process     string         `json:"process"`
	Location    *string        `json:"location"`
	Deadline    string         `json:"deadline"`
	StartDate   string         `json:"startDate"`
	Span        string         `json:"span"`
}
