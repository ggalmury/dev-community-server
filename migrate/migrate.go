package main

import (
	"dev_community_server/initializers"
	"dev_community_server/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DBConnection()
}

func main() {
	initializers.DB.AutoMigrate(&models.PartyArticleEntity{}, &models.UserEntity{})
}
