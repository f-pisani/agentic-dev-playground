package feedbin

import (
	"time"
)

// Subscription represents a Feedbin subscription
type Subscription struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int       `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"`
}

// JSONFeed represents additional metadata for JSON feeds
type JSONFeed struct {
	Favicon     string `json:"favicon,omitempty"`
	FeedURL     string `json:"feed_url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
	HomePageURL string `json:"home_page_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

// SubscriptionCreateRequest represents a request to create a subscription
type SubscriptionCreateRequest struct {
	FeedURL string `json:"feed_url"`
}

// SubscriptionUpdateRequest represents a request to update a subscription
type SubscriptionUpdateRequest struct {
	Title string `json:"title"`
}

// Entry represents a Feedbin entry
type Entry struct {
	ID                  int       `json:"id"`
	FeedID              int       `json:"feed_id"`
	Title               *string   `json:"title"`
	URL                 string    `json:"url"`
	ExtractedContentURL string    `json:"extracted_content_url"`
	Author              *string   `json:"author"`
	Content             *string   `json:"content"`
	Summary             *string   `json:"summary"`
	Published           time.Time `json:"published"`
	CreatedAt           time.Time `json:"created_at"`

	// Extended mode fields
	Original          *EntryOriginal     `json:"original,omitempty"`
	Images            *EntryImages       `json:"images,omitempty"`
	Enclosure         *EntryEnclosure    `json:"enclosure,omitempty"`
	TwitterID         *int64             `json:"twitter_id,omitempty"`
	TwitterThreadIDs  []int64            `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles []ExtractedArticle `json:"extracted_articles,omitempty"`
	JSONFeed          *JSONFeed          `json:"json_feed,omitempty"`
}

// EntryOriginal represents the original entry data before updates
type EntryOriginal struct {
	Author    *string     `json:"author"`
	Content   *string     `json:"content"`
	Title     *string     `json:"title"`
	URL       string      `json:"url"`
	EntryID   string      `json:"entry_id"`
	Published time.Time   `json:"published"`
	Data      interface{} `json:"data"`
}

// EntryImages represents images associated with an entry
type EntryImages struct {
	OriginalURL string         `json:"original_url"`
	Size1       EntryImageSize `json:"size_1"`
}

// EntryImageSize represents a specific image size
type EntryImageSize struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// EntryEnclosure represents podcast/RSS enclosure data
type EntryEnclosure struct {
	EnclosureURL    string `json:"enclosure_url"`
	EnclosureType   string `json:"enclosure_type"`
	EnclosureLength string `json:"enclosure_length"`
	ItunesDuration  string `json:"itunes_duration,omitempty"`
	ItunesImage     string `json:"itunes_image,omitempty"`
}

// ExtractedArticle represents an article extracted from a tweet
type ExtractedArticle struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Host    string `json:"host"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

// UnreadEntriesRequest represents a request to mark entries as unread
type UnreadEntriesRequest struct {
	UnreadEntries []int `json:"unread_entries"`
}

// StarredEntriesRequest represents a request to star entries
type StarredEntriesRequest struct {
	StarredEntries []int `json:"starred_entries"`
}

// Tag represents a Feedbin tag
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Tagging represents a Feedbin tagging (association between a tag and a feed)
type Tagging struct {
	ID     int `json:"id"`
	FeedID int `json:"feed_id"`
	TagID  int `json:"tag_id"`
}

// TaggingCreateRequest represents a request to create a tagging
type TaggingCreateRequest struct {
	FeedID int `json:"feed_id"`
	TagID  int `json:"tag_id"`
}

// SavedSearch represents a Feedbin saved search
type SavedSearch struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// SavedSearchCreateRequest represents a request to create a saved search
type SavedSearchCreateRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// SavedSearchUpdateRequest represents a request to update a saved search
type SavedSearchUpdateRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}
