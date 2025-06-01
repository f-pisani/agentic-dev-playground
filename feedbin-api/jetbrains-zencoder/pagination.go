package feedbin

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// PaginationLinks represents the pagination links from the Link header.
type PaginationLinks struct {
	First string
	Prev  string
	Next  string
	Last  string
}

// PaginationInfo represents pagination information from a response.
type PaginationInfo struct {
	Links       PaginationLinks
	TotalCount  int
	CurrentPage int
}

// parsePaginationLinks parses the Link header from an HTTP response.
func parsePaginationLinks(resp *http.Response) PaginationLinks {
	links := PaginationLinks{}

	// Get the Link header
	linkHeader := resp.Header.Get("Link")
	if linkHeader == "" {
		return links
	}

	// Regular expression to parse the Link header
	// Format: <url>; rel="relation"
	re := regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)

	// Find all matches
	matches := re.FindAllStringSubmatch(linkHeader, -1)
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		url := match[1]
		rel := match[2]

		switch rel {
		case "first":
			links.First = url
		case "prev":
			links.Prev = url
		case "next":
			links.Next = url
		case "last":
			links.Last = url
		}
	}

	return links
}

// GetPaginationInfo extracts pagination information from an HTTP response.
func GetPaginationInfo(resp *http.Response) PaginationInfo {
	info := PaginationInfo{
		Links: parsePaginationLinks(resp),
	}

	// Get the total count from the X-Feedbin-Record-Count header
	countStr := resp.Header.Get("X-Feedbin-Record-Count")
	if countStr != "" {
		count, err := strconv.Atoi(countStr)
		if err == nil {
			info.TotalCount = count
		}
	}

	// Try to determine the current page from the request URL
	if resp.Request != nil && resp.Request.URL != nil {
		query := resp.Request.URL.Query()
		pageStr := query.Get("page")
		if pageStr != "" {
			page, err := strconv.Atoi(pageStr)
			if err == nil {
				info.CurrentPage = page
			}
		} else {
			// If no page parameter, assume it's page 1
			info.CurrentPage = 1
		}
	}

	return info
}

// GetNextPageURL returns the URL for the next page of results.
func (p PaginationInfo) GetNextPageURL() string {
	return p.Links.Next
}

// GetPrevPageURL returns the URL for the previous page of results.
func (p PaginationInfo) GetPrevPageURL() string {
	return p.Links.Prev
}

// GetFirstPageURL returns the URL for the first page of results.
func (p PaginationInfo) GetFirstPageURL() string {
	return p.Links.First
}

// GetLastPageURL returns the URL for the last page of results.
func (p PaginationInfo) GetLastPageURL() string {
	return p.Links.Last
}

// HasNextPage returns true if there is a next page of results.
func (p PaginationInfo) HasNextPage() bool {
	return p.Links.Next != ""
}

// HasPrevPage returns true if there is a previous page of results.
func (p PaginationInfo) HasPrevPage() bool {
	return p.Links.Prev != ""
}

// GetPageFromURL extracts the page number from a URL.
func GetPageFromURL(urlStr string) (int, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return 0, err
	}

	query := u.Query()
	pageStr := query.Get("page")
	if pageStr == "" {
		return 1, nil // Default to page 1
	}

	return strconv.Atoi(pageStr)
}

// AddPaginationParams adds pagination parameters to a URL query.
func AddPaginationParams(query url.Values, page, perPage int) url.Values {
	if query == nil {
		query = url.Values{}
	}

	if page > 0 {
		query.Set("page", strconv.Itoa(page))
	}

	if perPage > 0 {
		query.Set("per_page", strconv.Itoa(perPage))
	}

	return query
}
