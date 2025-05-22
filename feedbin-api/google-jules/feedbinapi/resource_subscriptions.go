package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// SubscriptionsService handles operations related to feed subscriptions.
type SubscriptionsService struct {
	client *Client
}

// ListSubscriptions retrieves a list of all subscriptions for the authenticated user.
// Supports 'since' and 'mode' parameters.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/subscriptions.md#get-subscriptions
func (s *SubscriptionsService) ListSubscriptions(ctx context.Context, params *ListSubscriptionsParams) ([]Subscription, *Response, error) {
	path := "subscriptions.json"

	queryParams, err := structToURLValues(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert params to query values for ListSubscriptions: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListSubscriptions request: %w", err)
	}

	var subscriptions []Subscription
	resp, err := s.client.do(req, &subscriptions)
	if err != nil {
		return nil, resp, err
	}

	return subscriptions, resp, nil
}

// GetSubscription retrieves a single subscription by its ID.
// Supports 'mode' parameter (via GetSubscriptionParams, to be added to models.go).
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/subscriptions.md#get-subscription
func (s *SubscriptionsService) GetSubscription(ctx context.Context, subscriptionID int64, params *GetSubscriptionParams) (*Subscription, *Response, error) {
	if subscriptionID <= 0 {
		return nil, nil, fmt.Errorf("subscriptionID must be a positive integer")
	}
	path := fmt.Sprintf("subscriptions/%d.json", subscriptionID)

	queryParams, err := structToURLValues(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert params to query values for GetSubscription: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GetSubscription request: %w", err)
	}

	var sub Subscription
	resp, err := s.client.do(req, &sub)
	if err != nil {
		return nil, resp, err
	}

	return &sub, resp, nil
}

// CreateSubscription creates a new subscription to the specified feed URL.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/subscriptions.md#create-subscription
func (s *SubscriptionsService) CreateSubscription(ctx context.Context, feedURL string) (*Subscription, *Response, error) {
	if feedURL == "" {
		return nil, nil, fmt.Errorf("feedURL cannot be empty")
	}
	path := "subscriptions.json"
	requestBody := &CreateSubscriptionRequest{FeedURL: feedURL}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CreateSubscription request: %w", err)
	}

	var sub Subscription
	// Special handling for 300 Multiple Choices
	// The `do` method will return an APIError for non-2xx. We might need to inspect it.
	// For now, assume `do` handles unmarshalling for 201 and 302 based on Content-Type.
	// If a 300 is returned, `sub` will be zero, and `err` will be an `*APIError`.
	// The caller can then inspect `err.(*APIError).Body` for the choices.
	// Or, we can try to parse `[]MultipleFeedChoice{}` if status is 300.

	// A more advanced `do` method could potentially return the raw body or a pre-parsed
	// error structure for specific status codes like 300.
	// For now, the client will get an APIError and can inspect the body.

	resp, err := s.client.do(req, &sub)
	// If err is APIError and status is 300, the body might contain []MultipleFeedChoice
	// Caller should check:
	// if apiErr, ok := err.(*APIError); ok && apiErr.IsStatus(http.StatusMultipleChoices) {
	//     var choices []MultipleFeedChoice
	//     if json.Unmarshal(apiErr.Body, &choices) == nil { /* use choices */ }
	// }
	if err != nil {
		return nil, resp, err
	}

	// If status is 201 Created or 302 Found, `sub` should be populated.
	// If status is 300 Multiple Choices, `sub` will be empty, and `err` should have been non-nil.
	// This part might need refinement if the API returns an empty body for 300 but a valid JSON for choices
	// in the error body. The current `do` method only tries to unmarshal into `v` for 2xx.

	return &sub, resp, nil
}

// DeleteSubscription deletes a subscription by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/subscriptions.md#delete-subscription
func (s *SubscriptionsService) DeleteSubscription(ctx context.Context, subscriptionID int64) (*Response, error) {
	if subscriptionID <= 0 {
		return nil, fmt.Errorf("subscriptionID must be a positive integer")
	}
	path := fmt.Sprintf("subscriptions/%d.json", subscriptionID)

	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeleteSubscription request: %w", err)
	}

	resp, err := s.client.do(req, nil) // No body expected on success (204)
	return resp, err
}

// UpdateSubscription updates a subscription, typically its title.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/subscriptions.md#update-subscription
func (s *SubscriptionsService) UpdateSubscription(ctx context.Context, subscriptionID int64, newTitle string) (*Subscription, *Response, error) {
	if subscriptionID <= 0 {
		return nil, nil, fmt.Errorf("subscriptionID must be a positive integer")
	}
	if newTitle == "" { // Assuming title cannot be empty, adjust if API allows
		return nil, nil, fmt.Errorf("newTitle cannot be empty")
	}

	path := fmt.Sprintf("subscriptions/%d.json", subscriptionID)
	requestBody := &UpdateSubscriptionRequest{Title: newTitle}

	req, err := s.client.newRequest(ctx, http.MethodPatch, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create UpdateSubscription (PATCH) request: %w", err)
	}

	var updatedSub Subscription
	resp, err := s.client.do(req, &updatedSub)
	if err != nil {
		// Check if it's an error that suggests trying POST alternative
		// For now, just return the error. User can try s.client.patchViaPOST if needed.
		return nil, resp, err
	}

	return &updatedSub, resp, nil
}

// UpdateSubscriptionAlt handles updating a subscription using the POST alternative.
// `POST /v2/subscriptions/:id/update.json`
func (s *SubscriptionsService) UpdateSubscriptionAlt(ctx context.Context, subscriptionID int64, newTitle string) (*Subscription, *Response, error) {
	if subscriptionID <= 0 {
		return nil, nil, fmt.Errorf("subscriptionID must be a positive integer")
	}
	if newTitle == "" {
		return nil, nil, fmt.Errorf("newTitle cannot be empty")
	}

	path := fmt.Sprintf("subscriptions/%d/update.json", subscriptionID)
	requestBody := &UpdateSubscriptionRequest{Title: newTitle}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create UpdateSubscription (POST alt) request: %w", err)
	}

	var updatedSub Subscription
	resp, err := s.client.do(req, &updatedSub)
	if err != nil {
		return nil, resp, err
	}
	return &updatedSub, resp, nil
}
