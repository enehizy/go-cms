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

func GetPosts(c *gin.Context) {
	var posts []models.Post
	db, err := db.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err := db.Preload("Media").Find(&posts).Error; err != nil {
		errObj := utils.RequestErrors{Error: err, Message: "failed to fetch post"}
		c.JSON(http.StatusNotFound, errObj.GetError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}
func GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	db, err := db.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err := db.Preload("Media").Where("id = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errObj := utils.RequestErrors{
				Message: fmt.Sprintf(`post with %s doesnt exist`, id),
				Error:   err,
			}
			c.JSON(http.StatusInternalServerError, errObj.GetError())
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}
func CreatePost(c *gin.Context) {
	var body struct {
		Title   string         `json:"title" binding:"required"`
		Content string         `json:"content" binding:"required"`
		Author  string         `json:"author" binding:"required"`
		Media   []MediaRequest `json:"media"`
	}

	if !utils.BindAndValidate(c, &body) {
		return
	}
	post := models.Post{
		Title:   body.Title,
		Content: body.Content,
		Author:  body.Author,
	}
	if len(body.Media) > 0 {
		var records []models.Media
		for _, v := range body.Media {
			records = append(records, models.Media{
				Uri:      v.Uri,
				MimeType: v.MimeType,
			})
		}
		post.Media = records
	}
	db, err := db.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err := db.Model(&models.Post{}).Create(&post).Error; err != nil {
		errObj := utils.RequestErrors{
			Message: "failed to create post",
			Error:   err,
		}
		c.JSON(http.StatusInternalServerError, errObj.GetError())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "post sucessfully  created", "data": post})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("no post with is %v avaliable yet", id)})

}
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	db, err := db.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	var post models.Post
	if err := db.Model(&models.Post{}).Where("id = ?", id).Delete(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errObj := utils.RequestErrors{
				Message: fmt.Sprintf(`post with %s doesnt exist`, id),
				Error:   err,
			}
			c.JSON(http.StatusNotFound, errObj.GetError())
			return
		}
		errObj := utils.RequestErrors{
			Message: fmt.Sprintf(`failed to delete post with %s`, id),
			Error:   err,
		}
		c.JSON(http.StatusNotFound, errObj.GetError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("sucessfully deleted post with id :%s", id), "data": post})

}
