package models

import "gorm.io/gorm"

type Submission struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	ProblemID uint   `json:"problem_id"`
	Code      string `json:"code"`     // 代码内容
	Language  string `json:"language"` // cpp / python
	Result    string `json:"result"`   // Accepted / WA / TLE / ...
}
