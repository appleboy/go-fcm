# go-fcm

[![GoDoc](https://pkg.go.dev/badge/github.com/appleboy/go-fcm)](https://pkg.go.dev/github.com/appleboy/go-fcm)
[![Lint and Testing](https://github.com/appleboy/go-fcm/actions/workflows/testing.yml/badge.svg?branch=master)](https://github.com/appleboy/go-fcm/actions/workflows/testing.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/go-fcm)](https://goreportcard.com/report/github.com/appleboy/go-fcmm)

This project was forked from [github.com/edganiukov/fcm](https://github.com/edganiukov/fcm).

More information on [Firebase Cloud Messaging](https://firebase.google.com/docs/cloud-messaging/)

## Feature

* [x] Send messages to a single device
* [x] Send messages to a multiple devices
* [x] Send messages to a topic
* [x] Supports condition attribute

## Getting Started

To install fcm, use `go get`:

```bash
go get github.com/appleboy/go-fcm
```

## Provide credentials using ADC

Google Application Default Credentials (ADC) for Firebase projects support Google service accounts, which you can use to call Firebase server APIs from your app server or trusted environment. If you're developing code locally or deploying your application on-premises, you can use credentials obtained via this service account to authorize server requests.

To authenticate a service account and authorize it to access Firebase services, you must generate a private key file in JSON format.

**To generate a private key file for your service account:**

1. In the Firebase console, open **Settings > [Service Accounts][11]**.
2. Click **Generate New Private Key**, then confirm by clicking **Generate Key**.
3. Securely store the JSON file containing the key.

When authorizing via a service account, you have two choices for providing the credentials to your application. You can either set the **GOOGLE_APPLICATION_CREDENTIALS** environment variable, or you can explicitly pass the path to the service account key in code. The first option is more secure and is strongly recommended.

See the more detail information [here][12].

[11]: https://console.firebase.google.com/project/_/settings/serviceaccounts/adminsdk
[12]: https://firebase.google.com/docs/cloud-messaging/auth-server#provide-credentials-using-adc

## Usage

Here is a simple example illustrating how to use FCM library:

```go
package main

import (
  "context"
  "fmt"
  "log"

  "firebase.google.com/go/v4/messaging"
  fcm "github.com/appleboy/go-fcm"
)

func main() {
  ctx := context.Background()
  client, err := fcm.NewClient(
    ctx,
    fcm.WithCredentialsFile("path/to/serviceAccountKey.json"),
    // initial with service account
    // fcm.WithServiceAccount("my-client-id@my-project-id.iam.gserviceaccount.com"),
  )
  if err != nil {
    log.Fatal(err)
  }

  // Send to a single device
  token := "test"
  resp, err := client.Send(
    ctx,
    &messaging.Message{
      Token: token,
      Data: map[string]string{
        "foo": "bar",
      },
    },
  )
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("success count:", resp.SuccessCount)
  fmt.Println("failure count:", resp.FailureCount)
  fmt.Println("message id:", resp.Responses[0].MessageID)
  fmt.Println("error msg:", resp.Responses[0].Error)

  // Send to topic
  resp, err = client.Send(
    ctx,
    &messaging.Message{
      Data: map[string]string{
        "foo": "bar",
      },
      Topic: "highScores",
    },
  )
  if err != nil {
    log.Fatal(err)
  }

  // Send with condition
  // Define a condition which will send to devices which are subscribed
  // to either the Google stock or the tech industry topics.
  condition := "'stock-GOOG' in topics || 'industry-tech' in topics"

  // See documentation on defining a message payload.
  message := &messaging.Message{
    Data: map[string]string{
      "score": "850",
      "time":  "2:45",
    },
    Condition: condition,
  }

  resp, err = client.Send(
    ctx,
    message,
  )
  if err != nil {
    log.Fatal(err)
  }

  // Send multiple messages to device
  // Create a list containing up to 500 messages.
  registrationToken := "YOUR_REGISTRATION_TOKEN"
  messages := []*messaging.Message{
    {
      Notification: &messaging.Notification{
        Title: "Price drop",
        Body:  "5% off all electronics",
      },
      Token: registrationToken,
    },
    {
      Notification: &messaging.Notification{
        Title: "Price drop",
        Body:  "2% off all books",
      },
      Topic: "readers-club",
    },
  }
  resp, err = client.Send(
    ctx,
    messages...,
  )
  if err != nil {
    log.Fatal(err)
  }

  // Send multicast message
  // Create a list containing up to 500 registration tokens.
  // This registration tokens come from the client FCM SDKs.
  registrationTokens := []string{
    "YOUR_REGISTRATION_TOKEN_1",
    // ...
    "YOUR_REGISTRATION_TOKEN_n",
  }
  msg := &messaging.MulticastMessage{
    Data: map[string]string{
      "score": "850",
      "time":  "2:45",
    },
    Tokens: registrationTokens,
  }
  resp, err = client.SendMulticast(
    ctx,
    msg,
  )
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%d messages were sent successfully\n", resp.SuccessCount)
  if resp.FailureCount > 0 {
    var failedTokens []string
    for idx, resp := range resp.Responses {
      if !resp.Success {
        // The order of responses corresponds to the order of the registration tokens.
        failedTokens = append(failedTokens, registrationTokens[idx])
      }
    }

    fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
  }
}
```

### Custom HTTP Client

You can use a custom HTTP client by passing it to the `NewClient` function:

```go
func main() {
  httpTimeout := time.Duration(5) * time.Second
  tlsTimeout := time.Duration(5) * time.Second

  transport := &http2.Transport{
    DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
      return tls.DialWithDialer(&net.Dialer{Timeout: tlsTimeout}, network, addr, cfg)
    },
  }

  httpClient := &http.Client{
    Transport: transport,
    Timeout:   httpTimeout,
  }

  ctx := context.Background()
  client, err := fcm.NewClient(
    ctx,
    fcm.WithCredentialsFile("path/to/serviceAccountKey.json"),
    fcm.WithHTTPClient(httpClient),
  )
}
```

### Custom Proxy Server

You can use a custom proxy server by passing it to the `NewClient` function:

```go
func main() {
  ctx := context.Background()
  client, err := fcm.NewClient(
    ctx,
    fcm.WithCredentialsFile("path/to/serviceAccountKey.json"),
    fcm.WithHTTPProxy("http://localhost:8088"),
  )
}
```

### Mock Client for Testing

You can use a mock client for Unit Testing by initializing your `NewClient` with the following:

* `fcm.WithEndpoint` with a local `httptest.NewServer` URL
* `fcm.WithProjectID` any string
* `fcm.WithServiceAccount` any string
* `fcm.WithCustomClientOption` with option `option.WithoutAuthentication()`

```go
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
```
