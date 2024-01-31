package fcm

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

// Option configurates Client with defined option.
type Option func(*Client) error

// WithEndpoint returns Option to configure FCM Endpoint.
func WithEndpoint(endpoint string) Option {
	return func(c *Client) error {
		if endpoint == "" {
			return errors.New("invalid endpoint")
		}
		c.endpoint = endpoint
		return nil
	}
}

// WithHTTPClient returns Option to configure HTTP Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		c.client = httpClient
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
		c.client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
		return nil
	}
}
