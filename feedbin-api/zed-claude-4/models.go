package feedbin

import "time"

// Subscription represents a feed subscription
type Subscription struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int       `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"`
}

// JSONFeed contains additional metadata for subscriptions in extended mode
type JSONFeed struct {
	Favicon     string `json:"favicon,omitempty"`
	FeedURL     string `json:"feed_url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
	HomePageURL string `json:"home_page_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

// Entry represents a feed entry
type Entry struct {
	ID                  int        `json:"id"`
	FeedID              int        `json:"feed_id"`
	Title               *string    `json:"title"`
	URL                 string     `json:"url"`
	ExtractedContentURL string     `json:"extracted_content_url"`
	Author              *string    `json:"author"`
	Content             *string    `json:"content"`
	Summary             string     `json:"summary"`
	Published           time.Time  `json:"published"`
	CreatedAt           time.Time  `json:"created_at"`
	Original            *Original  `json:"original,omitempty"`
	TwitterID           *string    `json:"twitter_id,omitempty"`
	TwitterThreadIDs    []string   `json:"twitter_thread_ids,omitempty"`
	Images              []Image    `json:"images,omitempty"`
	Enclosure           *Enclosure `json:"enclosure,omitempty"`
	ExtractedArticles   []Article  `json:"extracted_articles,omitempty"`
}

// Original contains original entry data if entry has been updated
type Original struct {
	Author    *string   `json:"author"`
	Content   *string   `json:"content"`
	Title     *string   `json:"title"`
	URL       string    `json:"url"`
	EntryID   int       `json:"entry_id"`
	Published time.Time `json:"published"`
	Data      string    `json:"data"`
}

// Image represents image metadata for entries
type Image struct {
	OriginalURL string `json:"original_url"`
	Size1       *Size  `json:"size_1,omitempty"`
}

// Size represents image size information
type Size struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Enclosure represents podcast/RSS enclosure data
type Enclosure struct {
	EnclosureURL    string `json:"enclosure_url"`
	EnclosureType   string `json:"enclosure_type"`
	EnclosureLength *int   `json:"enclosure_length,omitempty"`
	ItunesDuration  string `json:"itunes_duration,omitempty"`
	ItunesImage     string `json:"itunes_image,omitempty"`
}

// Article represents extracted article content
type Article struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Host    string `json:"host"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

// ExtractedContent represents Mercury Parser response format
type ExtractedContent struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	DatePublished time.Time `json:"date_published"`
	LeadImageURL  *string   `json:"lead_image_url"`
	Dek           *string   `json:"dek"`
	NextPageURL   *string   `json:"next_page_url"`
	URL           string    `json:"url"`
	Domain        string    `json:"domain"`
	Excerpt       string    `json:"excerpt"`
	WordCount     int       `json:"word_count"`
	Direction     string    `json:"direction"`
	TotalPages    int       `json:"total_pages"`
	RenderedPages int       `json:"rendered_pages"`
}

// Tagging represents a feed tag assignment
type Tagging struct {
	ID     int    `json:"id"`
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// FeedChoice represents a feed option when multiple feeds are found
type FeedChoice struct {
	FeedURL string `json:"feed_url"`
	Title   string `json:"title"`
}

// SubscriptionOptions contains options for listing subscriptions
type SubscriptionOptions struct {
	Since *time.Time
	Mode  *string
}

// EntryOptions contains options for listing entries
type EntryOptions struct {
	Page               *int
	Since              *time.Time
	IDs                []int
	Read               *bool
	Starred            *bool
	PerPage            *int
	Mode               *string
	IncludeOriginal    *bool
	IncludeEnclosure   *bool
	IncludeContentDiff *bool
}

// CreateSubscriptionRequest represents a subscription creation request
type CreateSubscriptionRequest struct {
	FeedURL string `json:"feed_url"`
}

// UpdateSubscriptionRequest represents a subscription update request
type UpdateSubscriptionRequest struct {
	Title string `json:"title"`
}

// UnreadEntriesRequest represents a request to mark entries as unread/read
type UnreadEntriesRequest struct {
	UnreadEntries []int `json:"unread_entries"`
}

// StarredEntriesRequest represents a request to star/unstar entries
type StarredEntriesRequest struct {
	StarredEntries []int `json:"starred_entries"`
}

// CreateTaggingRequest represents a tagging creation request
type CreateTaggingRequest struct {
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// RenameTagRequest represents a tag rename request
type RenameTagRequest struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// DeleteTagRequest represents a tag deletion request
type DeleteTagRequest struct {
	Name string `json:"name"`
}

// Feed represents a feed object
type Feed struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	SiteURL string `json:"site_url"`
}

// SavedSearch represents a saved search
type SavedSearch struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// CreateSavedSearchRequest represents a saved search creation request
type CreateSavedSearchRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// UpdateSavedSearchRequest represents a saved search update request
type UpdateSavedSearchRequest struct {
	Name string `json:"name"`
}

// SavedSearchOptions contains options for saved search entry retrieval
type SavedSearchOptions struct {
	IncludeEntries *bool
	Page           *int
}

// RecentlyReadEntriesRequest represents a request to mark entries as recently read
type RecentlyReadEntriesRequest struct {
	RecentlyReadEntries []int `json:"recently_read_entries"`
}

// UpdatedEntriesRequest represents a request to mark updated entries as read
type UpdatedEntriesRequest struct {
	UpdatedEntries []int `json:"updated_entries"`
}

// UpdatedEntriesOptions contains options for updated entries
type UpdatedEntriesOptions struct {
	Since *time.Time
}

// Icon represents a feed icon
type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

// Import represents an OPML import
type Import struct {
	ID          int          `json:"id"`
	Complete    bool         `json:"complete"`
	CreatedAt   time.Time    `json:"created_at"`
	ImportItems []ImportItem `json:"import_items,omitempty"`
}

// ImportItem represents an individual item in an import
type ImportItem struct {
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	Status  string `json:"status"`
}

// CreatePageRequest represents a page creation request
type CreatePageRequest struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// PaginationInfo represents pagination information from Link headers
type PaginationInfo struct {
	First    string
	Previous string
	Next     string
	Last     string
	Total    int
}

// HasNext returns true if there is a next page
func (p *PaginationInfo) HasNext() bool {
	return p.Next != ""
}

// HasPrevious returns true if there is a previous page
func (p *PaginationInfo) HasPrevious() bool {
	return p.Previous != ""
}

// PaginatedResponse wraps responses that support pagination
type PaginatedResponse struct {
	Pagination *PaginationInfo
}

// Helper functions for creating pointers to basic types for optional fields
func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

func Bool(b bool) *bool {
	return &b
}

func Time(t time.Time) *time.Time {
	return &t
}
