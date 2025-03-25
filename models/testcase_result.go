package models

import "gorm.io/gorm"

type TestcaseResult struct {
	gorm.Model
	SubmissionID uint   `json:"submission_id"` // 所属提交
	CaseID       string `json:"case_id"`       // 测试点编号，如 "1", "2"
	Status       string `json:"status"`        // Accepted / Wrong Answer / Runtime Error
	Output       string `json:"output"`        // 用户输出
	Expected     string `json:"expected"`      // 预期输出
	RuntimeMs    int    `json:"runtime_ms"`    // 执行时间（毫秒）
}
