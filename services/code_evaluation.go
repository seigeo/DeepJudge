package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CaseResult struct {
	CaseID    string
	Status    string
	Output    string
	Expected  string
	RuntimeMs int
}

func EvaluateCode(code, lang string, problemID uint) ([]CaseResult, string, error) {
	// 生成临时代码文件
	tmpName := fmt.Sprintf("submissions/code_%d", time.Now().UnixNano())
	var codeFile string
	if lang == "cpp" {
		codeFile = tmpName + ".cpp"
	} else if lang == "python" {
		codeFile = tmpName + ".py"
	} else {
		return nil, "Unsupported Language", nil
	}

	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, "File Write Error", err
	}

	var results []CaseResult
	finalStatus := "Accepted"

	testDir := fmt.Sprintf("testcases/%d", problemID)
	files, err := os.ReadDir(testDir)
	if err != nil {
		return nil, "Missing Testcases", nil
	}

	// 获取所有 .in 文件名（去掉后缀）
	var testCases []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".in") {
			testCases = append(testCases, strings.TrimSuffix(f.Name(), ".in"))
		}
	}

	// 遍历每组测试点
	for _, caseID := range testCases {
		inputPath := filepath.Join(testDir, caseID+".in")
		outputPath := filepath.Join(testDir, caseID+".out")

		input, err1 := os.ReadFile(inputPath)
		expectedBytes, err2 := os.ReadFile(outputPath)
		if err1 != nil || err2 != nil {
			results = append(results, CaseResult{CaseID: caseID, Status: "Missing Testcase"})
			finalStatus = "Wrong Answer"
			continue
		}
		expected := strings.TrimSpace(string(expectedBytes))

		// 拼接 docker run 命令
		absCodePath, _ := filepath.Abs(codeFile)
		containerPath := "/app/code"

		dockerArgs := []string{
			"run", "--rm",
			"-v", fmt.Sprintf("%s:%s", absCodePath, containerPath),
			"deepjudge-runner",
			lang, containerPath, string(input),
		}

		// 设置超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "docker", dockerArgs...)
		start := time.Now()
		outputBytes, err := cmd.CombinedOutput()
		elapsed := time.Since(start).Milliseconds()

		actual := strings.TrimSpace(string(outputBytes))
		status := "Accepted"

		if ctx.Err() == context.DeadlineExceeded {
			status = "Time Limit Exceeded"
			finalStatus = "Wrong Answer"
		} else if err != nil {
			status = "Runtime Error"
			finalStatus = "Wrong Answer"
		} else if actual != expected {
			status = "Wrong Answer"
			finalStatus = "Wrong Answer"
		}

		results = append(results, CaseResult{
			CaseID:    caseID,
			Status:    status,
			Output:    actual,
			Expected:  expected,
			RuntimeMs: int(elapsed),
		})
	}

	// 清理临时文件
	_ = os.Remove(codeFile)

	return results, finalStatus, nil
}
