package feedbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

// TaggingsService handles operations related to entry taggings.
type TaggingsService struct {
	client *Client
}

// NewTaggingsService creates a new service for tagging related operations.
func NewTaggingsService(client *Client) *TaggingsService {
	return &TaggingsService{client: client}
}

// TaggingListOptions specifies optional parameters for listing taggings.
type TaggingListOptions struct {
	ListOptions // Embeds Page, PerPage, Since
}

// List retrieves all taggings. Each tagging object links an entry to a tag.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#get-taggings
func (s *TaggingsService) List(opts *TaggingListOptions) ([]Tagging, *http.Response, error) {
	path := "taggings.json"
	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		if params := v.Encode(); params != "" {
			path = fmt.Sprintf("%s?%s", path, params)
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var taggings []Tagging
	resp, err := s.client.Do(req, &taggings)
	if err != nil {
		return nil, resp, err
	}
	return taggings, resp, nil
}

// CreateTaggingOptions specifies the parameters for creating a new tagging.
type CreateTaggingOptions struct {
	EntryID int64  `json:"entry_id"`
	Name    string `json:"name"` // The name of the tag to apply
}

// Create applies a tag to an entry.
// Request body: {"entry_id": entry_id, "name": "tag_name"}
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#create-tagging
func (s *TaggingsService) Create(opts *CreateTaggingOptions) (*Tagging, *http.Response, error) {
	if opts == nil || opts.EntryID == 0 || opts.Name == "" {
		return nil, nil, fmt.Errorf("EntryID and Name are required to create a tagging")
	}
	path := "taggings.json"

	req, err := s.client.NewRequest(http.MethodPost, path, opts)
	if err != nil {
		return nil, nil, err
	}

	var createdTagging Tagging
	resp, err := s.client.Do(req, &createdTagging)
	if err != nil {
		return nil, resp, err
	}
	return &createdTagging, resp, nil
}

// Delete removes a tag from an entry using the tagging ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/taggings.md#delete-tagging
func (s *TaggingsService) Delete(taggingID int64) (*http.Response, error) {
	if taggingID == 0 {
		return nil, fmt.Errorf("taggingID is required to delete a tagging")
	}
	path := fmt.Sprintf("taggings/%d.json", taggingID)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil) // Expect 204 No Content
	if err != nil {
		return resp, err
	}
	return resp, nil
}
