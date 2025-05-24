package feedbin

// CreatePage saves a web page as an entry
func (c *Client) CreatePage(url string, title string) (*Entry, error) {
	body := struct {
		URL   string `json:"url"`
		Title string `json:"title,omitempty"`
	}{
		URL:   url,
		Title: title,
	}

	var entry Entry
	err := c.post("/pages.json", body, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// CreatePageFromURL is an alias for CreatePage for clarity
func (c *Client) CreatePageFromURL(url string) (*Entry, error) {
	return c.CreatePage(url, "")
}
