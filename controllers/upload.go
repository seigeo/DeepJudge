package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadTestcases(c *gin.Context) {
	// 获取题目 ID
	problemID := c.Param("id")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": "上传格式错误"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(400, gin.H{"error": "未选择文件"})
		return
	}

	// 构建目标目录
	targetDir := fmt.Sprintf("testcases/%s", problemID)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		c.JSON(500, gin.H{"error": "创建目录失败"})
		return
	}

	var saved []string
	for _, file := range files {
		dst := filepath.Join(targetDir, filepath.Base(file.Filename))

		src, err := file.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		out, err := os.Create(dst)
		if err != nil {
			continue
		}
		defer out.Close()

		_, _ = io.Copy(out, src)
		saved = append(saved, file.Filename)
	}

	c.JSON(200, gin.H{
		"message": "上传成功",
		"files":   saved,
	})
}
