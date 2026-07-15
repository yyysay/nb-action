package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	// ⚠️ 注意：如果你的 go.mod 模块名不是 github.com/yangtudou/nb-action
	// 请将 "github.com/yangtudou/nb-action" 替换为实际的模块名
	"github.com/yangtudou/nb-action/internal/logger"
	"github.com/yangtudou/nb-action/internal/parser"
	"github.com/yangtudou/nb-action/internal/resolver"
	"github.com/yangtudou/nb-action/internal/result"
	"github.com/yangtudou/nb-action/internal/retry"
	"github.com/yangtudou/nb-action/internal/syncer"
	"github.com/yangtudou/nb-action/internal/worker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: registry-sync sync [--src <image> | --base <file>] --dst-prefix <registry>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sync":
		syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)

		src := syncCmd.String("src", "", "Source image")
		base := syncCmd.String("base", "images.txt", "Base image manifest file")
		srcPrefix := syncCmd.String("src-prefix", "", "Source registry prefix")
		srcFlatten := syncCmd.Bool("src-flatten", false, "Flatten source image path")
		dstPrefix := syncCmd.String("dst-prefix", "", "Destination registry prefix")
		dstFlatten := syncCmd.Bool("dst-flatten", false, "Flatten destination image path")
		platform := syncCmd.String("platform", "", "Target platform (e.g., linux/amd64)")
		concurrency := syncCmd.Int("concurrency", 4, "Concurrent workers")
		retries := syncCmd.Int("retries", 3, "Retry attempts")
		dryRun := syncCmd.Bool("dry-run", false, "Only show actions")

		_ = syncCmd.Parse(os.Args[2:])

		// 1. 参数防御性检查
		if *dstPrefix == "" {
			syncCmd.Usage()
			os.Exit(1)
		}
		if *concurrency < 1 {
			*concurrency = 1
		}

		// 2. 加载镜像列表
		var images []string
		if *src != "" {
			images = []string{*src}
		} else {
			var err error
			images, err = parser.ReadImageList(*base)
			if err != nil {
				fmt.Printf("Error: Base manifest '%s' not found: %v\n", *base, err)
				os.Exit(1)
			}
		}

		if len(images) == 0 {
			fmt.Println("No images found to sync.")
			return
		}

		total := len(images)
		stats := result.New(total)
		var tasks []worker.Task

		for i, img := range images {
			index := i + 1
			image := img // 局部变量确保闭包安全

			tasks = append(tasks, func() error {
				source := resolver.Resolve(image, resolver.Rule{Prefix: *srcPrefix, Flatten: *srcFlatten})
				target := resolver.Resolve(image, resolver.Rule{Prefix: *dstPrefix, Flatten: *dstFlatten})

				if *dryRun {
					logger.Printf("[%d/%d] Dry Run: %s -> %s\n", index, total, source, target)
					stats.AddSuccess()
					return nil
				}

				err := retry.Do(func() error {
					return syncer.CopyWithPlatform(source, target, *platform)
				}, *retries, 2*time.Second)

				if err != nil {
					stats.AddFailed()
					logger.Printf("[%d/%d] Failed: %s (%v)\n", index, total, image, err)
					return err
				}

				stats.AddSuccess()
				logger.Printf("[%d/%d] Success: %s\n", index, total, image)
				return nil
			})
		}

		// 3. 执行任务
		_ = worker.Run(tasks, *concurrency)

		// 4. 打印报告
		mode := "Sync"
		if *dryRun {
			mode = "Dry Run"
		}
		logger.Printf("\n%s Summary\nTotal: %d | Success: %d | Failed: %d | Time: %s\n",
			mode, stats.Total, stats.Success, stats.Failed, stats.Duration())

		if stats.Failed > 0 {
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
