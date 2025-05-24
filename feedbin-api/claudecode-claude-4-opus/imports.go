package feedbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// CreateImport imports an OPML file
func (c *Client) CreateImport(opmlData []byte) (*Import, error) {
	// Create request with XML body
	resp, err := c.requestXML(http.MethodPost, "/imports.json", opmlData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var importResult Import
	if err := json.NewDecoder(resp.Body).Decode(&importResult); err != nil {
		return nil, fmt.Errorf("failed to decode import response: %w", err)
	}

	return &importResult, nil
}

// requestXML performs an HTTP request with XML content
func (c *Client) requestXML(method, path string, xmlData []byte) (*http.Response, error) {
	// Build URL
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("/%s%s", apiVersion, path)

	// Create request
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(xmlData))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.SetBasicAuth(c.Email, c.Password)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, newAPIError(resp, string(bodyBytes))
	}

	return resp, nil
}

// CreateImportFromReader imports an OPML file from an io.Reader
func (c *Client) CreateImportFromReader(reader io.Reader) (*Import, error) {
	opmlData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read OPML data: %w", err)
	}
	return c.CreateImport(opmlData)
}

// GetImport retrieves the status of an import
func (c *Client) GetImport(id int) (*Import, error) {
	var importResult Import
	path := fmt.Sprintf("/imports/%d.json", id)
	err := c.get(path, nil, &importResult)
	if err != nil {
		return nil, err
	}
	return &importResult, nil
}

// GetImports retrieves all imports for the user
func (c *Client) GetImports() ([]Import, error) {
	var imports []Import
	err := c.get("/imports.json", nil, &imports)
	return imports, err
}

// WaitForImport polls the import status until it's complete
func (c *Client) WaitForImport(id int) (*Import, error) {
	for {
		importResult, err := c.GetImport(id)
		if err != nil {
			return nil, err
		}

		if importResult.Complete {
			return importResult, nil
		}

		// Wait before polling again
		time.Sleep(2 * time.Second)
	}
}
