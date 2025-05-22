package feedbinapi

import "time"

// Entry represents a single feed entry.
type Entry struct {
	ID                  int64              `json:"id"`
	FeedID              int64              `json:"feed_id"`
	Title               *string            `json:"title"`
	URL                 *string            `json:"url"`
	ExtractedContentURL *string            `json:"extracted_content_url,omitempty"`
	Author              *string            `json:"author"`
	Content             *string            `json:"content"` // HTML content
	Summary             *string            `json:"summary"`
	Published           time.Time          `json:"published"`
	CreatedAt           time.Time          `json:"created_at"`
	Original            *OriginalEntry     `json:"original,omitempty"`           // Included with include_original=true
	Images              *EntryImages       `json:"images,omitempty"`             // Extended mode
	Enclosure           *EntryEnclosure    `json:"enclosure,omitempty"`          // Extended mode or include_enclosure=true
	TwitterID           *int64             `json:"twitter_id,omitempty"`         // Extended mode
	TwitterThreadIDs    []int64            `json:"twitter_thread_ids,omitempty"` // Extended mode
	ExtractedArticles   []ExtractedArticle `json:"extracted_articles,omitempty"` // Extended mode
	JSONFeed            *JSONFeedData      `json:"json_feed,omitempty"`          // Extended mode
	ContentDiff         *string            `json:"content_diff,omitempty"`       // Included with include_content_diff=true
}

// OriginalEntry represents the original version of an entry if it has been updated.
type OriginalEntry struct {
	Author    *string   `json:"author"`
	Content   *string   `json:"content"`
	Title     *string   `json:"title"`
	URL       *string   `json:"url"`
	EntryID   string    `json:"entry_id"` // Note: This is a string in the example
	Published time.Time `json:"published"`
	Data      *struct{} `json:"data"` // Example shows null or empty object, define if structure known
}

// EntryImages represents images associated with an entry.
type EntryImages struct {
	OriginalURL *string     `json:"original_url"`
	Size1       *SizedImage `json:"size_1,omitempty"` // Example shows "size_1", there might be others
	// Add other sizes if the API defines them (e.g., Size2, Size3)
}

// SizedImage represents a specific size of an image with CDN URL.
type SizedImage struct {
	CDNURL *string `json:"cdn_url"`
	Width  *int    `json:"width"`
	Height *int    `json:"height"`
}

// EntryEnclosure represents podcast/RSS enclosure data.
type EntryEnclosure struct {
	URL            *string `json:"enclosure_url"`
	Type           *string `json:"enclosure_type"`
	Length         *string `json:"enclosure_length"` // String because example "54103635" might be large for int32
	ITunesDuration *string `json:"itunes_duration,omitempty"`
	ITunesImage    *string `json:"itunes_image,omitempty"`
}

// ExtractedArticle represents an article extracted from a link in an entry (e.g., a tweet).
type ExtractedArticle struct {
	URL     *string `json:"url"`
	Title   *string `json:"title"`
	Host    *string `json:"host"`
	Author  *string `json:"author"`
	Content *string `json:"content"` // HTML content
}

// JSONFeedData holds additional metadata if the entry is from a JSON Feed.
// This structure might be shared between Entry and Subscription.
type JSONFeedData struct {
	Version     *string   `json:"version,omitempty"`
	UserComment *string   `json:"user_comment,omitempty"`
	NextURL     *string   `json:"next_url,omitempty"`
	Icon        *string   `json:"icon,omitempty"`
	Favicon     *string   `json:"favicon,omitempty"`
	Author      *struct { // JSONFeed author object
		Name   *string `json:"name,omitempty"`
		URL    *string `json:"url,omitempty"`
		Avatar *string `json:"avatar,omitempty"`
	} `json:"author,omitempty"`
	Expired *bool `json:"expired,omitempty"`
	Hubs    []struct {
		Type *string `json:"type,omitempty"`
		URL  *string `json:"url,omitempty"`
	} `json:"hubs,omitempty"`
	// Fields specific to Entry's JSONFeed data (if any)
	FeedURL     *string `json:"feed_url,omitempty"`      // Present in Subscription's JSONFeed
	HomePageURL *string `json:"home_page_url,omitempty"` // Present in Subscription's JSONFeed
	Title       *string `json:"title,omitempty"`         // Present in Subscription's JSONFeed
}

// Feed represents a single feed.
type Feed struct {
	ID      int64   `json:"id"`
	Title   *string `json:"title"`
	FeedURL *string `json:"feed_url"`
	SiteURL *string `json:"site_url"`
}

