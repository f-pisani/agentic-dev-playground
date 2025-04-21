package feedbin

import (
	"net/http"
)

// UnreadService handles communication with the unread entries related
// endpoints of the Feedbin API
type UnreadService struct {
	client *Client
}

// UnreadEntriesRequest represents a request to mark entries as read or unread
type UnreadEntriesRequest struct {
	UnreadEntries []int `json:"unread_entries"`
}

// List returns all unread entry IDs
// https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#get-unread-entries
func (s *UnreadService) List() ([]int, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "unread_entries.json", nil)
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

// MarkAsUnread marks entries as unread
// https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#create-unread-entries-mark-as-unread
func (s *UnreadService) MarkAsUnread(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) > 1000 {
		return nil, nil, ErrTooManyIDs
	}

	request := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "unread_entries.json", request)
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

// MarkAsRead marks entries as read
// https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#delete-unread-entries-mark-as-read
func (s *UnreadService) MarkAsRead(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) > 1000 {
		return nil, nil, ErrTooManyIDs
	}

	request := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodDelete, "unread_entries.json", request)
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

// MarkAsReadAlternative provides an alternative to DELETE for clients that don't support
// bodies in DELETE requests
// https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#delete-alternative
func (s *UnreadService) MarkAsReadAlternative(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) > 1000 {
		return nil, nil, ErrTooManyIDs
	}

	request := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "unread_entries/delete.json", request)
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
