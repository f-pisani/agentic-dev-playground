package feedbin

import (
	"net/http"
)

// PagesService handles communication with the pages related
// methods of the Feedbin API.
type PagesService struct {
	client *Client
}

// CreatePageOptions specifies the parameters to the
// PagesService.Create method.
type CreatePageOptions struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// Create creates a new page.
func (s *PagesService) Create(opts *CreatePageOptions) (*Entry, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/pages.json", opts)
	if err != nil {
		return nil, nil, err
	}

	entry := new(Entry)
	resp, err := s.client.Do(req, entry)
	if err != nil {
		return nil, resp, err
	}

	return entry, resp, nil
}
