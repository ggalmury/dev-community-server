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
	kakaoAccessToken, err := utils.GetBearerToken(authHeader)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Request Kakao user data using Kakao auth token
	req, err := http.NewRequest("GET", kakaoAPIURL, nil)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", "Bearer "+*kakaoAccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var (
		user      *models.UserEntity
		userDto   *dto.UserDto
		kakaoResp dto.KakaoResponse
	)

	if err = json.Unmarshal(respBody, &kakaoResp); err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	props := kakaoResp.Properties

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		log.Error(tx.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
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
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			log.Info("New Kakao user created:", user.Uuid)
		} else {
			log.Error(err)
			tx.Rollback()
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	tokenDto, err := crypto.GenerateTokens(user)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = crypto.SaveTokens(user.Uuid, tokenDto.RefreshToken); err != nil {
		log.Error(err)
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userDto = dto.NewUserDto(*user, *tokenDto)

	if err = tx.Commit().Error; err != nil {
		log.Error(err)
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(201, userDto)
	log.Info("Kakao user logged in / [uuid]:", userDto.Uuid)
}

func AutoLogin(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	refreshToken := c.Request.Header.Get("X-Refresh-Token")
	accessToken, err := utils.GetBearerToken(authHeader)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user *models.UserEntity

	atClaims, atErr := crypto.ValidateAccessToken(*accessToken)
	rtClaims, rtErr := crypto.ValidateRefreshToken(refreshToken)
	if atClaims.Uuid != rtClaims.Uuid {
		log.Errorf("Payload mismatch / [access payload]: $s [refresh payload]: $s", atClaims.Uuid, rtClaims.Uuid)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if atErr == nil {
		initializers.DB.Where("uuid = ?", atClaims.Uuid).First(&user)
		log.Info("Valid access token / [uuid]:", atClaims.Uuid)
	} else if atErr != nil && rtErr == nil {
		initializers.DB.Where("uuid = ?", rtClaims.Uuid).First(&user)
		log.Info("Valid refresh token / [uuid]:", rtClaims.Uuid)
	} else {
		if err = crypto.DeleteTokens(rtClaims.Uuid); err != nil {
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		log.Error(atErr, rtErr)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenDto, err := crypto.GenerateTokens(user)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = crypto.SaveTokens(user.Uuid, tokenDto.RefreshToken); err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userDto := dto.NewUserDto(*user, *tokenDto)

	c.JSON(201, userDto)
	log.Info("User auto logged in / [uuid]:", userDto.Uuid)
}

func Logout(c *gin.Context) {
	var body dto.LogoutDto
	if err := c.Bind(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := crypto.DeleteTokens(body.Uuid); err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(201, gin.H{})
	log.Info("Kakao user logged out / [uuid]:", body.Uuid)
}

func UseAuthRouter(g *gin.Engine) {
	sg := g.Group("/auth")

	sg.POST("/kakao", KakaoLogin)
	sg.POST("/auto-login", AutoLogin)
	sg.POST("/logout", Logout)
}
