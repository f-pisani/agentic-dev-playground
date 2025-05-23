package feedbin

import (
	"context"
	"fmt"
)

// GetPages retrieves all pages for the authenticated user.
// Returns a slice of pages.
func (c *Client) GetPages(ctx context.Context) ([]Page, error) {
	var pages []Page
	_, err := c.get(ctx, "pages.json", nil, &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// GetPage retrieves a specific page by ID.
func (c *Client) GetPage(ctx context.Context, id int) (*Page, error) {
	path := fmt.Sprintf("pages/%d.json", id)

	var page Page
	_, err := c.get(ctx, path, nil, &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// CreatePage creates a new page from a URL.
// Returns the created page or an error.
func (c *Client) CreatePage(ctx context.Context, url string) (*Page, error) {
	request := CreatePageRequest{
		URL: url,
	}

	var page Page
	_, err := c.post(ctx, "pages.json", nil, request, &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

// DeletePage deletes a page by ID.
func (c *Client) DeletePage(ctx context.Context, id int) error {
	path := fmt.Sprintf("pages/%d.json", id)
	_, err := c.delete(ctx, path, nil, nil)
	return err
}

// GetPageByURL retrieves a page by URL.
// Returns the first page with the matching URL, or nil if not found.
func (c *Client) GetPageByURL(ctx context.Context, url string) (*Page, error) {
	pages, err := c.GetPages(ctx)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		if page.URL == url {
			return &page, nil
		}
	}

	return nil, fmt.Errorf("page with URL '%s' not found", url)
}

// GetPagesByTitle retrieves all pages with a specific title.
func (c *Client) GetPagesByTitle(ctx context.Context, title string) ([]Page, error) {
	pages, err := c.GetPages(ctx)
	if err != nil {
		return nil, err
	}

	var matchingPages []Page
	for _, page := range pages {
		if page.Title == title {
			matchingPages = append(matchingPages, page)
		}
	}

	return matchingPages, nil
}

// GetPagesByAuthor retrieves all pages by a specific author.
func (c *Client) GetPagesByAuthor(ctx context.Context, author string) ([]Page, error) {
	pages, err := c.GetPages(ctx)
	if err != nil {
		return nil, err
	}

	var authorPages []Page
	for _, page := range pages {
		if page.Author == author {
			authorPages = append(authorPages, page)
		}
	}

	return authorPages, nil
}
