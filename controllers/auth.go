package controllers

import (
	"dev_community_server/dto"
	"dev_community_server/initializers"
	"dev_community_server/models"
	"dev_community_server/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/shyunku-libraries/go-logger"
	"io"
	"net/http"
)

const (
	kakaoAPIURL      = "https://kapi.kakao.com/v2/user/me"
	authHeaderPrefix = "Bearer "
)

func KakaoLogin(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	kakaoAccessToken, err := utils.GetBearerToken(&authHeader)
	if err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusUnauthorized, "Incorrect Kakao auth token")
		return
	}

	// access token -> body (x)
	// access token -> header/authorization: "bearer ..."
	// refresh token -> header/X-Refresh-Token: ""

	// Request Kakao user data using Kakao auth token
	req, err := http.NewRequest("GET", kakaoAPIURL, nil)
	if err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Failed to create Kakao request")
		return
	}

	req.Header.Add("Authorization", authHeaderPrefix+*kakaoAccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Failed to request Kakao account")
		return
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Failed to read Kakao account")
		return
	}

	var (
		user      *models.UserEntity
		userDto   *dto.UserDto
		kakaoResp dto.KakaoResponse
	)

	if err = json.Unmarshal(responseBody, &kakaoResp); err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Failed to unmarshal Kakao user data response")
		return
	}
	properties := kakaoResp.Properties

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Transaction error")
		return
	}

	// Kakao user not exists in database
	if err = tx.Where("kakao_id = ?", kakaoResp.ID).First(&user).Error; err != nil {
		createdUser := models.UserEntity{
			Uuid:                 uuid.New().String(),
			Email:                &kakaoResp.KakaoAccount.Email,
			Password:             nil,
			Nickname:             &properties.Nickname,
			ProfileImgUrl:        &properties.ProfileImage,
			KakaoId:              &kakaoResp.ID,
			KakaoEmail:           &kakaoResp.KakaoAccount.Email,
			KakaoNickname:        &properties.Nickname,
			KakaoProfileImgUrl:   &properties.ProfileImage,
			KakaoThumbnailImgUrl: &properties.ThumbnailImage,
		}

		if err = tx.Create(&createdUser).Error; err != nil {
			log.Error(err)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot save Kakao user account")
			return
		}

		// 한 함수로 합쳐
		tokenDto, accErr, refErr := utils.GenerateTokens(&createdUser)
		if accErr != nil || refErr != nil {
			log.Error(err)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred in generating tokens")
			return
		}

		createdToken := models.TokenEntity{
			Uuid:         createdUser.Uuid,
			RefreshToken: tokenDto.RefreshToken,
		}

		if err = tx.Create(&createdToken).Error; err != nil {
			log.Error(err)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot save user auth token")
			return
		}

		userDto = dto.NewUserDto(createdUser, *tokenDto)

		log.Info("New Kakao account created user logged in:", userDto.Uuid)
	} else { // Kakao user already exists in database
		tokenDto, accErr, refErr := utils.GenerateTokens(user)

		if accErr != nil || refErr != nil {
			log.Error(err)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred in generating tokens")
			return
		}

		createdToken := models.TokenEntity{
			Uuid:         user.Uuid,
			RefreshToken: tokenDto.RefreshToken,
		}

		if err = tx.Create(&createdToken).Error; err != nil {
			log.Error(err)
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot save user auth token")
			return
		}

		userDto = dto.NewUserDto(*user, *tokenDto)

		log.Info("Kakao user logged in:", userDto.Uuid)
	}

	if err = tx.Commit().Error; err != nil {
		log.Error(err)
		tx.Rollback()
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred in transaction")
		return
	}

	c.JSON(201, userDto)
}

func UseAuthRouter(g *gin.Engine) {
	sg := g.Group("/auth")

	sg.POST("/kakao", KakaoLogin)
}
