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
