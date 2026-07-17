package bark

import (
	"testing"
)

func TestEncrypt(t *testing.T) {

	// 设置 AES key
	// AES-256 需要 32 字节

	t.Setenv(
		"BARK_AES_KEY",
		"12345678901234567890123456789012",
	)

	// 设置 AES GCM nonce(iv)
	// GCM 推荐 12 字节

	t.Setenv(
		"BARK_AES_IV",
		"123456789012",
	)

	// 设置 Bark Server

	t.Setenv(
		"BARK_SERVER",
		"https://api.day.app",
	)

	// 设置 Bark DeviceKey

	t.Setenv(
		"BARK_DEVICE_KEY",
		"test-device-key",
	)

	// 创建 Bark 消息

	payload := BarkPayload{

		Title: "test",

		Body: "hello bark",
	}

	// 调用完整加密流程

	result, err := encrypt(
		payload,
	)

	// 检查错误

	if err != nil {

		t.Fatal(err)

	}

	// 检查结果

	if result == "" {

		t.Fatal(
			"encrypt result is empty",
		)

	}

	// 输出结果观察

	t.Log(
		"encrypted result:",
	)

	t.Log(
		result,
	)

}
