package dto

import (
	"dev_community_server/models"
	"dev_community_server/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type PartyArticleDto struct {
	Poster      string         `json:"poster"`
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

func PartyArticleDtoFromEntity(c *gin.Context, entity models.PartyArticleEntity) *PartyArticleDto {
	return &PartyArticleDto{
		Poster:      entity.Poster,
		Title:       entity.Title,
		Description: entity.Description,
		TechSkill:   utils.ErrHandledUnmarshal[[]string](c, entity.TechSkill),
		Position:    utils.ErrHandledUnmarshal[map[string]int](c, entity.Position),
		Process:     entity.Process,
		Category:    entity.Category,
		Deadline:    entity.Deadline,
		StartDate:   entity.StartDate,
		Span:        entity.Span,
		Location:    entity.Location,
		CreatedAt:   entity.CreatedAt,
	}
}

type PartyArticleCreateDto struct {
	Poster      string         `json:"poster"`
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
