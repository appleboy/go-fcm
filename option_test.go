package fcm

import (
	"net/http"
	"net/url"
	"testing"
)

func TestWithHTTPProxy(t *testing.T) {
	proxyURL := "http://example.com/proxy"

	c := &Client{
		client: &http.Client{
			Transport: &http.Transport{},
		},
	}

	err := WithHTTPProxy(proxyURL)(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	transport := c.client.Transport.(*http.Transport)
	req := &http.Request{
		URL: &url.URL{
			Scheme: "https",
			Host:   "fcm.googleapis.com",
		},
	}
	proxy, _ := transport.Proxy(req)

	if proxy.String() != proxyURL {
		t.Fatalf("expected proxy URL: %s\ngot: %s", proxyURL, proxy.String())
	}
}
