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
