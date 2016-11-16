package fcm

import "testing"

func TestValidate(t *testing.T) {
	t.Run("validate=to_many_reg_ids", func(t *testing.T) {
		msg := &Message{
			Token:           "test",
			RegistrationIDs: []string{"reg_id"},
			TimeToLive:      3600,
		}
		err := msg.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("validate=invalid_message", func(t *testing.T) {
		var msg *Message
		err := msg.Validate()
		if err == nil {
			t.Fatal("expected 'message is invalid' error but got nil")
		}
	})

	t.Run("validate=invalid_token", func(t *testing.T) {
		msg := &Message{}
		err := msg.Validate()
		if err == nil {
			t.Fatal("expected 'device token is invalid' error but got nil")
		}
	})

	t.Run("validaate=to_many_reg_ids", func(t *testing.T) {
		msg := &Message{
			Token:           "test",
			RegistrationIDs: make([]string, 2000),
		}
		err := msg.Validate()
		if err == nil {
			t.Fatal("expected 'too many registrations id' error but got nil")
		}
	})

	t.Run("validaate=to_many_reg_ids", func(t *testing.T) {
		msg := &Message{
			Token:           "test",
			RegistrationIDs: []string{"reg_id"},
			TimeToLive:      2500000,
		}
		err := msg.Validate()
		if err == nil {
			t.Fatal("expected 'message time-to-live is invali' error but got nil")
		}
	})
}
