package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrServerError = errors.New("server error")
)

type Client struct {
	httpClient *http.Client
	host       string
}

func NewClient(httpClient *http.Client, host string) *Client {
	return &Client{
		httpClient: httpClient,
		host:       host,
	}
}

func (c *Client) Encrypt(ctx context.Context, text []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.host, bytes.NewBuffer(text))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"error respons from encryptor status code [%d] message [%s]: %w",
			resp.StatusCode, string(body), ErrServerError,
		)
	}

	return body, nil
}
