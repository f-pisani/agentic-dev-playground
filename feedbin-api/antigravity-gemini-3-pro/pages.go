package feedbin

import "context"

type pageCreateRequest struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// CreatePage creates a new page (entry) from a URL.
func (c *Client) CreatePage(ctx context.Context, url, title string) (*Entry, *Response, error) {
	reqBody := &pageCreateRequest{URL: url, Title: title}
	req, err := c.NewRequest("POST", "pages.json", reqBody)
	if err != nil {
		return nil, nil, err
	}

	var entry Entry
	resp, err := c.Do(ctx, req, &entry)
	if err != nil {
		return nil, resp, err
	}

	return &entry, resp, nil
}
