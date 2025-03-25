package routes

import (
	"deepjudge/controllers"
	"deepjudge/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// 公共题目接口
	r.GET("/problems", controllers.GetProblems)
	r.GET("/problems/:id", controllers.GetProblemByID)

	// 受保护的题目接口（添加题目）
	auth := r.Group("/auth", middleware.AuthMiddleware())
	{
		auth.POST("/problems", controllers.CreateProblem)
		auth.PUT("/problems/:id", controllers.UpdateProblem)    // 编辑题目
		auth.DELETE("/problems/:id", controllers.DeleteProblem) // 删除题目

		auth.POST("/submit", controllers.SubmitCode)
		auth.GET("/submissions/:id", controllers.GetSubmission)
	}
}
