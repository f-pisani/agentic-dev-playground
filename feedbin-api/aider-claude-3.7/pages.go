package feedbin

import (
	"fmt"
	"net/http"
)

// GetPages retrieves all pages
func (c *Client) GetPages() ([]Page, error) {
	req, err := c.NewRequest(http.MethodGet, "/v2/pages.json", nil)
	if err != nil {
		return nil, err
	}
	
	var pages []Page
	_, err = c.Do(req, &pages)
	if err != nil {
		return nil, err
	}
	
	return pages, nil
}

// GetPage retrieves a specific page by ID
func (c *Client) GetPage(id int64) (*Page, error) {
	path := fmt.Sprintf("/v2/pages/%d.json", id)
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	page := new(Page)
	_, err = c.Do(req, page)
	if err != nil {
		return nil, err
	}
	
	return page, nil
}

// GetEntryPages retrieves all pages for a specific entry
func (c *Client) GetEntryPages(entryID int64) ([]Page, error) {
	path := fmt.Sprintf("/v2/entries/%d/pages.json", entryID)
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	var pages []Page
	_, err = c.Do(req, &pages)
	if err != nil {
		return nil, err
	}
	
	return pages, nil
}
