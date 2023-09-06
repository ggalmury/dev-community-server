package controllers

import "github.com/gin-gonic/gin"

func GetPartyArticle(c *gin.Context) {

}

func CreatePartyArticle(c *gin.Context) {

}

func UsePartyRouter(g *gin.Engine) {
	sg := g.Group("/party")

	sg.GET("/articles", GetPartyArticle)
	sg.POST("/create", CreatePartyArticle)
}
