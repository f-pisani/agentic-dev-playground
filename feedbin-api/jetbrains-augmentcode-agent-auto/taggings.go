package feedbin

import (
	"fmt"
	"net/http"
)

// TaggingsService handles communication with the taggings related
// methods of the Feedbin API.
type TaggingsService struct {
	client *Client
}

// List returns all taggings.
func (s *TaggingsService) List() ([]*Tagging, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/v2/taggings.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var taggings []*Tagging
	resp, err := s.client.Do(req, &taggings)
	if err != nil {
		return nil, resp, err
	}

	return taggings, resp, nil
}

// Get returns a single tagging.
func (s *TaggingsService) Get(id int64) (*Tagging, *http.Response, error) {
	u := fmt.Sprintf("/v2/taggings/%d.json", id)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	tagging := new(Tagging)
	resp, err := s.client.Do(req, tagging)
	if err != nil {
		return nil, resp, err
	}

	return tagging, resp, nil
}

// CreateTaggingOptions specifies the parameters to the
// TaggingsService.Create method.
type CreateTaggingOptions struct {
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// Create creates a new tagging.
func (s *TaggingsService) Create(opts *CreateTaggingOptions) (*Tagging, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/taggings.json", opts)
	if err != nil {
		return nil, nil, err
	}

	tagging := new(Tagging)
	resp, err := s.client.Do(req, tagging)
	if err != nil {
		// Check if we got a 302 Found (tagging exists)
		if errResp, ok := err.(*ErrorResponse); ok {
			if errResp.Response.StatusCode == http.StatusFound {
				// Get the tagging from the Location header
				location := errResp.Response.Header.Get("Location")
				if location != "" {
					req, err := s.client.NewRequest(http.MethodGet, location, nil)
					if err != nil {
						return nil, nil, err
					}

					resp, err := s.client.Do(req, tagging)
					if err != nil {
						return nil, resp, err
					}

					return tagging, resp, nil
				}
			}
		}

		return nil, resp, err
	}

	return tagging, resp, nil
}

// Delete deletes a tagging.
func (s *TaggingsService) Delete(id int64) (*http.Response, error) {
	u := fmt.Sprintf("/v2/taggings/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
