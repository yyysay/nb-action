package logger

import (
	"fmt"
	"os"
	"sync"
)

var mu sync.Mutex

// Println 输出调试日志到 stderr
func Println(
	args ...any,
) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintln(
		os.Stderr,
		args...,
	)
}

// Printf 输出格式日志到 stderr
func Printf(
	format string,
	args ...any,
) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintf(
		os.Stderr,
		format,
		args...,
	)
}
