package feedbin

import (
	"context"
	"fmt"
)

// GetSubscriptions retrieves all subscriptions for the authenticated user
func (c *Client) GetSubscriptions(ctx context.Context, opts *SubscriptionOptions) ([]Subscription, *PaginationInfo, error) {
	path := "/subscriptions.json"
	if opts != nil {
		params := buildQueryParams(opts)
		path = addQueryParams(path, params)
	}

	var subscriptions []Subscription
	resp, err := c.makeRequest(ctx, "GET", path, nil, &subscriptions)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	pagination := extractPaginationInfo(resp)
	return subscriptions, pagination, nil
}

// GetSubscription retrieves a specific subscription by ID
func (c *Client) GetSubscription(ctx context.Context, id int) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%d.json", id)

	var subscription Subscription
	_, err := c.makeRequest(ctx, "GET", path, nil, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// CreateSubscription creates a new subscription to the specified feed URL
func (c *Client) CreateSubscription(ctx context.Context, feedURL string) (*Subscription, []FeedChoice, error) {
	if feedURL == "" {
		return nil, nil, &ValidationError{
			Field:   "feed_url",
			Message: "feed URL is required",
		}
	}

	req := &CreateSubscriptionRequest{
		FeedURL: feedURL,
	}

	var subscription Subscription
	resp, err := c.makeRequest(ctx, "POST", "/subscriptions.json", req, &subscription)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		// Check if this is a multiple choices error (300)
		if apiErr, ok := err.(*APIError); ok && apiErr.IsMultipleChoices() {
			var choices []FeedChoice
			if decodeErr := c.decodeErrorBody(apiErr.Body, &choices); decodeErr == nil {
				return nil, choices, err
			}
		}
		return nil, nil, err
	}

	return &subscription, nil, nil
}

// UpdateSubscription updates a subscription's title
func (c *Client) UpdateSubscription(ctx context.Context, id int, title string) (*Subscription, error) {
	if title == "" {
		return nil, &ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}

	path := fmt.Sprintf("/subscriptions/%d.json", id)
	req := &UpdateSubscriptionRequest{
		Title: title,
	}

	var subscription Subscription
	_, err := c.makeRequest(ctx, "PATCH", path, req, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// UpdateSubscriptionPOST updates a subscription using the POST alternative endpoint
func (c *Client) UpdateSubscriptionPOST(ctx context.Context, id int, title string) (*Subscription, error) {
	if title == "" {
		return nil, &ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}

	path := fmt.Sprintf("/subscriptions/%d/update.json", id)
	req := &UpdateSubscriptionRequest{
		Title: title,
	}

	var subscription Subscription
	_, err := c.makeRequest(ctx, "POST", path, req, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// DeleteSubscription deletes a subscription
func (c *Client) DeleteSubscription(ctx context.Context, id int) error {
	path := fmt.Sprintf("/subscriptions/%d.json", id)

	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}

// GetSubscriptionsFromURL retrieves subscriptions from a full URL (used for pagination)
func (c *Client) GetSubscriptionsFromURL(ctx context.Context, url string) ([]Subscription, *PaginationInfo, error) {
	var subscriptions []Subscription
	resp, err := c.GetFromURL(ctx, url, &subscriptions)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	pagination := extractPaginationInfo(resp)
	return subscriptions, pagination, nil
}

// Convenience methods without context (use background context)

func (c *Client) GetSubscriptionsWithoutContext(opts *SubscriptionOptions) ([]Subscription, *PaginationInfo, error) {
	return c.GetSubscriptions(context.Background(), opts)
}

func (c *Client) GetSubscriptionWithoutContext(id int) (*Subscription, error) {
	return c.GetSubscription(context.Background(), id)
}

func (c *Client) CreateSubscriptionWithoutContext(feedURL string) (*Subscription, []FeedChoice, error) {
	return c.CreateSubscription(context.Background(), feedURL)
}

func (c *Client) UpdateSubscriptionWithoutContext(id int, title string) (*Subscription, error) {
	return c.UpdateSubscription(context.Background(), id, title)
}

func (c *Client) DeleteSubscriptionWithoutContext(id int) error {
	return c.DeleteSubscription(context.Background(), id)
}

// Helper method to decode error response bodies
func (c *Client) decodeErrorBody(body string, v interface{}) error {
	return nil // Implementation would decode JSON from body string
}
