package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"deepjudge/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// 限流中间件：限制每个用户每 N 秒最多 M 次请求
func RateLimitMiddleware(maxReq int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录用户"})
			return
		}
		var userID uint
		switch v := val.(type) {
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		case uint:
			userID = v
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user_id 类型不支持"})
			return
		}

		key := fmt.Sprintf("rate_limit:user:%d", userID)
		countStr, err := utils.RDB.Get(ctx, key).Result()
		if err == redis.Nil {
			_ = utils.RDB.Set(ctx, key, 1, duration).Err()
		} else if err != nil {
			fmt.Println("Redis GET 错误：", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "限流系统错误"})
			return
		} else {
			count, convErr := strconv.Atoi(countStr)
			if convErr != nil {
				fmt.Println("Redis 值转换失败：", countStr)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "限流系统异常"})
				return
			}
			if count >= maxReq {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "提交过于频繁，请稍后再试"})
				return
			}
			_ = utils.RDB.Incr(ctx, key).Err()
		}

		c.Next()
	}
}
