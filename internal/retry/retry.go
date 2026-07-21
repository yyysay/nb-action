package retry

import (
	"strings"
	"time"
)

func Do(
	image string,
	task func() error,
	attempts int,
	delay time.Duration,
) error {

	var err error

	if attempts < 1 {
		attempts = 1
	}

	for i := 0; i < attempts; i++ {

		err = task()

		if err == nil {
			return nil
		}

		if !retryable(err) {
			return err
		}

		if i == attempts-1 {
			break
		}

		time.Sleep(delay)

		delay *= 2
	}

	return err
}

func retryable(
	err error,
) bool {

	message := err.Error()

	nonRetryErrors := []string{
		"MANIFEST_UNKNOWN",
		"NAME_UNKNOWN",
		"denied",
		"unauthorized",
	}

	for _, item := range nonRetryErrors {

		if strings.Contains(
			message,
			item,
		) {
			return false
		}
	}

	return true
}
