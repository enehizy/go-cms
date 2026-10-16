package utils

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestErrors struct {
	Message string

	Error error
}

func (re *RequestErrors) GetError() any {
	return gin.H{

		"message": re.Message,
		"error":   re.Error.Error(),
	}
}
func BindAndValidate[T any](c *gin.Context, dst *T) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		message := "missing field in body"
		if errors.Is(err, io.EOF) {
			message = "no body provided for the request"
		}

		errObj := RequestErrors{
			Message: message,
			Error:   err,
		}
		c.JSON(http.StatusBadRequest, errObj.GetError())
		return false
	}
	return true
}
