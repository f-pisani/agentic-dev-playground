package feedbin

import (
	"fmt"
	"net/http"
)

// TaggingsService handles operations related to taggings
type TaggingsService struct {
	client *Client
}

// List returns all taggings
func (s *TaggingsService) List() ([]Tagging, error) {
	req, err := s.client.NewRequest(http.MethodGet, "taggings.json", nil)
	if err != nil {
		return nil, err
	}

	var taggings []Tagging
	_, err = s.client.Do(req, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// Get returns a single tagging by ID
func (s *TaggingsService) Get(id int) (*Tagging, error) {
	path := fmt.Sprintf("/taggings/%d.json", id)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	tagging := new(Tagging)
	_, err = s.client.Do(req, tagging)
	if err != nil {
		return nil, err
	}

	return tagging, nil
}

// Create creates a new tagging
func (s *TaggingsService) Create(feedID int, name string) (*Tagging, error) {
	params := &TaggingParams{
		FeedID: feedID,
		Name:   name,
	}

	req, err := s.client.NewRequest(http.MethodPost, "/taggings.json", params)
	if err != nil {
		return nil, err
	}

	tagging := new(Tagging)
	_, err = s.client.Do(req, tagging)
	if err != nil {
		return nil, err
	}

	return tagging, nil
}

// Delete deletes a tagging
func (s *TaggingsService) Delete(id int) error {
	path := fmt.Sprintf("/taggings/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}
