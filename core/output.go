package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// maskString 负责对敏感字符串进行局部打码，保留首尾各 4 位用于比对
func maskString(s string) string {
	runes := []rune(s)
	n := len(runes)

	if n == 0 {
		return ""
	}

	// 如果太短（小于等于 8 个字符），只保留第一位
	if n <= 8 {
		return string(runes[0:1]) + "..."
	}

	// 保留前 4 位和后 4 位，中间隐藏
	return string(runes[:4]) + "..." + string(runes[n-4:])
}

// maskSensitive 递归处理敏感字段
func maskSensitive(data interface{}) interface{} {

	// Map 类型递归处理
	if m, ok := data.(map[string]interface{}); ok {

		masked := make(map[string]interface{})

		// 真正需要保护的字段
		sensitiveKeys := []string{
			"password",
			"secret",
			"token",
			"device_key",
			"device_token",
			"access_key",
			"private_key",
		}

		for k, v := range m {

			isSensitive := false

			lowerKey := strings.ToLower(k)

			for _, s := range sensitiveKeys {

				if strings.Contains(lowerKey, s) {
					isSensitive = true
					break
				}
			}

			if isSensitive {

				strVal := fmt.Sprintf("%v", v)

				masked[k] = maskString(strVal)

			} else {

				masked[k] = maskSensitive(v)

			}
		}

		return masked
	}

	// Slice 类型递归处理
	if s, ok := data.([]interface{}); ok {

		maskedSlice := make([]interface{}, len(s))

		for i, val := range s {

			maskedSlice[i] = maskSensitive(val)

		}

		return maskedSlice
	}

	// 基础类型直接返回
	return data
}

func WriteOutput(data interface{}) {

	// 输出前统一脱敏
	safeData := maskSensitive(data)

	output, err := json.Marshal(safeData)

	if err != nil {

		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	fmt.Println(string(output))
}

func WriteError(err error) {

	output := map[string]interface{}{
		"error": err.Error(),
	}

	data, marshalErr := json.Marshal(output)

	if marshalErr != nil {

		fmt.Fprintln(os.Stderr, marshalErr)

		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, string(data))
}
