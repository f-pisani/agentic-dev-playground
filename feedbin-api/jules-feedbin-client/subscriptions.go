package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings" // Required for addOptions, if used.

	"github.com/google/go-querystring/query"
)

// SubscriptionsService handles operations related to feed subscriptions.
type SubscriptionsService struct {
	client *Client
}

// NewSubscriptionsService creates a new service for subscription related operations.
func NewSubscriptionsService(client *Client) *SubscriptionsService {
	return &SubscriptionsService{client: client}
}

// SubscriptionListOptions specifies the optional parameters to the SubscriptionsService.List method.
type SubscriptionListOptions struct {
	ListOptions      // Embeds Page, PerPage, Since
	Mode        string `url:"mode,omitempty"` // "extended" for more details
}

// List retrieves all subscriptions.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#get-subscriptions
func (s *SubscriptionsService) List(opts *SubscriptionListOptions) ([]Subscription, *http.Response, error) {
	path := "subscriptions.json"
	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		if params := v.Encode(); params != "" {
			path = fmt.Sprintf("%s?%s", path, params)
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
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

// SubscriptionGetOptions specifies the optional parameters to the SubscriptionsService.Get method.
type SubscriptionGetOptions struct {
	Mode string `url:"mode,omitempty"` // "extended" for more details
}

// Get retrieves a single subscription.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#get-subscription
func (s *SubscriptionsService) Get(id int64, opts *SubscriptionGetOptions) (*Subscription, *http.Response, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		if params := v.Encode(); params != "" {
			path = fmt.Sprintf("%s?%s", path, params)
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
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

// CreateSubscriptionOptions specifies the parameters to the SubscriptionsService.Create method.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#create-subscription
type CreateSubscriptionOptions struct {
	FeedURL string `json:"feed_url"`
}

// Create subscribes to a new feed.
// Returns the created Subscription.
func (s *SubscriptionsService) Create(opts *CreateSubscriptionOptions) (*Subscription, *http.Response, error) {
	if opts == nil || opts.FeedURL == "" {
		return nil, nil, fmt.Errorf("FeedURL is required to create a subscription")
	}
	path := "subscriptions.json"
	req, err := s.client.NewRequest(http.MethodPost, path, opts)
	if err != nil {
		return nil, nil, err
	}

	var createdSubscription Subscription
	resp, err := s.client.Do(req, &createdSubscription)
	if err != nil {
		return nil, resp, err
	}
	return &createdSubscription, resp, nil
}

// UpdateSubscriptionOptions specifies the parameters to the SubscriptionsService.Update method.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#update-subscription
type UpdateSubscriptionOptions struct {
	Title string `json:"title"` // The only updatable field mentioned is title.
}

// Update modifies an existing subscription. Currently, only the title can be updated.
// Returns the updated Subscription.
func (s *SubscriptionsService) Update(id int64, opts *UpdateSubscriptionOptions) (*Subscription, *http.Response, error) {
	if opts == nil || opts.Title == "" {
		return nil, nil, fmt.Errorf("Title is required to update a subscription")
	}
	path := fmt.Sprintf("subscriptions/%d.json", id)
	req, err := s.client.NewRequest(http.MethodPatch, path, opts)
	if err != nil {
		return nil, nil, err
	}

	var updatedSubscription Subscription
	resp, err := s.client.Do(req, &updatedSubscription)
	if err != nil {
		return nil, resp, err
	}
	return &updatedSubscription, resp, nil
}

// Delete unsubscribes from a feed.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/subscriptions.md#delete-subscription
func (s *SubscriptionsService) Delete(id int64) (*http.Response, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	// Successful deletion returns 204 No Content, so pass nil for the body interface.
	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Helper function to add query parameters to a URL path.
// Note: This is a more generic helper. The go-querystring library is used above,
// which is generally preferred for struct-to-querystring conversion.
// This function is not actively used in the current code but kept for potential future use.
func addOptions(path string, opts interface{}) (string, error) {
	if opts == nil {
		return path, nil
	}
	v, err := query.Values(opts)
	if err != nil {
		return path, err
	}
	if params := v.Encode(); params != "" {
		if strings.Contains(path, "?") {
			return fmt.Sprintf("%s&%s", path, params), nil
		}
		return fmt.Sprintf("%s?%s", path, params), nil
	}
	return path, nil
}

// Ensure you have `github.com/google/go-querystring/query` by running:
// go get github.com/google/go-querystring/query
// This will be handled later when we create go.mod and tidy dependencies.
