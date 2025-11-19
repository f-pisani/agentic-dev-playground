package feedbin

import (
	"context"
	"fmt"
)

// GetSubscriptions fetches all subscriptions.
func (c *Client) GetSubscriptions(ctx context.Context, opt *PaginationOptions) ([]Subscription, *Response, error) {
	u, err := addOptions("subscriptions.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var subs []Subscription
	resp, err := c.Do(ctx, req, &subs)
	if err != nil {
		return nil, resp, err
	}

	return subs, resp, nil
}

// GetSubscription fetches a single subscription by ID.
func (c *Client) GetSubscription(ctx context.Context, id int64) (*Subscription, *Response, error) {
	u := fmt.Sprintf("subscriptions/%d.json", id)
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var sub Subscription
	resp, err := c.Do(ctx, req, &sub)
	if err != nil {
		return nil, resp, err
	}

	return &sub, resp, nil
}

type subscriptionCreateRequest struct {
	FeedURL string `json:"feed_url"`
}

// CreateSubscription creates a new subscription.
func (c *Client) CreateSubscription(ctx context.Context, feedURL string) (*Subscription, *Response, error) {
	reqBody := &subscriptionCreateRequest{FeedURL: feedURL}
	req, err := c.NewRequest("POST", "subscriptions.json", reqBody)
	if err != nil {
		return nil, nil, err
	}

	var sub Subscription
	resp, err := c.Do(ctx, req, &sub)
	if err != nil {
		return nil, resp, err
	}

	return &sub, resp, nil
}

type subscriptionUpdateRequest struct {
	Title string `json:"title"`
}

// UpdateSubscription updates a subscription's title.
func (c *Client) UpdateSubscription(ctx context.Context, id int64, title string) (*Subscription, *Response, error) {
	u := fmt.Sprintf("subscriptions/%d.json", id)
	reqBody := &subscriptionUpdateRequest{Title: title}
	req, err := c.NewRequest("PATCH", u, reqBody)
	if err != nil {
		return nil, nil, err
	}

	var sub Subscription
	resp, err := c.Do(ctx, req, &sub)
	if err != nil {
		return nil, resp, err
	}

	return &sub, resp, nil
}

// DeleteSubscription deletes a subscription.
func (c *Client) DeleteSubscription(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("subscriptions/%d.json", id)
	req, err := c.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, req, nil)
}
