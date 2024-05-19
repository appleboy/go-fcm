package fcm

import (
	"net"
	"time"
)

const (
	minBackoff = 100 * time.Millisecond
	maxBackoff = 1 * time.Minute
)

func retry(fn func() error, attempts int) error {
	var attempt int
	for {
		err := fn()
		if err == nil {
			return nil
		}

		if _, ok := err.(net.Error); !ok {
			return err
		}

		attempt++
		backoff := minBackoff * time.Duration(attempt*attempt)
		if attempt > attempts || backoff > maxBackoff {
			return err
		}

		time.Sleep(backoff)
	}
}
