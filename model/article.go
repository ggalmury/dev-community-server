package model

import (
	"dev_community_server/entity"
	"encoding/json"
	"time"
)

type Party struct {
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
	Poster      Poster         `json:"poster"`
}

func PartyFromEntity(entity entity.PartyEntity) (*Party, error) {
	var (
		techSkill []string
		position  map[string]int
	)

	tsErr := json.Unmarshal(entity.TechSkill, &techSkill)
	posErr := json.Unmarshal(entity.Position, &position)

	if tsErr != nil || posErr != nil {
		return nil, tsErr
	}

	poster := PosterFromEntity(entity.Poster)

	return &Party{
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
		Poster:      *poster,
	}, nil
}

type Poster struct {
	Uuid          string `json:"uuid"`
	Nickname      string `json:"nickname"`
	ProfileImgUrl string `json:"profileImgUrl"`
}

func PosterFromEntity(entity entity.UserEntity) *Poster {
	var profileImgUrl string

	switch entity.Platform {
	case "email":
		profileImgUrl = *entity.ProfileImgUrl
	case "kakao":
		profileImgUrl = *entity.KakaoProfileImgUrl
	}

	return &Poster{
		Uuid:          entity.Uuid,
		Nickname:      *entity.Nickname,
		ProfileImgUrl: profileImgUrl,
	}
}

type Comment struct {
	Id        uint      `json:"id"`
	PostId    uint      `json:"postId"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	Group     *uint     `json:"group"`
	Depth     int       `json:"depth"`
	Poster    Poster    `json:"poster"`
}

func CommentFromEntity(entity entity.PartyCommentEntity) *Comment {
	poster := PosterFromEntity(entity.Poster)

	return &Comment{
		Id:        entity.ID,
		PostId:    entity.PostId,
		Comment:   entity.Comment,
		CreatedAt: entity.CreatedAt,
		Group:     entity.Group,
		Depth:     entity.Depth,
		Poster:    *poster,
	}
}

type CommentGroup struct {
	Comment    Comment   `json:"comment"`
	SubComment []Comment `json:"subComment"`
}