// Subscription represents a user's subscription to a feed.
type Subscription struct {
	ID        int64         `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	FeedID    int64         `json:"feed_id"`
	Title     string        `json:"title"`               // Docs show this as required
	FeedURL   string        `json:"feed_url"`            // Docs show this as required
	SiteURL   string        `json:"site_url"`            // Docs show this as required
	JSONFeed  *JSONFeedData `json:"json_feed,omitempty"` // Extended mode
}

// Tagging represents a tag applied to a feed.
type Tagging struct {
	ID     int64  `json:"id"`
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a user's saved search query.
type SavedSearch struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Import represents an OPML import process.
type Import struct {
	ID          int64         `json:"id"`
	Complete    bool          `json:"complete"`
	CreatedAt   time.Time     `json:"created_at"`
	ImportItems []*ImportItem `json:"import_items,omitempty"` // Only present when getting a specific import
}

// ImportItem represents a single item within an OPML import.
type ImportItem struct {
	Title   *string `json:"title"`
	FeedURL *string `json:"feed_url"`
	Status  *string `json:"status"` // "pending", "complete", "failed"
}

// Icon represents a favicon for a feed.
type Icon struct {
	Host *string `json:"host"`
	URL  *string `json:"url"`
}

// PageCreateRequest is used to create a new page (entry from URL).
type PageCreateRequest struct {
	URL   string  `json:"url"`
	Title *string `json:"title,omitempty"`
}

// --- Request/Response Structs for specific actions ---

// StarredEntryRequest is used for POST/DELETE bodies for starred entries.
type StarredEntryRequest struct {
	StarredEntries []int64 `json:"starred_entries"`
}

// UnreadEntryRequest is used for POST/DELETE bodies for unread entries.
type UnreadEntryRequest struct {
	UnreadEntries []int64 `json:"unread_entries"`
}

// RecentlyReadEntryRequest is used for POST bodies for recently_read entries.
type RecentlyReadEntryRequest struct {
	RecentlyReadEntries []int64 `json:"recently_read_entries"`
}

// UpdatedEntryRequest is used for DELETE bodies for updated_entries.
type UpdatedEntryRequest struct {
	UpdatedEntries []int64 `json:"updated_entries"`
}

// CreateSubscriptionRequest is used to create a new subscription.
type CreateSubscriptionRequest struct {
	FeedURL string `json:"feed_url"`
}

// UpdateSubscriptionRequest is used to update a subscription (e.g., title).
type UpdateSubscriptionRequest struct {
	Title string `json:"title"`
}

// CreateTaggingRequest is used to create a new tagging.
type CreateTaggingRequest struct {
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// RenameTagRequest is used to rename a tag.
type RenameTagRequest struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// DeleteTagRequest is used to delete a tag.
type DeleteTagRequest struct {
	Name string `json:"name"`
}

// CreateSavedSearchRequest is used to create a saved search.
type CreateSavedSearchRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// UpdateSavedSearchRequest is used to update a saved search.
type UpdateSavedSearchRequest struct {
	Name *string `json:"name,omitempty"` // Only name can be updated based on docs
}

// --- Parameter Structs for List operations ---

// ListEntriesParams holds parameters for listing entries.
type ListEntriesParams struct {
	Page               *int    `url:"page,omitempty"`
	Since              *string `url:"since,omitempty"` // ISO 8601 date
	IDs                []int64 `url:"ids,omitempty,comma"`
	Read               *bool   `url:"read,omitempty"`
	Starred            *bool   `url:"starred,omitempty"`
	PerPage            *int    `url:"per_page,omitempty"`
	Mode               *string `url:"mode,omitempty"` // "extended"
	IncludeOriginal    *bool   `url:"include_original,omitempty"`
	IncludeEnclosure   *bool   `url:"include_enclosure,omitempty"`
	IncludeContentDiff *bool   `url:"include_content_diff,omitempty"`
}

// ListSubscriptionsParams holds parameters for listing subscriptions.
type ListSubscriptionsParams struct {
	Since *string `url:"since,omitempty"` // ISO 8601 date
	Mode  *string `url:"mode,omitempty"`  // "extended"
}

// GetSubscriptionParams holds parameters for getting a single subscription.
type GetSubscriptionParams struct {
	Mode *string `url:"mode,omitempty"` // "extended"
}

// GetSavedSearchParams holds parameters for getting a saved search's entries.
type GetSavedSearchParams struct {
	IncludeEntries *bool `url:"include_entries,omitempty"`
	Page           *int  `url:"page,omitempty"` // Only if include_entries=true
}

// ListUpdatedEntryIDsParams holds parameters for listing updated entry IDs.
type ListUpdatedEntryIDsParams struct {
	Since *string `url:"since,omitempty"` // ISO 8601 date
}

// GetEntryParams holds parameters for getting a single entry.
// Similar to ListEntriesParams but contextually for a single entry.
type GetEntryParams struct {
	Mode               *string `url:"mode,omitempty"` // "extended"
	IncludeOriginal    *bool   `url:"include_original,omitempty"`
	IncludeEnclosure   *bool   `url:"include_enclosure,omitempty"`
	IncludeContentDiff *bool   `url:"include_content_diff,omitempty"`
}

// MultipleFeedChoice represents one choice when feed discovery yields multiple feeds.
type MultipleFeedChoice struct {
	FeedURL string `json:"feed_url"`
	Title   string `json:"title"`
}

// Page represents a single "page" resource (entry created from a URL)
// The API returns an Entry struct for a page, so we can reuse Entry.
// However, the create operation uses PageCreateRequest.
// No specific Page struct is needed for responses if it's identical to Entry.

// TaggingsOnEntry represents the tags associated with a specific entry.
// The API returns an array of strings (tag names).
// Example: ["tag1", "tag2"]
// So, a `[]string` will suffice for this.

// PresignedS3Upload represents the response from POST /v2/presigned_s3_uploads.json
// This is used for OPML import.
type PresignedS3Upload struct {
	URL            string            `json:"url"`            // The URL to PUT the file to
	Fields         map[string]string `json:"fields"`         // Form fields to include in the PUT request
	Path           string            `json:"path"`           // The path of the file on S3 (used to create import)
	Method         string            `json:"method"`         // Should be "put"
	ACL            string            `json:"acl"`            // e.g., "private"
	Key            string            `json:"key"`            // The key (filename) on S3
	AWSAccessKeyID string            `json:"AWSAccessKeyId"` // Note the casing
	Policy         string            `json:"Policy"`
	Signature      string            `json:"Signature"`
	ContentType    string            `json:"Content-Type"`
}

// CreateImportRequest is used to create an import after uploading an OPML file.
type CreateImportRequest struct {
	Path string `json:"path"` // The 'path' from the PresignedS3Upload response
}
