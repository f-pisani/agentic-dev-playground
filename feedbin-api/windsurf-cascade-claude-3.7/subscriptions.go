package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// SubscriptionsService handles operations related to subscriptions
type SubscriptionsService struct {
	client *Client
}

// List returns all subscriptions
// If options is nil, default options will be used
func (s *SubscriptionsService) List(options *ListOptions) ([]Subscription, error) {
	path := "subscriptions.json"
	if options != nil {
		path = addOptionsToPath(path, options)
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var subscriptions []Subscription
	_, err = s.client.Do(req, &subscriptions)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Get returns a single subscription by ID
func (s *SubscriptionsService) Get(id int) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	subscription := new(Subscription)
	_, err = s.client.Do(req, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// Create creates a new subscription
func (s *SubscriptionsService) Create(feedURL string) (*Subscription, error) {
	params := &SubscriptionParams{
		FeedURL: feedURL,
	}

	req, err := s.client.NewRequest(http.MethodPost, "subscriptions.json", params)
	if err != nil {
		return nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.Do(req, subscription)
	if err != nil {
		// Check if it's a 300 Multiple Choices error
		if apiErr, ok := err.(*Error); ok && apiErr.StatusCode == http.StatusMultipleChoices {
			// Parse the response to get the available feeds
			var feeds []struct {
				FeedURL string `json:"feed_url"`
				Title   string `json:"title"`
			}

			decodeErr := json.NewDecoder(resp.Body).Decode(&feeds)
			if decodeErr != nil {
				return nil, decodeErr
			}

			// Return a custom error with available feeds
			return nil, &Error{
				StatusCode: http.StatusMultipleChoices,
				Message:    fmt.Sprintf("Multiple feeds found. Choose one of: %v", feeds),
			}
		}
		return nil, err
	}

	return subscription, nil
}

// Update updates a subscription
func (s *SubscriptionsService) Update(id int, title string) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)

	params := struct {
		Title string `json:"title"`
	}{
		Title: title,
	}

	req, err := s.client.NewRequest(http.MethodPatch, path, params)
	if err != nil {
		return nil, err
	}

	subscription := new(Subscription)
	_, err = s.client.Do(req, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// Delete deletes a subscription
func (s *SubscriptionsService) Delete(id int) error {
	path := fmt.Sprintf("subscriptions/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Helper function to add options to the path
func addOptionsToPath(path string, options *ListOptions) string {
	values := url.Values{}

	if options.Page != nil {
		values.Add("page", strconv.Itoa(*options.Page))
	}

	if options.PerPage != nil {
		values.Add("per_page", strconv.Itoa(*options.PerPage))
	}

	if options.Since != nil {
		values.Add("since", options.Since.Format(time.RFC3339Nano))
	}

	if len(values) > 0 {
		return path + "?" + values.Encode()
	}

	return path
}
