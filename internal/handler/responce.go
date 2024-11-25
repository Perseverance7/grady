package handler

import (
	"github.com/gin-gonic/gin"
)

func newErrorResponce(c *gin.Context, statusCode int, err error) {
	c.Error(err)
	c.JSON(statusCode, gin.H{"error": err.Error()})
}