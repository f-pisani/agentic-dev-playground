package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SubscriptionsService handles subscription-related operations
type SubscriptionsService struct {
	client *Client
}

// List returns all subscriptions
func (s *SubscriptionsService) List(since *time.Time, extended bool) ([]Subscription, error) {
	path := "/subscriptions.json"

	// Add query parameters
	if since != nil || extended {
		query := url.Values{}

		if since != nil {
			query.Set("since", since.Format(time.RFC3339Nano))
		}

		if extended {
			query.Set("mode", "extended")
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var subscriptions []Subscription
	_, err = s.client.do(req, &subscriptions)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Get returns a specific subscription
func (s *SubscriptionsService) Get(id int) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%d.json", id)

	req, err := s.client.newRequest(http.MethodGet, path, nil)
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

// Create creates a new subscription
func (s *SubscriptionsService) Create(feedURL string) (*Subscription, error) {
	createReq := SubscriptionCreateRequest{
		FeedURL: feedURL,
	}

	req, err := s.client.newRequest(http.MethodPost, "/subscriptions.json", createReq)
	if err != nil {
		return nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.do(req, subscription)
	if err != nil {
		// Check if it's a multiple choices error
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusMultipleChoices {
			var choices []struct {
				FeedURL string `json:"feed_url"`
				Title   string `json:"title"`
			}

			if err := json.Unmarshal([]byte(apiErr.Body), &choices); err != nil {
				return nil, fmt.Errorf("multiple feeds found but could not parse response: %v", err)
			}

			return nil, fmt.Errorf("multiple feeds found: %v", choices)
		}

		return nil, err
	}

	// Check if it's a redirect (subscription already exists)
	if resp.StatusCode == http.StatusFound {
		// Extract the subscription ID from the Location header
		location := resp.Header.Get("Location")
		if location == "" {
			return nil, fmt.Errorf("subscription already exists but no Location header found")
		}

		// Parse the ID from the Location header
		parts := strings.Split(location, "/")
		if len(parts) == 0 {
			return nil, fmt.Errorf("invalid Location header: %s", location)
		}

		idPart := parts[len(parts)-1]
		idPart = strings.TrimSuffix(idPart, ".json")

		// Get the subscription by ID
		return s.Get(parseInt(idPart))
	}

	return subscription, nil
}

// Delete deletes a subscription
func (s *SubscriptionsService) Delete(id int) error {
	path := fmt.Sprintf("/subscriptions/%d.json", id)

	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}

// Update updates a subscription
func (s *SubscriptionsService) Update(id int, title string) (*Subscription, error) {
	path := fmt.Sprintf("/subscriptions/%d.json", id)

	updateReq := SubscriptionUpdateRequest{
		Title: title,
	}

	req, err := s.client.newRequest(http.MethodPatch, path, updateReq)
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

// Helper function to parse an integer from a string
func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
