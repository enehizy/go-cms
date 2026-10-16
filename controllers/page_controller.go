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

//	type Page struct {
//		BaseModel
//		Title   string `gorm:"not null"`
//		Slug    string `gorm:"size:255;unique;not null"`
//		Content string `gorm:"not null"`
//	}
func GetPages(c *gin.Context) {
	var page []models.Page
	db, err := db.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err := db.Find(&page).Error; err != nil {
		errObj := utils.RequestErrors{Error: err, Message: "failed to fetch page"}
		c.JSON(http.StatusNotFound, errObj.GetError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": page})
}
func GetPage(c *gin.Context) {
	id := c.Param("id")

	var page models.Page
	db, err := db.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err := db.Where("id = ?", id).First(&page).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errObj := utils.RequestErrors{
				Message: fmt.Sprintf(`page with %s doesnt exist`, id),
				Error:   err,
			}
			c.JSON(http.StatusInternalServerError, errObj.GetError())
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": page})
}
func CreatePage(c *gin.Context) {
	var body struct {
		Title   string `json:"title" binding:"required"`
		Slug    string `json:"slug" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if !utils.BindAndValidate(c, &body) {
		return
	}

	db, err := db.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	page := models.Page{
		Title:   body.Title,
		Slug:    body.Slug,
		Content: body.Content,
	}
	if err := db.Model(&models.Page{}).Create(&page).Error; err != nil {
		errObj := utils.RequestErrors{
			Message: "failed to create a page resource",
			Error:   err,
		}
		c.JSON(http.StatusNotFound, errObj.GetError())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "successfully created a page resource", "data": page})

}
func UpdatePage(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("no page with is %v avaliable yet", id)})
}
func DeletePage(c *gin.Context) {
	id := c.Param("id")
	db, err := db.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	var page models.Page
	if err := db.Model(&models.Page{}).Where("id = ?", id).Delete(&page).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errObj := utils.RequestErrors{
				Message: fmt.Sprintf(`page with %s doesnt exist`, id),
				Error:   err,
			}
			c.JSON(http.StatusNotFound, errObj.GetError())
			return
		}
		errObj := utils.RequestErrors{
			Message: fmt.Sprintf(`failed to delete page with %s`, id),
			Error:   err,
		}
		c.JSON(http.StatusNotFound, errObj.GetError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("sucessfully deleted medai with id :%s", id), "data": page})
}
