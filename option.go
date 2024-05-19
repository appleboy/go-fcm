package fcm

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"google.golang.org/api/option"
)

// Option configurates Client with defined option.
type Option func(*Client) error

// WithHTTPClient returns Option to configure HTTP Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout returns Option to configure HTTP Client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d.Nanoseconds() <= 0 {
			return errors.New("invalid timeout duration")
		}
		c.timeout = d
		return nil
	}
}

// WithHTTPProxy returns Option to configure HTTP Client proxy.
func WithHTTPProxy(proxyURL string) Option {
	return func(c *Client) error {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return err
		}
		c.httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
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
