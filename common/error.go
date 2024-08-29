package common

import (
	"github.com/gin-gonic/gin"
)

/*
	{
		status: number,
		message: ".....",
		error: error
	}
*/
type ErrorRes struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewErrorResponse(c *gin.Context, status int, message string, err interface{}) {
	errorResponse := ErrorRes{
		Status:  status,
		Message: message,
		Error:   err,
	}
	c.JSON(status, errorResponse)
}
