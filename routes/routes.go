package routes

import (
	"deepjudge/controllers"
	"deepjudge/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register) //注册
	r.POST("/login", controllers.Login)       //登陆

	// 公共题目接口
	r.GET("/problems", controllers.GetProblems)               //查看题目列表
	r.GET("/problems/:id", controllers.GetProblemByID)        //查看题目
	r.GET("/problems/:id/stats", controllers.GetProblemStats) //查看题目统计
	r.GET("/leaderboard", controllers.GetLeaderboard)         //查看排行榜

	// 受保护接口
	auth := r.Group("/auth", middleware.AuthMiddleware())
	{
		problems := auth.Group("/problems")
		{
			problems.POST("/:id/submit", middleware.RateLimitMiddleware(2, 10*time.Second), // 10 秒最多 2 次
				controllers.SubmitCode) //提交代码
			problems.GET("/:id/all_submissions", controllers.GetProblemSubmissions)
			problems.GET("/:id/submissions", controllers.GetProblemSubmissions)
		}

		edit := auth.Group("/edit")
		{
			edit.POST("/add", controllers.CreateProblem)          //创建题目
			edit.PUT("/:id", controllers.UpdateProblem)           //编辑题目
			edit.DELETE("/:id", controllers.DeleteProblem)        //删除题目
			edit.POST("/:id/upload", controllers.UploadTestcases) //上传测试用例
		}

		// 与用户相关的通用接口
		auth.GET("/submissions", controllers.GetUserSubmissions)             // 获取当前用户所有提交记录
		auth.GET("/submissions/:id/results", controllers.GetTestcaseResults) // 查看测试点结果
		auth.GET("/rank", controllers.GetUserRank)                           // 获取当前用户排名
		auth.GET("/profile", controllers.GetUserProfile)                     // 获取用户资料
		auth.PUT("/profile", controllers.UpdateUserProfile)                  // 更新用户资料
		// 分析结果接口
		auth.POST("/submissions/:id/analyze", controllers.TriggerAnalysis) // 提交分析请求
		auth.GET("/submissions/:id/analysis", controllers.GetAllAnalyses)  // 查看分析结果
	}
}
