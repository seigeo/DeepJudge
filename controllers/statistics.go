package controllers

import (
	"deepjudge/models"
	"deepjudge/services"
	"deepjudge/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取题目统计信息
func GetProblemStats(c *gin.Context) {
	problemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	var problem models.Problem
	if err := utils.DB.First(&problem, problemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accepted_count":   problem.AcceptedCount,
		"submission_count": problem.SubmissionCount,
		"pass_rate":        problem.PassRate,
	})
}

// 获取用户排名
func GetUserRank(c *gin.Context) {
	userID := c.GetUint("user_id")
	rank, solvedCount, err := services.GetUserRanking(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取排名失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rank":         rank,
		"solved_count": solvedCount,
	})
}

// 获取排行榜
func GetLeaderboard(c *gin.Context) {
	limit := 50 // 默认显示前50名
	if limitStr := c.DefaultQuery("limit", "50"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	rankings, err := services.GetTopNUsers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取排行榜失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rankings": rankings,
	})
}
