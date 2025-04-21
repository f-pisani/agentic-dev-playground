package feedbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TaggingsService handles communication with the tagging related
// methods of the Feedbin API.
type TaggingsService service

// List retrieves all taggings for the authenticated user.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#get-taggings
func (s *TaggingsService) List() ([]*Tagging, *http.Response, error) {
	req, err := s.client.newRequest(http.MethodGet, "taggings.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var taggings []*Tagging
	resp, err := s.client.do(req, &taggings)
	if err != nil {
		return nil, resp, err
	}

	return taggings, resp, nil
}

// taggingRequest is used for Create requests.
type taggingRequest struct {
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// Create applies a tag to a specific feed.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#create-taggings
func (s *TaggingsService) Create(feedID int64, tagName string) (*Tagging, *http.Response, error) {
	body := taggingRequest{FeedID: feedID, Name: tagName}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPost, "taggings.json", buf)
	if err != nil {
		return nil, nil, err
	}

	var tagging Tagging
	resp, err := s.client.do(req, &tagging)
	if err != nil {
		return nil, resp, err
	}

	return &tagging, resp, nil
}

// Delete removes a tagging (removes a tag from a feed) by its ID.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#delete-taggings
func (s *TaggingsService) Delete(taggingID int64) (*http.Response, error) {
	path := fmt.Sprintf("taggings/%d.json", taggingID)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	// Expect 204 No Content
	resp, err := s.client.do(req, nil)
	return resp, err
}
