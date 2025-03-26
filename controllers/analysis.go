// controllers/analysis.go
package controllers

import (
	"deepjudge/models"
	"deepjudge/services"
	"deepjudge/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /auth/submissions/:id/analyze
func TriggerAnalysis(c *gin.Context) {
	// 获取提交 ID
	sidStr := c.Param("id")
	sid, err := strconv.Atoi(sidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "提交ID非法"})
		return
	}

	// 创建分析记录
	analysis := models.AnalysisResult{
		SubmissionID: uint(sid),
		Status:       "pending", // 初始状态为 pending
	}

	if err := utils.DB.Create(&analysis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库写入失败"})
		return
	}

	// 异步启动分析任务
	go func() {
		// 更新分析状态为 analyzing
		_ = services.UpdateAnalysisStatus(analysis.ID, "analyzing")

		// 启动分析任务
		result, err := services.AnalyzeCodeWithDeepSeek(uint(sid)) // 传入题目和用户代码进行分析
		if err != nil {
			// 如果分析失败，更新状态为 failed
			_ = services.UpdateAnalysisStatus(analysis.ID, "failed")
			return
		}

		// 保存分析结果到数据库
		_ = services.SaveAnalysisResult(analysis.ID, result)

		// 更新分析状态为 completed
		_ = services.UpdateAnalysisStatus(analysis.ID, "completed")
	}()

	// 返回分析任务已提交
	c.JSON(http.StatusAccepted, gin.H{
		"message": "分析任务已提交",
		"status":  "pending",
	})
}

// GET /auth/submissions/:id/analyses
func GetAllAnalyses(c *gin.Context) {
	// 获取提交 ID
	sidStr := c.Param("id")
	sid, err := strconv.Atoi(sidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "提交ID非法"})
		return
	}

	// 查找该提交记录下的所有分析结果
	var analyses []models.AnalysisResult
	if err := utils.DB.Where("submission_id = ?", sid).Find(&analyses).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到分析结果"})
		return
	}

	// 返回所有分析记录
	c.JSON(http.StatusOK, gin.H{
		"submission_id": sid,
		"analyses":      analyses,
	})
}
