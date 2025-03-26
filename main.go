package main

import (
	"deepjudge/routes"
	"deepjudge/services"
	"deepjudge/utils"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitRedis()
	utils.InitDB()
	services.StartJudgeWorkerPool(4) // 启动 4 个评测 worker
	// services.ProcessProblemUpdateQueue()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.SetupRoutes(r)

	r.Run(":8080")
}
