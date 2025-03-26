// models/analysis_result.go
package models

import "gorm.io/gorm"

type AnalysisResult struct {
	gorm.Model
	SubmissionID     uint   `json:"submission_id"`
	ReasoningContent string `gorm:"type:text" json:"reasoning_content"`
	Content          string `gorm:"type:text" json:"content"`
	Status           string `gorm:"type:varchar(50)" json:"status"`
}
