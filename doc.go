// Package fcm provides Firebase Cloud Messaging functionality for Golang
//
// Here is a simple example illustrating how to use FCM library:
//
//	func main() {
//		// Create the message to be sent.
//		msg := &fcm.Message{
//			To: "sample_device_token",
//			Data: map[string]interface{}{
//				"foo": "bar",
//			},
//		}
//
//		// Create a FCM client to send the message.
//		client, err := fcm.NewClient("sample_api_key")
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		// Send the message and receive the response without retries.
//		response, err := client.Send(msg)
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		log.Printf("%#v\n", response)
//	}
package fcm
