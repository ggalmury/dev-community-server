package model

import (
	"dev_community_server/entity"
	"time"
)

type Poster struct {
	Uuid          string `json:"uuid"`
	Nickname      string `json:"nickname"`
	ProfileImgUrl string `json:"profileImgUrl"`
}

func NewPoster(entity entity.UserEntity) *Poster {
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
	PostId    uint      `json:"postId"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	Group     *uint     `json:"group"`
	Depth     int       `json:"depth"`
	Poster    Poster    `json:"poster"`
}

func NewComment(entity entity.PartyCommentEntity) *Comment {
	poster := NewPoster(entity.Poster)

	return &Comment{
		PostId:    entity.PostId,
		Comment:   entity.Comment,
		CreatedAt: entity.CreatedAt,
		Group:     entity.Group,
		Depth:     entity.Depth,
		Poster:    *poster,
	}
}

type CommentGroup struct {
	Comment    Comment `json:"comment"`
	SubComment Comment `json:"subComment"`
}
