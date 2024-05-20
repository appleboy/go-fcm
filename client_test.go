package fcm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/messaging"
	"golang.org/x/oauth2"
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
}

func checkSuccessfulBatchResponseForSendEach(t *testing.T, resp *messaging.BatchResponse) {
	if resp.SuccessCount != 1 {
		t.Fatalf("expected 1 successes\ngot: %d sucesses", resp.SuccessCount)
	}
	if resp.FailureCount != 0 {
		t.Fatalf("expected 0 failures\ngot: %d failures", resp.FailureCount)
	}
}
