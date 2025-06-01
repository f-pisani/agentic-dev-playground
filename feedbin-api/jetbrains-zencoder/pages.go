package feedbin

import (
	"fmt"
	"net/http"
)

// GetPages retrieves all pages for the authenticated user.
func (c *Client) GetPages() ([]Page, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/pages.json", nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var pages []Page
	if err := parseResponse(resp, &pages); err != nil {
		return nil, err
	}

	return pages, nil
}

// GetPage retrieves a specific page by ID.
func (c *Client) GetPage(id int) (*Page, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/pages/%d.json", id), nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var page Page
	if err := parseResponse(resp, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// CreatePage creates a new page from a URL.
func (c *Client) CreatePage(url string) (*Page, error) {
	// Create the request body
	reqBody := PageRequest{
		URL: url,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/pages.json", reqBody, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var page Page
	if err := parseResponse(resp, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// DeletePage deletes a page by ID.
func (c *Client) DeletePage(id int) error {
	// Make the request
	resp, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/pages/%d.json", id), nil, nil)
	if err != nil {
		return err
	}

	// Check for success (204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
