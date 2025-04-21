package client

import (
	"fmt"
	"io"
	"os"
)

// ImportsService handles communication with the imports related
// methods of the Feedbin API
type ImportsService struct {
	client *Client
}

// Create initiates an import from an OPML file
func (s *ImportsService) Create(filePath string) (*ImportStatus, error) {
	// Open the OPML file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening OPML file: %v", err)
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading OPML file: %v", err)
	}

	// Create the request
	type importRequest struct {
		OPML string `json:"opml"`
	}

	req, err := s.client.newRequest("POST", "imports.json", &importRequest{
		OPML: string(content),
	})
	if err != nil {
		return nil, err
	}

	importStatus := new(ImportStatus)
	_, err = s.client.do(req, importStatus)
	if err != nil {
		return nil, err
	}

	return importStatus, nil
}

// Get returns the status of an import
func (s *ImportsService) Get(id int64) (*ImportStatus, error) {
	u := fmt.Sprintf("imports/%d.json", id)

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	importStatus := new(ImportStatus)
	_, err = s.client.do(req, importStatus)
	if err != nil {
		return nil, err
	}

	return importStatus, nil
}
