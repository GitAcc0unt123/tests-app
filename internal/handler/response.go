package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

type idResponse struct {
	Id int `json:"id"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		log.Print(err.Error())
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Message: message,
	})
}
