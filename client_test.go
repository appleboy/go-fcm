package fcm

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	t.Run("send=success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"q1w2e3r4",
					"registration_id": "t5y6u7i8o9",
					"error": ""
				}]
			}`)
		}))
		defer server.Close()

		client, err := NewClient("test", WithEndpoint(server.URL), WithTimeout(10*time.Second))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.Send(&Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Success != 1 {
			t.Fatalf("expected 1 successes\ngot: %d sucesses", resp.Success)
		}
		if resp.Failure != 0 {
			t.Fatalf("expected 0 failures\ngot: %d failures", resp.Failure)
		}
	})

	t.Run("send=success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"q1w2e3r4",
					"registration_id": "t5y6u7i8o9",
					"error": ""
				}]
			}`)
		}))
		defer server.Close()

		client, err := NewClient("test", WithEndpoint(server.URL), WithTimeout(10*time.Second))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.Send(&Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Success != 1 {
			t.Fatalf("expected 1 successes\ngot: %d sucesses", resp.Success)
		}
		if resp.Failure != 0 {
			t.Fatalf("expected 0 failures\ngot: %d failures", resp.Failure)
		}
	})

	t.Run("send=failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			rw.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client, err := NewClient("test", WithEndpoint(server.URL))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.Send(&Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if resp != nil {
			t.Fatalf("expected nil response\ngot: %v response", resp)
		}
	})

	t.Run("send=invalid_token", func(t *testing.T) {
		_, err := NewClient("test", WithEndpoint(""))
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("send=invalid_message", func(t *testing.T) {
		c, err := NewClient("test", WithEndpoint("test"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = c.Send(&Message{})
		if err == nil {
			t.Fatal("expected error but go nil")
		}
	})

	t.Run("send=invalid-response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": {
					"message_id":"q1w2e3r4",
					"registration_id": "t5y6u7i8o9",
					"error": ""
				}
			}`)
		}))
		defer server.Close()

		client, err := NewClient("test",
			WithEndpoint(server.URL),
			WithHTTPClient(&http.Client{}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.Send(&Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		})
		if err == nil {
			t.Fatal("expected error but go nil")
		}

		if resp != nil {
			t.Fatalf("expected nil\ngot response: %v", resp)
		}
	})
}

func TestSendWithRetry(t *testing.T) {
	t.Run("send_with_retry=success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"q1w2e3r4",
					"registration_id": "t5y6u7i8o9",
					"error": ""
				}]
			}`)
		}))
		defer server.Close()

		client, err := NewClient("test",
			WithEndpoint(server.URL),
			WithHTTPClient(&http.Client{}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.SendWithRetry(&Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		}, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Success != 1 {
			t.Fatalf("expected 1 successes\ngot: %d successes", resp.Success)
		}
		if resp.Failure != 0 {
			t.Fatalf("expected 0 failures\ngot: %d failures", resp.Failure)
		}
	})

	t.Run("send_with_retry=failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			rw.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client, err := NewClient("test",
			WithEndpoint(server.URL),
			WithHTTPClient(&http.Client{}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.SendWithRetry(&Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		}, 2)

		if err == nil {
			t.Fatal("expected error\ngot nil")
		}
		if resp != nil {
			t.Fatalf("expected nil response\ngot: %v response", resp)
		}
	})

	t.Run("send_with_retry=success_retry", func(t *testing.T) {
		var attempts int
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			attempts++
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			if attempts < 3 {
				rw.WriteHeader(http.StatusInternalServerError)
			} else {
				rw.WriteHeader(http.StatusOK)
			}
			rw.Header().Set("Content-Type", "application/json")

			fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"q1w2e3r4",
					"registration_id": "t5y6u7i8o9",
					"error": ""
				}]
			}`)
		}))
		defer server.Close()

		client, err := NewClient("test",
			WithEndpoint(server.URL),
			WithHTTPClient(&http.Client{}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.SendWithRetry(&Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		}, 4)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if attempts != 3 {
			t.Fatalf("expected 3 attempts\ngot: %d attempts", attempts)
		}
		if resp.Success != 1 {
			t.Fatalf("expected 1 successes\ngot: %d successes", resp.Success)
		}
		if resp.Failure != 0 {
			t.Fatalf("expected 0 failures\ngot: %d failures", resp.Failure)
		}
	})

	t.Run("send_with_retry=failure_retry", func(t *testing.T) {
		client, err := NewClient("test",
			WithEndpoint("127.0.0.1:80"),
			WithHTTPClient(&http.Client{

				Timeout: time.Nanosecond,
			}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.SendWithRetry(&Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		}, 3)

		if err == nil {
			t.Fatal("expected error\ngot nil")
		}
		if resp != nil {
			t.Fatalf("expected nil response\ngot: %v response", resp)
		}
	})
}

func TestSendWithContext(t *testing.T) {
	t.Run("send_context=success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"q1w2e3r4",
					"registration_id": "t5y6u7i8o9",
					"error": ""
				}]
			}`)
		}))
		defer server.Close()

		client, err := NewClient("test", WithEndpoint(server.URL), WithTimeout(10*time.Second))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ctx := context.Background()
		resp, err := client.SendWithContext(ctx, &Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Success != 1 {
			t.Fatalf("expected 1 successes\ngot: %d sucesses", resp.Success)
		}
		if resp.Failure != 0 {
			t.Fatalf("expected 0 failures\ngot: %d failures", resp.Failure)
		}
	})

	t.Run("send_context=timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "key=test" {
				t.Fatalf("expected: key=test\ngot: %s", req.Header.Get("Authorization"))
			}
			time.Sleep(time.Millisecond * 100)
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
				"success": 1,
				"failure":0,
				"results": [{
					"message_id":"q1w2e3r4",
					"registration_id": "t5y6u7i8o9",
					"error": ""
				}]
			}`)
		}))
		defer server.Close()

		client, err := NewClient("test", WithEndpoint(server.URL), WithTimeout(10*time.Second))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		defer cancel()
		_, err = client.SendWithContext(ctx, &Message{
			To: "test",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		})
		if err == nil {
			t.Fatalf("no context timeout")
		}

		_, ok := err.(connectionError)
		if !ok {
			t.Fatalf("error is not fcm.connectionError \ngot: %T", err)
		}
	})
}
