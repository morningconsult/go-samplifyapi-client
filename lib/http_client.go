package samplify

import (
	"net/http"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

type defaultHTTPClient struct {
	client *retryablehttp.Client
}

func newDefaultHTTPClient() defaultHTTPClient {
	return defaultHTTPClient{client: retryablehttp.NewClient()}
}

// Do makes an HTTP request.
func (d defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rreq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}

	return d.client.Do(rreq)
}

func (d defaultHTTPClient) SetTimeout(timeout time.Duration) {
	d.client.HTTPClient.Timeout = timeout
}
