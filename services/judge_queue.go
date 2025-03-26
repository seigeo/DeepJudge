package services

import (
	"deepjudge/models"
	"deepjudge/utils"
	"log"
)

type JudgeTask struct {
	Submission models.Submission
}

var judgeQueue = make(chan JudgeTask, 100)

// 外部调用：将任务加入评测队列
func EnqueueSubmission(sub models.Submission) {
	judgeQueue <- JudgeTask{Submission: sub}
}

// 启动评测 worker 池
func StartJudgeWorkerPool(n int) {
	for i := 0; i < n; i++ {
		go func(workerID int) {
			for task := range judgeQueue {
				sub := task.Submission
				log.Printf("[Worker %d] 评测 submission %d 开始", workerID, sub.ID)

				caseResults, result, err := EvaluateCode(sub.Code, sub.Language, sub.ProblemID)
				if err != nil {
					log.Printf("[Worker %d] 评测 submission %d 失败: %v", workerID, sub.ID, err)
					result = "System Error"
				}

				// 统计通过数
				passCount := 0
				for _, r := range caseResults {
					if r.Status == "Accepted" {
						passCount++
					}
				}

				// 更新提交记录
				utils.DB.Model(&models.Submission{}).
					Where("id = ?", sub.ID).
					Updates(map[string]interface{}{
						"result":       result,
						"passed_count": passCount,
						"total_count":  len(caseResults),
					})

				// 写入每组测试点评测结果
				for _, r := range caseResults {
					utils.DB.Create(&models.TestcaseResult{
						SubmissionID: sub.ID,
						CaseID:       r.CaseID,
						Status:       r.Status,
						Output:       r.Output,
						Expected:     r.Expected,
						RuntimeMs:    r.RuntimeMs,
					})
				}

				log.Printf("[Worker %d] 评测 submission %d 完成，结果: %s", workerID, sub.ID, result)
			}
		}(i)
	}
}
