// Package models contains the data models for the Feedbin API.
package models

import (
	"time"
)

// Entry represents a Feedbin entry
type Entry struct {
	// ID is the entry ID
	ID int64 `json:"id"`

	// FeedID is the ID of the feed
	FeedID int64 `json:"feed_id"`

	// Title is the title of the entry
	Title string `json:"title,omitempty"`

	// URL is the URL of the entry
	URL string `json:"url"`

	// ExtractedContentURL is the URL for the extracted content
	ExtractedContentURL string `json:"extracted_content_url,omitempty"`

	// Author is the author of the entry
	Author string `json:"author,omitempty"`

	// Content is the content of the entry
	Content string `json:"content,omitempty"`

	// Summary is the summary of the entry
	Summary string `json:"summary,omitempty"`

	// Published is the time the entry was published
	Published time.Time `json:"published"`

	// CreatedAt is the time the entry was created
	CreatedAt time.Time `json:"created_at"`
}

// EntryIDs represents a list of entry IDs
type EntryIDs struct {
	// EntryIDs is the list of entry IDs
	EntryIDs []int64 `json:"entry_ids,omitempty"`

	// UnreadEntries is the list of unread entry IDs
	UnreadEntries []int64 `json:"unread_entries,omitempty"`

	// StarredEntries is the list of starred entry IDs
	StarredEntries []int64 `json:"starred_entries,omitempty"`
}
