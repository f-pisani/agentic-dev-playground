package client

import (
	"fmt"
	"time"
)

// TagsService handles communication with the tag related
// methods of the Feedbin API
type TagsService struct {
	client *Client
}

// List returns all tags for the user
func (s *TagsService) List() ([]*Tag, error) {
	req, err := s.client.newRequest("GET", "tags.json", nil)
	if err != nil {
		return nil, err
	}

	var tags []*Tag
	_, err = s.client.do(req, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// ListWithSince returns all tags created after the specified time
func (s *TagsService) ListWithSince(since time.Time) ([]*Tag, error) {
	u := fmt.Sprintf("tags.json?since=%s", since.Format(time.RFC3339Nano))

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var tags []*Tag
	_, err = s.client.do(req, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Get returns a single tag
func (s *TagsService) Get(id int64) (*Tag, error) {
	u := fmt.Sprintf("tags/%d.json", id)

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tag := new(Tag)
	_, err = s.client.do(req, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// Delete removes a tag
func (s *TagsService) Delete(id int64) error {
	u := fmt.Sprintf("tags/%d.json", id)

	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
