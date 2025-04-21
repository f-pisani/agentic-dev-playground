package feedbin

import (
	"fmt"
	"net/http"
	"strconv"
)

// SubscriptionsService handles communication with the subscription related
// endpoints of the Feedbin API
type SubscriptionsService struct {
	client *Client
}

// SubscriptionListOptions specifies the optional parameters to the
// SubscriptionsService.List method
type SubscriptionListOptions struct {
	ListOptions
	Mode string // "extended" mode is available
}

// SubscriptionCreateOptions specifies the parameters to the
// SubscriptionsService.Create method
type SubscriptionCreateOptions struct {
	FeedURL string `json:"feed_url"`
}

// SubscriptionUpdateOptions specifies the parameters to the
// SubscriptionsService.Update method
type SubscriptionUpdateOptions struct {
	Title string `json:"title"`
}

// List lists all subscriptions for the authenticated user
func (s *SubscriptionsService) List(opts *SubscriptionListOptions) ([]*Subscription, *http.Response, error) {
	url := "subscriptions.json"
	if opts != nil {
		params := make(map[string]string)
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if !opts.Since.IsZero() {
			params["since"] = opts.Since.Format("2006-01-02T15:04:05.000000Z")
		}
		if opts.Mode != "" {
			params["mode"] = opts.Mode
		}

		// Add parameters
		if len(params) > 0 {
			first := true
			for k, v := range params {
				if first {
					url += "?" + k + "=" + v
					first = false
				} else {
					url += "&" + k + "=" + v
				}
			}
		}
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

// Get gets a single subscription by ID
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
func (s *SubscriptionsService) Create(opts *SubscriptionCreateOptions) (*Subscription, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "subscriptions.json", opts)
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

// Delete deletes a subscription
func (s *SubscriptionsService) Delete(id int) (*http.Response, error) {
	url := fmt.Sprintf("subscriptions/%d.json", id)
	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Update updates a subscription
func (s *SubscriptionsService) Update(id int, opts *SubscriptionUpdateOptions) (*Subscription, *http.Response, error) {
	url := fmt.Sprintf("subscriptions/%d.json", id)
	req, err := s.client.NewRequest(http.MethodPatch, url, opts)
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

// UpdateWithPost updates a subscription using POST instead of PATCH
// Some proxies may block PATCH requests
func (s *SubscriptionsService) UpdateWithPost(id int, opts *SubscriptionUpdateOptions) (*Subscription, *http.Response, error) {
	url := fmt.Sprintf("subscriptions/%d/update.json", id)
	req, err := s.client.NewRequest(http.MethodPost, url, opts)
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
