// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"strconv"
)

// IconService handles communication with the icon related
// methods of the Feedbin API.
type IconService struct {
	client *Client
}

// GetIcons retrieves icons for the specified feed IDs.
func (s *IconService) GetIcons(feedIDs []int) ([]Icon, error) {
	// Convert feed IDs to strings and join with commas
	feedIDsParam := ""
	for i, id := range feedIDs {
		if i > 0 {
			feedIDsParam += ","
		}
		feedIDsParam += strconv.Itoa(id)
	}

	path := fmt.Sprintf("/icons.json?feed_ids=%s", feedIDsParam)
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var icons []Icon
	_, err = s.client.Do(req, &icons)
	if err != nil {
		return nil, err
	}

	return icons, nil
}
