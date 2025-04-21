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

// CreateTaggingRequest represents a request to create a tagging
type CreateTaggingRequest struct {
	FeedID int    `json:"feed_id"`
	TagID  int    `json:"tag_id,omitempty"`
	Name   string `json:"name,omitempty"`
}

// List returns all taggings
// https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#get-taggings
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

// Create adds a new tag to a feed
// https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#create-tagging
func (s *TaggingsService) Create(feedID, tagID int, name string) (*Tagging, *http.Response, error) {
	request := &CreateTaggingRequest{
		FeedID: feedID,
	}

	// Either tagID or name should be set, not both
	if tagID > 0 {
		request.TagID = tagID
	} else if name != "" {
		request.Name = name
	} else {
		return nil, nil, fmt.Errorf("either tagID or name must be provided")
	}

	req, err := s.client.NewRequest(http.MethodPost, "taggings.json", request)
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

// Delete removes a tagging
// https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#delete-tagging
func (s *TaggingsService) Delete(id int) (*http.Response, error) {
	url := fmt.Sprintf("taggings/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
