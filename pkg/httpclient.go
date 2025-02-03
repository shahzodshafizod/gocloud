//go:generate mockgen -source=httpclient.go -package=mocks -destination=mocks/httpclient.go
package pkg

import "io"

type HTTPClient interface {
	// Sends an HTTP request using the specified method (e.g., GET, POST), URL, and optional data payload (`io.Reader`).
	// It returns the response body as a byte array or an error if the request fails.
	SendRequest(method string, url string, data io.Reader) ([]byte, error)
}
