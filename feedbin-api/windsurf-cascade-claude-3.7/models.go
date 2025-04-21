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
	// Extended mode fields
	JSONFeed *JSONFeed `json:"json_feed,omitempty"`
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

// Entry represents a Feedbin entry
type Entry struct {
	ID                  int       `json:"id"`
	FeedID              int       `json:"feed_id"`
	Title               *string   `json:"title"`
	Author              *string   `json:"author"`
	Summary             *string   `json:"summary"`
	Content             *string   `json:"content"`
	URL                 string    `json:"url"`
	ExtractedContentURL string    `json:"extracted_content_url"`
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
	ContentDiff       *string            `json:"content_diff,omitempty"`
}

// EntryOriginal represents the original version of an entry
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
	OriginalURL string          `json:"original_url"`
	Size1       *EntryImageSize `json:"size_1,omitempty"`
}

// EntryImageSize represents an image size
type EntryImageSize struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// EntryEnclosure represents enclosure data for podcasts
type EntryEnclosure struct {
	EnclosureURL    string `json:"enclosure_url"`
	EnclosureType   string `json:"enclosure_type"`
	EnclosureLength string `json:"enclosure_length"`
	ITunesDuration  string `json:"itunes_duration,omitempty"`
	ITunesImage     string `json:"itunes_image,omitempty"`
}

// ExtractedArticle represents an article extracted from a tweet
type ExtractedArticle struct {
	URL     string  `json:"url"`
	Title   string  `json:"title"`
	Host    string  `json:"host"`
	Author  *string `json:"author"`
	Content *string `json:"content"`
}

// Tagging represents a Feedbin tagging
type Tagging struct {
	ID     int    `json:"id"`
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a Feedbin saved search
type SavedSearch struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Icon represents a Feedbin feed icon
type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

// Import represents a Feedbin OPML import
type Import struct {
	ID          int          `json:"id"`
	Complete    bool         `json:"complete"`
	CreatedAt   time.Time    `json:"created_at"`
	ImportItems []ImportItem `json:"import_items,omitempty"`
}

// ImportItem represents an item in a Feedbin import
type ImportItem struct {
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	Status  string `json:"status"`
}

// SubscriptionParams represents the parameters for creating a subscription
type SubscriptionParams struct {
	FeedURL string `json:"feed_url"`
}

// TaggingParams represents the parameters for creating a tagging
type TaggingParams struct {
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearchParams represents the parameters for creating a saved search
type SavedSearchParams struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// PageParams represents the parameters for creating a page
type PageParams struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// TagParams represents the parameters for updating a tag
type TagParams struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// TagDeleteParams represents the parameters for deleting a tag
type TagDeleteParams struct {
	Name string `json:"name"`
}

// EntryIDs represents a list of entry IDs
type EntryIDs struct {
	IDs []int `json:"ids,omitempty"`
}

// UnreadEntryIDs represents a list of unread entry IDs
type UnreadEntryIDs struct {
	UnreadEntries []int `json:"unread_entries"`
}

// StarredEntryIDs represents a list of starred entry IDs
type StarredEntryIDs struct {
	StarredEntries []int `json:"starred_entries"`
}

// UpdatedEntryIDs represents a list of updated entry IDs
type UpdatedEntryIDs struct {
	UpdatedEntries []int `json:"updated_entries"`
}

// RecentlyReadEntryIDs represents a list of recently read entry IDs
type RecentlyReadEntryIDs struct {
	RecentlyReadEntries []int `json:"recently_read_entries"`
}

// ListOptions represents options for list endpoints
type ListOptions struct {
	Page    *int       `url:"page,omitempty"`
	PerPage *int       `url:"per_page,omitempty"`
	Since   *time.Time `url:"since,omitempty"`
}

// EntryListOptions represents options for listing entries
type EntryListOptions struct {
	ListOptions
	IDs                []int  `url:"ids,omitempty,comma"`
	Read               *bool  `url:"read,omitempty"`
	Starred            *bool  `url:"starred,omitempty"`
	Mode               string `url:"mode,omitempty"`
	IncludeOriginal    *bool  `url:"include_original,omitempty"`
	IncludeEnclosure   *bool  `url:"include_enclosure,omitempty"`
	IncludeContentDiff *bool  `url:"include_content_diff,omitempty"`
}
