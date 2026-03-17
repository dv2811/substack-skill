package koyfin

import (
	"errors"
	"io"
	"net/http"
)

func (c *Client) getResponse(method, url string, requestBody io.Reader, headers map[string][]string) (io.ReadCloser, error) {
	// prepare request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		req.Header = headers
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// HTTP error handling
	if resp.StatusCode > 399 {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		return nil, errors.Join(KoyfinAPIErr, errors.New(resp.Status))
	}

	return resp.Body, nil
}
