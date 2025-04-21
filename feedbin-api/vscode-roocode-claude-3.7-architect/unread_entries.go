package feedbin

import (
	"net/http"
)

// UnreadEntriesService handles unread entry operations
type UnreadEntriesService struct {
	client *Client
}

// List returns all unread entry IDs
func (s *UnreadEntriesService) List() ([]int, error) {
	req, err := s.client.newRequest(http.MethodGet, "/unread_entries.json", nil)
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

// MarkAsUnread marks entries as unread
func (s *UnreadEntriesService) MarkAsUnread(ids []int) ([]int, error) {
	if len(ids) > 1000 {
		return nil, ErrTooManyIDs
	}

	req := UnreadEntriesRequest{
		UnreadEntries: ids,
	}

	httpReq, err := s.client.newRequest(http.MethodPost, "/unread_entries.json", req)
	if err != nil {
		return nil, err
	}

	var markedIDs []int
	_, err = s.client.do(httpReq, &markedIDs)
	if err != nil {
		return nil, err
	}

	return markedIDs, nil
}

// MarkAsRead marks entries as read
func (s *UnreadEntriesService) MarkAsRead(ids []int) ([]int, error) {
	if len(ids) > 1000 {
		return nil, ErrTooManyIDs
	}

	req := UnreadEntriesRequest{
		UnreadEntries: ids,
	}

	httpReq, err := s.client.newRequest(http.MethodDelete, "/unread_entries.json", req)
	if err != nil {
		return nil, err
	}

	var markedIDs []int
	_, err = s.client.do(httpReq, &markedIDs)
	if err != nil {
		return nil, err
	}

	return markedIDs, nil
}

// MarkAsReadAlternative marks entries as read using the alternative POST endpoint
// This is useful for clients that don't support DELETE requests with a body
func (s *UnreadEntriesService) MarkAsReadAlternative(ids []int) ([]int, error) {
	if len(ids) > 1000 {
		return nil, ErrTooManyIDs
	}

	req := UnreadEntriesRequest{
		UnreadEntries: ids,
	}

	httpReq, err := s.client.newRequest(http.MethodPost, "/unread_entries/delete.json", req)
	if err != nil {
		return nil, err
	}

	var markedIDs []int
	_, err = s.client.do(httpReq, &markedIDs)
	if err != nil {
		return nil, err
	}

	return markedIDs, nil
}

// ErrTooManyIDs is returned when too many IDs are provided in a request
var ErrTooManyIDs = &APIError{
	StatusCode: 400,
	Message:    "Too many IDs provided (maximum 1000)",
}
