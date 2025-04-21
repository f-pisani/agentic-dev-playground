package feedbin

import (
	"net/http"
)

// UnreadEntriesService handles communication with the unread entries related
// endpoints of the Feedbin API
type UnreadEntriesService struct {
	client *Client
}

// UnreadEntriesRequest is used to mark entries as read/unread
type UnreadEntriesRequest struct {
	UnreadEntries []int `json:"unread_entries"`
}

// List returns a list of unread entry IDs
func (s *UnreadEntriesService) List() ([]int, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "unread_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var unreadEntries []int
	resp, err := s.client.Do(req, &unreadEntries)
	if err != nil {
		return nil, resp, err
	}

	return unreadEntries, resp, nil
}

// MarkAsUnread marks entries as unread
func (s *UnreadEntriesService) MarkAsUnread(entryIDs []int) ([]int, *http.Response, error) {
	unreadRequest := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "unread_entries.json", unreadRequest)
	if err != nil {
		return nil, nil, err
	}

	var markedEntries []int
	resp, err := s.client.Do(req, &markedEntries)
	if err != nil {
		return nil, resp, err
	}

	return markedEntries, resp, nil
}

// MarkAsRead marks entries as read
func (s *UnreadEntriesService) MarkAsRead(entryIDs []int) ([]int, *http.Response, error) {
	unreadRequest := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodDelete, "unread_entries.json", unreadRequest)
	if err != nil {
		return nil, nil, err
	}

	var markedEntries []int
	resp, err := s.client.Do(req, &markedEntries)
	if err != nil {
		return nil, resp, err
	}

	return markedEntries, resp, nil
}

// MarkAsReadWithPost marks entries as read using POST instead of DELETE
// Some clients may not allow sending a body with DELETE
func (s *UnreadEntriesService) MarkAsReadWithPost(entryIDs []int) ([]int, *http.Response, error) {
	unreadRequest := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "unread_entries/delete.json", unreadRequest)
	if err != nil {
		return nil, nil, err
	}

	var markedEntries []int
	resp, err := s.client.Do(req, &markedEntries)
	if err != nil {
		return nil, resp, err
	}

	return markedEntries, resp, nil
}
