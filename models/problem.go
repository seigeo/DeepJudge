package models

import "gorm.io/gorm"

type Problem struct {
	gorm.Model

	ID           uint   `json:"id"` // 👈 添加这一行，或者在匿名结构体上标注
	Title        string `json:"title"`
	Description  string `json:"description"` // 题目描述
	Input        string `json:"input"`       // 输入描述
	Output       string `json:"output"`      // 输出描述
	SampleInput  string `json:"sample_input"`
	SampleOutput string `json:"sample_output"`
	Difficulty   string `json:"difficulty"` // easy / medium / hard
}
