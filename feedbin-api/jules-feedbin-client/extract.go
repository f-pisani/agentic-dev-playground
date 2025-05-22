package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// ExtractBaseURL is the base URL for the Feedbin Full Content Extraction API.
	ExtractBaseURL = "https://extract.feedbin.com/parser"
)

// ExtractService handles operations for the Full Content Extraction API.
type ExtractService struct {
	httpClient *http.Client // Use a specific httpClient, can be shared from main client
	apiToken   string       // API token, typically the Feedbin username
	userAgent  string       // User agent, can be shared from main client
}

// NewExtractService creates a new service for content extraction.
func NewExtractService(sharedHttpClient *http.Client, apiToken string, userAgent string) *ExtractService {
	return &ExtractService{
		httpClient: sharedHttpClient,
		apiToken:   apiToken,
		userAgent:  userAgent,
	}
}

// SetApiToken allows updating the API token for the extraction service.
func (s *ExtractService) SetApiToken(apiToken string) {
	s.apiToken = apiToken
}

// Extract fetches the parsed content of a given URL.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/extract-full-content.md
// The endpoint is GET /parser?token=<token>&url=<url_to_parse>
func (s *ExtractService) Extract(urlToParse string) (*ExtractResult, *http.Response, error) {
	if s.apiToken == "" {
		return nil, nil, fmt.Errorf("API token for ExtractService is not set")
	}
	if urlToParse == "" {
		return nil, nil, fmt.Errorf("URL to parse is required")
	}

	extractURL, err := url.Parse(ExtractBaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse ExtractBaseURL: %w", err)
	}

	q := extractURL.Query()
	q.Set("token", s.apiToken)
	q.Set("url", urlToParse)
	extractURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, extractURL.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("creating extraction request: %w", err)
	}
	if s.userAgent != "" {
		req.Header.Set("User-Agent", s.userAgent)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("executing extraction request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Attempt to read body for more detailed error message, similar to main client's CheckResponse
		// errorBody, _ := io.ReadAll(resp.Body) // Be careful with consuming body
		// return nil, resp, fmt.Errorf("extraction API error: status %s, body: %s", resp.Status, string(errorBody))
		return nil, resp, fmt.Errorf("extraction API error: status %s", resp.Status)
	}

	var result ExtractResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, resp, fmt.Errorf("decoding extraction response: %w", err)
	}

	return &result, resp, nil
}
