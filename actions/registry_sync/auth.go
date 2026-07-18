package registry_sync

import (
	"fmt"
	"os"
	"path/filepath"
)

// PrepareAuth
//
// 认证准备流程:
//
// 1. 检查本地 Docker 配置
// 2. 没有配置时检查环境变量
// 3. 后续可以在这里扩展 docker login
func PrepareAuth() error {

	if hasDockerConfig() {
		fmt.Println("docker config found")
		return nil
	}

	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")

	if username == "" || password == "" {
		return fmt.Errorf(
			"registry credentials not found",
		)
	}

	fmt.Println("registry credentials found")

	return nil
}

// hasDockerConfig
//
// 判断本机是否存在 Docker 登录配置
//
// 默认路径:
// ~/.docker/config.json
func hasDockerConfig() bool {

	home, err := os.UserHomeDir()

	if err != nil {
		return false
	}

	path := filepath.Join(
		home,
		".docker",
		"config.json",
	)

	_, err = os.Stat(path)

	return err == nil
}
