// Package models contains the data models for the Feedbin API.
package models

import (
	"time"
)

// Subscription represents a Feedbin subscription
type Subscription struct {
	// ID is the subscription ID
	ID int64 `json:"id"`

	// CreatedAt is the time the subscription was created
	CreatedAt time.Time `json:"created_at"`

	// FeedID is the ID of the feed
	FeedID int64 `json:"feed_id"`

	// Title is the title of the subscription
	Title string `json:"title"`

	// FeedURL is the URL of the feed
	FeedURL string `json:"feed_url"`

	// SiteURL is the URL of the site
	SiteURL string `json:"site_url"`

	// Extended fields (only available in extended mode)
	LastPublishedEntry time.Time `json:"last_published_entry,omitempty"`
	IconURL            string    `json:"icon_url,omitempty"`
	FaviconURL         string    `json:"favicon_url,omitempty"`
}

// SubscriptionCreateRequest represents a request to create a subscription
type SubscriptionCreateRequest struct {
	// FeedURL is the URL of the feed to subscribe to
	FeedURL string `json:"feed_url"`
}

// SubscriptionUpdateRequest represents a request to update a subscription
type SubscriptionUpdateRequest struct {
	// Title is the new title for the subscription
	Title string `json:"title"`
}
