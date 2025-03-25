package utils

import (
	"deepjudge/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("deepjudge.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(&models.User{}, &models.Problem{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
}
