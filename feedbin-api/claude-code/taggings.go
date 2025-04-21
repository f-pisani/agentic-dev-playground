package feedbin

import (
	"fmt"
	"net/http"
)

// TaggingsService handles communication with the taggings related
// endpoints of the Feedbin API
type TaggingsService struct {
	client *Client
}

// TaggingCreateOptions specifies the parameters to the
// TaggingsService.Create method
type TaggingCreateOptions struct {
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// List returns all taggings for the authenticated user
func (s *TaggingsService) List() ([]*Tagging, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "taggings.json", nil)
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

// Create creates a new tagging (assigns a tag to a feed)
func (s *TaggingsService) Create(opts *TaggingCreateOptions) (*Tagging, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "taggings.json", opts)
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

// Delete deletes a tagging
func (s *TaggingsService) Delete(id int) (*http.Response, error) {
	url := fmt.Sprintf("taggings/%d.json", id)
	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
