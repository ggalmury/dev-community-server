package controllers

import (
	"dev_community_server/dto"
	"dev_community_server/initializers"
	"dev_community_server/models"
	"dev_community_server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/shyunku-libraries/go-logger"
	"io"
	"net/http"
)

func KakaoLogin(c *gin.Context) {
	var body dto.KakaoLoginDto

	if bindErr := c.Bind(&body); bindErr != nil {
		utils.AbortWithStrJson(c, 400, "Cannot bind request body")
		return
	}

	kakaoAPIURL := "https://kapi.kakao.com/v2/user/me"

	req, reqErr := http.NewRequest("GET", kakaoAPIURL, nil)

	if reqErr != nil {
		utils.AbortWithStrJson(c, 500, "Failed to create kakao request")
		return
	}

	req.Header.Add("Authorization", "Bearer "+body.AccessToken)

	client := &http.Client{}
	res, resErr := client.Do(req)

	if resErr != nil {
		utils.AbortWithStrJson(c, 500, "Failed to request kakao account info")
		return
	}

	defer res.Body.Close()

	responseBody, responseBodyErr := io.ReadAll(res.Body)

	if responseBodyErr != nil {
		utils.AbortWithStrJson(c, 500, "Failed to read kakao account info")
		return
	}

	var user *models.UserEntity

	kakaoResponse := utils.ErrHandledUnmarshal[dto.KakaoResponse](c, responseBody)
	properties := kakaoResponse.Properties

	findErr := initializers.DB.Where("kakao_id = ?", kakaoResponse.ID).First(&user).Error

	log.Fatal(user)

	if findErr != nil {
		createdUser := models.UserEntity{
			Uuid:                 uuid.New().String(),
			Email:                nil,
			Password:             nil,
			Nickname:             nil,
			ProfileImgUrl:        nil,
			KakaoId:              &kakaoResponse.ID,
			KakaoEmail:           &kakaoResponse.KakaoAccount.Email,
			KakaoNickname:        &properties.Nickname,
			KakaoProfileImgUrl:   &properties.ProfileImage,
			KakaoThumbnailImgUrl: &properties.ThumbnailImage,
		}

		create := initializers.DB.Create(&createdUser)

		if create.Error != nil {
			utils.AbortWithStrJson(c, 500, "Cannot create article")
			return
		}

		userDto := dto.UserDtoFromEntity(createdUser)

		c.JSON(201, gin.H{
			"result": userDto,
		})

		log.Info("New Kakao account created")
	} else {
		c.JSON(201, gin.H{
			"result": user,
		})

		log.Info("Kakao account has sent to the client")
	}

}

func UseAuthRouter(g *gin.Engine) {
	sg := g.Group("/auth")

	sg.POST("/kakao", KakaoLogin)
}
