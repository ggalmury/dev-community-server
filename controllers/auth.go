package controllers

import (
	"dev_community_server/crypto"
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
	kakaoAPIURL = "https://kapi.kakao.com/v2/user/me"
)

func KakaoLogin(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	kakaoAccessToken, err := utils.GetBearerToken(&authHeader)
	if err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusUnauthorized, "Incorrect Kakao auth token")
		return
	}

	// access token -> header/authorization: "bearer ..."
	// refresh token -> header/X-Refresh-Token: ""

	// Request Kakao user data using Kakao auth token
	req, err := http.NewRequest("GET", kakaoAPIURL, nil)
	if err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Failed to create Kakao request")
		return
	}

	req.Header.Add("Authorization", "Bearer "+*kakaoAccessToken)

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
		log.Error(tx.Error)
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
			Platform:             utils.Kakao,
		}

		if err = tx.Create(&createdUser).Error; err != nil {
			log.Error(err)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot save Kakao user account")
			return
		}

		tokenDto, tokenDtoErr := crypto.GenerateTokens(&createdUser)
		if tokenDtoErr != nil {
			log.Error(tokenDtoErr)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred in generating tokens")
			return
		}

		if saveTokenErr := crypto.SaveTokens(createdUser.Uuid, tokenDto.RefreshToken); saveTokenErr != nil {
			log.Error(saveTokenErr)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot save user auth token")
			return
		}

		userDto = dto.NewUserDto(createdUser, *tokenDto)

		log.Info("New Kakao user created:", userDto.Uuid)
	} else { // Kakao user already exists in database
		tokenDto, tokenDtoErr := crypto.GenerateTokens(user)
		if tokenDtoErr != nil {
			log.Error(tokenDtoErr)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred in generating tokens")
			return
		}

		if saveTokenErr := crypto.SaveTokens(user.Uuid, tokenDto.RefreshToken); saveTokenErr != nil {
			log.Error(saveTokenErr)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot save user auth token")
			return
		}

		userDto = dto.NewUserDto(*user, *tokenDto)
	}

	if err = tx.Commit().Error; err != nil {
		log.Error(err)
		tx.Rollback()
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred in transaction")
		return
	}

	log.Info("Kakao user logged in:", userDto.Uuid)
	c.JSON(201, userDto)
}

func AutoLogin(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	refreshToken := c.Request.Header.Get("X-Refresh-Token")
	accessToken, err := utils.GetBearerToken(&authHeader)
	if err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusUnauthorized, "Incorrect access token format")
		return
	}

	var user *models.UserEntity

	atClaims, atErr := crypto.ValidateAccessToken(*accessToken)
	rtClaims, rtErr := crypto.ValidateRefreshToken(refreshToken)
	if atErr == nil {
		initializers.DB.Where("uuid = ?", atClaims.Uuid).First(&user)
	} else if atErr != nil && rtErr == nil {
		initializers.DB.Where("uuid = ?", rtClaims.Uuid).First(&user)
	} else {
		log.Error(atErr, rtErr)
		utils.AbortWithStrJson(c, http.StatusUnauthorized, "Invalid tokens")
		return
	}

	tokenDto, tokenDtoErr := crypto.GenerateTokens(user)
	if tokenDtoErr != nil {
		log.Error(tokenDtoErr)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred in generating tokens")
		return
	}

	if saveTokenErr := crypto.SaveTokens(user.Uuid, tokenDto.RefreshToken); saveTokenErr != nil {
		log.Error(saveTokenErr)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot save user auth token")
		return
	}

	userDto := dto.NewUserDto(*user, *tokenDto)
	c.JSON(201, userDto)

}

func UseAuthRouter(g *gin.Engine) {
	sg := g.Group("/auth")

	sg.POST("/kakao", KakaoLogin)
	sg.POST("/auto-login", AutoLogin)
}
