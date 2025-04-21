package models

type Entry struct {
	ID                  int    `json:"id"`
	FeedID              int    `json:"feed_id"`
	Title               string `json:"title"`
	URL                 string `json:"url"`
	Author              string `json:"author"`
	Content             string `json:"content"`
	Published           string `json:"published"`
	CreatedAt           string `json:"created_at"`
	ExtractedContentURL string `json:"extracted_content_url"`
}
