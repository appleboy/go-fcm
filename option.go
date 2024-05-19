package fcm

import (
	"net/http"
	"net/url"

	"google.golang.org/api/option"
)

// Option configurates Client with defined option.
type Option func(*Client) error

// WithHTTPClient returns Option to configure HTTP Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		c.options = append(c.options, option.WithHTTPClient(httpClient))
		return nil
	}
}

// WithHTTPProxy returns Option to configure HTTP Client proxy.
func WithHTTPProxy(proxyURL string) Option {
	return func(c *Client) error {
		httpClient := http.DefaultClient
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return err
		}
		httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
		c.options = append(c.options, option.WithHTTPClient(httpClient))
		return nil
	}
}

// WithCredentialsFile returns a ClientOption that authenticates
// API calls with the given service account or refresh token JSON
// credentials file.
func WithCredentialsFile(filename string) Option {
	return func(c *Client) error {
		c.options = append(c.options, option.WithCredentialsFile(filename))
		return nil
	}
}

// WithCredentialsJSON returns a ClientOption that authenticates
// API calls with the given service account or refresh token JSON
// credentials.
func WithCredentialsJSON(json []byte) Option {
	return func(c *Client) error {
		c.options = append(c.options, option.WithCredentialsJSON(json))
		return nil
	}
}

// WithCustomClientOption is an option function that allows you to provide custom client options.
// It appends the provided custom options to the client's options list.
// The custom options are applied when sending requests to the FCM server.
// If no custom options are provided, this function does nothing.
//
// Parameters:
//   - opts: The custom client options to be appended to the client's options list.
//
// Returns:
//   - An error if there was an issue appending the custom options to the client's options list, or nil otherwise.
func WithCustomClientOption(opts ...option.ClientOption) Option {
	return func(c *Client) error {
		if len(opts) == 0 {
			return nil
		}
		c.options = append(c.options, opts...)
		return nil
	}
}
