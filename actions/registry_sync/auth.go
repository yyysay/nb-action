package registry_sync

import (
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/authn"
)

type AuthMode string

const (
	AuthModeAnonymous AuthMode = "anonymous"

	AuthModeDocker AuthMode = "docker"

	AuthModeEnv AuthMode = "env"
)

// AuthConfig
//
// source 和 destination 分开认证
type AuthConfig struct {

	// source registry
	SrcMode AuthMode

	SrcKeychain authn.Keychain

	// destination registry
	DstMode AuthMode

	DstKeychain authn.Keychain
}

// envKeychain
//
// 使用环境变量认证
type envKeychain struct {
	auth authn.Authenticator
}

func (k envKeychain) Resolve(
	target authn.Resource,
) (authn.Authenticator, error) {

	return k.auth, nil
}

// PrepareAuth
//
// 认证优先级:
//
// source:
// 1. SRC_REGISTRY_USERNAME/PASSWORD
// 2. ~/.docker/config.json
// 3. anonymous
//
// destination:
// 1. DST_REGISTRY_USERNAME/PASSWORD
// 2. ~/.docker/config.json
// 3. anonymous
func PrepareAuth() (*AuthConfig, error) {

	config := &AuthConfig{}

	// ==========================
	// source auth
	// ==========================

	srcUsername := os.Getenv(
		"SRC_REGISTRY_USERNAME",
	)

	srcPassword := os.Getenv(
		"SRC_REGISTRY_PASSWORD",
	)

	if srcUsername != "" && srcPassword != "" {

		config.SrcMode = AuthModeEnv

		config.SrcKeychain = envKeychain{
			auth: authn.FromConfig(
				authn.AuthConfig{
					Username: srcUsername,
					Password: srcPassword,
				},
			),
		}

	} else if hasDockerConfig() {

		config.SrcMode = AuthModeDocker

		config.SrcKeychain = authn.DefaultKeychain

	} else {

		config.SrcMode = AuthModeAnonymous

		config.SrcKeychain = authn.DefaultKeychain
	}

	// ==========================
	// destination auth
	// ==========================

	dstUsername := os.Getenv(
		"DST_REGISTRY_USERNAME",
	)

	dstPassword := os.Getenv(
		"DST_REGISTRY_PASSWORD",
	)

	if dstUsername != "" && dstPassword != "" {

		config.DstMode = AuthModeEnv

		config.DstKeychain = envKeychain{
			auth: authn.FromConfig(
				authn.AuthConfig{
					Username: dstUsername,
					Password: dstPassword,
				},
			),
		}

	} else if hasDockerConfig() {

		config.DstMode = AuthModeDocker

		config.DstKeychain = authn.DefaultKeychain

	} else {

		config.DstMode = AuthModeAnonymous

		config.DstKeychain = authn.DefaultKeychain
	}

	return config, nil
}

// hasDockerConfig
//
// 判断本机 Docker config 是否存在
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
