package feedbin

import (
	"fmt"
	"net/http"
)

// GetImports retrieves all imports
func (c *Client) GetImports() ([]Import, error) {
	req, err := c.NewRequest(http.MethodGet, "/v2/imports.json", nil)
	if err != nil {
		return nil, err
	}
	
	var imports []Import
	_, err = c.Do(req, &imports)
	if err != nil {
		return nil, err
	}
	
	return imports, nil
}

// GetImport retrieves a specific import by ID
func (c *Client) GetImport(id int64) (*Import, error) {
	path := fmt.Sprintf("/v2/imports/%d.json", id)
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	importObj := new(Import)
	_, err = c.Do(req, importObj)
	if err != nil {
		return nil, err
	}
	
	return importObj, nil
}

// CreateImport creates a new import from OPML data
func (c *Client) CreateImport(opml string) (*Import, error) {
	importReq := &ImportRequest{
		OPML: opml,
	}
	
	req, err := c.NewRequest(http.MethodPost, "/v2/imports.json", importReq)
	if err != nil {
		return nil, err
	}
	
	importObj := new(Import)
	_, err = c.Do(req, importObj)
	if err != nil {
		return nil, err
	}
	
	return importObj, nil
}
