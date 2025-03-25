package models

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	UserID      uint      `json:"user_id"`
	ProblemID   uint      `json:"problem_id"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	Result      string    `json:"result"` // 总体结果
	SubmitTime  time.Time `json:"submit_time"`
	PassedCount int       `json:"passed_count"` // ✅ 通过数量
	TotalCount  int       `json:"total_count"`  // ✅ 总数量
}
