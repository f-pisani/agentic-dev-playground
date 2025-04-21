package client

import (
	"fmt"
)

// UnreadEntriesService handles communication with the unread entries related
// methods of the Feedbin API
type UnreadEntriesService struct {
	client *Client
}

// List returns all unread entry IDs
func (s *UnreadEntriesService) List() ([]int64, error) {
	req, err := s.client.newRequest("GET", "unread_entries.json", nil)
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

// MarkAsUnread marks the specified entry IDs as unread
func (s *UnreadEntriesService) MarkAsUnread(entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs can be marked as unread in a single request")
	}

	type unreadRequest struct {
		UnreadEntries []int64 `json:"unread_entries"`
	}

	req, err := s.client.newRequest("POST", "unread_entries.json", &unreadRequest{
		UnreadEntries: entryIDs,
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

// MarkAsRead marks the specified entry IDs as read
func (s *UnreadEntriesService) MarkAsRead(entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs can be marked as read in a single request")
	}

	type unreadRequest struct {
		UnreadEntries []int64 `json:"unread_entries"`
	}

	req, err := s.client.newRequest("DELETE", "unread_entries.json", &unreadRequest{
		UnreadEntries: entryIDs,
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

// MarkAsReadWithPOST marks the specified entry IDs as read using POST instead of DELETE
// Some clients like Android don't easily allow a body with a DELETE request
func (s *UnreadEntriesService) MarkAsReadWithPOST(entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs can be marked as read in a single request")
	}

	type unreadRequest struct {
		UnreadEntries []int64 `json:"unread_entries"`
	}

	req, err := s.client.newRequest("POST", "unread_entries/delete.json", &unreadRequest{
		UnreadEntries: entryIDs,
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
