package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/messaging"
	fcm "github.com/appleboy/go-fcm"
	"google.golang.org/api/option"
)

func TestMockClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
      "name": "q1w2e3r4"
    }`))
	}))
	defer server.Close()

	client, err := fcm.NewClient(
		context.Background(),
		fcm.WithEndpoint(server.URL),
		fcm.WithProjectID("test"),
		fcm.WithCustomClientOption(option.WithoutAuthentication()),
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
}

func checkSuccessfulBatchResponseForSendEach(t *testing.T, resp *messaging.BatchResponse) {
	if resp.SuccessCount != 1 {
		t.Fatalf("expected 1 successes\ngot: %d sucesses", resp.SuccessCount)
	}
	if resp.FailureCount != 0 {
		t.Fatalf("expected 0 failures\ngot: %d failures", resp.FailureCount)
	}
}
