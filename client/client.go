package client

import (
	"io/ioutil"
	"net/http"
)

type RequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	httpClient RequestDoer
}

func NewClient(client RequestDoer) *Client {
	return &Client{client}
}

func (c *Client) Get(url string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Spike_Crawler/0.1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
