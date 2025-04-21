// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ParseLinkHeader parses the Link header from an HTTP response and returns a PaginationLinks struct.
func ParseLinkHeader(resp *http.Response) *PaginationLinks {
	link := resp.Header.Get("Link")
	if link == "" {
		return nil
	}

	links := &PaginationLinks{}

	// Regular expression to extract URL and rel from each link
	re := regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)
	matches := re.FindAllStringSubmatch(link, -1)

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

// GetTotalPages extracts the total number of pages from the last page URL.
func GetTotalPages(links *PaginationLinks) (int, error) {
	if links == nil || links.Last == "" {
		return 0, nil
	}

	// Extract page parameter from the last page URL
	re := regexp.MustCompile(`page=(\d+)`)
	match := re.FindStringSubmatch(links.Last)
	if len(match) != 2 {
		return 0, nil
	}

	return strconv.Atoi(match[1])
}

// GetTotalRecords extracts the total number of records from the X-Feedbin-Record-Count header.
func GetTotalRecords(resp *http.Response) (int, error) {
	countStr := resp.Header.Get("X-Feedbin-Record-Count")
	if countStr == "" {
		return 0, nil
	}

	return strconv.Atoi(countStr)
}

// PageOptions represents options for paginated requests.
type PageOptions struct {
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Since   string `url:"since,omitempty"`
}

// AddPageParams adds pagination parameters to a URL.
func AddPageParams(url string, options *PageOptions) string {
	if options == nil {
		return url
	}

	params := []string{}

	if options.Page > 0 {
		params = append(params, "page="+strconv.Itoa(options.Page))
	}

	if options.PerPage > 0 {
		params = append(params, "per_page="+strconv.Itoa(options.PerPage))
	}

	if options.Since != "" {
		params = append(params, "since="+options.Since)
	}

	if len(params) == 0 {
		return url
	}

	if strings.Contains(url, "?") {
		return url + "&" + strings.Join(params, "&")
	}

	return url + "?" + strings.Join(params, "&")
}
