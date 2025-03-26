package services

import (
	"context"
	"deepjudge/models"
	"deepjudge/utils"
	"encoding/json"
	"log"
	"time"
)

var ctx = context.Background()
var redisQueueKey = "judge_queue"

// 👇 入队
func EnqueueSubmission(sub models.Submission) {
	data, _ := json.Marshal(sub)
	utils.RDB.LPush(ctx, redisQueueKey, data)
}

// 👇 Worker Pool 出队并处理
func StartJudgeWorkerPool(n int) {
	for i := 0; i < n; i++ {
		go func(workerID int) {
			for {
				res, err := utils.RDB.BRPop(ctx, 0*time.Second, redisQueueKey).Result()
				if err != nil || len(res) < 2 {
					log.Printf("[Worker %d] 拉任务失败: %v", workerID, err)
					continue
				}

				var sub models.Submission
				if err := json.Unmarshal([]byte(res[1]), &sub); err != nil {
					log.Printf("[Worker %d] JSON解析失败: %v", workerID, err)
					continue
				}

				log.Printf("[Worker %d] 开始评测 submission %d", workerID, sub.ID)
				caseResults, result, _ := EvaluateCode(sub.Code, sub.Language, sub.ProblemID)

				passCount := 0
				for _, r := range caseResults {
					if r.Status == "Accepted" {
						passCount++
					}
				}

				utils.DB.Model(&models.Submission{}).
					Where("id = ?", sub.ID).
					Updates(map[string]interface{}{
						"result":       result,
						"passed_count": passCount,
						"total_count":  len(caseResults),
					})

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

				log.Printf("[Worker %d] 评测完成 submission %d: %s", workerID, sub.ID, result)
			}
		}(i)
	}
}
