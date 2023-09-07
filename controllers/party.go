package controllers

import (
	"dev_community_server/dto"
	"dev_community_server/initializers"
	"dev_community_server/models"
	"dev_community_server/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPartyArticle(c *gin.Context) {

}

func CreatePartyArticle(c *gin.Context) {
	var body dto.PartyArticleCreateDto

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deadline, deadlineErr := utils.StringToTime(body.Deadline)
	startDate, startDateErr := utils.StringToTime(body.StartDate)
	position, positionErr := json.Marshal(body.Position)
	techSkill, techSkillErr := json.Marshal(body.TechSkill)

	if deadlineErr != nil || startDateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"datetime error": deadlineErr.Error()})
		return
	}

	if positionErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"position error": positionErr.Error()})
		return
	}

	if techSkillErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"techSkill error": positionErr.Error()})
		return
	}

	partyArticle := models.PartyArticle{
		Poster:      body.Poster,
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
		c.JSON(http.StatusInternalServerError, gin.H{"Cannot create article": result.Error})
		return
	}

	//c.Status(200)
	c.JSON(200, gin.H{
		"success": true,
	})
}

func UsePartyRouter(g *gin.Engine) {
	sg := g.Group("/party")

	sg.GET("/articles", GetPartyArticle)
	sg.POST("/create", CreatePartyArticle)
}
