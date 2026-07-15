package result

import (
	"sync"
	"time"
)

type Result struct {
	Total   int
	Success int
	Failed  int
	Start   time.Time

	mu sync.Mutex
}

// New 创建统计对象
func New(total int) *Result {
	return &Result{
		Total: total,
		Start: time.Now(),
	}
}

// AddSuccess 成功+1
func (r *Result) AddSuccess() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Success++
}

// AddFailed 失败+1
func (r *Result) AddFailed() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Failed++
}

// Duration 执行耗时
func (r *Result) Duration() time.Duration {
	return time.Since(r.Start)
}
