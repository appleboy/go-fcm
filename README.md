# go-fcm

[![GoDoc](https://godoc.org/github.com/appleboy/go-fcm?status.svg)](https://godoc.org/github.com/appleboy/go-fcm)
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

## Sample Usage

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
  client, err := fcm.NewClient(
    context.Background(),
    fcm.WithCredentialsFile("path/to/serviceAccountKey.json"),
  )
  if err != nil {
    log.Fatal(err)
  }

  // Send to a single device
  token := "test"
  resp, err := client.Send(
    context.Background(),
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
}
```
