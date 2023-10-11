package controllers

import (
	"dev_community_server/dto"
	"dev_community_server/entity"
	"dev_community_server/initializers"
	"dev_community_server/middlewares"
	"dev_community_server/model"
	"dev_community_server/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/shyunku-libraries/go-logger"
	"net/http"
	"strconv"
)

func GetParty(c *gin.Context) {
	var entities []entity.PartyEntity

	if err := initializers.DB.Preload("Poster").Find(&entities).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	parties := make([]dto.PartyDto, len(entities))

	for idx, e := range entities {
		pad, err := dto.NewPartyDto(e)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		parties[idx] = *pad
	}

	c.JSON(http.StatusOK, parties)
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

	partyArticle := entity.PartyEntity{
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

	if err := initializers.DB.Create(&partyArticle).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
	log.Info("Article successfully created")
}

func GetPartyComment(c *gin.Context) {
	var entities []entity.PartyCommentEntity

	query := c.DefaultQuery("postId", "")
	postId, err := strconv.Atoi(query)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = initializers.DB.Preload("Poster").Where("post_id = ?", postId).Find(&entities).Order("created_at").Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	comments := make([]model.Comment, len(entities))

	for idx, e := range entities {
		nc := model.NewComment(e)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		comments[idx] = *nc
	}

	c.JSON(http.StatusOK, comments)
}

func CreatePartyComment(c *gin.Context) {
	var body dto.PartyCommentCreateDto
	uuid, ok := c.Get("uuid")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := c.Bind(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// main comment
	if body.Depth == 0 {
		partyComment := entity.PartyCommentEntity{
			PosterUuid: uuid.(string),
			PostId:     body.PostId,
			Comment:    body.Comment,
			Group:      nil,
			Depth:      body.Depth,
		}

		tx := initializers.DB.Begin()

		if err := tx.Create(&partyComment).Scan(&partyComment).Error; err != nil {
			tx.Rollback()
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		partyComment.Group = &partyComment.ID
		if err := tx.Save(&partyComment).Error; err != nil {
			tx.Rollback()
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		tx.Commit()

		c.JSON(http.StatusCreated, gin.H{})
		log.Info("Comment successfully created")
	} else { // sub comment
		partyComment := entity.PartyCommentEntity{
			PosterUuid: uuid.(string),
			PostId:     body.PostId,
			Comment:    body.Comment,
			Group:      body.Group,
			Depth:      body.Depth,
		}

		if err := initializers.DB.Create(&partyComment).Error; err != nil {
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusCreated, gin.H{})
		log.Info("Sub comment successfully created")

	}

}

func UsePartyRouter(g *gin.Engine) {
	sg := g.Group("/party")
	sg.Use(middlewares.TokenMiddleWare())

	sg.GET("/articles", GetParty)
	sg.POST("/create", CreateParty)
	sg.GET("/comment", GetPartyComment)
	sg.POST("comment-create", CreatePartyComment)
}
