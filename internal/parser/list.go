package parser

import (
	"bufio"
	"os"
	"strings"
)

// ReadImageList 从文件读取镜像列表
//
// 支持:
// - 空行
// - # 注释
func ReadImageList(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var images []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		images = append(images, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return images, nil
}
