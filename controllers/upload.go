package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 验证测试用例文件名格式
func validateTestcaseFilename(filename string) bool {
	// 文件名必须是 数字.in 或 数字.out 格式
	base := filepath.Base(filename)
	if !strings.HasSuffix(base, ".in") && !strings.HasSuffix(base, ".out") {
		return false
	}
	numStr := strings.TrimSuffix(strings.TrimSuffix(base, ".in"), ".out")
	for _, c := range numStr {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func UploadTestcases(c *gin.Context) {
	// 获取题目 ID
	problemID := c.Param("id")

	// 构建目标目录
	targetDir := fmt.Sprintf("testcases/%s", problemID)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		c.JSON(500, gin.H{"error": "创建目录失败"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "上传文件失败"})
		return
	}

	// 检查是否是zip文件
	if filepath.Ext(file.Filename) == ".zip" {
		// 保存zip文件到临时位置
		tempZip := filepath.Join(os.TempDir(), "upload.zip")
		if err := c.SaveUploadedFile(file, tempZip); err != nil {
			c.JSON(500, gin.H{"error": "保存文件失败"})
			return
		}
		defer os.Remove(tempZip)

		// 打开zip文件
		reader, err := zip.OpenReader(tempZip)
		if err != nil {
			c.JSON(500, gin.H{"error": "无法打开zip文件"})
			return
		}
		defer reader.Close()

		// 记录成功解压的文件
		var extracted []string
		inFiles := make(map[string]bool)
		outFiles := make(map[string]bool)

		// 解压文件
		for _, file := range reader.File {
			// 跳过目录
			if file.FileInfo().IsDir() {
				continue
			}

			// 验证文件名格式
			if !validateTestcaseFilename(file.Name) {
				continue
			}

			// 构建目标路径
			filename := filepath.Base(file.Name)
			targetPath := filepath.Join(targetDir, filename)

			// 记录.in和.out文件
			if strings.HasSuffix(filename, ".in") {
				inFiles[strings.TrimSuffix(filename, ".in")] = true
			} else if strings.HasSuffix(filename, ".out") {
				outFiles[strings.TrimSuffix(filename, ".out")] = true
			}

			// 解压文件
			src, err := file.Open()
			if err != nil {
				continue
			}

			dst, err := os.Create(targetPath)
			if err != nil {
				src.Close()
				continue
			}

			_, err = io.Copy(dst, src)
			src.Close()
			dst.Close()

			if err == nil {
				extracted = append(extracted, filename)
			}
		}

		// 验证每个测试点都有对应的.in和.out文件
		var validPairs []string
		for num := range inFiles {
			if outFiles[num] {
				validPairs = append(validPairs, num)
			}
		}

		if len(validPairs) == 0 {
			c.JSON(400, gin.H{
				"error": "未找到有效的测试用例对（需要对应的.in和.out文件）",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "zip文件解压成功",
			"files":   extracted,
			"pairs":   validPairs,
		})
		return
	}

	// 处理单个文件上传
	if !validateTestcaseFilename(file.Filename) {
		c.JSON(400, gin.H{"error": "文件名格式错误，应为 数字.in 或 数字.out"})
		return
	}

	dst := filepath.Join(targetDir, filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败"})
		return
	}

	c.JSON(200, gin.H{
		"message": "文件上传成功",
		"file":    file.Filename,
	})
}
