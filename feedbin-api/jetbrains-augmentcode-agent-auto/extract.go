package feedbin

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// ExtractService handles communication with the full content extraction service.
type ExtractService struct {
	client   *Client
	username string
	secret   string
}

// SetSecret sets the secret key for the extract service.
func (s *ExtractService) SetSecret(secret string) {
	s.secret = secret
}

// NewExtractService creates a new extract service.
func NewExtractService(username, secret string) *ExtractService {
	return &ExtractService{
		username: username,
		secret:   secret,
	}
}

// Article represents an extracted article.
type ExtractedArticle struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Author        string `json:"author"`
	DatePublished string `json:"date_published"`
	LeadImageURL  string `json:"lead_image_url"`
	Dek           string `json:"dek"`
	NextPageURL   string `json:"next_page_url"`
	URL           string `json:"url"`
	Domain        string `json:"domain"`
	Excerpt       string `json:"excerpt"`
	WordCount     int    `json:"word_count"`
	Direction     string `json:"direction"`
	TotalPages    int    `json:"total_pages"`
	RenderedPages int    `json:"rendered_pages"`
}

// Extract extracts the full content of a webpage.
func (s *ExtractService) Extract(pageURL string) (*ExtractedArticle, error) {
	// Create HMAC-SHA1 signature
	h := hmac.New(sha1.New, []byte(s.secret))
	h.Write([]byte(pageURL))
	signature := hex.EncodeToString(h.Sum(nil))

	// Base64 encode the URL
	base64URL := base64.URLEncoding.EncodeToString([]byte(pageURL))

	// Create the request URL
	extractURL := fmt.Sprintf("https://extract.feedbin.com/parser/%s/%s?base64_url=%s",
		s.username, signature, base64URL)

	// Create a new HTTP client
	client := &http.Client{}

	// Create the request
	req, err := http.NewRequest(http.MethodGet, extractURL, nil)
	if err != nil {
		return nil, err
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("extract service returned status code %d", resp.StatusCode)
	}

	// Decode the response
	var article ExtractedArticle
	err = json.NewDecoder(resp.Body).Decode(&article)
	if err != nil {
		return nil, err
	}

	return &article, nil
}
