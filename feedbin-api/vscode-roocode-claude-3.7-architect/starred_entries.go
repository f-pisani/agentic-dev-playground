package feedbin

import (
	"net/http"
)

// StarredEntriesService handles starred entry operations
type StarredEntriesService struct {
	client *Client
}

// List returns all starred entry IDs
func (s *StarredEntriesService) List() ([]int, error) {
	req, err := s.client.newRequest(http.MethodGet, "/starred_entries.json", nil)
	if err != nil {
		return nil, err
	}

	var ids []int
	_, err = s.client.do(req, &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// Star stars entries
func (s *StarredEntriesService) Star(ids []int) ([]int, error) {
	if len(ids) > 1000 {
		return nil, ErrTooManyIDs
	}

	req := StarredEntriesRequest{
		StarredEntries: ids,
	}

	httpReq, err := s.client.newRequest(http.MethodPost, "/starred_entries.json", req)
	if err != nil {
		return nil, err
	}

	var starredIDs []int
	_, err = s.client.do(httpReq, &starredIDs)
	if err != nil {
		return nil, err
	}

	return starredIDs, nil
}

// Unstar unstars entries
func (s *StarredEntriesService) Unstar(ids []int) ([]int, error) {
	if len(ids) > 1000 {
		return nil, ErrTooManyIDs
	}

	req := StarredEntriesRequest{
		StarredEntries: ids,
	}

	httpReq, err := s.client.newRequest(http.MethodDelete, "/starred_entries.json", req)
	if err != nil {
		return nil, err
	}

	var unstarredIDs []int
	_, err = s.client.do(httpReq, &unstarredIDs)
	if err != nil {
		return nil, err
	}

	return unstarredIDs, nil
}

// UnstarAlternative unstars entries using the alternative POST endpoint
// This is useful for clients that don't support DELETE requests with a body
func (s *StarredEntriesService) UnstarAlternative(ids []int) ([]int, error) {
	if len(ids) > 1000 {
		return nil, ErrTooManyIDs
	}

	req := StarredEntriesRequest{
		StarredEntries: ids,
	}

	httpReq, err := s.client.newRequest(http.MethodPost, "/starred_entries/delete.json", req)
	if err != nil {
		return nil, err
	}

	var unstarredIDs []int
	_, err = s.client.do(httpReq, &unstarredIDs)
	if err != nil {
		return nil, err
	}

	return unstarredIDs, nil
}
