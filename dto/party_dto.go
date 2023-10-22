package dto

import (
	"dev_community_server/entity"
	"dev_community_server/model"
)

type PartyListDto struct {
	Parties []model.Party `json:"parties"`
}

func PartyListDtoFromEntity(entity []entity.PartyEntity) (*PartyListDto, error) {
	parties := make([]model.Party, len(entity))

	for idx, e := range entity {
		p, err := model.PartyFromEntity(e)
		if err != nil {
			return nil, err
		}

		parties[idx] = *p
	}

	return &PartyListDto{Parties: parties}, nil
}

type PartyCommentDto struct {
	Comment model.Comment `json:"comment"`
}

func PartyCommentDtoFromEntity(entity entity.PartyCommentEntity) *PartyCommentDto {
	comment := model.CommentFromEntity(entity)

	return &PartyCommentDto{Comment: *comment}
}

type PartyCommentListDto struct {
	Comments []model.CommentGroup `json:"comments"`
}

func PartyCommentListDtoFromEntity(entity []entity.PartyCommentEntity) *PartyCommentListDto {
	if len(entity) == 0 {
		return &PartyCommentListDto{Comments: []model.CommentGroup{}}
	}

	comments := make([]model.Comment, len(entity))

	for idx, e := range entity {
		c := model.CommentFromEntity(e)
		comments[idx] = *c
	}

	var commentGroup []model.CommentGroup

	for _, comment := range comments {
		if comment.Depth == 0 {
			cg := model.CommentGroup{Comment: comment}
			cg.Reply = make([]model.Comment, 0)

			for _, subComment := range comments {
				if subComment.Depth == 1 && *subComment.Group == *comment.Group {
					cg.Reply = append(cg.Reply, subComment)
				}
			}
			commentGroup = append(commentGroup, cg)
		}
	}

	return &PartyCommentListDto{Comments: commentGroup}
}

type PartyCreateDto struct {
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

type PartyCommentCreateDto struct {
	PostId  uint   `json:"postId"`
	Comment string `json:"comment"`
	Depth   int    `json:"depth"`
	Group   *uint  `json:"group"`
}
