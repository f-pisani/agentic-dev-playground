package feedbin

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// GetSubscriptions retrieves all subscriptions for the authenticated user.
// Returns a slice of subscriptions and pagination information.
func (c *Client) GetSubscriptions(ctx context.Context, opts *SubscriptionOptions) ([]Subscription, *PaginationInfo, error) {
	params := url.Values{}

	if opts != nil {
		if opts.Since != nil {
			params.Set("since", opts.Since.Format("2006-01-02T15:04:05.000000Z"))
		}
		if opts.Mode != "" {
			params.Set("mode", opts.Mode)
		}
	}

	var subscriptions []Subscription
	resp, err := c.get(ctx, "subscriptions.json", params, &subscriptions)
	if err != nil {
		return nil, nil, err
	}

	pagination := c.parsePagination(resp)
	return subscriptions, pagination, nil
}

// GetSubscription retrieves a specific subscription by ID.
func (c *Client) GetSubscription(ctx context.Context, id int) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)

	var subscription Subscription
	_, err := c.get(ctx, path, nil, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// CreateSubscription creates a new subscription to the specified feed URL.
// Returns the created subscription or an error.
// If the subscription already exists, returns the existing subscription.
// If multiple feeds are found at the URL, returns a MultipleChoicesError.
func (c *Client) CreateSubscription(ctx context.Context, feedURL string) (*Subscription, error) {
	request := CreateSubscriptionRequest{
		FeedURL: feedURL,
	}

	var subscription Subscription
	resp, err := c.post(ctx, "subscriptions.json", nil, request, &subscription)
	if err != nil {
		// Check for multiple choices response (300)
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == 300 {
			// Parse multiple choices response
			var choices MultipleChoicesResponse
			if _, parseErr := c.get(ctx, "subscriptions.json", nil, &choices); parseErr == nil {
				return nil, &MultipleChoicesError{
					Choices: choices,
				}
			}
		}
		return nil, err
	}

	// Handle 302 Found (existing subscription)
	if resp.StatusCode == 302 {
		// The subscription already exists, return it
		return &subscription, nil
	}

	return &subscription, nil
}

// UpdateSubscription updates a subscription's title.
func (c *Client) UpdateSubscription(ctx context.Context, id int, title string) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	request := UpdateSubscriptionRequest{
		Title: title,
	}

	var subscription Subscription
	_, err := c.patch(ctx, path, nil, request, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// UpdateSubscriptionPOST updates a subscription's title using POST method.
// This is an alternative to PATCH for clients that don't support PATCH requests.
func (c *Client) UpdateSubscriptionPOST(ctx context.Context, id int, title string) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d/update.json", id)
	request := UpdateSubscriptionRequest{
		Title: title,
	}

	var subscription Subscription
	_, err := c.post(ctx, path, nil, request, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// DeleteSubscription deletes a subscription by ID.
func (c *Client) DeleteSubscription(ctx context.Context, id int) error {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	_, err := c.delete(ctx, path, nil, nil)
	return err
}

// GetSubscriptionsExtended retrieves all subscriptions with extended metadata.
// This is a convenience method that calls GetSubscriptions with mode="extended".
func (c *Client) GetSubscriptionsExtended(ctx context.Context, since *time.Time) ([]Subscription, *PaginationInfo, error) {
	opts := &SubscriptionOptions{
		Since: since,
		Mode:  "extended",
	}
	return c.GetSubscriptions(ctx, opts)
}

// MultipleChoicesError is returned when multiple feeds are found at a URL
type MultipleChoicesError struct {
	Choices MultipleChoicesResponse
}

// Error implements the error interface
func (e *MultipleChoicesError) Error() string {
	return fmt.Sprintf("multiple feeds found: %d choices available", len(e.Choices))
}

// GetChoices returns the available feed choices
func (e *MultipleChoicesError) GetChoices() MultipleChoicesResponse {
	return e.Choices
}
