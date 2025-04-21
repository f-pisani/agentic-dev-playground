package feedbin

import (
	"net/http"
)

// PagesService handles operations related to pages
type PagesService struct {
	client *Client
}

// Create creates a new page from a URL
// If title is provided, it will be used if Feedbin cannot find the title of the content
func (s *PagesService) Create(url string, title string) (*Entry, error) {
	params := &PageParams{
		URL: url,
	}

	// Only include title if it's not empty
	if title != "" {
		params.Title = title
	}

	req, err := s.client.NewRequest(http.MethodPost, "/pages.json", params)
	if err != nil {
		return nil, err
	}

	entry := new(Entry)
	_, err = s.client.Do(req, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}
