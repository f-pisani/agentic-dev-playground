package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
)

// GetSubscriptions retrieves all subscriptions
func (c *Client) GetSubscriptions() ([]Subscription, error) {
	req, err := c.NewRequest(http.MethodGet, "/v2/subscriptions.json", nil)
	if err != nil {
		return nil, err
	}
	
	var subscriptions []Subscription
	_, err = c.Do(req, &subscriptions)
	if err != nil {
		return nil, err
	}
	
	return subscriptions, nil
}

// GetSubscription retrieves a specific subscription by ID
func (c *Client) GetSubscription(id int64) (*Subscription, error) {
	path := fmt.Sprintf("/v2/subscriptions/%d.json", id)
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	subscription := new(Subscription)
	_, err = c.Do(req, subscription)
	if err != nil {
		return nil, err
	}
	
	return subscription, nil
}

// CreateSubscription creates a new subscription
func (c *Client) CreateSubscription(feedURL string) (*Subscription, error) {
	subscriptionReq := &SubscriptionRequest{
		FeedURL: feedURL,
	}
	
	req, err := c.NewRequest(http.MethodPost, "/v2/subscriptions.json", subscriptionReq)
	if err != nil {
		return nil, err
	}
	
	subscription := new(Subscription)
	_, err = c.Do(req, subscription)
	if err != nil {
		return nil, err
	}
	
	return subscription, nil
}

// UpdateSubscription updates a subscription's title
func (c *Client) UpdateSubscription(id int64, title string) (*Subscription, error) {
	path := fmt.Sprintf("/v2/subscriptions/%d.json", id)
	updateReq := &SubscriptionUpdateRequest{
		Title: title,
	}
	
	req, err := c.NewRequest(http.MethodPatch, path, updateReq)
	if err != nil {
		return nil, err
	}
	
	subscription := new(Subscription)
	_, err = c.Do(req, subscription)
	if err != nil {
		return nil, err
	}
	
	return subscription, nil
}

// DeleteSubscription deletes a subscription
func (c *Client) DeleteSubscription(id int64) error {
	path := fmt.Sprintf("/v2/subscriptions/%d.json", id)
	req, err := c.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// GetFeedSubscriptions retrieves all subscriptions for a specific feed
func (c *Client) GetFeedSubscriptions(feedID int64) ([]Subscription, error) {
	path := fmt.Sprintf("/v2/feeds/%d/subscriptions.json", feedID)
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	var subscriptions []Subscription
	_, err = c.Do(req, &subscriptions)
	if err != nil {
		return nil, err
	}
	
	return subscriptions, nil
}

// GetSubscriptionsWithParams retrieves subscriptions with optional parameters
func (c *Client) GetSubscriptionsWithParams(params url.Values) ([]Subscription, error) {
	path := "/v2/subscriptions.json"
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}
	
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	var subscriptions []Subscription
	_, err = c.Do(req, &subscriptions)
	if err != nil {
		return nil, err
	}
	
	return subscriptions, nil
}

// GetSubscriptionCount returns the total number of subscriptions
func (c *Client) GetSubscriptionCount() (int, error) {
	req, err := c.NewRequest(http.MethodGet, "/v2/subscriptions.json", nil)
	if err != nil {
		return 0, err
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("API error: %s", resp.Status)
	}
	
	return GetTotalCount(resp), nil
}
