package feedbin

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// TagsService handles communication with the tag related
// methods of the Feedbin API (renaming and deleting).
type TagsService service

// renameTagRequest is used for the Rename request.
type renameTagRequest struct {
	OldName string `json:"from"`
	NewName string `json:"to"`
}

// Rename changes the name of an existing tag across all feeds it's applied to.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/tags.md#rename-tag
func (s *TagsService) Rename(oldName, newName string) (*http.Response, error) {
	body := renameTagRequest{OldName: oldName, NewName: newName}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest(http.MethodPatch, "tags/rename.json", buf)
	if err != nil {
		return nil, err
	}

	// Expect 204 No Content
	resp, err := s.client.do(req, nil)
	return resp, err
}

// deleteTagRequest is used for the Delete request.
type deleteTagRequest struct {
	Name string `json:"name"`
}

// Delete removes a tag entirely from all feeds.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/tags.md#delete-tag
func (s *TagsService) Delete(tagName string) (*http.Response, error) {
	body := deleteTagRequest{Name: tagName}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest(http.MethodDelete, "tags/destroy.json", buf)
	if err != nil {
		return nil, err
	}

	// Expect 204 No Content
	resp, err := s.client.do(req, nil)
	return resp, err
}
