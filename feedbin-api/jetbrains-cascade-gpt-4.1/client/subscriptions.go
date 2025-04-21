package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// GetSubscriptions fetches all subscriptions.
func (c *Client) GetSubscriptions() ([]models.Subscription, error) {
	req, err := c.newRequest(http.MethodGet, "/subscriptions.json", nil)
	if err != nil {
		return nil, err
	}
	var subs models.SubscriptionsResponse
	err = c.do(req, &subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

// GetSubscription fetches a subscription by ID.
func (c *Client) GetSubscription(id int64) (*models.Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%d.json", id)
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var sub models.Subscription
	err = c.do(req, &sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// UpdateSubscription updates a subscription (e.g., custom title).
func (c *Client) UpdateSubscription(id int64, payload map[string]interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/subscriptions/%d.json", id)
	req, err := c.newRequest(http.MethodPatch, path, strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
