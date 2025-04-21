// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
	"time"
)

// PageService handles communication with the page related
// methods of the Feedbin API.
type PageService struct {
	client *Client
}

// GetPages retrieves all pages.
func (s *PageService) GetPages() ([]Page, error) {
	req, err := s.client.NewRequest("GET", "/pages.json", nil)
	if err != nil {
		return nil, err
	}

	var pages []Page
	_, err = s.client.Do(req, &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// GetPage retrieves a specific page.
func (s *PageService) GetPage(id int) (*Page, error) {
	path := fmt.Sprintf("/pages/%d.json", id)
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var page Page
	_, err = s.client.Do(req, &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// CreatePage creates a new page.
func (s *PageService) CreatePage(title, url string, published time.Time) (*Page, error) {
	body := map[string]interface{}{
		"title":     title,
		"url":       url,
		"published": published.Format(time.RFC3339),
	}

	req, err := s.client.NewRequest("POST", "/pages.json", body)
	if err != nil {
		return nil, err
	}

	var page Page
	_, err = s.client.Do(req, &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// UpdatePage updates an existing page.
func (s *PageService) UpdatePage(id int, title, url string, published time.Time) (*Page, error) {
	body := map[string]interface{}{
		"title":     title,
		"url":       url,
		"published": published.Format(time.RFC3339),
	}

	path := fmt.Sprintf("/pages/%d.json", id)
	req, err := s.client.NewRequest("PATCH", path, body)
	if err != nil {
		return nil, err
	}

	var page Page
	_, err = s.client.Do(req, &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// DeletePage deletes a page.
func (s *PageService) DeletePage(id int) error {
	path := fmt.Sprintf("/pages/%d.json", id)
	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}