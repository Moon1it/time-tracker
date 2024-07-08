package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Description string `json:"description"`
}

type statusResponse struct {
	Description string `json:"description"`
}

func newErrorResponse(c *gin.Context, statusCode int, description string) {
	logrus.Error(description)
	c.AbortWithStatusJSON(statusCode, errorResponse{description})
}
