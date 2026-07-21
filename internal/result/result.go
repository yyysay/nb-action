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

	Items []Item

	mu sync.Mutex
}

// New 创建统计对象
func New(total int) *Result {

	return &Result{
		Total: total,
		Start: time.Now(),

		Items: make(
			[]Item,
			0,
			total,
		),
	}
}

// AddSuccess 成功+1
func (r *Result) AddSuccess(
	item Item,
) {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.Success++

	item.Status = "success"

	r.Items = append(
		r.Items,
		item,
	)
}

// AddFailed 失败+1
func (r *Result) AddFailed(
	item Item,
) {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.Failed++

	item.Status = "failed"

	r.Items = append(
		r.Items,
		item,
	)
}

// Duration 执行耗时
func (r *Result) Duration() time.Duration {
	return time.Since(r.Start)
}
