package feedbin

import (
	"fmt"
	"net/http"

	// "github.com/google/go-querystring/query" // Not needed if no options
)

// PagesService handles fetching processed page content for entries.
type PagesService struct {
	client *Client
}

// NewPagesService creates a new service for page operations.
func NewPagesService(client *Client) *PagesService {
	return &PagesService{client: client}
}

// PageGetOptions specifies optional parameters for getting a page.
// The API doc for "GET /v2/pages/{entry_id}.json" does not specify any query parameters.
type PageGetOptions struct {
	// No options defined in spec for this endpoint.
}

// Get retrieves the processed page content for a given entry ID.
// The response is a Page object.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/pages.md#get-page
func (s *PagesService) Get(entryID int64, opts *PageGetOptions) (*Page, *http.Response, error) {
	if entryID == 0 {
		return nil, nil, fmt.Errorf("entryID is required to get a page")
	}

	path := fmt.Sprintf("pages/%d.json", entryID)
	// No query options specified in docs for this endpoint.
	// if opts != nil {
	//  v, err := query.Values(opts)
	//  if err != nil {
	//   return nil, nil, err
	//  }
	//  if params := v.Encode(); params != "" {
	//   path = fmt.Sprintf("%s?%s", path, params)
	//  }
	// }

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var page Page
	resp, err := s.client.Do(req, &page)
	if err != nil {
		return nil, resp, err
	}
	return &page, resp, nil
}
