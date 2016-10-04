package fcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	// FCMEndpoint contains endpoint URL of FCM service.
	FCMEndpoint = "https://fcm.googleapis.com/fcm/send"
)

var (
	// ErrInvalidAPIKey occurs if API key is not set.
	ErrInvalidAPIKey = errors.New("client API Key is invalid")
)

// HTTPClient abstracts the interaction between the application server and the
// FCM server via HTTP protocol. The developer must obtain an API key from the
// Google APIs Console page and pass it to the `HTTPClient` so that it can
// perform authorized requests on the application server's behalf.
// To send a message to one or more devices use the Client's Send.
//
// If the `HTTP` field is nil, a zeroed http.Client will be allocated and used
// to send messages.
type HTTPClient struct {
	apiKey   string
	client   *http.Client
	endpoint string
}

// NewClient creates new HTTP FCM Client based on API key and
// with default endpoint and http client.
func NewClient(apiKey string) (*HTTPClient, error) {
	if apiKey == "" {
		return nil, ErrInvalidAPIKey
	}

	return &HTTPClient{
		apiKey:   apiKey,
		endpoint: FCMEndpoint,
		client:   &http.Client{},
	}, nil
}

// NewClientWithHTTP creates new HTTP FCM Client based on API key, endpoint and http client.
func NewClientWithHTTP(httpClient *http.Client, apiKey, endpoint string) (*HTTPClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if endpoint == "" {
		endpoint = FCMEndpoint
	}
	if apiKey == "" {
		return nil, ErrInvalidAPIKey
	}

	return &HTTPClient{
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   httpClient,
	}, nil
}

// Send sends a message to the FCM server without retrying in case of
// service unavailability. A non-nil error is returned if a non-recoverable
// error occurs (i.e. if the response status is not "200 OK").
func (c *HTTPClient) Send(msg *Message) (*Response, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return c.send(data)
}

// SendWithRetry sends a message to the FCM server with defined number of retrying
// in case of temporary error.
func (c *HTTPClient) SendWithRetry(msg *Message, retryAttempts int) (*Response, error) {
	resp := new(Response)
	err := retry(func() error {
		var err error
		resp, err = c.Send(msg)
		return err
	}, retryAttempts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HTTPClient) send(data []byte) (*Response, error) {
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", c.apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, connectionError(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, serverError(fmt.Sprintf("%d error: %s", resp.StatusCode, resp.Status))
		}
		return nil, fmt.Errorf("%d error: %s", resp.StatusCode, resp.Status)
	}

	response := new(Response)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
}
