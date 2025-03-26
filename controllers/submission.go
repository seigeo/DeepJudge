package controllers

import (
	"deepjudge/models"
	"deepjudge/services"
	"deepjudge/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SubmitCode(c *gin.Context) {
	var req struct {
		Code     string `json:"code"`
		Language string `json:"language"`
	}

	fmt.Println(req)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式错误"})
		return
	}

	problemIDStr := c.Param("id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	userID := c.GetUint("user_id")

	var problem models.Problem
	if err := utils.DB.First(&problem, problemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	// 创建提交记录（先写入数据库，评测后再更新）
	sub := models.Submission{
		UserID:     userID,
		ProblemID:  uint(problemID),
		Code:       req.Code,
		Language:   req.Language,
		Result:     "Pending",
		SubmitTime: time.Now(),
	}
	utils.DB.Create(&sub)

	// 异步评测（入队）
	services.EnqueueSubmission(sub)

	c.JSON(http.StatusOK, gin.H{
		"message":       "提交成功，已加入评测队列",
		"submission_id": sub.ID,
	})
}

func GetSubmission(c *gin.Context) {
	id := c.Param("id")
	var sub models.Submission

	if err := utils.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提交不存在"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func GetTestcaseResults(c *gin.Context) {
	id := c.Param("id")
	var sub models.Submission
	var results []models.TestcaseResult

	if err := utils.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提交不存在"})
		return
	}

	if err := utils.DB.Where("submission_id = ?", id).Find(&results).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评测记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"submission": sub,
		"results":    results,
	})
}

func GetUserSubmissions(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 读取分页参数
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(limit)

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	var submissions []models.Submission

	// 计算总数
	utils.DB.Model(&models.Submission{}).Where("user_id = ?", userID).Count(&total)

	// 查询分页数据
	utils.DB.
		Where("user_id = ?", userID).
		Order("submit_time desc").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&submissions)

	// 返回分页结果
	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"page":        pageNum,
		"limit":       pageSize,
		"submissions": submissions,
	})
}
