package registry_sync

const helpText = `
registry-sync

Docker 镜像仓库同步工具

Usage:
  nb-action registry-sync [options]

Required:
  --dst-prefix <registry>
        目标镜像仓库

Options:
  --src-prefix <registry>
        源镜像仓库

  --base <file>
        镜像列表文件

  --src-flatten
        源仓库路径扁平化

  --dst-flatten
        目标仓库路径扁平化

  --platform <platform>
        目标平台

  --concurrency <number>
        并发数量

  --retries <number>
        重试次数

  --dry-run
        仅显示同步计划

Examples:
  nb-action registry-sync \
    --src-prefix docker.io \
    --dst-prefix registry.example.com \
    --base images.txt
`
