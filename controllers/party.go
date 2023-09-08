package controllers

import (
	"dev_community_server/dto"
	"dev_community_server/initializers"
	"dev_community_server/models"
	"dev_community_server/utils"
	"github.com/gin-gonic/gin"
	log "github.com/shyunku-libraries/go-logger"
)

func GetPartyArticle(c *gin.Context) {
	var entities []models.PartyArticle

	initializers.DB.Find(&entities)

	result := make([]dto.PartyArticleDto, len(entities))

	for idx, entity := range entities {
		result[idx] = dto.PartyArticleDto{
			Id:          entity.Id,
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

	c.JSON(200, gin.H{
		"result": result,
	})

	log.Info("Article list successfully sent")
}

func CreatePartyArticle(c *gin.Context) {
	var body dto.PartyArticleCreateDto

	if err := c.Bind(&body); err != nil {
		utils.AbortWithStrJson(c, 400, "Cannot bind request body")
		return
	}

	partyArticle := models.PartyArticle{
		Poster:      body.Poster,
		Title:       body.Title,
		Description: body.Description,
		TechSkill:   utils.ErrHandledMarshal(c, body.TechSkill),
		Position:    utils.ErrHandledMarshal(c, body.Position),
		Process:     body.Process,
		Category:    body.Category,
		Deadline:    utils.StringToTime(c, body.Deadline),
		StartDate:   utils.StringToTime(c, body.StartDate),
		Span:        body.Span,
		Location:    body.Location,
	}

	result := initializers.DB.Create(&partyArticle)

	if result.Error != nil {
		utils.AbortWithStrJson(c, 500, "Cannot create article")
		return

	}

	//c.Status(200)
	c.JSON(200, gin.H{})
}

func UsePartyRouter(g *gin.Engine) {
	sg := g.Group("/party")

	sg.GET("/articles", GetPartyArticle)
	sg.POST("/create", CreatePartyArticle)
}
