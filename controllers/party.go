package controllers

import (
	"dev_community_server/dto"
	"dev_community_server/initializers"
	"dev_community_server/middlewares"
	"dev_community_server/models"
	"dev_community_server/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/shyunku-libraries/go-logger"
	"net/http"
)

func GetPartyArticle(c *gin.Context) {
	var entities []models.PartyArticleEntity

	if err := initializers.DB.Preload("Poster").Find(&entities).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Info(entities)

	result := make([]dto.PartyArticleDto, len(entities))

	for idx, entity := range entities {
		pad, err := dto.NewPartyArticleDto(entity)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		result[idx] = *pad
	}

	c.JSON(http.StatusOK, result)
	log.Info("Article list successfully sent")
}

func CreatePartyArticle(c *gin.Context) {
	var body dto.PartyArticleCreateDto
	uuid, ok := c.Get("uuid")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := c.Bind(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	techSkill, tsErr := json.Marshal(body.TechSkill)
	position, posErr := json.Marshal(body.Position)
	if tsErr != nil || posErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	deadline, dlErr := utils.StringToTime(body.Deadline)
	startDate, sdErr := utils.StringToTime(body.StartDate)
	if dlErr != nil || sdErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	partyArticle := models.PartyArticleEntity{
		PosterUuid:  uuid.(string),
		Title:       body.Title,
		Description: body.Description,
		TechSkill:   techSkill,
		Position:    position,
		Process:     body.Process,
		Category:    body.Category,
		Deadline:    *deadline,
		StartDate:   *startDate,
		Span:        body.Span,
		Location:    body.Location,
	}

	result := initializers.DB.Create(&partyArticle)

	if result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
	log.Info("Article successfully created")
}

func UsePartyRouter(g *gin.Engine) {
	sg := g.Group("/party")
	sg.Use(middlewares.TokenMiddleWare())

	sg.GET("/articles", GetPartyArticle)
	sg.POST("/create", CreatePartyArticle)
}
