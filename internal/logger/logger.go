package logger

import (
	"fmt"
	"sync"
)

var mu sync.Mutex

// Println 线程安全输出
func Println(args ...any) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println(args...)
}

// Printf 线程安全格式化输出
func Printf(format string, args ...any) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Printf(format, args...)
}
