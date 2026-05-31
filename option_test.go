package fcm

import (
	"errors"
	"io/fs"
	"path/filepath"
	"testing"

	"google.golang.org/api/option"
)

func TestWithHTTPProxyInvalidURL(t *testing.T) {
	c := &Client{}
	if err := WithHTTPProxy("http://\x7f")(c); err == nil {
		t.Fatal("expected error for invalid proxy URL, got nil")
	}
}

func TestWithCredentialsFileMissingWrapsError(t *testing.T) {
	c := &Client{}
	// Use a path under the test's temp dir so non-existence is guaranteed
	// regardless of the working directory.
	missing := filepath.Join(t.TempDir(), "does-not-exist.json")
	err := WithCredentialsFile(missing)(c)
	if err == nil {
		t.Fatal("expected error for missing credentials file, got nil")
	}
	// The %w wrapping must preserve the underlying os error so callers can
	// classify it with errors.Is.
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected error wrapping fs.ErrNotExist, got: %v", err)
	}
}

func TestWithCustomClientOptionEmptyIsNoop(t *testing.T) {
	c := &Client{}
	if err := WithCustomClientOption()(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.options) != 0 {
		t.Fatalf("expected no options appended, got %d", len(c.options))
	}
}

func TestWithCustomClientOptionAppends(t *testing.T) {
	c := &Client{}
	if err := WithCustomClientOption(option.WithoutAuthentication())(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.options) != 1 {
		t.Fatalf("expected 1 option appended, got %d", len(c.options))
	}
}
