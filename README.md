# fcm
[![GoDoc](https://godoc.org/github.com/edganiukov/fcm?status.svg)](https://godoc.org/github.com/edganiukov/fcm)

Golang client library for Firebase Cloud Messaging. Implemented only [HTTP client](https://firebase.google.com/docs/cloud-messaging/http-server-ref#downstream).

More information on [Firebase Cloud Messaging](https://firebase.google.com/docs/cloud-messaging/)

### Getting Started
-------------------
To install fcm, use `go get`:

```bash
go get github.com/edganiukov/fcm
```
or `govendor`:

```bash
govendor fetch github.com/edganiukov/fcm
```

### Sample Usage
----------------
Here is a simple example illustrating how to use FCM library:
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/edganiukov/fcm"
)

func main() {
	// Create the message to be sent.
	msg := &gcm.Message{
     		Token: "sample_device_token",
      		Data: &fcm.Data{
         		"foo": "bar",
      		}
  	}

	// Create a Sender to send the message.
	client := fcm.NewHTTPClient("sample_api_key")

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	if err != nil {
		/* ... */
    	}
	/* ... */
}
```
