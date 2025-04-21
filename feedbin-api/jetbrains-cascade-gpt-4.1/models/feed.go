package models

type Feed struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	SiteURL string `json:"site_url"`
}
