package models

type Tagging struct {
	ID        int64  `json:"id"`
	FeedID    int64  `json:"feed_id"`
	TagID     int64  `json:"tag_id"`
	CreatedAt string `json:"created_at"`
}

type TaggingsResponse []Tagging
