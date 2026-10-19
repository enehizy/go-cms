package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/enehizy/go-cms/database"
	"github.com/enehizy/go-cms/models"
	"github.com/enehizy/go-cms/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type MediaRequest struct {
	Uri      string `json:"uri" binding:"required"`
	MimeType string `json:"mime_type" binding:"required"`
}

func GetMedias(c *gin.Context) {
	db, err := database.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var mediaArray []models.Media
	if err := db.Model(&models.Media{}).Preload("Post").Find(&mediaArray).Error; err != nil {
		errObj := utils.RequestErrors{
			Error:   err,
			Message: "unable to fetch media data",
		}
		c.JSON(http.StatusNotFound, errObj.GetError())
		return

	}
	c.JSON(http.StatusOK, gin.H{"data": mediaArray})
}
func GetMedia(c *gin.Context) {
	id := c.Param("id")

	var media models.Media
	db, err := database.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if err := db.Preload("Post").Where("id = ?", id).First(&media).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errObj := utils.RequestErrors{
				Message: fmt.Sprintf(`media with %s doesnt exist`, id),
				Error:   err,
			}
			c.JSON(http.StatusNotFound, errObj.GetError())
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": media})
}
func CreateMedia(c *gin.Context) {
	var body MediaRequest
	if !utils.BindAndValidate(c, &body) {
		return
	}
	db, err := database.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	media := models.Media{
		Uri:      body.Uri,
		MimeType: body.MimeType,
	}
	if err := db.Model(&models.Media{}).Create(&media).Error; err != nil {
		errObj := utils.RequestErrors{
			Message: "failed to create a media resource",
			Error:   err,
		}
		c.JSON(http.StatusNotFound, errObj.GetError())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "successfully created a media resource", "data": media})
}

// type Media struct {
// 	BaseModel
// 	Uri      string `gorm:"not null; unique" json:"uri"`
// 	MimeType string `gorm:"column:mime_type;not null;check:mime_type in ('image/jpg','image/png','video/mp4','video/webp')" json:"mime_type"`
// 	Post     []Post `gorm:"many2many" json:"post"`
// }

func UpdateMedia(c *gin.Context) {
	var body struct {
		Uri      *string          `json:"uri"`
		MimeType *models.Mimetype `json:"mime_type"`
		Post any `json:"post"`
	}
	if !utils.BindAndValidate(c, &body) {
		return
	}
	db, dbErr := database.GetDBFromContext(c)
	if dbErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if body.Uri == nil && body.MimeType == nil {
		errObj := utils.RequestErrors{
			Message: "at least one field must be provided",
			Error:   errors.New("missing fields"),
		}
		c.JSON(http.StatusBadRequest, errObj.GetError())
		return
	}
	if body.Post!=nil{
		errObj := utils.RequestErrors{
			Message: utils.ERR_BAD_REQUEST.Error(),
			Error:   errors.New("can only update post in the /post endpoint"),
		}
		c.JSON(http.StatusBadRequest, errObj.GetError())
		return
	}
	updates := make(map[string]any)
	if body.Uri != nil {
		updates["uri"] = *body.Uri
	}
	if body.MimeType != nil {
		updates["mime_type"] = *body.MimeType
	}
	var media models.Media
	id := c.Param("id")
	if err := db.Model(models.Media{}).Where("id = ?", id).Updates(updates).Scan(&media).Error; err != nil {
		message := fmt.Sprintf("failed to update media with id %s", id)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				message = "uri already exist ,must be unique"
			}
			if pgErr.Code == "23514" {
				message = "unsupported mime_type use ('image/jpg','image/png','video/mp4','video/webp')"
			}
		}
		errObj := utils.RequestErrors{
			Error:   err,
			Message: message,
		}
		c.JSON(http.StatusInternalServerError, errObj.GetError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"media": media})
}
func DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	var media models.Media
	if err := db.Model(&models.Media{}).Where("id = ?", id).Delete(&media).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errObj := utils.RequestErrors{
				Message: fmt.Sprintf(`media with %s doesnt exist`, id),
				Error:   err,
			}
			c.JSON(http.StatusNotFound, errObj.GetError())
			return
		}
		errObj := utils.RequestErrors{
			Message: fmt.Sprintf(`failed to delete media with %s`, id),
			Error:   err,
		}
		c.JSON(http.StatusNotFound, errObj.GetError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("sucessfully deleted medai with id :%s", id), "data": media})
}
