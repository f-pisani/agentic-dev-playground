package feedbin

import (
	"net/http"
)

// StarredService handles communication with the starred entries related
// endpoints of the Feedbin API
type StarredService struct {
	client *Client
}

// StarredEntriesRequest represents a request to star or unstar entries
type StarredEntriesRequest struct {
	StarredEntries []int `json:"starred_entries"`
}

// List returns all starred entry IDs
// https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#get-starred-entries
func (s *StarredService) List() ([]int, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "starred_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var entryIDs []int
	resp, err := s.client.Do(req, &entryIDs)
	if err != nil {
		return nil, resp, err
	}

	return entryIDs, resp, nil
}

// Star marks entries as starred
// https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#create-starred-entries
func (s *StarredService) Star(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) > 1000 {
		return nil, nil, ErrTooManyIDs
	}

	request := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "starred_entries.json", request)
	if err != nil {
		return nil, nil, err
	}

	var markedIDs []int
	resp, err := s.client.Do(req, &markedIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedIDs, resp, nil
}

// Unstar removes the star from entries
// https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#delete-starred-entries-unstar
func (s *StarredService) Unstar(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) > 1000 {
		return nil, nil, ErrTooManyIDs
	}

	request := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodDelete, "starred_entries.json", request)
	if err != nil {
		return nil, nil, err
	}

	var markedIDs []int
	resp, err := s.client.Do(req, &markedIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedIDs, resp, nil
}

// UnstarAlternative provides an alternative to DELETE for clients that don't support
// bodies in DELETE requests
// https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#delete-alternative
func (s *StarredService) UnstarAlternative(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) > 1000 {
		return nil, nil, ErrTooManyIDs
	}

	request := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "starred_entries/delete.json", request)
	if err != nil {
		return nil, nil, err
	}

	var markedIDs []int
	resp, err := s.client.Do(req, &markedIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedIDs, resp, nil
}
