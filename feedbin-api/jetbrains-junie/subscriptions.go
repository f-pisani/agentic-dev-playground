// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
	"time"
)

// SubscriptionService provides access to the subscription-related endpoints of the Feedbin API.
type SubscriptionService struct {
	client *Client
}

// SubscriptionCreateRequest represents a request to create a subscription.
type SubscriptionCreateRequest struct {
	FeedURL string `json:"feed_url"`
}

// SubscriptionUpdateRequest represents a request to update a subscription.
type SubscriptionUpdateRequest struct {
	Title string `json:"title"`
}

// GetSubscriptions retrieves all subscriptions for the authenticated user.
func (s *SubscriptionService) GetSubscriptions(options *PageOptions) ([]Subscription, *http.Response, error) {
	path := "/subscriptions.json"
	if options != nil {
		path = AddPageParams(path, options)
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var subscriptions []Subscription
	resp, err := s.client.Do(req, &subscriptions)
	if err != nil {
		return nil, resp, err
	}

	return subscriptions, resp, nil
}

// GetSubscription retrieves a specific subscription by ID.
func (s *SubscriptionService) GetSubscription(id int) (*Subscription, *http.Response, error) {
	path := fmt.Sprintf("/subscriptions/%d.json", id)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var subscription Subscription
	resp, err := s.client.Do(req, &subscription)
	if err != nil {
		return nil, resp, err
	}

	return &subscription, resp, nil
}

// CreateSubscription creates a new subscription.
func (s *SubscriptionService) CreateSubscription(feedURL string) (*Subscription, *http.Response, error) {
	req, err := s.client.NewRequest("POST", "/subscriptions.json", &SubscriptionCreateRequest{
		FeedURL: feedURL,
	})
	if err != nil {
		return nil, nil, err
	}

	var subscription Subscription
	resp, err := s.client.Do(req, &subscription)
	if err != nil {
		// Check if it's a multiple choices error (300)
		if IsMultipleChoices(err) {
			// Return the error as is, caller can handle the multiple choices
			return nil, resp, err
		}
		return nil, resp, err
	}

	return &subscription, resp, nil
}

// DeleteSubscription deletes a subscription.
func (s *SubscriptionService) DeleteSubscription(id int) (*http.Response, error) {
	path := fmt.Sprintf("/subscriptions/%d.json", id)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateSubscription updates a subscription's title.
func (s *SubscriptionService) UpdateSubscription(id int, title string) (*Subscription, *http.Response, error) {
	path := fmt.Sprintf("/subscriptions/%d.json", id)

	req, err := s.client.NewRequest("PATCH", path, &SubscriptionUpdateRequest{
		Title: title,
	})
	if err != nil {
		return nil, nil, err
	}

	var subscription Subscription
	resp, err := s.client.Do(req, &subscription)
	if err != nil {
		return nil, resp, err
	}

	return &subscription, resp, nil
}

// UpdateSubscriptionAlternative updates a subscription's title using the POST alternative.
// Some proxies block or filter PATCH requests, so this method provides an alternative.
func (s *SubscriptionService) UpdateSubscriptionAlternative(id int, title string) (*Subscription, *http.Response, error) {
	path := fmt.Sprintf("/subscriptions/%d/update.json", id)

	req, err := s.client.NewRequest("POST", path, &SubscriptionUpdateRequest{
		Title: title,
	})
	if err != nil {
		return nil, nil, err
	}

	var subscription Subscription
	resp, err := s.client.Do(req, &subscription)
	if err != nil {
		return nil, resp, err
	}

	return &subscription, resp, nil
}

// GetSubscriptionsSince retrieves all subscriptions created after the specified time.
func (s *SubscriptionService) GetSubscriptionsSince(since time.Time) ([]Subscription, *http.Response, error) {
	options := &PageOptions{
		Since: since.Format(time.RFC3339Nano),
	}

	return s.GetSubscriptions(options)
}

// GetSubscriptionsExtended retrieves all subscriptions with extended metadata.
func (s *SubscriptionService) GetSubscriptionsExtended() ([]Subscription, *http.Response, error) {
	path := "/subscriptions.json?mode=extended"

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var subscriptions []Subscription
	resp, err := s.client.Do(req, &subscriptions)
	if err != nil {
		return nil, resp, err
	}

	return subscriptions, resp, nil
}
