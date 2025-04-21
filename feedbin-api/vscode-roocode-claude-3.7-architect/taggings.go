package feedbin

import (
	"fmt"
	"net/http"
)

// TaggingsService handles tagging-related operations
type TaggingsService struct {
	client *Client
}

// List returns all taggings
func (s *TaggingsService) List() ([]Tagging, error) {
	req, err := s.client.newRequest(http.MethodGet, "/taggings.json", nil)
	if err != nil {
		return nil, err
	}

	var taggings []Tagging
	_, err = s.client.do(req, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// Create creates a new tagging
func (s *TaggingsService) Create(feedID, tagID int) (*Tagging, error) {
	createReq := TaggingCreateRequest{
		FeedID: feedID,
		TagID:  tagID,
	}

	req, err := s.client.newRequest(http.MethodPost, "/taggings.json", createReq)
	if err != nil {
		return nil, err
	}

	tagging := new(Tagging)
	_, err = s.client.do(req, tagging)
	if err != nil {
		return nil, err
	}

	return tagging, nil
}

// Delete deletes a tagging
func (s *TaggingsService) Delete(id int) error {
	path := fmt.Sprintf("/taggings/%d.json", id)

	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
