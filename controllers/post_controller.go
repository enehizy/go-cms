package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/enehizy/go-cms/database"
	"github.com/enehizy/go-cms/models"
	"github.com/enehizy/go-cms/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPosts(c *gin.Context) {
	var posts []models.Post
	db, err := database.GetDBFromContext(c)
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
	db, err := database.GetDBFromContext(c)
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
	db, err := database.GetDBFromContext(c)
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

	var body struct {
		Title     *string        `json:"title" `
		Content   *string        `json:"content" `
		Author    *string        `json:"author" `
		Published *bool          `json:"published"`
		Media     []MediaRequest `json:"media"`
	}
	if !utils.BindAndValidate(c, &body) {
		return
	}
	db, dbErr := database.GetDBFromContext(c)
	if dbErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if body.Title == nil && body.Author == nil && body.Content == nil && body.Published == nil {
		errObj := utils.RequestErrors{
			Message: "at least one field must be provided",
			Error:   errors.New("missing fields"),
		}
		c.JSON(http.StatusBadRequest, errObj.GetError())
		return
	}
	if body.Media != nil {
		errObj := utils.RequestErrors{
			Message: utils.ERR_BAD_REQUEST.Error(),
			Error:   errors.New("can only update media in the /media endpoint"),
		}
		c.JSON(http.StatusBadRequest, errObj.GetError())
		return
	}
	updates := make(map[string]any)
	if body.Title != nil {
		updates["title"] = *body.Title
	}
	if body.Author != nil {
		updates["author"] = *body.Author
	}
	if body.Content != nil {
		updates["content"] = *body.Content
	}
	if body.Published != nil {
		updates["published"] = *body.Published
	}

	id := c.Param("id")

	var post models.Post
	if err := db.Model(&models.Post{}).Where("id = ?", id).Updates(updates).Scan(&post).Error; err != nil {
		errObj := utils.RequestErrors{
			Error:   err,
			Message: fmt.Sprintf("failed to update post with id %s", id),
		}
		c.JSON(http.StatusInternalServerError, errObj.GetError())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})

}
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDBFromContext(c)
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
