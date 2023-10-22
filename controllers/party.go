package controllers

import (
	"dev_community_server/dto"
	"dev_community_server/entity"
	"dev_community_server/initializers"
	"dev_community_server/middlewares"
	"dev_community_server/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/shyunku-libraries/go-logger"
	"net/http"
	"strconv"
)

func GetParty(c *gin.Context) {
	var partyEntity []entity.PartyEntity

	if err := initializers.DB.Preload("Poster").Find(&partyEntity).Error; err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	partyListDto, err := dto.PartyListDtoFromEntity(partyEntity)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, partyListDto)
	log.Info("Article list successfully sent")
}

func CreateParty(c *gin.Context) {
	var body dto.PartyCreateDto
	uuid, ok := c.Get("uuid")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := c.Bind(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	techSkill, tsErr := json.Marshal(body.TechSkill)
	position, posErr := json.Marshal(body.Position)
	if tsErr != nil || posErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	deadline, dlErr := utils.StringToTime(body.Deadline)
	startDate, sdErr := utils.StringToTime(body.StartDate)
	if dlErr != nil || sdErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	partyEntity := entity.PartyEntity{
		PosterUuid:  uuid.(string),
		Title:       body.Title,
		Description: body.Description,
		TechSkill:   techSkill,
		Position:    position,
		Process:     body.Process,
		Category:    body.Category,
		Deadline:    *deadline,
		StartDate:   *startDate,
		Span:        body.Span,
		Location:    body.Location,
	}

	if err := initializers.DB.Create(&partyEntity).Error; err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
	log.Info("Article successfully created")
}

func GetPartyComment(c *gin.Context) {
	var partyCommentEntity []entity.PartyCommentEntity

	query := c.DefaultQuery("postId", "")
	postId, err := strconv.Atoi(query)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = initializers.DB.Preload("Poster").Where("post_id = ?", postId).Order("created_at desc").Find(&partyCommentEntity).Error; err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	partyCommentListDto := dto.PartyCommentListDtoFromEntity(partyCommentEntity)

	c.JSON(http.StatusOK, partyCommentListDto)
	log.Info("Party comment list responded to the client")
}

func CreatePartyComment(c *gin.Context) {
	var (
		body               dto.PartyCommentCreateDto
		partyCommentEntity entity.PartyCommentEntity
	)

	uuid, ok := c.Get("uuid")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := c.Bind(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tx := initializers.DB.Begin()

	// main comment
	if body.Depth == 0 {
		partyCommentEntity = entity.PartyCommentEntity{
			PosterUuid: uuid.(string),
			PostId:     body.PostId,
			Comment:    body.Comment,
			Group:      nil,
			Depth:      body.Depth,
		}

		if err := tx.Create(&partyCommentEntity).Error; err != nil {
			log.Error(err)
			tx.Rollback()
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		partyCommentEntity.Group = &partyCommentEntity.ID
		if err := tx.Save(&partyCommentEntity).Error; err != nil {
			log.Error(err)
			tx.Rollback()
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	} else { // sub comment
		partyCommentEntity = entity.PartyCommentEntity{
			PosterUuid: uuid.(string),
			PostId:     body.PostId,
			Comment:    body.Comment,
			Group:      body.Group,
			Depth:      body.Depth,
		}

		if err := tx.Create(&partyCommentEntity); err != nil {
			log.Error(err)
			tx.Rollback()
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Model(&partyCommentEntity).Preload("Poster").First(&partyCommentEntity).Error; err != nil {
		log.Error(err)
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tx.Commit()

	partyCommentDto := dto.PartyCommentDtoFromEntity(partyCommentEntity)

	c.JSON(http.StatusCreated, partyCommentDto)
	log.Info("Comment successfully created")
}

func UsePartyRouter(g *gin.Engine) {
	sg := g.Group("/party")
	sg.Use(middlewares.TokenMiddleWare())

	sg.GET("/articles", GetParty)
	sg.POST("/create", CreateParty)
	sg.GET("/comment", GetPartyComment)
	sg.POST("/comment-create", CreatePartyComment)
}
