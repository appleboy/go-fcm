package fcm

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type debugTransport struct {
	t http.RoundTripper
}

func (d debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	log.Printf("%s", string(reqDump)) //nolint:gosec // debug output from local HTTP dump, not user-controlled

	resp, err := d.t.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
		return nil, err
	}
	log.Printf("%s", string(respDump)) //nolint:gosec // debug output from local HTTP dump, not user-controlled
	return resp, nil
}
