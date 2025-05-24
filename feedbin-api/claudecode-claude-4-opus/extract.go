package feedbin

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// ExtractContent retrieves the full content for an entry using Mercury Parser
func (c *Client) ExtractContent(entryID int) (*ExtractedContent, error) {
	// Get the entry to get its extracted_content_url
	entry, err := c.GetEntry(entryID)
	if err != nil {
		return nil, err
	}

	if entry.ExtractedContentURL == "" {
		return nil, fmt.Errorf("no extracted content URL available for entry %d", entryID)
	}

	// Parse the URL to extract the path
	parts := strings.Split(entry.ExtractedContentURL, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid extracted content URL format")
	}

	// The URL format is: https://extract.feedbin.com/parser/{signature}/{encoded_url}
	// We need to make the request directly to the extract service
	var content ExtractedContent

	// Since we're using the standard library only, we'll make a direct HTTP request
	resp, err := c.HTTPClient.Get(entry.ExtractedContentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch extracted content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch extracted content: status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, fmt.Errorf("failed to decode extracted content: %w", err)
	}

	return &content, nil
}

// GenerateExtractedContentURL generates a signed URL for content extraction
// This is useful if you want to construct the URL yourself
func (c *Client) GenerateExtractedContentURL(url string, secret string) (string, error) {
	// Base64 URL encode the URL (RFC 4648 url-safe variant)
	encodedURL := base64.RawURLEncoding.EncodeToString([]byte(url))

	// Generate HMAC-SHA1 signature
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(encodedURL))
	signature := hex.EncodeToString(h.Sum(nil))

	// Construct the full URL
	extractURL := fmt.Sprintf("https://extract.feedbin.com/parser/%s/%s", signature, encodedURL)

	return extractURL, nil
}

// ExtractContentFromURL extracts content from any URL
// Note: This requires a secret key which is not provided via the standard API
func (c *Client) ExtractContentFromURL(url string, secret string) (*ExtractedContent, error) {
	extractURL, err := c.GenerateExtractedContentURL(url, secret)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Get(extractURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch extracted content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch extracted content: status %d", resp.StatusCode)
	}

	var content ExtractedContent
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, fmt.Errorf("failed to decode extracted content: %w", err)
	}

	return &content, nil
}
