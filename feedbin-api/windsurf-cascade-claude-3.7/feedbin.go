// Package feedbin provides a client for the Feedbin REST API v2.
package feedbin

// Feedbin is the main client for the Feedbin API
type Feedbin struct {
	// Client is the underlying HTTP client
	Client *Client

	// Authentication provides methods for authentication
	Authentication *AuthenticationService

	// Subscriptions provides methods for managing subscriptions
	Subscriptions *SubscriptionsService

	// Entries provides methods for accessing entries
	Entries *EntriesService

	// UnreadEntries provides methods for managing unread entries
	UnreadEntries *UnreadEntriesService

	// StarredEntries provides methods for managing starred entries
	StarredEntries *StarredEntriesService

	// Taggings provides methods for managing taggings
	Taggings *TaggingsService

	// Tags provides methods for managing tags
	Tags *TagsService

	// SavedSearches provides methods for managing saved searches
	SavedSearches *SavedSearchesService

	// RecentlyReadEntries provides methods for managing recently read entries
	RecentlyReadEntries *RecentlyReadEntriesService

	// UpdatedEntries provides methods for managing updated entries
	UpdatedEntries *UpdatedEntriesService

	// Icons provides methods for accessing feed icons
	Icons *IconsService

	// Imports provides methods for managing OPML imports
	Imports *ImportsService

	// Pages provides methods for creating pages
	Pages *PagesService
}

// New creates a new Feedbin API client
func New(username, password string) *Feedbin {
	client := NewClient(username, password)

	return &Feedbin{
		Client:              client,
		Authentication:      &AuthenticationService{client: client},
		Subscriptions:       &SubscriptionsService{client: client},
		Entries:             &EntriesService{client: client},
		UnreadEntries:       &UnreadEntriesService{client: client},
		StarredEntries:      &StarredEntriesService{client: client},
		Taggings:            &TaggingsService{client: client},
		Tags:                &TagsService{client: client},
		SavedSearches:       &SavedSearchesService{client: client},
		RecentlyReadEntries: &RecentlyReadEntriesService{client: client},
		UpdatedEntries:      &UpdatedEntriesService{client: client},
		Icons:               &IconsService{client: client},
		Imports:             &ImportsService{client: client},
		Pages:               &PagesService{client: client},
	}
}
