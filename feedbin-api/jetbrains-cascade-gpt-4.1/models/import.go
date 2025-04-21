package models

type Import struct {
	ID          int64        `json:"id"`
	Complete    bool         `json:"complete"`
	CreatedAt   string       `json:"created_at"`
	ImportItems []ImportItem `json:"import_items,omitempty"`
}

type ImportItem struct {
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	Status  string `json:"status"`
}

type ImportsResponse []Import
