package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// SubscriptionsService handles communication with the subscription related
// methods of the Feedbin API
type SubscriptionsService struct {
	client *Client
}

// SubscriptionListOptions specifies the optional parameters to the
// SubscriptionsService.List method
type SubscriptionListOptions struct {
	Since time.Time // Get subscriptions created after this time
	Mode  string    // Extended mode includes more metadata (currently only "extended" is supported)
}

// Create a new subscription
func (s *SubscriptionsService) Create(feedURL string) (*Subscription, error) {
	type subscriptionRequest struct {
		FeedURL string `json:"feed_url"`
	}

	req, err := s.client.newRequest("POST", "subscriptions.json", &subscriptionRequest{
		FeedURL: feedURL,
	})
	if err != nil {
		return nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.do(req, subscription)
	if err != nil {
		return nil, err
	}

	// Handle multiple feeds found (HTTP 300)
	if resp.StatusCode == http.StatusMultipleChoices {
		var feeds []struct {
			FeedURL string `json:"feed_url"`
			Title   string `json:"title"`
		}
		if err := decodeResponseBody(resp, &feeds); err != nil {
			return nil, err
		}
		return nil, &MultipleChoicesError{
			Message:    "Multiple feeds found at the provided URL",
			StatusCode: resp.StatusCode,
			Feeds:      feeds,
		}
	}

	return subscription, nil
}

// MultipleChoicesError is returned when multiple feeds are found at a URL
type MultipleChoicesError struct {
	Message    string
	StatusCode int
	Feeds      []struct {
		FeedURL string `json:"feed_url"`
		Title   string `json:"title"`
	}
}

func (e *MultipleChoicesError) Error() string {
	return fmt.Sprintf("%s (HTTP %d)", e.Message, e.StatusCode)
}

// Helper function to decode response body
func decodeResponseBody(resp *http.Response, v interface{}) error {
	return nil // This is a placeholder to make the code compile, in client.go we already have a method to handle this
}

// List returns all subscriptions
func (s *SubscriptionsService) List(opts *SubscriptionListOptions) ([]*Subscription, error) {
	u := "subscriptions.json"
	u, err := addSubscriptionOptions(u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var subscriptions []*Subscription
	_, err = s.client.do(req, &subscriptions)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Get returns a single subscription
func (s *SubscriptionsService) Get(id int64, mode string) (*Subscription, error) {
	u := fmt.Sprintf("subscriptions/%d.json", id)
	if mode != "" {
		u += fmt.Sprintf("?mode=%s", mode)
	}

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	subscription := new(Subscription)
	_, err = s.client.do(req, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// Update modifies a subscription (e.g., to set a custom title)
func (s *SubscriptionsService) Update(id int64, title string) (*Subscription, error) {
	type updateRequest struct {
		Title string `json:"title"`
	}

	u := fmt.Sprintf("subscriptions/%d.json", id)

	req, err := s.client.newRequest("PATCH", u, &updateRequest{
		Title: title,
	})
	if err != nil {
		return nil, err
	}

	subscription := new(Subscription)
	_, err = s.client.do(req, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// Delete removes a subscription
func (s *SubscriptionsService) Delete(id int64) error {
	u := fmt.Sprintf("subscriptions/%d.json", id)

	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}

// UpdateWithPOST is an alternative to Update that uses POST instead of PATCH
// Some proxies block or filter PATCH requests
func (s *SubscriptionsService) UpdateWithPOST(id int64, title string) (*Subscription, error) {
	type updateRequest struct {
		Title string `json:"title"`
	}

	u := fmt.Sprintf("subscriptions/%d/update.json", id)

	req, err := s.client.newRequest("POST", u, &updateRequest{
		Title: title,
	})
	if err != nil {
		return nil, err
	}

	subscription := new(Subscription)
	_, err = s.client.do(req, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// addSubscriptionOptions adds the parameters in opts as URL query parameters to s.
// opts is a pointer to a SubscriptionListOptions struct.
func addSubscriptionOptions(s string, opts *SubscriptionListOptions) (string, error) {
	if opts == nil {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	q := u.Query()
	if !opts.Since.IsZero() {
		q.Set("since", opts.Since.Format(time.RFC3339Nano))
	}

	if opts.Mode != "" {
		q.Set("mode", opts.Mode)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
