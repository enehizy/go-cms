package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/enehizy/go-cms/db"
	"github.com/enehizy/go-cms/models"
	"github.com/enehizy/go-cms/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MediaRequest struct {
	Uri      string `json:"uri" binding:"required"`
	MimeType string `json:"mime_type" binding:"required"`
}

func GetMedias(c *gin.Context) {
	db, err := db.GetDBFromContext(c)
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
	db, err := db.GetDBFromContext(c)
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
	db, err := db.GetDBFromContext(c)
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
func UpdateMedia(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("no media with is %v avaliable yet", id)})
}
func DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	db, err := db.GetDBFromContext(c)
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
