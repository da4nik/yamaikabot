package telegram

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	token  string
	client *http.Client
}

const botURL = "https://api.telegram.org"

func NewClient(botToken string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		token:  botToken,
		client: &http.Client{Transport: tr},
	}
}

func (c *Client) do(method string, params interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.request(method, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected return code: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) request(method string, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest("POST", c.url(method), body)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}

	return req, nil
}

func (c *Client) url(method string) string {
	return fmt.Sprintf("%s/bot%s/%s",
		botURL,
		c.token,
		method,
	)
}
