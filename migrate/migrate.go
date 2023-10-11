package main

import (
	"dev_community_server/entity"
	"dev_community_server/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DBConnection()
}

func main() {
	initializers.DB.AutoMigrate(&entity.PartyEntity{}, &entity.UserEntity{}, &entity.PartyCommentEntity{})
}
