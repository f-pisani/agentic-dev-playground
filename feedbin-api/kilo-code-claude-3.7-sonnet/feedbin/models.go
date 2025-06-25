package feedbin

import "time"

type Entry struct {
	ID                  int64               `json:"id"`
	FeedID              int64               `json:"feed_id"`
	Title               *string             `json:"title"`
	URL                 string              `json:"url"`
	Author              *string             `json:"author"`
	Content             *string             `json:"content"`
	Summary             string              `json:"summary"`
	Published           time.Time           `json:"published"`
	CreatedAt           time.Time           `json:"created_at"`
	ExtractedContentURL string              `json:"extracted_content_url"`
	Original            *Entry              `json:"original,omitempty"`
	Images              *Images             `json:"images,omitempty"`
	Enclosure           *Enclosure          `json:"enclosure,omitempty"`
	TwitterID           *int64              `json:"twitter_id,omitempty"`
	TwitterThreadIDs    []int64             `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles   []*ExtractedArticle `json:"extracted_articles,omitempty"`
	JSONFeed            *map[string]any     `json:"json_feed,omitempty"`
}

type Feed struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	SiteURL string `json:"site_url"`
}

type Subscription struct {
	ID        int64           `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	FeedID    int64           `json:"feed_id"`
	Title     string          `json:"title"`
	FeedURL   string          `json:"feed_url"`
	SiteURL   string          `json:"site_url"`
	JSONFeed  *map[string]any `json:"json_feed,omitempty"`
}

type Tagging struct {
	ID     int64  `json:"id"`
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

type SavedSearch struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

type Images struct {
	OriginalURL string    `json:"original_url"`
	Size1       ImageSize `json:"size_1"`
}

type ImageSize struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Enclosure struct {
	URL      string `json:"enclosure_url"`
	Type     string `json:"enclosure_type"`
	Length   string `json:"enclosure_length"`
	Duration string `json:"itunes_duration"`
	Image    string `json:"itunes_image"`
}

type ExtractedArticle struct {
	URL     string  `json:"url"`
	Title   string  `json:"title"`
	Host    string  `json:"host"`
	Author  *string `json:"author"`
	Content string  `json:"content"`
}

type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

type Import struct {
	ID          int64         `json:"id"`
	Complete    bool          `json:"complete"`
	CreatedAt   time.Time     `json:"created_at"`
	ImportItems []*ImportItem `json:"import_items,omitempty"`
}

type ImportItem struct {
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	Status  string `json:"status"`
}

type Page struct {
	URL   string  `json:"url"`
	Title *string `json:"title,omitempty"`
}

type PageContent struct {
	Content string `json:"content"`
}

type TwitterAccount struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	ScreenName      string    `json:"screen_name"`
	ProfileImageURL string    `json:"profile_image_url"`
}

type Tweet struct {
	ID         int64     `json:"id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
	ScreenName string    `json:"screen_name"`
	TweetID    int64     `json:"tweet_id"`
}