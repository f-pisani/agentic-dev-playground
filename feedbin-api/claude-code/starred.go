package feedbin

import (
	"net/http"
)

// StarredEntriesService handles communication with the starred entries related
// endpoints of the Feedbin API
type StarredEntriesService struct {
	client *Client
}

// StarredEntriesRequest is used to star/unstar entries
type StarredEntriesRequest struct {
	StarredEntries []int `json:"starred_entries"`
}

// List returns a list of starred entry IDs
func (s *StarredEntriesService) List() ([]int, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "starred_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var starredEntries []int
	resp, err := s.client.Do(req, &starredEntries)
	if err != nil {
		return nil, resp, err
	}

	return starredEntries, resp, nil
}

// Star marks entries as starred
func (s *StarredEntriesService) Star(entryIDs []int) ([]int, *http.Response, error) {
	starRequest := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "starred_entries.json", starRequest)
	if err != nil {
		return nil, nil, err
	}

	var starredEntries []int
	resp, err := s.client.Do(req, &starredEntries)
	if err != nil {
		return nil, resp, err
	}

	return starredEntries, resp, nil
}

// Unstar marks entries as unstarred
func (s *StarredEntriesService) Unstar(entryIDs []int) ([]int, *http.Response, error) {
	starRequest := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodDelete, "starred_entries.json", starRequest)
	if err != nil {
		return nil, nil, err
	}

	var unstarredEntries []int
	resp, err := s.client.Do(req, &unstarredEntries)
	if err != nil {
		return nil, resp, err
	}

	return unstarredEntries, resp, nil
}

// UnstarWithPost marks entries as unstarred using POST instead of DELETE
// Some clients may not allow sending a body with DELETE
func (s *StarredEntriesService) UnstarWithPost(entryIDs []int) ([]int, *http.Response, error) {
	starRequest := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "starred_entries/delete.json", starRequest)
	if err != nil {
		return nil, nil, err
	}

	var unstarredEntries []int
	resp, err := s.client.Do(req, &unstarredEntries)
	if err != nil {
		return nil, resp, err
	}

	return unstarredEntries, resp, nil
}
