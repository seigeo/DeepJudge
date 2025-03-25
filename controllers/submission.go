package controllers

import (
	"deepjudge/models"
	"deepjudge/services"
	"deepjudge/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SubmitCode(c *gin.Context) {
	var req struct {
		ProblemID uint   `json:"problem_id"`
		Code      string `json:"code"`
		Language  string `json:"language"` // cpp / python
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式错误"})
		return
	}

	userID := c.GetUint("user_id")

	// 查询题目
	var problem models.Problem
	if err := utils.DB.First(&problem, req.ProblemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	// 评测代码
	result, _ := services.EvaluateCode(req.Code, req.Language, problem.SampleInput, problem.SampleOutput)

	// 保存记录
	sub := models.Submission{
		UserID:    userID,
		ProblemID: req.ProblemID,
		Code:      req.Code,
		Language:  req.Language,
		Result:    result,
	}
	utils.DB.Create(&sub)

	c.JSON(http.StatusOK, gin.H{
		"message": "提交成功",
		"result":  result,
		"id":      sub.ID,
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
