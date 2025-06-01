package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GetSubscriptions retrieves all subscriptions for the authenticated user.
func (c *Client) GetSubscriptions(since *time.Time, extended bool) ([]Subscription, *PaginationInfo, error) {
	// Build query parameters
	query := url.Values{}

	if since != nil {
		query.Set("since", since.Format(time.RFC3339Nano))
	}

	if extended {
		query.Set("mode", "extended")
	}

	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/subscriptions.json", nil, query)
	if err != nil {
		return nil, nil, err
	}

	// Parse pagination info
	paginationInfo := GetPaginationInfo(resp)

	// Parse the response
	var subscriptions []Subscription
	if err := parseResponse(resp, &subscriptions); err != nil {
		return nil, nil, err
	}

	return subscriptions, &paginationInfo, nil
}

// GetSubscription retrieves a specific subscription by ID.
func (c *Client) GetSubscription(id int) (*Subscription, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/subscriptions/%d.json", id), nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var subscription Subscription
	if err := parseResponse(resp, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// CreateSubscription creates a new subscription.
func (c *Client) CreateSubscription(feedURL string) (*Subscription, error) {
	// Create the request body
	reqBody := SubscriptionRequest{
		FeedURL: feedURL,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/subscriptions.json", reqBody, nil)
	if err != nil {
		return nil, err
	}

	// Handle multiple choices (status code 300)
	if resp.StatusCode == http.StatusMultipleChoices {
		var options []map[string]string
		if err := parseResponse(resp, &options); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("multiple feeds found at URL, options: %v", options)
	}

	// Parse the response
	var subscription Subscription
	if err := parseResponse(resp, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// DeleteSubscription deletes a subscription by ID.
func (c *Client) DeleteSubscription(id int) error {
	// Make the request
	resp, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/subscriptions/%d.json", id), nil, nil)
	if err != nil {
		return err
	}

	// Check for success (204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// UpdateSubscription updates a subscription's title.
func (c *Client) UpdateSubscription(id int, title string) (*Subscription, error) {
	// Create the request body
	reqBody := SubscriptionUpdateRequest{
		Title: title,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPatch, fmt.Sprintf("/subscriptions/%d.json", id), reqBody, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var subscription Subscription
	if err := parseResponse(resp, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// UpdateSubscriptionAlternative updates a subscription's title using the POST alternative.
// This is useful for proxies that block or filter PATCH requests.
func (c *Client) UpdateSubscriptionAlternative(id int, title string) (*Subscription, error) {
	// Create the request body
	reqBody := SubscriptionUpdateRequest{
		Title: title,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, fmt.Sprintf("/subscriptions/%d/update.json", id), reqBody, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var subscription Subscription
	if err := parseResponse(resp, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}
