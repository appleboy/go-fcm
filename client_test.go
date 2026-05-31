package fcm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"firebase.google.com/go/v4/messaging"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

// MockTokenSource is a TokenSource implementation that can be used for testing.
type MockTokenSource struct {
	AccessToken string
}

// Token returns the test token associated with the TokenSource.
func (ts *MockTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: ts.AccessToken}, nil
}

func TestSendEach(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
      "name": "q1w2e3r4"
    }`))
	}))
	defer server.Close()
	t.Run("send each and send each dry run success", func(t *testing.T) {
		client, err := NewClient(
			context.Background(),
			WithEndpoint(server.URL),
			WithProjectID("test"),
			WithTokenSource(&MockTokenSource{AccessToken: "test-token"}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.Send(
			context.Background(),
			&messaging.Message{
				Token: "test",
				Data: map[string]string{
					"foo": "bar",
				},
			})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkSuccessfulBatchResponseForSendEach(t, resp)

		resp, err = client.SendDryRun(
			context.Background(),
			&messaging.Message{
				Token: "test",
				Data: map[string]string{
					"foo": "bar",
				},
			})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkSuccessfulBatchResponseForSendEach(t, resp)
	})

	t.Run("missing multicast message", func(t *testing.T) {
		client, err := NewClient(
			context.Background(),
			WithEndpoint(server.URL),
			WithProjectID("test"),
			WithTokenSource(&MockTokenSource{AccessToken: "12345-token"}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.SendMulticast(
			context.Background(),
			nil,
		)
		if err == nil {
			t.Fatalf("expected error\n got: %v", err)
		}
		if resp != nil {
			t.Fatalf("expected nil response\ngot: %v", resp)
		}
	})

	t.Run("send multicast and send multicast dry run success", func(t *testing.T) {
		client, err := NewClient(
			context.Background(),
			WithEndpoint(server.URL),
			WithProjectID("test"),
			WithTokenSource(&MockTokenSource{AccessToken: "12345-token"}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.SendMulticast(
			context.Background(),
			&messaging.MulticastMessage{
				Tokens: []string{"test01"},
				Data: map[string]string{
					"foo": "bar",
				},
			})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkSuccessfulBatchResponseForSendEach(t, resp)

		resp, err = client.SendMulticastDryRun(
			context.Background(),
			&messaging.MulticastMessage{
				Tokens: []string{"test01"},
				Data: map[string]string{
					"foo": "bar",
				},
			})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkSuccessfulBatchResponseForSendEach(t, resp)
	})

	t.Run("send message without token using custom client option", func(t *testing.T) {
		client, err := NewClient(
			context.Background(),
			WithEndpoint(server.URL),
			WithProjectID("test"),
			WithCustomClientOption(option.WithoutAuthentication()),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.Send(
			context.Background(),
			&messaging.Message{
				Token: "test",
				Data: map[string]string{
					"foo": "bar",
				},
			})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkSuccessfulBatchResponseForSendEach(t, resp)
	})
}

func checkSuccessfulBatchResponseForSendEach(t *testing.T, resp *messaging.BatchResponse) {
	if resp.SuccessCount != 1 {
		t.Fatalf("expected 1 successes\ngot: %d successes", resp.SuccessCount)
	}
	if resp.FailureCount != 0 {
		t.Fatalf("expected 0 failures\ngot: %d failures", resp.FailureCount)
	}
}

// TestTransportAuthCombos verifies that attaching a custom transport (debug
// logging or a custom http.Client) stays compatible with the non-JSON auth
// methods. Previously these combinations forced a service-account-JSON path and
// failed at construction with "unexpected end of JSON input".
func TestTransportAuthCombos(t *testing.T) {
	var (
		mu      sync.Mutex
		gotAuth string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotAuth = r.Header.Get("Authorization")
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name": "q1w2e3r4"}`))
	}))
	defer server.Close()

	msg := &messaging.Message{
		Token: "test",
		Data:  map[string]string{"foo": "bar"},
	}

	// sendAndCheck sends one message, asserts a successful batch response, and
	// returns the Authorization header the server observed so callers can verify
	// the auth transport actually attached (or omitted) the bearer token.
	sendAndCheck := func(t *testing.T, client *Client) string {
		t.Helper()
		mu.Lock()
		gotAuth = ""
		mu.Unlock()
		resp, err := client.Send(context.Background(), msg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkSuccessfulBatchResponseForSendEach(t, resp)
		mu.Lock()
		defer mu.Unlock()
		return gotAuth
	}

	t.Run("debug with token source", func(t *testing.T) {
		client, err := NewClient(
			context.Background(),
			WithEndpoint(server.URL),
			WithProjectID("test"),
			WithTokenSource(&MockTokenSource{AccessToken: "test-token"}),
			WithDebug(true),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if auth := sendAndCheck(t, client); auth != "Bearer test-token" {
			t.Fatalf("expected bearer token to be attached, got %q", auth)
		}
	})

	t.Run("debug without authentication", func(t *testing.T) {
		client, err := NewClient(
			context.Background(),
			WithEndpoint(server.URL),
			WithProjectID("test"),
			WithCustomClientOption(option.WithoutAuthentication()),
			WithDebug(true),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if auth := sendAndCheck(t, client); auth != "" {
			t.Fatalf("expected no Authorization header, got %q", auth)
		}
	})

	t.Run("custom http client with token source", func(t *testing.T) {
		client, err := NewClient(
			context.Background(),
			WithEndpoint(server.URL),
			WithProjectID("test"),
			WithTokenSource(&MockTokenSource{AccessToken: "test-token"}),
			WithHTTPClient(&http.Client{}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if auth := sendAndCheck(t, client); auth != "Bearer test-token" {
			t.Fatalf("expected bearer token to be attached, got %q", auth)
		}
	})
}
