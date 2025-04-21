package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// SubscriptionsService handles communication with the subscription related
// methods of the Feedbin API.
type SubscriptionsService struct {
	client *Client
}

// SubscriptionListOptions specifies the optional parameters to the
// SubscriptionsService.List method.
type SubscriptionListOptions struct {
	Since time.Time
	Mode  string // "extended" is the only available mode
}

// List returns all subscriptions.
func (s *SubscriptionsService) List(opts *SubscriptionListOptions) ([]*Subscription, *http.Response, error) {
	u := "/v2/subscriptions.json"

	if opts != nil {
		params := url.Values{}

		if !opts.Since.IsZero() {
			params.Add("since", opts.Since.Format(time.RFC3339Nano))
		}

		if opts.Mode != "" {
			params.Add("mode", opts.Mode)
		}

		if len(params) > 0 {
			u += "?" + params.Encode()
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
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

// Get returns a single subscription.
func (s *SubscriptionsService) Get(id int64, mode string) (*Subscription, *http.Response, error) {
	u := fmt.Sprintf("/v2/subscriptions/%d.json", id)

	if mode != "" {
		u += fmt.Sprintf("?mode=%s", mode)
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
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

// CreateSubscriptionOptions specifies the parameters to the
// SubscriptionsService.Create method.
type CreateSubscriptionOptions struct {
	FeedURL string `json:"feed_url"`
}

// Create creates a new subscription.
func (s *SubscriptionsService) Create(opts *CreateSubscriptionOptions) (*Subscription, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/subscriptions.json", opts)
	if err != nil {
		return nil, nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.Do(req, subscription)
	if err != nil {
		// Check if we got a 302 Found (subscription exists) or 300 Multiple Choices (multiple feeds found)
		if errResp, ok := err.(*ErrorResponse); ok {
			if errResp.Response.StatusCode == http.StatusFound {
				// Get the subscription from the Location header
				location := errResp.Response.Header.Get("Location")
				if location != "" {
					req, err := s.client.NewRequest(http.MethodGet, location, nil)
					if err != nil {
						return nil, nil, err
					}

					resp, err := s.client.Do(req, subscription)
					if err != nil {
						return nil, resp, err
					}

					return subscription, resp, nil
				}
			} else if errResp.Response.StatusCode == http.StatusMultipleChoices {
				// Multiple feeds found, return the error
				var feeds []*struct {
					FeedURL string `json:"feed_url"`
					Title   string `json:"title"`
				}

				if err := json.NewDecoder(errResp.Response.Body).Decode(&feeds); err != nil {
					return nil, errResp.Response, err
				}

				return nil, errResp.Response, fmt.Errorf("multiple feeds found: %v", feeds)
			}
		}

		return nil, resp, err
	}

	return subscription, resp, nil
}

// Delete deletes a subscription.
func (s *SubscriptionsService) Delete(id int64) (*http.Response, error) {
	u := fmt.Sprintf("/v2/subscriptions/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UpdateSubscriptionOptions specifies the parameters to the
// SubscriptionsService.Update method.
type UpdateSubscriptionOptions struct {
	Title string `json:"title"`
}

// Update updates a subscription.
func (s *SubscriptionsService) Update(id int64, opts *UpdateSubscriptionOptions) (*Subscription, *http.Response, error) {
	u := fmt.Sprintf("/v2/subscriptions/%d.json", id)

	req, err := s.client.NewRequest(http.MethodPatch, u, opts)
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

// UpdateAlternative updates a subscription using the POST alternative for PATCH.
func (s *SubscriptionsService) UpdateAlternative(id int64, opts *UpdateSubscriptionOptions) (*Subscription, *http.Response, error) {
	u := fmt.Sprintf("/v2/subscriptions/%d/update.json", id)

	req, err := s.client.NewRequest(http.MethodPost, u, opts)
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
