package client

import (
	"fmt"
	"time"
)

// TaggingsService handles communication with the tagging related
// methods of the Feedbin API
type TaggingsService struct {
	client *Client
}

// List returns all taggings for the user
func (s *TaggingsService) List() ([]*Tagging, error) {
	req, err := s.client.newRequest("GET", "taggings.json", nil)
	if err != nil {
		return nil, err
	}

	var taggings []*Tagging
	_, err = s.client.do(req, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// ListWithSince returns all taggings created after the specified time
func (s *TaggingsService) ListWithSince(since time.Time) ([]*Tagging, error) {
	u := fmt.Sprintf("taggings.json?since=%s", since.Format(time.RFC3339Nano))

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var taggings []*Tagging
	_, err = s.client.do(req, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// Get returns a single tagging
func (s *TaggingsService) Get(id int64) (*Tagging, error) {
	u := fmt.Sprintf("taggings/%d.json", id)

	req, err := s.client.newRequest("GET", u, nil)
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

// Create adds a tag to a feed
func (s *TaggingsService) Create(feedID int64, name string) (*Tagging, error) {
	type taggingRequest struct {
		FeedID int64  `json:"feed_id"`
		Name   string `json:"name"`
	}

	req, err := s.client.newRequest("POST", "taggings.json", &taggingRequest{
		FeedID: feedID,
		Name:   name,
	})
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

// Delete removes a tagging
func (s *TaggingsService) Delete(id int64) error {
	u := fmt.Sprintf("taggings/%d.json", id)

	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
