package feedbinapi

import (
	"net/http" // Added import
	"regexp"
	"strconv"
	"strings"
)

// PaginationInfo holds pagination details from API responses.
type PaginationInfo struct {
	NextPageURL  string
	PrevPageURL  string
	FirstPageURL string
	LastPageURL  string
	TotalRecords int
}

// Regular expression to parse Link headers, handles multiple links.
var linkHeaderRegex = regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)

// parseLinkHeader parses a Link HTTP header and returns a map of rel -> url.
func parseLinkHeader(header string) map[string]string {
	links := make(map[string]string)
	if header == "" {
		return links
	}

	matches := linkHeaderRegex.FindAllStringSubmatch(header, -1)
	for _, match := range matches {
		if len(match) == 3 {
			url := strings.TrimSpace(match[1])
			rel := strings.TrimSpace(match[2])
			links[rel] = url
		}
	}
	return links
}

// extractPaginationInfo creates a PaginationInfo object from HTTP headers.
func extractPaginationInfo(respHeaders http.Header) *PaginationInfo {
	if respHeaders == nil {
		return nil
	}

	linkHeader := respHeaders.Get("Link")
	recordCountHeader := respHeaders.Get("X-Feedbin-Record-Count")

	if linkHeader == "" && recordCountHeader == "" {
		return nil // No pagination headers present
	}

	pageInfo := &PaginationInfo{}
	parsedLinks := parseLinkHeader(linkHeader)

	pageInfo.NextPageURL = parsedLinks["next"]
	pageInfo.PrevPageURL = parsedLinks["prev"]
	pageInfo.FirstPageURL = parsedLinks["first"]
	pageInfo.LastPageURL = parsedLinks["last"]

	if count, err := strconv.Atoi(recordCountHeader); err == nil {
		pageInfo.TotalRecords = count
	}

	// If any field is populated, return the struct
	if pageInfo.NextPageURL != "" || pageInfo.PrevPageURL != "" ||
		pageInfo.FirstPageURL != "" || pageInfo.LastPageURL != "" || pageInfo.TotalRecords > 0 {
		return pageInfo
	}

	return nil
}
