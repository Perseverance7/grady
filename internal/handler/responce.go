package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/Perseverance7/grady/pkg/logging"
)

type errorResponce struct {
	Message string `json:"message"`
}

func newErrorResponce(c *gin.Context, statusCode int, message string) {
	logger := logging.GetLogger()
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponce{Message: message})
}