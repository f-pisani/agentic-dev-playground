// Package endpoints contains the API endpoint implementations for the Feedbin API.
package endpoints

import (
	"fmt"
	"net/http"
	"time"

	"z/repos/ai-editor-vibing/feedbin-api/trae-claude-3.7/models"
)

// SubscriptionsService handles communication with the subscription related
// endpoints of the Feedbin API
type SubscriptionsService struct {
	client *Client
}

// List returns a list of subscriptions
// If since is not nil, only subscriptions created after the specified time will be returned
// If mode is "extended", additional metadata will be included
func (s *SubscriptionsService) List(since *time.Time, mode string) ([]models.Subscription, error) {
	params := make(map[string]string)

	if since != nil {
		params["since"] = since.Format(time.RFC3339)
	}

	if mode == "extended" {
		params["mode"] = "extended"
	}

	url := "/subscriptions.json"
	if len(params) > 0 {
		url = AddQueryParams(url, params)
	}

	resp, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var subscriptions []models.Subscription
	if err := s.client.ParseResponse(resp, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Get returns a subscription by ID
func (s *SubscriptionsService) Get(id int64) (*models.Subscription, error) {
	url := fmt.Sprintf("/subscriptions/%d.json", id)

	resp, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var subscription models.Subscription
	if err := s.client.ParseResponse(resp, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// Create creates a new subscription
func (s *SubscriptionsService) Create(feedURL string) (*models.Subscription, error) {
	req := models.SubscriptionCreateRequest{FeedURL: feedURL}

	resp, err := s.client.NewRequest(http.MethodPost, "/subscriptions.json", req)
	if err != nil {
		return nil, err
	}

	var subscription models.Subscription
	if err := s.client.ParseResponse(resp, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// Update updates a subscription
func (s *SubscriptionsService) Update(id int64, title string) (*models.Subscription, error) {
	req := models.SubscriptionUpdateRequest{Title: title}

	url := fmt.Sprintf("/subscriptions/%d.json", id)
	resp, err := s.client.NewRequest(http.MethodPatch, url, req)
	if err != nil {
		return nil, err
	}

	var subscription models.Subscription
	if err := s.client.ParseResponse(resp, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// Delete deletes a subscription
func (s *SubscriptionsService) Delete(id int64) error {
	url := fmt.Sprintf("/subscriptions/%d.json", id)

	resp, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return s.client.ParseResponse(resp, nil)
}
