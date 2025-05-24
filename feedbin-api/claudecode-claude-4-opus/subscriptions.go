package feedbin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetSubscriptions retrieves all subscriptions
func (c *Client) GetSubscriptions(options ...QueryOption) ([]Subscription, error) {
	var subscriptions []Subscription
	query := buildQuery(options...)
	err := c.get("/subscriptions.json", query, &subscriptions)
	return subscriptions, err
}

// GetSubscription retrieves a single subscription by ID
func (c *Client) GetSubscription(id int) (*Subscription, error) {
	var subscription Subscription
	path := fmt.Sprintf("/subscriptions/%d.json", id)
	err := c.get(path, nil, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// CreateSubscription creates a new subscription to a feed
func (c *Client) CreateSubscription(feedURL string) (*Subscription, error) {
	body := CreateSubscriptionRequest{
		FeedURL: feedURL,
	}

	var subscription Subscription
	err := c.post("/subscriptions.json", body, &subscription)
	if err != nil {
		// Check if it's a 300 Multiple Choices response
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == 300 {
			return nil, fmt.Errorf("multiple feeds found for URL: %s", feedURL)
		}
		// Check if it's a 302 Found response (already subscribed)
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == 302 {
			return nil, fmt.Errorf("already subscribed to feed: %s", feedURL)
		}
		return nil, err
	}
	return &subscription, nil
}

// CreateSubscriptionWithChoice creates a subscription when multiple feeds are available
func (c *Client) CreateSubscriptionWithChoice(feedURL string) ([]*MultipleChoiceItem, error) {
	body := CreateSubscriptionRequest{
		FeedURL: feedURL,
	}

	// First attempt to create the subscription
	var subscription Subscription
	err := c.post("/subscriptions.json", body, &subscription)

	// If successful, return the subscription
	if err == nil {
		return []*MultipleChoiceItem{{
			FeedURL: subscription.FeedURL,
			Title:   subscription.Title,
		}}, nil
	}

	// Check if it's a 300 Multiple Choices response
	if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == 300 {
		// Parse the multiple choices from the response
		var choices []*MultipleChoiceItem
		if apiErr.Response != nil {
			defer apiErr.Response.Body.Close()
			if err := decodeJSON(apiErr.Response.Body, &choices); err == nil {
				return choices, nil
			}
		}
		return nil, fmt.Errorf("multiple feeds found but could not parse choices")
	}

	return nil, err
}

// decodeJSON is a helper to decode JSON from an io.Reader
func decodeJSON(r interface{}, v interface{}) error {
	switch reader := r.(type) {
	case *http.Response:
		return json.NewDecoder(reader.Body).Decode(v)
	case io.Reader:
		return json.NewDecoder(reader).Decode(v)
	default:
		return fmt.Errorf("unsupported reader type")
	}
}

// UpdateSubscription updates a subscription's title
func (c *Client) UpdateSubscription(id int, title string) (*Subscription, error) {
	body := UpdateSubscriptionRequest{
		Title: title,
	}

	var subscription Subscription
	path := fmt.Sprintf("/subscriptions/%d.json", id)
	err := c.patch(path, body, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// DeleteSubscription deletes a subscription
func (c *Client) DeleteSubscription(id int) error {
	path := fmt.Sprintf("/subscriptions/%d.json", id)
	return c.delete(path, nil)
}

// UpdateSubscriptionAlt updates a subscription's title using POST (alternative to PATCH)
func (c *Client) UpdateSubscriptionAlt(id int, title string) (*Subscription, error) {
	body := UpdateSubscriptionRequest{
		Title: title,
	}

	var subscription Subscription
	path := fmt.Sprintf("/subscriptions/%d/update.json", id)
	err := c.post(path, body, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetSubscriptionsByTag retrieves subscriptions with a specific tag
func (c *Client) GetSubscriptionsByTag(tagName string) ([]Subscription, error) {
	// First get all taggings for the tag
	taggings, err := c.GetTaggingsByName(tagName)
	if err != nil {
		return nil, err
	}

	// Then get all subscriptions and filter by feed IDs
	allSubs, err := c.GetSubscriptions()
	if err != nil {
		return nil, err
	}

	// Create a map of feed IDs from taggings
	feedIDMap := make(map[int]bool)
	for _, tagging := range taggings {
		feedIDMap[tagging.FeedID] = true
	}

	// Filter subscriptions
	var filteredSubs []Subscription
	for _, sub := range allSubs {
		if feedIDMap[sub.FeedID] {
			filteredSubs = append(filteredSubs, sub)
		}
	}

	return filteredSubs, nil
}
