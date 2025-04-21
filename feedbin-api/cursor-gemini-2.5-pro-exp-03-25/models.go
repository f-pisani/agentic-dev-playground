package feedbin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIError represents an error returned by the Feedbin API.
// It includes the HTTP status code and the response body for more context.
type APIError struct {
	StatusCode int
	Message    string // Parsed error message, if available
	Body       string // Raw response body
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Feedbin API error: status %d, message: %s, body: %s", e.StatusCode, e.Message, e.Body)
}

// newAPIError creates a new APIError from an HTTP response.
func newAPIError(resp *http.Response) *APIError {
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	// Try to parse a standard error structure if possible, otherwise use the full body
	// Example: {"message": "error detail"}
	var jsonErr struct {
		Message string `json:"message"`
	}
	message := bodyString
	_ = json.Unmarshal(bodyBytes, &jsonErr) // Ignore error, if it's not JSON or doesn't match
	if jsonErr.Message != "" {
		message = jsonErr.Message
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    message,
		Body:       bodyString,
	}
}

// TimeRFC3339Nano is a wrapper around time.Time to handle Feedbin's specific date format.
type TimeRFC3339Nano time.Time

const feedbinTimeFormat = "2006-01-02T15:04:05.999999Z07:00"

func (t *TimeRFC3339Nano) UnmarshalJSON(b []byte) error {
	// Feedbin uses ISO 8601 with microsecond precision
	s := string(b)
	// Trim quotes
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	parsedTime, err := time.Parse(feedbinTimeFormat, s)
	if err != nil {
		// Fallback to standard RFC3339Nano if the custom one fails
		parsedTime, err = time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return fmt.Errorf("cannot parse time %q: %w", s, err)
		}
	}
	*t = TimeRFC3339Nano(parsedTime)
	return nil
}

func (t TimeRFC3339Nano) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(t).Format(feedbinTimeFormat) + "\""), nil
}

func (t TimeRFC3339Nano) Time() time.Time {
	return time.Time(t)
}

// Subscription represents a feed subscription.
type Subscription struct {
	ID        int64           `json:"id"`
	CreatedAt TimeRFC3339Nano `json:"created_at"`
	FeedID    int64           `json:"feed_id"`
	Title     string          `json:"title"`
	FeedURL   string          `json:"feed_url"`
	SiteURL   string          `json:"site_url"`
	// Extended mode fields
	JSONFeed *JSONFeedSubscription `json:"json_feed,omitempty"`
}

// JSONFeedSubscription represents the JSON Feed specific metadata within a Subscription.
type JSONFeedSubscription struct {
	Favicon     string `json:"favicon,omitempty"`
	FeedURL     string `json:"feed_url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
	HomePageURL string `json:"home_page_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

// FeedChoice represents one of the choices returned when subscribing to a URL
// that hosts multiple feeds.
type FeedChoice struct {
	FeedURL string `json:"feed_url"`
	Title   string `json:"title"`
}

// Entry represents a single feed entry.
type Entry struct {
	ID                  int64           `json:"id"`
	FeedID              int64           `json:"feed_id"`
	Title               *string         `json:"title"` // Pointer to handle potential null
	URL                 string          `json:"url"`
	ExtractedContentURL string          `json:"extracted_content_url,omitempty"`
	Author              *string         `json:"author"`  // Pointer to handle potential null
	Content             *string         `json:"content"` // Pointer to handle potential null
	Summary             string          `json:"summary"`
	Published           TimeRFC3339Nano `json:"published"`
	CreatedAt           TimeRFC3339Nano `json:"created_at"`
	// Extended mode / optional include fields
	Original          *OriginalEntry     `json:"original,omitempty"`
	Images            *Images            `json:"images,omitempty"`
	Enclosure         *Enclosure         `json:"enclosure,omitempty"`
	TwitterID         *int64             `json:"twitter_id,omitempty"`         // Pointer for optional field
	TwitterThreadIDs  []int64            `json:"twitter_thread_ids,omitempty"` // Slice, empty if not present
	ExtractedArticles []ExtractedArticle `json:"extracted_articles,omitempty"`
	JSONFeed          *JSONFeedEntry     `json:"json_feed,omitempty"`
	ContentDiff       *string            `json:"content_diff,omitempty"` // HTML diff if requested
}

// OriginalEntry represents the original state of an entry if it has been updated.
type OriginalEntry struct {
	Author    string          `json:"author"`
	Content   string          `json:"content"`
	Title     string          `json:"title"`
	URL       string          `json:"url"`
	EntryID   string          `json:"entry_id"`
	Published TimeRFC3339Nano `json:"published"`
	// Data map[string]interface{} `json:"data"` // Not strictly defined, using interface{}
}

// Images represents image data associated with an entry.
type Images struct {
	OriginalURL string     `json:"original_url"`
	Size1       *ImageSize `json:"size_1,omitempty"` // Assuming size_1 is the primary one, others might exist
	// Potentially other sizes like size_2, size_3 etc.
}

// ImageSize represents a specific size variant of an image.
type ImageSize struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Enclosure represents podcast/RSS enclosure data.
type Enclosure struct {
	EnclosureURL    string `json:"enclosure_url"`
	EnclosureType   string `json:"enclosure_type"`
	EnclosureLength string `json:"enclosure_length"` // Often a string in feeds
	ITunesDuration  string `json:"itunes_duration,omitempty"`
	ITunesImage     string `json:"itunes_image,omitempty"`
}

// ExtractedArticle represents article content extracted from links within an entry (e.g., a tweet).
type ExtractedArticle struct {
	URL     string  `json:"url"`
	Title   string  `json:"title"`
	Host    string  `json:"host"`
	Author  *string `json:"author"` // Pointer to handle potential null
	Content string  `json:"content"`
}

// JSONFeedEntry represents JSON Feed specific metadata within an Entry.
// Structure might vary based on actual JSON Feed implementation, this is a guess.
type JSONFeedEntry map[string]interface{}

// PaginationInfo holds pagination details extracted from response headers.
type PaginationInfo struct {
	NextPageURL  string
	PrevPageURL  string
	FirstPageURL string
	LastPageURL  string
	TotalRecords int // From X-Feedbin-Record-Count
}

// Helper functions for creating pointers to primitive types, useful for optional parameters.

// Bool returns a pointer to a bool value.
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to an int value.
func Int(v int) *int {
	return &v
}

// Int64 returns a pointer to an int64 value.
func Int64(v int64) *int64 {
	return &v
}

// String returns a pointer to a string value.
func String(v string) *string {
	return &v
}

// Time returns a pointer to a time.Time value.
func Time(v time.Time) *time.Time {
	return &v
}

// Tagging represents an association between a feed and a tag name.
type Tagging struct {
	ID     int64  `json:"id"`
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a user-defined search query.
type SavedSearch struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Icon represents a favicon associated with a host.
type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

// Import represents an OPML import job.
type Import struct {
	ID          int64           `json:"id"`
	Complete    bool            `json:"complete"`
	CreatedAt   TimeRFC3339Nano `json:"created_at"`
	ImportItems []ImportItem    `json:"import_items,omitempty"` // Only included when getting a specific import
}

// ImportItem represents the status of a single feed within an import job.
type ImportItem struct {
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	Status  string `json:"status"` // e.g., "pending", "complete", "failed"
}
