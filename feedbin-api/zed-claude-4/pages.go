package feedbin

import (
	"context"
)

// CreatePage creates a new page from a URL
func (c *Client) CreatePage(ctx context.Context, url string, title string) (*Entry, error) {
	if url == "" {
		return nil, &ValidationError{
			Field:   "url",
			Message: "URL is required",
		}
	}

	req := &CreatePageRequest{
		URL:   url,
		Title: title,
	}

	var entry Entry
	_, err := c.makeRequest(ctx, "POST", "/pages.json", req, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// CreatePageWithoutTitle creates a new page from a URL without specifying a title
func (c *Client) CreatePageWithoutTitle(ctx context.Context, url string) (*Entry, error) {
	return c.CreatePage(ctx, url, "")
}

// Convenience methods without context (use background context)

func (c *Client) CreatePageWithoutContext(url string, title string) (*Entry, error) {
	return c.CreatePage(context.Background(), url, title)
}

func (c *Client) CreatePageWithoutTitleWithoutContext(url string) (*Entry, error) {
	return c.CreatePageWithoutTitle(context.Background(), url)
}
