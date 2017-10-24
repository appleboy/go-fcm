package fcm

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	t.Run("unmarshal=success", func(t *testing.T) {
		data := []byte(`{
  "multicast_id":10,
  "success": 0,
  "failure":1,
  "canonical_ids":10,
  "results": [{
    "message_id":"q1w2e3r4",
    "registration_id": "t5y6u7i8o9",
    "error": "NotRegistered"
  }]
}`)

		var response Response
		err := json.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := Response{
			MulticastID:  10,
			Success:      0,
			Failure:      1,
			CanonicalIDs: 10,
			Results: []Result{{
				MessageID:      "q1w2e3r4",
				RegistrationID: "t5y6u7i8o9",
				Error:          ErrNotRegistered,
			}},
		}
		if !reflect.DeepEqual(response, expected) {
			t.Fatalf("expected: %v\ngot: %v", expected, response)
		}

		if !response.Results[0].Unregistered() {
			t.Fatalf("expected: true\ngot: %t", response.Results[0].Unregistered())
		}
	})

	t.Run("Device Group HTTP Response", func(t *testing.T) {
		data := []byte(`{
  "success":1,
  "failure":2,
  "failed_registration_ids":[
    "regId1",
    "regId2"
  ]
}`)

		var response Response
		err := json.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := Response{
			Success: 1,
			Failure: 2,
			FailedRegistrationIDs: []string{
				"regId1",
				"regId2",
			},
		}
		if !reflect.DeepEqual(response, expected) {
			t.Fatalf("expected: %v\ngot: %v", expected, response)
		}
	})

	t.Run("Topic HTTP response", func(t *testing.T) {
		data := []byte(`{
  "message_id": 1023456,
  "error": "TopicsMessageRateExceeded"
}`)

		var response Response
		err := json.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := Response{
			MessageID: int64(1023456),
			Error:     errMap["TopicsMessageRateExceeded"],
		}
		if !reflect.DeepEqual(response, expected) {
			t.Fatalf("expected: %v\ngot: %v", expected, response)
		}
	})

	t.Run("unmarshal=success_timeout", func(t *testing.T) {
		data := []byte(`{
  "multicast_id":10,
  "success": 0,
  "failure":1,
  "canonical_ids":10,
  "results": [{
    "message_id":"q1w2e3r4",
    "registration_id": "t5y6u7i8o9",
    "error": "Unavailable"
  }]
}`)

		var response Response
		err := json.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if response.Results[0].Unregistered() {
			t.Fatalf("expected: false\ngot: %t", response.Results[0].Unregistered())
		}
	})

	t.Run("unmarshal=failed", func(t *testing.T) {
		data := []byte(`{
  "multicast_id":10,
  "success": 0,
  "failure":1,
  "canonical_ids":10,
  "results": [{
    "message_id":["q1w2e3r4"],
    "registration_id": "t5y6u7i8o9",
    "error": "NotRegistered"
  }]
}`)

		var response Response
		err := json.Unmarshal(data, &response)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
