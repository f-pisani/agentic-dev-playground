package client

import (
	"bytes"
	"fmt"
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// CreateImport uploads an OPML file (as XML string) to Feedbin.
func (c *Client) CreateImport(xmlData string) (*models.Import, error) {
	req, err := c.newRequest(http.MethodPost, "/imports.json", bytes.NewReader([]byte(xmlData)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml")
	var imp models.Import
	err = c.do(req, &imp)
	if err != nil {
		return nil, err
	}
	return &imp, nil
}

// GetImports lists all imports for the user.
func (c *Client) GetImports() ([]models.Import, error) {
	req, err := c.newRequest(http.MethodGet, "/imports.json", nil)
	if err != nil {
		return nil, err
	}
	var imps models.ImportsResponse
	err = c.do(req, &imps)
	if err != nil {
		return nil, err
	}
	return imps, nil
}

// GetImport fetches the status/details of a specific import by ID.
func (c *Client) GetImport(id int64) (*models.Import, error) {
	path := "/imports/" + fmt.Sprint(id) + ".json"
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var imp models.Import
	err = c.do(req, &imp)
	if err != nil {
		return nil, err
	}
	return &imp, nil
}
