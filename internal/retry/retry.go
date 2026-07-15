package retry

import (
	"time"
)

// Do 执行带重试的任务
//
// attempts: 最大尝试次数
// delay: 初始等待时间
func Do(task func() error, attempts int, delay time.Duration) error {
	var err error

	for i := 0; i < attempts; i++ {
		err = task()

		if err == nil {
			return nil
		}

		// 最后一次失败，不等待
		if i == attempts-1 {
			break
		}

		time.Sleep(delay)

		// 简单指数退避
		delay *= 2
	}

	return err
}
