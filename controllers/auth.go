package controllers

import (
	"dev_community_server/crypto"
	"dev_community_server/dto"
	"dev_community_server/initializers"
	"dev_community_server/models"
	"dev_community_server/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/shyunku-libraries/go-logger"
	"gorm.io/gorm"
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
	props := kakaoResp.Properties

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		log.Error(tx.Error)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Transaction error")
		return
	}

	// Kakao user not exists in database
	if err = tx.Where("kakao_id = ?", kakaoResp.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			*user = models.UserEntity{
				Uuid:                 uuid.New().String(),
				Email:                &kakaoResp.KakaoAccount.Email,
				Password:             nil,
				Nickname:             &props.Nickname,
				ProfileImgUrl:        &props.ProfileImage,
				KakaoId:              &kakaoResp.ID,
				KakaoEmail:           &kakaoResp.KakaoAccount.Email,
				KakaoNickname:        &props.Nickname,
				KakaoProfileImgUrl:   &props.ProfileImage,
				KakaoThumbnailImgUrl: &props.ThumbnailImage,
				Platform:             utils.Kakao,
			}

			if err = tx.Create(&user).Error; err != nil {
				log.Error(err)
				tx.Rollback()
				utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot save Kakao user account")
				return
			}

			log.Info("New Kakao user created:", user.Uuid)
		} else {
			log.Error(err)
			tx.Rollback()
			utils.AbortWithStrJson(c, http.StatusInternalServerError, "Error occurred in find user")
			return
		}
	}

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

func Logout(c *gin.Context) {
	var body dto.LogoutDto
	if bindErr := c.Bind(&body); bindErr != nil {
		utils.AbortWithStrJson(c, http.StatusBadRequest, "Cannot bind request body")
		return
	}

	if err := crypto.DeleteTokens(body.Uuid); err != nil {
		log.Error(err)
		utils.AbortWithStrJson(c, http.StatusInternalServerError, "Cannot delete user auth token")
		return
	}

	log.Info("Kakao user logged out:", body.Uuid)
	c.JSON(201, gin.H{})
}

func UseAuthRouter(g *gin.Engine) {
	sg := g.Group("/auth")

	sg.POST("/kakao", KakaoLogin)
	sg.POST("/auto-login", AutoLogin)
	sg.POST("/logout", Logout)
}
