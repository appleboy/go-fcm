package fcm

import "errors"

var (
	// ErrInvalidMessage occurs if push notitication message is nil.
	ErrInvalidMessage = errors.New("message is invalid")
	// ErrInvalidToken occurs if device token is empty.
	ErrInvalidToken = errors.New("device token is invalid")
	// ErrToManyRegIDs occurs when registration ids more then 1000.
	ErrToManyRegIDs = errors.New("too many registrations ids")
	// ErrInvalidTimeToLive occurs if TimeToLive more then 2419200.
	ErrInvalidTimeToLive = errors.New("messages time-to-live is invalid")
)

// Message represents list of targets, options, and payload for HTTP JSON messages.
type Message struct {
	Token                    string        `json:"to"`
	RegistrationIDs          []string      `json:"registration_ids,omitempty"`
	Condition                string        `json:"condition,omitempty"`
	CollapseKey              string        `json:"collapse_key,omitempty"`
	Priority                 string        `json:"priority,omitempty"`
	ContentAvailable         bool          `json:"content_available,omitempty"`
	DelayWhileIdle           bool          `json:"delay_while_idle,omitempty"`
	TimeToLive               int           `json:"time_to_live,omitempty"`
	DeliveryReceiptRequested bool          `json:"delivery_receipt_requested,omitempty"`
	DryRun                   bool          `json:"dry_run,omitempty"`
	Notification             *Notification `json:"notification,omitempty"`
	Data                     *Data         `json:"data,omitempty"`
}

// Notification specifies the predefined, user-visible key-value pairs
// of the notification payload
type Notification struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Sound        string `json:"sound,omitempty"`
	Badge        string `json:"badge,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Color        string `json:"color,omitempty"`
	ClickAction  string `json:"click_action,omitempty"`
	BodyLocKey   string `json:"body_loc_key,omitempty"`
	BodyLocArgs  string `json:"body_loc_args,omitempty"`
	TitleLocKey  string `json:"title_loc_key,omitempty"`
	TitleLocArgs string `json:"title_loc_args,omitempty"`
}

// Data specifies the custom key-value pairs of the message's payload.
type Data map[string]interface{}

// Validate returns an error if the message is not well-formed.
func (msg *Message) Validate() error {
	switch {
	case msg == nil:
		return ErrInvalidMessage
	case msg.Token == "":
		return ErrInvalidToken
	case len(msg.RegistrationIDs) > 1000:
		return ErrToManyRegIDs
	case msg.TimeToLive > 2419200:
		return ErrInvalidTimeToLive
	default:
		return nil
	}
}
