package feedbin

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// LinkHeader represents parsed link header relations
type LinkHeader struct {
	FirstURL *url.URL
	PrevURL  *url.URL
	NextURL  *url.URL
	LastURL  *url.URL
	Count    int
}

// ParseLinkHeader parses the Link header from an HTTP response
func ParseLinkHeader(resp *http.Response) (*LinkHeader, error) {
	link := resp.Header.Get("Link")
	if link == "" {
		return nil, nil
	}

	linkHeader := &LinkHeader{}

	// Parse record count if available
	countStr := resp.Header.Get("X-Feedbin-Record-Count")
	if countStr != "" {
		count, err := strconv.Atoi(countStr)
		if err == nil {
			linkHeader.Count = count
		}
	}

	// Parse link relations
	linkRegex := regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)
	matches := linkRegex.FindAllStringSubmatch(link, -1)

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		uri := match[1]
		rel := match[2]

		parsedURL, err := url.Parse(uri)
		if err != nil {
			continue
		}

		switch rel {
		case "first":
			linkHeader.FirstURL = parsedURL
		case "prev":
			linkHeader.PrevURL = parsedURL
		case "next":
			linkHeader.NextURL = parsedURL
		case "last":
			linkHeader.LastURL = parsedURL
		}
	}

	return linkHeader, nil
}

// PageOptions contains pagination options
type PageOptions struct {
	Page    *int   `url:"page,omitempty"`
	PerPage *int   `url:"per_page,omitempty"`
	Since   string `url:"since,omitempty"`
}

// AddQueryParams adds query parameters to the URL
func AddQueryParams(baseURL string, opts interface{}) (string, error) {
	if opts == nil {
		return baseURL, nil
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	params := u.Query()

	// We'll use reflection to extract query parameters from the struct
	switch v := opts.(type) {
	case *PageOptions:
		if v.Page != nil {
			params.Set("page", strconv.Itoa(*v.Page))
		}
		if v.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*v.PerPage))
		}
		if v.Since != "" {
			params.Set("since", v.Since)
		}
	}

	// Handle additional option types
	handleEntryOptions(opts, params)
	handleSubscriptionOptions(opts, params)

	u.RawQuery = params.Encode()
	return u.String(), nil
}

// handleEntryOptions adds entry-specific query parameters
func handleEntryOptions(opts interface{}, params url.Values) {
	type EntryOptions struct {
		IDs                []int  `url:"ids,omitempty"`
		Read               *bool  `url:"read,omitempty"`
		Starred            *bool  `url:"starred,omitempty"`
		Mode               string `url:"mode,omitempty"`
		IncludeOriginal    *bool  `url:"include_original,omitempty"`
		IncludeEnclosure   *bool  `url:"include_enclosure,omitempty"`
		IncludeContentDiff *bool  `url:"include_content_diff,omitempty"`
	}

	if entryOpts, ok := opts.(*EntryOptions); ok {
		if len(entryOpts.IDs) > 0 {
			idStrs := make([]string, len(entryOpts.IDs))
			for i, id := range entryOpts.IDs {
				idStrs[i] = strconv.Itoa(id)
			}
			params.Set("ids", strings.Join(idStrs, ","))
		}
		if entryOpts.Read != nil {
			params.Set("read", strconv.FormatBool(*entryOpts.Read))
		}
		if entryOpts.Starred != nil {
			params.Set("starred", strconv.FormatBool(*entryOpts.Starred))
		}
		if entryOpts.Mode != "" {
			params.Set("mode", entryOpts.Mode)
		}
		if entryOpts.IncludeOriginal != nil {
			params.Set("include_original", strconv.FormatBool(*entryOpts.IncludeOriginal))
		}
		if entryOpts.IncludeEnclosure != nil {
			params.Set("include_enclosure", strconv.FormatBool(*entryOpts.IncludeEnclosure))
		}
		if entryOpts.IncludeContentDiff != nil {
			params.Set("include_content_diff", strconv.FormatBool(*entryOpts.IncludeContentDiff))
		}
	}
}

// handleSubscriptionOptions adds subscription-specific query parameters
func handleSubscriptionOptions(opts interface{}, params url.Values) {
	type SubscriptionOptions struct {
		Since string `url:"since,omitempty"`
		Mode  string `url:"mode,omitempty"`
	}

	if subOpts, ok := opts.(*SubscriptionOptions); ok {
		if subOpts.Since != "" {
			params.Set("since", subOpts.Since)
		}
		if subOpts.Mode != "" {
			params.Set("mode", subOpts.Mode)
		}
	}
}
