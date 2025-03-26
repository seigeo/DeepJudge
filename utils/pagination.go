package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page  int
	Limit int
}

// Paginate 用于 GORM 的 Scopes，自动从 Gin 的 Query 中解析 page/limit 参数
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	// 读取参数
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}
