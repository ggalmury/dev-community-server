package controllers

import (
	"dev_community_server/dto"
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

	print(body.Title)
}

func UsePartyRouter(g *gin.Engine) {
	sg := g.Group("/party")

	sg.GET("/articles", GetPartyArticle)
	sg.POST("/create", CreatePartyArticle)
}
