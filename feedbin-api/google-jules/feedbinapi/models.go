package feedbinapi

import "time"

// Entry represents a single feed entry.
// Based on specs/content/entries.md and existing clients.
type Entry struct {
	ID          int64     `json:"id"`
	FeedID      int64     `json:"feed_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Author      string    `json:"author"`
	Content     string    `json:"content"` // HTML content
	Summary     string    `json:"summary"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"` // Added based on common practice

	// Fields from "extended" mode or other contexts
	Image       *EntryImage `json:"image,omitempty"`
	Source      string      `json:"source,omitempty"`      // e.g. name of the website
	Twitter     *TwitterData `json:"twitter,omitempty"`   // If entry is from Twitter
	ExtractedAt *time.Time   `json:"extracted_at,omitempty"` // If content was extracted
	Original    *OriginalData `json:"original,omitempty"`  // Original entry data if updated
}

// EntryImage represents an image associated with an entry.
type EntryImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// TwitterData holds Twitter-specific information for an entry.
// Based on specs/content/supporting-twitter.md
type TwitterData struct {
	TweetID         string `json:"twitter_tweet_id,omitempty"`
	ScreenName      string `json:"twitter_screen_name,omitempty"`
	Name            string `json:"twitter_name,omitempty"`
	ProfileImageURL string `json:"twitter_profile_image_url,omitempty"`
	RetweetedBy     string `json:"twitter_retweeted_by_screen_name,omitempty"` // Screen name of retweeter
}

// OriginalData holds the original content of an entry if it has been updated.
type OriginalData struct {
	URL         string `json:"url,omitempty"`
	Title       string `json:"title,omitempty"`
	Author      string `json:"author,omitempty"`
	Content     string `json:"content,omitempty"`
	Summary     string `json:"summary,omitempty"`
	PublishedAt string `json:"published_at,omitempty"` // Keep as string if unsure about parsing non-standard original dates
}


// Subscription (Feed) represents a subscribed feed.
// Based on specs/content/subscriptions.md
type Subscription struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int64     `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	// Extended fields
	IsSpark    *bool   `json:"is_spark,omitempty"`
	IsPodcast  *bool   `json:"is_podcast,omitempty"`
	IsTwitter  *bool   `json:"is_twitter,omitempty"`
	HasIcon    *bool   `json:"has_icon,omitempty"`
	Icon       *Icon   `json:"icon,omitempty"` // Potentially link to Icon struct
	FeedStats  *FeedStats `json:"feed_stats,omitempty"`
}

// FeedStats contains statistics for a feed.
// Part of extended subscription data.
type FeedStats struct {
	Subscribers      int    `json:"subscribers"`
	EntriesPerDay    float64 `json:"entries_per_day"`
	LastEntryAt      *time.Time `json:"last_entry_at,omitempty"`
	LastStatus       string `json:"last_status,omitempty"` // e.g., "ok", "error"
	LastStatusAt     *time.Time `json:"last_status_at,omitempty"`
	ConsecutiveErrors int   `json:"consecutive_errors"`
}

// Tag represents a tag.
// Based on specs/content/tags.md
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// FeedIDs []int64 `json:"feed_ids,omitempty"` // This seems to be part of Taggings, not the Tag object itself.
}

// Tagging represents a tagging of an entry with a tag.
// Based on specs/content/taggings.md
type Tagging struct {
	ID        int64     `json:"id"`
	FeedID    int64     `json:"feed_id"`
	EntryID   int64     `json:"entry_id"`
	Name      string    `json:"name"`      // Tag name
	CreatedAt time.Time `json:"created_at"` // Added, as it's common
}

// SavedSearch represents a saved search query.
// Based on specs/content/saved-searches.md
type SavedSearch struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
	// FeedIDs   []int64   `json:"feed_ids,omitempty"` // Optional: feeds to search within
}

// RecentlyReadEntry indicates an entry that was recently read.
// Based on specs/content/recently-read-entries.md
type RecentlyReadEntry struct {
	EntryID      int64     `json:"entry_id"`
	Interaction  string    `json:"interaction"` // e.g., "read", "scrolled"
	InteractedAt time.Time `json:"interacted_at"`
}

// UpdatedEntryID refers to an entry ID that has been updated.
// Based on specs/content/updated-entries.md
// This endpoint returns an array of entry IDs, so a specific struct might not be needed
// unless there's more structure to it. For now, assume it's []int64.

// Icon represents a feed icon.
// Based on specs/content/icons.md
type Icon struct {
	Host      string `json:"host"`
	Data      string `json:"data"` // Base64 encoded image data
	Extension string `json:"extension"` // e.g., "png", "ico"
}

// Import represents an OPML import job.
// Based on specs/content/imports.md
type Import struct {
	ID          int64      `json:"id"`
	Status      string     `json:"status"` // e.g., "pending", "running", "complete", "error"
	Message     string     `json:"message,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	TotalFeeds  int        `json:"total_feeds,omitempty"`
	Imported    int        `json:"imported_feeds,omitempty"` // Corrected field name based on common sense
	Unchanged   int        `json:"unchanged_feeds,omitempty"`
	Errored     int        `json:"errored_feeds,omitempty"`
}

// Page represents a downloaded page.
// Based on specs/content/pages.md
type Page struct {
	EntryID       int64     `json:"entry_id"`
	Body          string    `json:"body"` // HTML content of the page
	ExtractedFrom string    `json:"extracted_from"` // URL it was extracted from
	WordCount     int       `json:"word_count"`
	Processed     bool      `json:"processed"`
	PublishedAt   time.Time `json:"published_at"`
	Domain        string    `json:"domain"`
	Path          string    `json:"path"`
	// Image URLs and other extracted info might be here too, TBD from actual API response if needed.
}

// ExtractResult represents the result of a content extraction.
// Based on specs/content/extract-full-content.md
type ExtractResult struct {
	URL         string `json:"url"`
	Title       string `json:"title,omitempty"`
	Author      string `json:"author,omitempty"`
	PublishedAt string `json:"published_at,omitempty"` // String because format isn't guaranteed ISO 8601
	Dek         string `json:"dek,omitempty"`          // Subtitle or summary
	LeadImageURL string `json:"lead_image_url,omitempty"`
	Content     string `json:"content"` // HTML content
	NextPageURL string `json:"next_page_url,omitempty"`
	Excerpt     string `json:"excerpt,omitempty"`
	WordCount   int    `json:"word_count,omitempty"`
	Direction   string `json:"direction,omitempty"` // e.g., "ltr", "rtl"
	// Other fields like domain, total_pages, rendered_pages, etc.
}

// General API Options Structs

// ListOptions provides general options for list methods.
type ListOptions struct {
	Page    int  `url:"page,omitempty"`
	PerPage int  `url:"per_page,omitempty"`
	Since   string `url:"since,omitempty"` // ISO 8601 date string
}

// ModeOption allows specifying the "mode" parameter (e.g., "extended").
type ModeOption struct {
	Mode string `url:"mode,omitempty"` // e.g., "extended"
}

// IDsOption is for endpoints that accept a list of IDs.
type IDsOption struct {
	IDs []int64 `url:"ids,comma"` // Will be sent as "ids=1,2,3"
}
