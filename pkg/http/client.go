package http

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type httpClient struct {
	client *http.Client
}

func NewHTTPClient() pkg.HTTPClient {
	return &httpClient{
		client: &http.Client{
			Timeout: time.Second * 40,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (h *httpClient) SendRequest(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "io.ReadAll")
	}

	return data, nil
}
