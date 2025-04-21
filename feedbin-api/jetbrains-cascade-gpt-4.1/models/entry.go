package models

type Entry struct {
	ID                  int64   `json:"id"`
	FeedID              int64   `json:"feed_id"`
	Title               *string `json:"title"`
	URL                 string  `json:"url"`
	ExtractedContentURL string  `json:"extracted_content_url"`
	Author              *string `json:"author"`
	Content             *string `json:"content"`
	Summary             string  `json:"summary"`
	Published           string  `json:"published"`
	CreatedAt           string  `json:"created_at"`
}

type EntriesResponse []Entry
