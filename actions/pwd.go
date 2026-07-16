package actions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"

	"golang.org/x/crypto/curve25519"
)

type Password struct{}

func NewPassword() *Password {
	return &Password{}
}

func (p *Password) Name() string {
	return "pwd"
}

func (p *Password) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (map[string]interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("pwd requires a subcommand: wg-keypair or rand")
	}

	subcommand := args[0]

	switch subcommand {
	case "wg-keypair":
		return p.generateWGKeypair()
	case "rand":
		if len(args) < 2 {
			return nil, fmt.Errorf("rand requires size argument")
		}
		return p.generateRand(args[1:])
	default:
		return nil, fmt.Errorf("unknown subcommand: %s (supported: wg-keypair, rand)", subcommand)
	}
}

// generateWGKeypair 生成 WireGuard 密钥对
func (p *Password) generateWGKeypair() (map[string]interface{}, error) {
	var privateKeyBytes [32]byte
	_, err := rand.Read(privateKeyBytes[:])
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// WireGuard 私钥处理（按照 RFC 7748 - Curve25519 标准）
	privateKeyBytes[0] &= 248
	privateKeyBytes[31] = (privateKeyBytes[31] & 127) | 64

	// 计算公钥
	publicKeyBytes, err := curve25519.X25519(privateKeyBytes[:], curve25519.Basepoint[:])
	if err != nil {
		return nil, fmt.Errorf("failed to compute public key: %w", err)
	}

	privateKeyHex := hex.EncodeToString(privateKeyBytes[:])
	publicKeyHex := hex.EncodeToString(publicKeyBytes)

	return map[string]interface{}{
		"private_key": privateKeyHex,
		"public_key":  publicKeyHex,
	}, nil
}

// generateRand 生成随机密钥
func (p *Password) generateRand(args []string) (map[string]interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("rand requires size argument")
	}

	// 解析大小参数
	size, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid size: %s", args[0])
	}

	if size <= 0 {
		return nil, fmt.Errorf("size must be positive")
	}

	// 检查是否有 --base64 flag
	useBase64 := false
	for i := 1; i < len(args); i++ {
		if args[i] == "--base64" {
			useBase64 = true
			break
		}
	}

	// 生成随机数据
	buf := make([]byte, size)
	_, err = rand.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random data: %w", err)
	}

	var value string
	if useBase64 {
		// Base64 编码
		value = base64.StdEncoding.EncodeToString(buf)
	} else {
		// 十六进制编码
		value = hex.EncodeToString(buf)
	}

	return map[string]interface{}{
		"value": value,
	}, nil
}
