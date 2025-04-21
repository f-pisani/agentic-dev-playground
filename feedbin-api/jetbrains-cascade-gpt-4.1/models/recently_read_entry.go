package models

type RecentlyReadEntry struct {
	ID int64 `json:"id"`
}

type RecentlyReadEntriesResponse []RecentlyReadEntry
