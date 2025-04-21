package models

type UnreadEntry struct {
	ID int64 `json:"id"`
}

type UnreadEntriesResponse []UnreadEntry
