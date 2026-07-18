package syncer

import (
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Copy 镜像同步
func Copy(src string, dst string) error {
	return CopyWithPlatform(src, dst, "")
}

// CopyWithPlatform 指定架构同步
func CopyWithPlatform(src string, dst string, targetPlatform string) error {
	if targetPlatform == "" {
		return crane.Copy(src, dst)
	}

	p := parsePlatform(targetPlatform)

	return crane.Copy(
		src,
		dst,
		crane.WithPlatform(p),
	)
}

// parsePlatform
//
// 格式:
// linux/amd64
// linux/arm64/v8
func parsePlatform(value string) *v1.Platform {
	parts := strings.Split(value, "/")

	p := &v1.Platform{}

	if len(parts) > 0 {
		p.OS = parts[0]
	}

	if len(parts) > 1 {
		p.Architecture = parts[1]
	}

	if len(parts) > 2 {
		p.Variant = parts[2]
	}

	return p
}
