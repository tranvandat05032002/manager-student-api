package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func Pagination(c *gin.Context) (int, int, int) {
	limitStr := c.Query("limit")
	pageStr := c.Query("page")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	skip := limit * (page - 1)
	return page, limit, skip
}
