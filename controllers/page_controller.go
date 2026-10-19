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

//	type Page struct {
//		BaseModel
//		Title   string `gorm:"not null"`
//		Slug    string `gorm:"size:255;unique;not null"`
//		Content string `gorm:"not null"`
//	}
func GetPages(c *gin.Context) {
	var page []models.Page
	db, err := database.GetDBFromContext(c)
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
	db, err := database.GetDBFromContext(c)
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

	db, err := database.GetDBFromContext(c)
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
	var body struct {
		Title   *string `json:"title"`
		Slug    *string `json:"slug"`
		Content *string `json:"content"`
	}
	if !utils.BindAndValidate(c, &body) {
		return
	}
	db, err := database.GetDBFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if body.Title == nil && body.Slug == nil && body.Content == nil {
		errObj := utils.RequestErrors{
			Message: "at least one field must be provided",
			Error:   errors.New("missing fields"),
		}
		c.JSON(http.StatusBadRequest, errObj.GetError())
		return
	}
	updates := make(map[string]any)
	if body.Title != nil {
		updates["title"] = *body.Title
	}
	if body.Slug != nil {
		updates["slug"] = *body.Slug
	}
	if body.Content != nil {
		updates["content"] = *body.Content
	}
	id := c.Param("id")
	var page models.Page
	if err := db.Model(&models.Page{}).Where("id = ?", id).Updates(updates).Scan(&page).Error; err != nil {
		message := fmt.Sprintf("failed to update post with id %s", id)

		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				message = "slug already exist ,must be unique"
			}
		}

		errObj := utils.RequestErrors{
			Error:   err,
			Message: message,
		}
		c.JSON(http.StatusInternalServerError, errObj.GetError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": page})
}
func DeletePage(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDBFromContext(c)
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
