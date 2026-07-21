package output

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

type Record struct {
	Image  string
	Target string
	Status string
	Error  string
}

func Start(
	total int,
	concurrency int,
) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintln(
		os.Stderr,
		"registry-sync",
	)

	fmt.Fprintln(
		os.Stderr,
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
	)

	fmt.Fprintf(
		os.Stderr,
		"total:       %d\n",
		total,
	)

	fmt.Fprintf(
		os.Stderr,
		"concurrency: %d\n",
		concurrency,
	)

	fmt.Fprintln(
		os.Stderr,
	)
}

func PrintItem(
	index int,
	total int,
	item Record,
) {
	mu.Lock()
	defer mu.Unlock()

	status := "✓"

	if item.Status == "failed" {
		status = "✗"
	}

	fmt.Fprintf(
		os.Stderr,
		"[%d/%d] %s %s\n",
		index,
		total,
		status,
		item.Image,
	)

	if item.Target != "" {

		fmt.Fprintf(
			os.Stderr,
			"      target: %s\n",
			item.Target,
		)
	}

	if item.Error != "" {

		fmt.Fprintf(
			os.Stderr,
			"      ! %s\n",
			shortError(item.Error),
		)
	}

	fmt.Fprintln(
		os.Stderr,
	)
}

func Finish(
	success int,
	failed int,
	duration time.Duration,
) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintln(
		os.Stderr,
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
	)

	fmt.Fprintf(
		os.Stderr,
		"result:   %d success, %d failed\n",
		success,
		failed,
	)

	fmt.Fprintf(
		os.Stderr,
		"time:     %s\n",
		duration,
	)
}

func shortError(
	err string,
) string {

	err = strings.TrimSpace(err)

	// MANIFEST_UNKNOWN
	if index := strings.Index(
		err,
		"MANIFEST_UNKNOWN:",
	); index >= 0 {

		err = err[index+len("MANIFEST_UNKNOWN:"):]
	}

	err = strings.TrimSpace(err)

	// 去掉 crane registry error 附带 map
	if index := strings.Index(
		err,
		"; map[",
	); index >= 0 {

		err = err[:index]
	}

	err = strings.TrimSpace(err)

	const limit = 80

	if len(err) > limit {
		return err[:limit] + "..."
	}

	return err
}
