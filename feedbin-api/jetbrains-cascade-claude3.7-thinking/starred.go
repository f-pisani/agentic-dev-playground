package client

import (
	"fmt"
)

// StarredEntriesService handles communication with the starred entries related
// methods of the Feedbin API
type StarredEntriesService struct {
	client *Client
}

// List returns all starred entry IDs
func (s *StarredEntriesService) List() ([]int64, error) {
	req, err := s.client.newRequest("GET", "starred_entries.json", nil)
	if err != nil {
		return nil, err
	}

	var entryIDs []int64
	_, err = s.client.do(req, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// Star marks the specified entry IDs as starred
func (s *StarredEntriesService) Star(entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs can be starred in a single request")
	}

	type starredRequest struct {
		StarredEntries []int64 `json:"starred_entries"`
	}

	req, err := s.client.newRequest("POST", "starred_entries.json", &starredRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, err
	}

	var processedIDs []int64
	_, err = s.client.do(req, &processedIDs)
	if err != nil {
		return nil, err
	}

	return processedIDs, nil
}

// Unstar removes the star from the specified entry IDs
func (s *StarredEntriesService) Unstar(entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs can be unstarred in a single request")
	}

	type starredRequest struct {
		StarredEntries []int64 `json:"starred_entries"`
	}

	req, err := s.client.newRequest("DELETE", "starred_entries.json", &starredRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, err
	}

	var processedIDs []int64
	_, err = s.client.do(req, &processedIDs)
	if err != nil {
		return nil, err
	}

	return processedIDs, nil
}

// UnstarWithPOST removes the star from the specified entry IDs using POST instead of DELETE
// Some clients like Android don't easily allow a body with a DELETE request
func (s *StarredEntriesService) UnstarWithPOST(entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs can be unstarred in a single request")
	}

	type starredRequest struct {
		StarredEntries []int64 `json:"starred_entries"`
	}

	req, err := s.client.newRequest("POST", "starred_entries/delete.json", &starredRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, err
	}

	var processedIDs []int64
	_, err = s.client.do(req, &processedIDs)
	if err != nil {
		return nil, err
	}

	return processedIDs, nil
}
