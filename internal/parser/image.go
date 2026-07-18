package parser

import "strings"

// ImageName 从完整镜像地址中提取仓库路径
//
// 示例:
// docker.io/library/alpine
// -> library/alpine
//
// ghcr.io/sagernet/sing-box
// -> sagernet/sing-box
func ImageName(image string) string {
	image = strings.TrimSpace(image)

	// 去掉协议（如果存在）
	image = strings.TrimPrefix(image, "docker://")
	image = strings.TrimPrefix(image, "daemon://")

	parts := strings.Split(image, "/")

	if len(parts) == 1 {
		// alpine
		return "library/" + parts[0]
	}

	// 判断第一段是不是 registry
	// docker.io/library/alpine
	// ghcr.io/user/image
	if strings.Contains(parts[0], ".") ||
		strings.Contains(parts[0], ":") ||
		parts[0] == "localhost" {

		return strings.Join(parts[1:], "/")
	}

	// user/image
	return image
}
