// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
)

// TaggingService handles communication with the tagging related
// methods of the Feedbin API.
type TaggingService struct {
	client *Client
}

// GetTaggings retrieves all taggings.
func (s *TaggingService) GetTaggings() ([]Tagging, error) {
	req, err := s.client.NewRequest("GET", "/taggings.json", nil)
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

// CreateTagging creates a new tagging.
func (s *TaggingService) CreateTagging(feedID int, name string) (*Tagging, error) {
	body := map[string]interface{}{
		"feed_id": feedID,
		"name":    name,
	}

	req, err := s.client.NewRequest("POST", "/taggings.json", body)
	if err != nil {
		return nil, err
	}

	var tagging Tagging
	_, err = s.client.Do(req, &tagging)
	if err != nil {
		return nil, err
	}

	return &tagging, nil
}

// DeleteTagging deletes a tagging.
func (s *TaggingService) DeleteTagging(id int) error {
	path := fmt.Sprintf("/taggings/%d.json", id)
	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}