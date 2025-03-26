package services

import (
	"bytes"
	"deepjudge/models"
	"deepjudge/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AnalysisResult struct {
	ReasoningContent string
	Content          string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Role             string `json:"role"`
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"message"`
	} `json:"choices"`
}

// AnalyzeCodeWithDeepSeek 使用 DeepSeek-R1 分析代码，返回分析结果
func AnalyzeCodeWithDeepSeek(submissionID uint) (AnalysisResult, error) {
	// fmt.Println ("I m working")
	// 获取提交记录
	var submission models.Submission
	if err := utils.DB.Where("id = ?", submissionID).First(&submission).Error; err != nil {
		return AnalysisResult{}, fmt.Errorf("提交记录未找到：%v", err)
	}

	// 获取题目信息
	var problem models.Problem
	if err := utils.DB.Where("id = ?", submission.ProblemID).First(&problem).Error; err != nil {
		return AnalysisResult{}, fmt.Errorf("题目信息未找到：%v", err)
	}

	// 构造消息
	messages := []Message{
		{Role: "user", Content: fmt.Sprintf("假如你是一名优秀的算法竞赛教练，请从以下三个层面评价该学生的代码：结构摘要、可读性建议、问题提示。下面是题目描述：\n\n%s\n\n以及用户代码：\n\n%s", problem.Description, submission.Code)},
	}

	// fmt.Println(messages)
	// 调用 DeepSeek API
	reqBody := ChatRequest{
		Model:    "deepseek-reasoner", // 使用 DeepSeek 的推理模型
		Messages: messages,
	}

	data, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return AnalysisResult{}, err
	}

	// fmt.Println(req)

	// 从环境变量中读取 API Key
	apiKey := "sk-2ae026614a634e00ba1cb0078056ac90"

	// fmt.Println("aliving")

	// fmt.Println(apiKey)

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// fmt.Println(req)

	// 发送请求并处理响应
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return AnalysisResult{}, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("读取响应体失败: %v", err)
	}

	// fmt.Println("Response Body:", string(body)) // 打印响应体内容

	// 解析响应体
	var parsed ChatResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return AnalysisResult{}, fmt.Errorf("解析响应失败: %v", err)
	}

	// fmt.Println("Parsed Response:", parsed) // 打印解析后的响应内容

	// 获取推理过程和结论
	if len(parsed.Choices) == 0 {
		return AnalysisResult{}, fmt.Errorf("DeepSeek 返回结果为空")
	}

	msg := parsed.Choices[0].Message
	analysisResult := AnalysisResult{
		ReasoningContent: msg.ReasoningContent,
		Content:          msg.Content,
	}

	return analysisResult, nil
}

// services/analysis.go
func SaveAnalysisResult(analysisID uint, result AnalysisResult) error {
	// 更新分析结果
	return utils.DB.Model(&models.AnalysisResult{}).
		Where("id = ?", analysisID).
		Update("reasoning_content", result.ReasoningContent).
		Update("content", result.Content).Error
}

// UpdateAnalysisStatus 更新分析状态（pending, analyzing, completed, failed）
func UpdateAnalysisStatus(analysisID uint, status string) error {
	return utils.DB.Model(&models.AnalysisResult{}).
		Where("id = ?", analysisID).
		Update("status", status).Error
}
