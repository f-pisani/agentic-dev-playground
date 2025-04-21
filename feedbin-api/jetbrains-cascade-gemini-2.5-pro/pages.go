package feedbin

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// PagesService handles communication with the page related
// methods of the Feedbin API (for saving webpages).
type PagesService service

// createPageRequest is the structure for the POST request body.
type createPageRequest struct {
	URL   string  `json:"url"`
	Title *string `json:"title,omitempty"` // Optional
}

// Create saves a webpage URL to Feedbin, returning it as an Entry.
// The title is optional and only used if Feedbin cannot determine it from the content.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/pages.md#post-v2pagesjson
func (s *PagesService) Create(pageURL string, title *string) (*Entry, *http.Response, error) {
	body := createPageRequest{
		URL:   pageURL,
		Title: title,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPost, "pages.json", buf)
	if err != nil {
		return nil, nil, err
	}

	var entry Entry
	resp, err := s.client.do(req, &entry)
	if err != nil {
		return nil, resp, err
	}

	return &entry, resp, nil
}
