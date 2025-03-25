package controllers

import (
	"deepjudge/models"
	"deepjudge/services"
	"deepjudge/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SubmitCode(c *gin.Context) {
	var req struct {
		ProblemID uint   `json:"problem_id"`
		Code      string `json:"code"`
		Language  string `json:"language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式错误"})
		return
	}

	userID := c.GetUint("user_id")

	var problem models.Problem
	if err := utils.DB.First(&problem, req.ProblemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	// 多测试点评测
	caseResults, finalResult, _ := services.EvaluateCode(req.Code, req.Language, req.ProblemID)

	passCount := 0
	for _, r := range caseResults {
		if r.Status == "Accepted" {
			passCount++
		}
	}
	totalCount := len(caseResults)

	// 创建提交记录
	sub := models.Submission{
		UserID:      userID,
		ProblemID:   req.ProblemID,
		Code:        req.Code,
		Language:    req.Language,
		Result:      finalResult,
		SubmitTime:  time.Now(),
		PassedCount: passCount,
		TotalCount:  totalCount,
	}
	utils.DB.Create(&sub)

	// 保存每组测试点结果
	for _, res := range caseResults {
		utils.DB.Create(&models.TestcaseResult{
			SubmissionID: sub.ID,
			CaseID:       res.CaseID,
			Status:       res.Status,
			Output:       res.Output,
			Expected:     res.Expected,
			RuntimeMs:    res.RuntimeMs,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "提交成功",
		"submission_id": sub.ID,
		"result":        finalResult,
		"passed":        passCount,
		"total":         totalCount,
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
	var results []models.TestcaseResult
	if err := utils.DB.Where("submission_id = ?", id).Find(&results).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评测记录不存在"})
		return
	}
	c.JSON(http.StatusOK, results)
}
