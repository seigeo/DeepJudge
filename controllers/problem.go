package controllers

import (
	"deepjudge/models"
	"deepjudge/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取所有题目
func GetProblems(c *gin.Context) {
	var problems []models.Problem
	utils.DB.Find(&problems)
	c.JSON(http.StatusOK, problems)
}

// 获取指定题目
func GetProblemByID(c *gin.Context) {
	id := c.Param("id")
	var problem models.Problem
	if err := utils.DB.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}
	c.JSON(http.StatusOK, problem)
}

// 添加题目（需要认证）
func CreateProblem(c *gin.Context) {
	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式错误"})
		return
	}
	utils.DB.Create(&problem)
	c.JSON(http.StatusOK, gin.H{"message": "题目创建成功"})
}

func UpdateProblem(c *gin.Context) {
	id := c.Param("id")
	var problem models.Problem

	if err := utils.DB.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	var updated models.Problem
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式错误"})
		return
	}

	// 更新字段
	problem.Title = updated.Title
	problem.Description = updated.Description
	problem.Input = updated.Input
	problem.Output = updated.Output
	problem.SampleInput = updated.SampleInput
	problem.SampleOutput = updated.SampleOutput
	problem.Difficulty = updated.Difficulty

	utils.DB.Save(&problem)
	c.JSON(http.StatusOK, gin.H{"message": "题目更新成功"})
}

func DeleteProblem(c *gin.Context) {
	id := c.Param("id")
	var problem models.Problem

	if err := utils.DB.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	utils.DB.Delete(&problem)
	c.JSON(http.StatusOK, gin.H{"message": "题目删除成功"})
}
