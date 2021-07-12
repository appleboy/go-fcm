package main

import (
	"log"

	"github.com/appleboy/go-fcm"
)

func main() {
	// Create the message to be sent.
	msg := &fcm.Message{
		To: "e-udVMWkc0fshHO52uosbh:APA91bFn2WPHd_xZM69QxGL21On6wDZIQKMIXKo5ChBYw-3aiaMnzxQkgTLWXVIJg1U0vuGuNPSyaEJS6rlNckJKrOK5crtmqcLzU_vHOHiTzzIbIgNfpCXH2uR-qJNTICREZSug6YWH",
		Data: map[string]interface{}{
			"foo": "bar",
		},
		Notification: &fcm.Notification{
			Title: "title",
			Body:  "body",
		},
	}

	// Create a FCM client to send the message.
	client, err := fcm.NewClient("AAAAr3cjMuo:APA91bGk-e4Oe2FB43yjySy9KoMZk6qo-cY7G-uyrwB6XoV14cMgAEo7jkJyz4t33oovHu4RWyr-HYq4z1voN8N17QdHTnLWPrHxEhYAT1gR4RmDUIMpzilRJ6FDYDNbRW5b-Xtwy2nt")
	if err != nil {
		log.Fatalln(err)
	}

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", response)
}
