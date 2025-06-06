package services

import (
	"context"
	"deepjudge/models"
	"deepjudge/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	ProblemStatsKeyPrefix = "problem_stats:"
	UserRankKeyPrefix     = "user_rank:"
	RankingZSetKey        = "user_ranking"
	RankingExpiration     = 24 * time.Hour
)

// 将用户ID转换为Redis中使用的格式
func userIDToRedisKey(userID uint) string {
	return fmt.Sprintf("user:%d", userID)
}

// 从Redis key中提取用户ID
func redisKeyToUserID(key string) (uint, error) {
	// 去掉 "user:" 前缀
	if len(key) <= 5 || key[:5] != "user:" {
		return 0, fmt.Errorf("invalid user key format: %s", key)
	}
	id, err := strconv.ParseUint(key[5:], 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// 更新题目统计信息
func UpdateProblemStats(problemID uint, isAccepted bool) error {
	key := fmt.Sprintf("%s%d", ProblemStatsKeyPrefix, problemID)

	// 更新数据库
	var problem models.Problem
	if err := utils.DB.First(&problem, problemID).Error; err != nil {
		return err
	}

	problem.SubmissionCount++
	if isAccepted {
		problem.AcceptedCount++
	}
	problem.PassRate = float64(problem.AcceptedCount) / float64(problem.SubmissionCount)

	if err := utils.DB.Save(&problem).Error; err != nil {
		return err
	}

	// 更新Redis缓存
	statsData, _ := json.Marshal(map[string]interface{}{
		"accepted_count":   problem.AcceptedCount,
		"submission_count": problem.SubmissionCount,
		"pass_rate":        problem.PassRate,
	})

	return utils.RDB.Set(context.Background(), key, string(statsData), 24*time.Hour).Err()
}

// 更新用户排名
func UpdateUserRanking(userID uint, isAccepted bool) error {
	if !isAccepted {
		return nil
	}

	ctx := context.Background()

	// 检查用户是否存在
	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在: %d", userID)
		}
		return err
	}

	// 直接从数据库获取用户通过的题目数量
	var solvedCount int64
	if err := utils.DB.Model(&models.Submission{}).
		Where("user_id = ? AND result = ?", userID, "Accepted").
		Distinct("problem_id").
		Count(&solvedCount).Error; err != nil {
		return err
	}

	// 更新Redis排行榜
	if err := utils.RDB.ZAdd(ctx, RankingZSetKey, redis.Z{
		Score:  float64(solvedCount),
		Member: userIDToRedisKey(userID),
	}).Err(); err != nil {
		return err
	}

	// 设置过期时间
	utils.RDB.Expire(ctx, RankingZSetKey, RankingExpiration)

	return nil
}

// 获取用户排名
func GetUserRanking(userID uint) (int64, int, error) {
	ctx := context.Background()

	// 检查用户是否存在
	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, 0, fmt.Errorf("用户不存在: %d", userID)
		}
		return 0, 0, err
	}

	// 获取用户排名
	rank, err := utils.RDB.ZRevRank(ctx, RankingZSetKey, userIDToRedisKey(userID)).Result()
	if err == redis.Nil {
		// 如果Redis中没有数据，重新计算并更新
		if err := UpdateUserRanking(userID, true); err != nil {
			return 0, 0, err
		}
		rank, err = utils.RDB.ZRevRank(ctx, RankingZSetKey, userIDToRedisKey(userID)).Result()
		if err != nil {
			return 0, 0, err
		}
	} else if err != nil {
		return 0, 0, err
	}

	// 获取用户解题数
	score, err := utils.RDB.ZScore(ctx, RankingZSetKey, userIDToRedisKey(userID)).Result()
	if err != nil {
		return 0, 0, err
	}

	return rank + 1, int(score), nil
}

// 获取排行榜前N名
func GetTopNUsers(n int) ([]map[string]interface{}, error) {
	ctx := context.Background()

	// 获取前N名用户ID和分数
	results, err := utils.RDB.ZRevRangeWithScores(ctx, RankingZSetKey, 0, int64(n-1)).Result()
	if err != nil {
		return nil, err
	}

	// 如果没有数据，初始化排行榜
	if len(results) == 0 {
		// 从数据库获取所有用户的解题数据
		var submissions []struct {
			UserID    uint
			ProblemID uint
		}
		if err := utils.DB.Model(&models.Submission{}).
			Where("result = ?", "Accepted").
			Select("DISTINCT user_id, problem_id").
			Find(&submissions).Error; err != nil {
			return nil, err
		}

		// 统计每个用户的解题数
		userSolved := make(map[uint]int)
		for _, sub := range submissions {
			userSolved[sub.UserID]++
		}

		// 更新Redis排行榜
		for userID, count := range userSolved {
			// 检查用户是否存在
			var user models.User
			if err := utils.DB.First(&user, userID).Error; err != nil {
				continue // 跳过不存在的用户
			}

			utils.RDB.ZAdd(ctx, RankingZSetKey, redis.Z{
				Score:  float64(count),
				Member: userIDToRedisKey(userID),
			})
		}

		// 重新获取排行榜数据
		results, err = utils.RDB.ZRevRangeWithScores(ctx, RankingZSetKey, 0, int64(n-1)).Result()
		if err != nil {
			return nil, err
		}
	}

	var rankings []map[string]interface{}
	for i, z := range results {
		userKey := z.Member.(string)
		userID, err := redisKeyToUserID(userKey)
		if err != nil {
			// 如果用户ID格式不正确，从排行榜中移除
			utils.RDB.ZRem(ctx, RankingZSetKey, userKey)
			continue
		}

		var user models.User
		if err := utils.DB.First(&user, userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 如果用户不存在，从排行榜中移除
				utils.RDB.ZRem(ctx, RankingZSetKey, userKey)
				continue
			}
			return nil, err
		}

		rankings = append(rankings, map[string]interface{}{
			"rank":         i + 1,
			"user_id":      userID,
			"username":     user.Username,
			"signature":    user.Signature,
			"solved_count": int(z.Score),
		})
	}

	return rankings, nil
}
