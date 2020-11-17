package http

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	writer     io.Writer
	url        string
}

func NewClient(writer io.Writer, url string) *Client {
	client := Client{
		httpClient: newHttpClient(),
		writer:     writer,
		url:        url,
	}
	return &client
}

func newHttpClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &http.Client{
		Timeout:   time.Minute * 3,
		Transport: transport,
	}
}

func (c *Client) Read() error {
	requestURL := c.url
	req, err := http.NewRequest(http.MethodGet, requestURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()

	err = c.request(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) request(req *http.Request) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	_, err = io.Copy(c.writer, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Url() string {
	return c.url
}
