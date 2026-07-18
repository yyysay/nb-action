package bark

import (
	"fmt"
	"os"
)

type Config struct {

	// AES 配置

	Key []byte

	IV []byte

	// Bark 推送配置

	Server string

	DeviceKey string
}

func LoadConfig() (
	Config,
	error,
) {

	keyString := os.Getenv(
		"BARK_AES_KEY",
	)

	ivString := os.Getenv(
		"BARK_AES_IV",
	)

	server := os.Getenv(
		"BARK_SERVER",
	)

	deviceKey := os.Getenv(
		"BARK_DEVICE_KEY",
	)

	if len(keyString) != 32 {

		return Config{},
			fmt.Errorf(
				"BARK_AES_KEY 长度错误，需要32字节",
			)

	}

	if len(ivString) != 12 {

		return Config{},
			fmt.Errorf(
				"BARK_AES_IV 长度错误，需要12字节",
			)

	}

	if server == "" {

		return Config{},
			fmt.Errorf(
				"BARK_SERVER 未配置",
			)

	}

	if deviceKey == "" {

		return Config{},
			fmt.Errorf(
				"BARK_DEVICE_KEY 未配置",
			)

	}

	config := Config{

		Key: []byte(keyString),

		IV: []byte(ivString),

		Server: server,

		DeviceKey: deviceKey,
	}

	return config, nil
}
