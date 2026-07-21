package registry_sync

import (
	"context"
	"sync"
	"time"

	"github.com/yangtudou/nb-action/internal/logger"
	"github.com/yangtudou/nb-action/internal/output"
	"github.com/yangtudou/nb-action/internal/parser"
	"github.com/yangtudou/nb-action/internal/resolver"
	"github.com/yangtudou/nb-action/internal/result"
	"github.com/yangtudou/nb-action/internal/retry"
	"github.com/yangtudou/nb-action/internal/syncer"
	"github.com/yangtudou/nb-action/internal/worker"
)

func Run(
	ctx context.Context,
	opt *Options,
) (map[string]interface{}, error) {

	authConfig, err := PrepareAuth()

	if err != nil {
		return nil, err
	}

	if opt.Concurrency < 1 {
		opt.Concurrency = 1
	}

	var images []string

	if opt.Src != "" {

		images = []string{
			opt.Src,
		}

	} else {

		images, err = parser.ReadImageList(
			opt.Base,
		)

		if err != nil {
			return nil, err
		}
	}

	output.Start(
		len(images),
		opt.Concurrency,
	)

	stats := result.New(
		len(images),
	)

	var mu sync.Mutex

	targets := make(
		[]string,
		0,
		len(images),
	)

	mappings := make(
		[]map[string]string,
		0,
		len(images),
	)

	records := make(
		[]output.Record,
		len(images),
	)

	tasks := make(
		[]worker.Task,
		0,
		len(images),
	)

	for index, img := range images {

		image := img
		itemIndex := index

		tasks = append(
			tasks,
			func() error {

				select {

				case <-ctx.Done():

					return ctx.Err()

				default:

				}

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

				item := result.Item{
					Image:  image,
					Source: source,
					Target: target,
				}

				record := output.Record{
					Image:  image,
					Target: target,
				}

				if opt.DryRun {

					item.Status = "success"

					record.Status = "success"

					stats.AddSuccess(
						item,
					)

					mu.Lock()

					targets = append(
						targets,
						target,
					)

					mappings = append(
						mappings,
						map[string]string{
							"source": source,
							"target": target,
						},
					)

					records[itemIndex] = record

					mu.Unlock()

					return nil
				}

				err := retry.Do(
					image,
					func() error {

						return syncer.CopyWithPlatform(
							source,
							target,
							opt.Platform,
							authConfig.SrcKeychain,
							authConfig.DstKeychain,
						)
					},
					opt.Retries,
					2*time.Second,
				)

				if err != nil {

					item.Status = "failed"
					item.Error = err.Error()

					record.Status = "failed"
					record.Error = err.Error()

					stats.AddFailed(
						item,
					)

					mu.Lock()

					records[itemIndex] = record

					mu.Unlock()

					return nil
				}

				item.Status = "success"

				record.Status = "success"

				stats.AddSuccess(
					item,
				)

				mu.Lock()

				targets = append(
					targets,
					target,
				)

				records[itemIndex] = record

				mu.Unlock()

				return nil
			},
		)
	}

	err = worker.Run(
		tasks,
		opt.Concurrency,
	)

	if err != nil {
		return nil, err
	}

	for index := range records {

		output.PrintItem(
			index+1,
			len(records),
			records[index],
		)
	}

	output.Finish(
		stats.Success,
		stats.Failed,
		stats.Duration(),
	)

	logger.Printf(
		"registry-sync finished success=%d failed=%d\n",
		stats.Success,
		stats.Failed,
	)

	return map[string]interface{}{

		"status": "ok",

		"auth_mode": map[string]interface{}{
			"source":      authConfig.SrcMode,
			"destination": authConfig.DstMode,
		},

		"dry_run": opt.DryRun,

		"total": stats.Total,

		"success": stats.Success,

		"failed": stats.Failed,

		"duration_ms": stats.Duration().Milliseconds(),

		"value": map[string]interface{}{
			"images":  targets,
			"mapping": mappings,
			"items":   stats.Items,
		},
	}, nil
}
