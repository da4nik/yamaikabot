package telegram

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	token  string
	client *http.Client
}

const botURL = "https://api.telegram.org/bot"

func NewClient(botToken string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		token:  botToken,
		client: &http.Client{Transport: tr},
	}
}

func (c *Client) do(method string, params map[string]interface{}) ([]byte, error) {

	json, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := c.request(method)
	res, err := client.Do(req)
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

func (c *Client) request(method string) *http.Request {
	req := &http.NewRequest("POST", c.url(method), nil)
	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}
	return req
}

func (c *Client) url(method string) string {
	return fmt.Sprintf("%s/bot%s/%s",
		botURL,
		c.token,
		method,
	)
}
