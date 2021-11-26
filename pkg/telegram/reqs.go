package telegram

import (
	"encoding/json"
	"fmt"
)

type SetWebhookRequest struct {
	URL                string   `json:"url"`
	AllowedUpdates     []string `json:"allowed_updates"` // message, edited_channel_post, callback_query
	DropPendingUpdates bool     `json:"drop_pending_updates"`
}

type SetWebhookResponse struct {
}

func (c *Client) SetWebhook(params SetWebhookRequest) error {
	respBody, err := c.do("setWebhook", params)
	if err != nil {
		return err
	}

	// TODO: Remove debug
	fmt.Printf("[SetWebhook] response = %s", respBody)

	return fmt.Errorf("not implemented")
}

func (c *Client) DeleteWebhook(dropPendingUpdates bool) error {
	respBody, err := c.do("deleteWebhook", map[string]bool{
		"drop_pending_updates": dropPendingUpdates,
	})
	if err != nil {
		return err
	}

	// TODO: Remove debug
	fmt.Printf("[DeleteWebhook] response = %s", respBody)

	return fmt.Errorf("not implmented")
}

func (c *Client) GetMe() (*User, error) {
	respBody, err := c.do("getMe", nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

type GetUpdatesRequest struct {
	Offset         int      `json:"offset,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	Timeout        int      `json:"timeout,omitempty"`         // Seconds, defaults to 0, short pulling
	AllowedUpdates []string `json:"allowed_updates,omitempty"` // message, edited_channel_post, callback_query
}

type GetUpdatesResponse struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

func (c *Client) GetUpdates(params GetUpdatesRequest) (*GetUpdatesResponse, error) {
	respBody, err := c.do("getUpdates", params)
	if err != nil {
		return nil, err
	}

	var response GetUpdatesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type SendMessageRequest struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *Client) SendMessage(params SendMessageRequest) (*Message, error) {
	respBody, err := c.do("sendMessage", params)
	if err != nil {
		return nil, err
	}

	var response Message
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
