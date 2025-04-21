package feedbin

import (
	"fmt"
	"net/http"
)

// SubscriptionsService handles communication with the subscription related
// endpoints of the Feedbin API
type SubscriptionsService struct {
	client *Client
}

// SubscriptionOptions specifies optional parameters for subscription requests
type SubscriptionOptions struct {
	PageOptions
	Since string `url:"since,omitempty"`
	Mode  string `url:"mode,omitempty"` // only valid value is "extended"
}

// CreateSubscriptionRequest is used to create a new subscription
type CreateSubscriptionRequest struct {
	FeedURL string `json:"feed_url"`
}

// UpdateSubscriptionRequest is used to update a subscription
type UpdateSubscriptionRequest struct {
	Title string `json:"title"`
}

// List returns all subscriptions
// https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#get-subscriptions
func (s *SubscriptionsService) List(opts *SubscriptionOptions) ([]*Subscription, *http.Response, error) {
	url := "subscriptions.json"
	url, err := AddQueryParams(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	var subscriptions []*Subscription
	resp, err := s.client.Do(req, &subscriptions)
	if err != nil {
		return nil, resp, err
	}

	return subscriptions, resp, nil
}

// Get returns a specific subscription
// https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#get-subscription
func (s *SubscriptionsService) Get(id int) (*Subscription, *http.Response, error) {
	url := fmt.Sprintf("subscriptions/%d.json", id)

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.Do(req, subscription)
	if err != nil {
		return nil, resp, err
	}

	return subscription, resp, nil
}

// Create creates a new subscription
// https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#create-subscription
func (s *SubscriptionsService) Create(feedURL string) (*Subscription, *http.Response, error) {
	url := "subscriptions.json"

	request := &CreateSubscriptionRequest{
		FeedURL: feedURL,
	}

	req, err := s.client.NewRequest(http.MethodPost, url, request)
	if err != nil {
		return nil, nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.Do(req, subscription)
	if err != nil {
		// Check if we received a 302 Found or a 300 Multiple Choices
		if errResp, ok := err.(*ErrorResponse); ok {
			if errResp.Response.StatusCode == http.StatusFound {
				location := errResp.Response.Header.Get("Location")
				if location != "" {
					return s.Get(subscription.ID)
				}
			} else if errResp.Response.StatusCode == http.StatusMultipleChoices {
				// Handle multiple feeds found
				if errResp.Response.Body != nil {
					err = nil
					return nil, errResp.Response, fmt.Errorf("multiple feeds found, inspect response body for options")
				}
			}
		}
		return nil, resp, err
	}

	return subscription, resp, nil
}

// Delete removes a subscription
// https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#delete-subscription
func (s *SubscriptionsService) Delete(id int) (*http.Response, error) {
	url := fmt.Sprintf("subscriptions/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Update updates a subscription
// https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#update-subscription
func (s *SubscriptionsService) Update(id int, title string) (*Subscription, *http.Response, error) {
	url := fmt.Sprintf("subscriptions/%d.json", id)

	request := &UpdateSubscriptionRequest{
		Title: title,
	}

	req, err := s.client.NewRequest(http.MethodPatch, url, request)
	if err != nil {
		return nil, nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.Do(req, subscription)
	if err != nil {
		return nil, resp, err
	}

	return subscription, resp, nil
}

// UpdateAlternative provides an alternative to PATCH for proxies that filter PATCH requests
// https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#patch-alternative
func (s *SubscriptionsService) UpdateAlternative(id int, title string) (*Subscription, *http.Response, error) {
	url := fmt.Sprintf("subscriptions/%d/update.json", id)

	request := &UpdateSubscriptionRequest{
		Title: title,
	}

	req, err := s.client.NewRequest(http.MethodPost, url, request)
	if err != nil {
		return nil, nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.Do(req, subscription)
	if err != nil {
		return nil, resp, err
	}

	return subscription, resp, nil
}
