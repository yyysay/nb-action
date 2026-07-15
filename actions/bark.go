package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Bark struct {
	Server    string
	DeviceKey string
}

func NewBark(server string, deviceKey string) *Bark {
	return &Bark{
		Server:    server,
		DeviceKey: deviceKey,
	}
}

func (b *Bark) Name() string {
	return "bark"
}

// parseKeyValue 精准识别 key=value 格式，彻底排除 Base64 等等号填充的干扰
func parseKeyValue(arg string) (string, string, bool) {
	idx := strings.Index(arg, "=")
	if idx <= 0 || idx == len(arg)-1 {
		return "", "", false
	}

	k := strings.TrimSpace(arg[:idx])
	v := strings.TrimSpace(arg[idx+1:])

	// 1. 验证 Key：必须是合法的标识符（只允许英文、数字、下划线、中划线）
	if len(k) == 0 {
		return "", "", false
	}
	for _, r := range k {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			return "", "", false
		}
	}

	// 2. 验证 Value：不能为空，且不能纯由等号 "=" 组成（防止 Base64 尾部填充如 "Hw==" 被误判）
	if len(v) == 0 {
		return "", "", false
	}
	isAllEquals := true
	for _, r := range v {
		if r != '=' {
			isAllEquals = false
			break
		}
	}
	if isAllEquals {
		return "", "", false
	}

	return k, v, true
}

func (b *Bark) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (
	map[string]interface{},
	error,
) {
	payload := map[string]interface{}{
		"device_key": b.DeviceKey,
	}

	// 1. 优先载入上游 input 传入的参数，并自动做字段名归一化
	keys := []string{
		"title", "body", "subtitle", "subTitle", "group", "sound", "icon", "url",
		"ciphertext", "level", "copy", "badge", "isArchive", "isarchive",
		"automaticallyCopy", "automatically_copy", "category",
	}
	for _, key := range keys {
		if value, ok := input[key]; ok {
			officialKey := key
			switch key {
			case "subTitle":
				officialKey = "subtitle"
			case "isarchive":
				officialKey = "isArchive"
			case "automatically_copy":
				officialKey = "automaticallyCopy"
			}
			payload[officialKey] = value
		}
	}

	// 2. 解析命令行 args。支持 positional (title/body) 和 key=value 混排
	var positionalArgs []string
	for _, arg := range args {
		if k, v, ok := parseKeyValue(arg); ok {
			switch k {
			case "subtitle", "subTitle":
				payload["subtitle"] = v
			case "group":
				payload["group"] = v
			case "sound":
				payload["sound"] = v
			case "icon":
				payload["icon"] = v
			case "url":
				payload["url"] = v
			case "level":
				payload["level"] = v
			case "copy":
				payload["copy"] = v
			case "ciphertext":
				payload["ciphertext"] = v
			case "category":
				payload["category"] = v
			case "badge":
				if intVal, err := strconv.Atoi(v); err == nil {
					payload["badge"] = intVal
				} else {
					payload["badge"] = v
				}
			case "isArchive", "archive", "isarchive":
				if intVal, err := strconv.Atoi(v); err == nil {
					payload["isArchive"] = intVal
				} else if v == "true" || v == "1" {
					payload["isArchive"] = 1
				} else if v == "false" || v == "0" {
					payload["isArchive"] = 0
				} else {
					payload["isArchive"] = v
				}
			case "automaticallyCopy", "auto_copy", "autocopy":
				if intVal, err := strconv.Atoi(v); err == nil {
					payload["automaticallyCopy"] = intVal
				} else if v == "true" || v == "1" {
					payload["automaticallyCopy"] = 1
				} else if v == "false" || v == "0" {
					payload["automaticallyCopy"] = 0
				} else {
					payload["automaticallyCopy"] = v
				}
			default:
				// 无法直接识别的自定义 key 直接透传
				payload[k] = v
			}
		} else {
			positionalArgs = append(positionalArgs, arg)
		}
	}

	// 3. 将没有等号的普通参数映射为 title 和 body
	if len(positionalArgs) > 0 {
		if len(positionalArgs) >= 2 {
			payload["title"] = positionalArgs[0]
			payload["body"] = positionalArgs[1]
		} else {
			payload["body"] = positionalArgs[0]
		}
	}

	// 4. 兼容老逻辑：如果还是没有 body 字段，但 input 里包含 value，则使用 value 补齐
	if _, ok := payload["body"]; !ok {
		if value, ok := input["value"]; ok {
			payload["body"] = value
		}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		b.Server+"/push",
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bark request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bark push failed: %s", resp.Status)
	}

	return map[string]interface{}{
		"sent": true,
	}, nil
}
