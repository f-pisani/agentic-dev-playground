package feedbin

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// parseLinkHeader parses the Link header to extract pagination URLs
func parseLinkHeader(linkHeader string) *PaginationInfo {
	if linkHeader == "" {
		return &PaginationInfo{}
	}

	pagination := &PaginationInfo{}

	// Regular expression to parse Link header format:
	// <https://api.feedbin.com/v2/entries.json?page=2>; rel="next"
	linkRegex := regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)
	matches := linkRegex.FindAllStringSubmatch(linkHeader, -1)

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		url := match[1]
		rel := match[2]

		switch rel {
		case "first":
			pagination.First = url
		case "prev":
			pagination.Previous = url
		case "next":
			pagination.Next = url
		case "last":
			pagination.Last = url
		}
	}

	return pagination
}

// parseRecordCount parses the X-Feedbin-Record-Count header
func parseRecordCount(header string) int {
	if header == "" {
		return 0
	}

	count, err := strconv.Atoi(strings.TrimSpace(header))
	if err != nil {
		return 0
	}

	return count
}

// extractPaginationInfo extracts pagination information from HTTP response headers
func extractPaginationInfo(resp *http.Response) *PaginationInfo {
	pagination := parseLinkHeader(resp.Header.Get("Link"))
	pagination.Total = parseRecordCount(resp.Header.Get("X-Feedbin-Record-Count"))
	return pagination
}

// PaginatedEntries represents a paginated list of entries
type PaginatedEntries struct {
	Entries    []Entry
	Pagination *PaginationInfo
}

// PaginatedSubscriptions represents a paginated list of subscriptions
type PaginatedSubscriptions struct {
	Subscriptions []Subscription
	Pagination    *PaginationInfo
}

// PaginatedTaggings represents a paginated list of taggings
type PaginatedTaggings struct {
	Taggings   []Tagging
	Pagination *PaginationInfo
}
