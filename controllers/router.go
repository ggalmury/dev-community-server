package controllers

import "github.com/gin-gonic/gin"

func UseRouter(r *gin.Engine) {
	UsePartyRouter(r)
}
