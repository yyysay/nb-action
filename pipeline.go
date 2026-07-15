package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// parseByComma 负责把 ["random", "64", ",", "bark", "Secret: {value}"]
// 切割成 [["random", "64"], ["bark", "Secret: {value}"]]
func parseByComma(args []string) [][]string {
	var steps [][]string
	var current []string

	for _, arg := range args {
		trimmed := strings.TrimSpace(arg)
		if trimmed == "," {
			if len(current) > 0 {
				steps = append(steps, current)
				current = []string{} // 重置，准备收集下一个 Action
			}
		} else if trimmed != "" {
			current = append(current, arg)
		}
	}
	if len(current) > 0 {
		steps = append(steps, current)
	}
	return steps
}

// resolveArgs 负责在执行前，将参数中形如 {key} 的占位符，动态替换为上游输出的值
func resolveArgs(args []string, input map[string]interface{}) []string {
	resolved := make([]string, len(args))

	for i, arg := range args {
		temp := arg
		// 遍历上游输出的每一个 key，动态替换占位符
		for k, v := range input {
			placeholder := "{" + k + "}"
			if strings.Contains(temp, placeholder) {
				strVal := fmt.Sprintf("%v", v)
				temp = strings.ReplaceAll(temp, placeholder, strVal)
			}
		}
		resolved[i] = temp
	}

	return resolved
}

// runPipeline 内存直接传递 Map，并且【合并】每一步的输出，防止数据在传递中丢失
func runPipeline(ctx context.Context, steps [][]string, initialInput map[string]interface{}) {
	// 初始化累积 Context，防止空指针，并将初始输入深拷贝进去
	currentInput := make(map[string]interface{})
	if initialInput != nil {
		for k, v := range initialInput {
			currentInput[k] = v
		}
	}

	for _, step := range steps {
		if len(step) == 0 {
			continue
		}
		actionName := step[0]
		actionArgs := step[1:]

		action, ok := GetAction(actionName)
		if !ok {
			WriteError(fmt.Errorf("action not found in pipeline: %s", actionName))
			os.Exit(1)
		}

		// 1. 用累积的 currentInput 动态渲染当前步骤的命令行参数
		resolvedArgs := resolveArgs(actionArgs, currentInput)

		// 2. 接力执行（此时的 currentInput 包含了前面所有步骤累积下来的所有 key-value 对）
		result, err := action.Execute(ctx, resolvedArgs, currentInput)
		if err != nil {
			WriteError(fmt.Errorf("pipeline step [%s] failed: %w", actionName, err))
			os.Exit(1)
		}

		// 3. 【核心改进】：增量合并，而不是直接覆盖！
		// 这样前面所有步骤产生的字段（比如 random 产生的 value）都会保留在 currentInput 中，一路传到底！
		if result != nil {
			for k, v := range result {
				currentInput[k] = v
			}
		}
	}

	// 管道执行完毕，打印最终累积的所有结果
	WriteOutput(currentInput)
}
