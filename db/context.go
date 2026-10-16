package db

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDBFromContext(c *gin.Context) (*gorm.DB, error) {
	_db, exists := c.Get("db")
	if !exists {
		return nil, errors.New("failed to get db from context")
	}
	db, ok := _db.(*gorm.DB)
	if !ok {
		return nil, errors.New("failed to get db from context")
	}
	return db, nil
}
