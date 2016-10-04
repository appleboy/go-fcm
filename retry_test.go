package fcm

import (
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	t.Run("retry=succcess", func(t *testing.T) {
		var attempts int
		err := retry(func() error {
			attempts++
			if attempts < 3 {
				return connectionError("error")
			}
			return nil
		}, 4)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if attempts != 3 {
			t.Fatalf("expected 3 attempts\ngot: %d attempts", attempts)
		}
	})

	t.Run("retry=false", func(t *testing.T) {
		err := retry(func() error {
			return errors.New("error")
		}, 4)
		if err == nil {
			t.Fatalf("expected error: %v\ngot nil", err)
		}
	})

	t.Run("retry=maxAttempts", func(t *testing.T) {
		err := retry(func() error {
			return connectionError("error")
		}, 1)
		if err == nil {
			t.Fatalf("expected error: %v\ngot nil", err)
		}
	})
}
