package models

type StarredEntry struct {
	ID int64 `json:"id"`
}

type StarredEntriesResponse []StarredEntry
