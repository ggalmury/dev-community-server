package controllers

import (
	"dev_community_server/dto"
	"dev_community_server/initializers"
	"dev_community_server/models"
	"dev_community_server/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/shyunku-libraries/go-logger"
	"net/http"
)

func GetPartyArticle(c *gin.Context) {
	var entities []models.PartyArticleEntity

	if err := initializers.DB.Find(&entities).Error; err != nil {
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot found party articles")
		return
	}

	result := make([]dto.PartyArticleDto, len(entities))

	for idx, entity := range entities {
		pad, err := dto.NewPartyArticleDto(entity)
		if err != nil {
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred while binding party article entity to dto")
			return
		}

		result[idx] = *pad
	}

	c.JSON(http.StatusOK, result)

	log.Info("Article list successfully sent")
}

func CreatePartyArticle(c *gin.Context) {
	var body dto.PartyArticleCreateDto

	if err := c.Bind(&body); err != nil {
		utils.AbortWithStrJson(c, http.StatusBadRequest, "Cannot bind request body")
		return
	}

	techSkill, tsErr := json.Marshal(body.TechSkill)
	position, posErr := json.Marshal(body.Position)
	if tsErr != nil || posErr != nil {
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Failed to marshal article properties")
		return
	}

	partyArticle := models.PartyArticleEntity{
		Poster:      body.Poster,
		Title:       body.Title,
		Description: body.Description,
		TechSkill:   techSkill,
		Position:    position,
		Process:     body.Process,
		Category:    body.Category,
		Deadline:    utils.StringToTime(c, body.Deadline),
		StartDate:   utils.StringToTime(c, body.StartDate),
		Span:        body.Span,
		Location:    body.Location,
	}

	result := initializers.DB.Create(&partyArticle)

	if result.Error != nil {
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot create article")
		return
	}

	c.JSON(http.StatusCreated, gin.H{})

	log.Info("Article successfully created")
}

func UsePartyRouter(g *gin.Engine) {
	sg := g.Group("/party")

	sg.GET("/articles", GetPartyArticle)
	sg.POST("/create", CreatePartyArticle)
}
