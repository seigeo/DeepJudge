package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func EvaluateCode(code, lang string, input, expectedOutput string) (string, error) {
	tmpFile := fmt.Sprintf("submissions/tmp_%d", time.Now().UnixNano())
	var sourceFile, execFile string

	// 保存源码
	if lang == "cpp" {
		sourceFile = tmpFile + ".cpp"
		execFile = tmpFile
	} else if lang == "python" {
		sourceFile = tmpFile + ".py"
		execFile = sourceFile
	} else {
		return "Unsupported Language", nil
	}

	_ = os.WriteFile(sourceFile, []byte(code), 0644)

	var compileErr error
	if lang == "cpp" {
		cmd := exec.Command("g++", sourceFile, "-o", execFile)
		compileErr = cmd.Run()
		if compileErr != nil {
			return "Compilation Error", nil
		}
	}

	// 运行程序
	var cmd *exec.Cmd
	if lang == "cpp" {
		cmd = exec.Command(execFile)
	} else {
		cmd = exec.Command("python3", sourceFile)
	}

	stdin, _ := cmd.StdinPipe()
	// stdout, _ := cmd.Output()

	// 输入数据
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(input))
	}()

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return "Runtime Error", nil
	}

	output := strings.TrimSpace(string(outputBytes))
	expected := strings.TrimSpace(expectedOutput)

	if output == expected {
		return "Accepted", nil
	} else {
		return fmt.Sprintf("Wrong Answer\n你的输出: %s\n期望输出: %s", output, expected), nil
	}
}
