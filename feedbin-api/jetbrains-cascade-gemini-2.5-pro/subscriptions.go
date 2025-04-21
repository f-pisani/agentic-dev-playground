package feedbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SubscriptionsService handles communication with the subscription related
// methods of the Feedbin API.
type SubscriptionsService service

// ListSubscriptionsOptions specifies the optional parameters to the SubscriptionsService.List method.
type ListSubscriptionsOptions struct {
	Mode *string `url:"mode,omitempty"` // Only "extended" is valid
}

// List retrieves all subscriptions for the authenticated user.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#get-subscriptions
func (s *SubscriptionsService) List(opt *ListSubscriptionsOptions) ([]*Subscription, *http.Response, error) {
	u, err := addOptions("subscriptions.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var subs []*Subscription
	resp, err := s.client.do(req, &subs)
	if err != nil {
		return nil, resp, err
	}

	return subs, resp, nil
}

// Get retrieves a single subscription by its ID.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#get-single-subscription
func (s *SubscriptionsService) Get(subscriptionID int64) (*Subscription, *http.Response, error) {
	path := fmt.Sprintf("subscriptions/%d.json", subscriptionID)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var sub Subscription
	resp, err := s.client.do(req, &sub)
	if err != nil {
		return nil, resp, err
	}

	return &sub, resp, nil
}

// createSubscriptionRequest is the structure for the POST request body.
type createSubscriptionRequest struct {
	FeedURL string `json:"feed_url"`
}

// Create subscribes the user to a new feed.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#create-subscription
func (s *SubscriptionsService) Create(feedURL string) (*Subscription, *http.Response, error) {
	body := createSubscriptionRequest{FeedURL: feedURL}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPost, "subscriptions.json", buf)
	if err != nil {
		return nil, nil, err
	}

	var sub Subscription
	resp, err := s.client.do(req, &sub)
	if err != nil {
		// Handle 300 Multiple Choices specifically
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusMultipleChoices {
			return nil, resp, NewMultipleChoicesError(apiErr)
		}
		return nil, resp, err
	}

	return &sub, resp, nil
}

// updateSubscriptionRequest is the structure for the PATCH request body.
type updateSubscriptionRequest struct {
	Title string `json:"title"`
}

// Update modifies the title of an existing subscription.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#update-subscription
func (s *SubscriptionsService) Update(subscriptionID int64, newTitle string) (*Subscription, *http.Response, error) {
	path := fmt.Sprintf("subscriptions/%d.json", subscriptionID)
	body := updateSubscriptionRequest{Title: newTitle}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPatch, path, buf)
	if err != nil {
		return nil, nil, err
	}

	var sub Subscription
	resp, err := s.client.do(req, &sub)
	if err != nil {
		return nil, resp, err
	}

	return &sub, resp, nil
}

// Delete unsubscribes the user from a feed.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#delete-subscription
func (s *SubscriptionsService) Delete(subscriptionID int64) (*http.Response, error) {
	path := fmt.Sprintf("subscriptions/%d.json", subscriptionID)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	// Expect 204 No Content
	resp, err := s.client.do(req, nil)
	return resp, err
}
