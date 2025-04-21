package models

type Subscription struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	FeedID    int64  `json:"feed_id"`
	Title     string `json:"title"`
	FeedURL   string `json:"feed_url"`
	SiteURL   string `json:"site_url"`
}

type SubscriptionsResponse []Subscription
