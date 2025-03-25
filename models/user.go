package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"` // 注意：后续我们将加密存储
	Email    string `json:"email"`
}
