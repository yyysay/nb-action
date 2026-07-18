package registry_sync

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yangtudou/nb-action/internal/logger"
	"github.com/yangtudou/nb-action/internal/parser"
	"github.com/yangtudou/nb-action/internal/resolver"
	"github.com/yangtudou/nb-action/internal/result"
	"github.com/yangtudou/nb-action/internal/retry"
	"github.com/yangtudou/nb-action/internal/syncer"
	"github.com/yangtudou/nb-action/internal/worker"
)

func (r *RegistrySync) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (map[string]interface{}, error) {

	opt, err := ParseOptions(args)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid registry-sync arguments: %w",
			err,
		)
	}

	if opt.Src == "" {

		if v, ok := input["value"].(string); ok {
			opt.Src = v
		}

		if v, ok := input["image"].(string); ok {
			opt.Src = v
		}
	}

	if opt.DstPrefix == "" {
		return nil, fmt.Errorf(
			"missing required flag: --dst-prefix",
		)
	}

	// 增加私有仓库登入认证
	err = PrepareAuth()

	if err != nil {
		return nil, fmt.Errorf(
			"registry auth failed: %w",
			err,
		)
	}

	if opt.Concurrency < 1 {
		opt.Concurrency = 1
	}

	var images []string

	if opt.Src != "" {
		images = []string{opt.Src}
	} else {

		images, err = parser.ReadImageList(opt.Base)

		if err != nil {
			return nil, err
		}
	}

	stats := result.New(len(images))

	var mu sync.Mutex
	targets := make([]string, 0, len(images))

	tasks := make([]worker.Task, 0, len(images))

	for i, img := range images {

		index := i + 1
		image := img

		tasks = append(tasks, func() error {

			source := resolver.Resolve(
				image,
				resolver.Rule{
					Prefix:  opt.SrcPrefix,
					Flatten: opt.SrcFlatten,
				},
			)

			target := resolver.Resolve(
				image,
				resolver.Rule{
					Prefix:  opt.DstPrefix,
					Flatten: opt.DstFlatten,
				},
			)

			if opt.DryRun {

				logger.Printf(
					"[%d/%d] %s -> %s",
					index,
					len(images),
					source,
					target,
				)

			} else {

				err := retry.Do(
					func() error {
						return syncer.CopyWithPlatform(
							source,
							target,
							opt.Platform,
						)
					},
					opt.Retries,
					2*time.Second,
				)

				if err != nil {
					stats.AddFailed()
					return err
				}

			}

			stats.AddSuccess()

			mu.Lock()
			targets = append(targets, target)
			mu.Unlock()

			return nil
		})
	}

	worker.Run(
		tasks,
		opt.Concurrency,
	)

	return map[string]interface{}{
		"status":      "ok",
		"total":       stats.Total,
		"success":     stats.Success,
		"failed":      stats.Failed,
		"duration_ms": stats.Duration().Milliseconds(),
		"value": map[string]interface{}{
			"images": targets,
		},
	}, nil
}
