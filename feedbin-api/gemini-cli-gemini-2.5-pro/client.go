package feedbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL = "https://api.feedbin.com/v2/"
)

// Client is a client for the Feedbin API.
type Client struct {
	User     string
	Password string
	client   *http.Client
}

// New creates a new Feedbin API client.
func New(user, password string) *Client {
	return &Client{
		User:     user,
		Password: password,
		client:   http.DefaultClient,
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = u.Path + rel.Path
	u.RawQuery = rel.RawQuery

	var reqBody []byte
	if body != nil {
		// check if body is a string reader
		if _, ok := body.(*strings.Reader); !ok {
			reqBody, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		}
	}

	var req *http.Request
	if _, ok := body.(*strings.Reader); ok {
		req, err = http.NewRequest(method, u.String(), body.(*strings.Reader))
	} else {
		req, err = http.NewRequest(method, u.String(), bytes.NewBuffer(reqBody))
	}
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.User, c.Password)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return req, nil
}

type Response struct {
	Links       *Links
	RecordCount int
}

type Links struct {
	First string
	Prev  string
	Next  string
	Last  string
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, *Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return resp, nil, &APIError{Response: resp, Message: "request failed"}
	}

	response := &Response{}
	if linkHeader := resp.Header.Get("Link"); linkHeader != "" {
		response.Links = parseLinkHeader(linkHeader)
	}
	if recordCountHeader := resp.Header.Get("X-Feedbin-Record-Count"); recordCountHeader != "" {
		response.RecordCount, _ = strconv.Atoi(recordCountHeader)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, nil, err
		}
	}

	return resp, response, nil
}

func parseLinkHeader(header string) *Links {
	links := &Links{}
	parts := strings.Split(header, ",")
	for _, part := range parts {
		segments := strings.Split(strings.TrimSpace(part), ";")
		if len(segments) < 2 {
			continue
		}

		linkURL := strings.Trim(segments[0], "<>")

		for _, segment := range segments[1:] {
			relParts := strings.Split(strings.TrimSpace(segment), "=")
			if len(relParts) < 2 || relParts[0] != "rel" {
				continue
			}

			relValue := strings.Trim(relParts[1], `"`)
			switch relValue {
			case "first":
				links.First = linkURL
			case "prev":
				links.Prev = linkURL
			case "next":
				links.Next = linkURL
			case "last":
				links.Last = linkURL
			}
		}
	}
	return links
}

// Authenticate checks if the user's credentials are valid.
func (c *Client) Authenticate() (bool, error) {
	req, err := c.newRequest("GET", "authentication.json", nil)
	if err != nil {
		return false, err
	}

	httpResp, _, err := c.do(req, nil)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			if apiErr.Response.StatusCode == http.StatusUnauthorized {
				return false, nil // Correctly indicates authentication failure without an unexpected error.
			}
		}
		return false, err // Return other errors (network, etc.)
	}

	return httpResp.StatusCode == http.StatusOK, nil
}

type GetSubscriptionsOptions struct {
	Since time.Time
	Mode  string
}

// GetSubscriptions returns all subscriptions.
func (c *Client) GetSubscriptions(opt *GetSubscriptionsOptions) ([]Subscription, *Response, error) {
	path := "subscriptions.json"
	params := url.Values{}
	if opt != nil {
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339Nano))
		}
		if opt.Mode != "" {
			params.Add("mode", opt.Mode)
		}
	}
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var subscriptions []Subscription
	_, resp, err := c.do(req, &subscriptions)
	return subscriptions, resp, err
}

// GetSubscription returns a single subscription.
func (c *Client) GetSubscription(id int64) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	_, _, err = c.do(req, &subscription)
	return &subscription, err
}

// CreateSubscription creates a new subscription.
func (c *Client) CreateSubscription(feedURL string) (*Subscription, error) {
	body := map[string]string{"feed_url": feedURL}
	req, err := c.newRequest("POST", "subscriptions.json", body)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	_, _, err = c.do(req, &subscription)
	return &subscription, err
}

// DeleteSubscription deletes a subscription.
func (c *Client) DeleteSubscription(id int64) error {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, _, err = c.do(req, nil)
	return err
}

// UpdateSubscription updates a subscription.
func (c *Client) UpdateSubscription(id int64, title string, usePostAlternative ...bool) (*Subscription, error) {
	body := map[string]string{"title": title}
	path := fmt.Sprintf("subscriptions/%d.json", id)
	method := "PATCH"
	if len(usePostAlternative) > 0 && usePostAlternative[0] {
		method = "POST"
		path = fmt.Sprintf("subscriptions/%d/update.json", id)
	}
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	_, _, err = c.do(req, &subscription)
	return &subscription, err
}

type GetEntriesOptions struct {
	Page                 int
	Since                time.Time
	IDs                  []int64
	Read                 *bool
	Starred              *bool
	PerPage              int
	Mode                 string
	IncludeOriginal      bool
	IncludeEnclosure     bool
	IncludeContentDiff   bool
}

// GetEntries returns all entries.
func (c *Client) GetEntries(opt *GetEntriesOptions) ([]Entry, *Response, error) {
	path := "entries.json"
	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Add("page", strconv.Itoa(opt.Page))
		}
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339Nano))
		}
		if len(opt.IDs) > 0 {
			var idStrs []string
			for _, id := range opt.IDs {
				idStrs = append(idStrs, strconv.FormatInt(id, 10))
			}
			params.Add("ids", strings.Join(idStrs, ","))
		}
		if opt.Read != nil {
			params.Add("read", strconv.FormatBool(*opt.Read))
		}
		if opt.Starred != nil {
			params.Add("starred", strconv.FormatBool(*opt.Starred))
		}
		if opt.PerPage > 0 {
			params.Add("per_page", strconv.Itoa(opt.PerPage))
		}
		if opt.Mode != "" {
			params.Add("mode", opt.Mode)
		}
		if opt.IncludeOriginal {
			params.Add("include_original", "true")
		}
		if opt.IncludeEnclosure {
			params.Add("include_enclosure", "true")
		}
		if opt.IncludeContentDiff {
			params.Add("include_content_diff", "true")
		}
	}
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []Entry
	_, resp, err := c.do(req, &entries)
	return entries, resp, err
}

// GetFeedEntries returns all entries for a specific feed.
func (c *Client) GetFeedEntries(feedID int64, opt *GetEntriesOptions) ([]Entry, *Response, error) {
	path := fmt.Sprintf("feeds/%d/entries.json", feedID)
	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Add("page", strconv.Itoa(opt.Page))
		}
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339Nano))
		}
		if opt.Read != nil {
			params.Add("read", strconv.FormatBool(*opt.Read))
		}
		if opt.Starred != nil {
			params.Add("starred", strconv.FormatBool(*opt.Starred))
		}
		if opt.PerPage > 0 {
			params.Add("per_page", strconv.Itoa(opt.PerPage))
		}
		if opt.Mode != "" {
			params.Add("mode", opt.Mode)
		}
		if opt.IncludeOriginal {
			params.Add("include_original", "true")
		}
		if opt.IncludeEnclosure {
			params.Add("include_enclosure", "true")
		}
		if opt.IncludeContentDiff {
			params.Add("include_content_diff", "true")
		}
	}
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []Entry
	_, resp, err := c.do(req, &entries)
	return entries, resp, err
}

// GetEntry returns a single entry.
func (c *Client) GetEntry(id int64) (*Entry, error) {
	path := fmt.Sprintf("entries/%d.json", id)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var entry Entry
	_, _, err = c.do(req, &entry)
	return &entry, err
}

// GetUnreadEntries returns all unread entry IDs.
func (c *Client) GetUnreadEntries() ([]int64, *Response, error) {
	req, err := c.newRequest("GET", "unread_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var unreadEntries []int64
	_, resp, err := c.do(req, &unreadEntries)
	return unreadEntries, resp, err
}

// MarkAsUnread marks the given entry IDs as unread.
func (c *Client) MarkAsUnread(ids []int64) ([]int64, error) {
	body := map[string][]int64{"unread_entries": ids}
	req, err := c.newRequest("POST", "unread_entries.json", body)
	if err != nil {
		return nil, err
	}

	var unreadEntries []int64
	_, _, err = c.do(req, &unreadEntries)
	return unreadEntries, err
}

// MarkAsRead marks the given entry IDs as read.
func (c *Client) MarkAsRead(ids []int64, usePostAlternative ...bool) ([]int64, error) {
	body := map[string][]int64{"unread_entries": ids}
	path := "unread_entries.json"
	method := "DELETE"
	if len(usePostAlternative) > 0 && usePostAlternative[0] {
		method = "POST"
		path = "unread_entries/delete.json"
	}
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	var readEntries []int64
	_, _, err = c.do(req, &readEntries)
	return readEntries, err
}

// GetStarredEntries returns all starred entry IDs.
func (c *Client) GetStarredEntries() ([]int64, *Response, error) {
	req, err := c.newRequest("GET", "starred_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var starredEntries []int64
	_, resp, err := c.do(req, &starredEntries)
	return starredEntries, resp, err
}

// StarEntries stars the given entry IDs.
func (c *Client) StarEntries(ids []int64) ([]int64, error) {
	body := map[string][]int64{"starred_entries": ids}
	req, err := c.newRequest("POST", "starred_entries.json", body)
	if err != nil {
		return nil, err
	}

	var starredEntries []int64
	_, _, err = c.do(req, &starredEntries)
	return starredEntries, err
}

// UnstarEntries unstars the given entry IDs.
func (c *Client) UnstarEntries(ids []int64, usePostAlternative ...bool) ([]int64, error) {
	body := map[string][]int64{"starred_entries": ids}
	path := "starred_entries.json"
	method := "DELETE"
	if len(usePostAlternative) > 0 && usePostAlternative[0] {
		method = "POST"
		path = "starred_entries/delete.json"
	}
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	var unstarredEntries []int64
	_, _, err = c.do(req, &unstarredEntries)
	return unstarredEntries, err
}

// GetTaggings returns all taggings.
func (c *Client) GetTaggings() ([]Tagging, *Response, error) {
	req, err := c.newRequest("GET", "taggings.json", nil)
	if err != nil {
		return nil, nil, err
	}
	var taggings []Tagging
	_, resp, err := c.do(req, &taggings)
	return taggings, resp, err
}

// GetTagging returns a single tagging.
func (c *Client) GetTagging(id int64) (*Tagging, error) {
	path := fmt.Sprintf("taggings/%d.json", id)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var tagging Tagging
	_, _, err = c.do(req, &tagging)
	return &tagging, err
}

// CreateTagging creates a new tagging.
func (c *Client) CreateTagging(feedID int64, name string) (*Tagging, error) {
	body := map[string]interface{}{"feed_id": feedID, "name": name}
	req, err := c.newRequest("POST", "taggings.json", body)
	if err != nil {
		return nil, err
	}
	var tagging Tagging
	_, _, err = c.do(req, &tagging)
	return &tagging, err
}

// DeleteTagging deletes a tagging.
func (c *Client) DeleteTagging(id int64) error {
	path := fmt.Sprintf("taggings/%d.json", id)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, _, err = c.do(req, nil)
	return err
}

// RenameTag renames a tag.
func (c *Client) RenameTag(oldName, newName string) ([]Tagging, error) {
	body := map[string]string{"old_name": oldName, "new_name": newName}
	req, err := c.newRequest("POST", "tags.json", body)
	if err != nil {
		return nil, err
	}
	var taggings []Tagging
	_, _, err = c.do(req, &taggings)
	return taggings, err
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(name string) ([]Tagging, error) {
	body := map[string]string{"name": name}
	req, err := c.newRequest("DELETE", "tags.json", body)
	if err != nil {
		return nil, err
	}
	var taggings []Tagging
	_, _, err = c.do(req, &taggings)
	return taggings, err
}

// GetSavedSearches returns all saved searches.
func (c *Client) GetSavedSearches() ([]SavedSearch, *Response, error) {
	req, err := c.newRequest("GET", "saved_searches.json", nil)
	if err != nil {
		return nil, nil, err
	}
	var searches []SavedSearch
	_, resp, err := c.do(req, &searches)
	return searches, resp, err
}

type GetSavedSearchEntriesOptions struct {
	IncludeEntries bool
	Page           int
}

// GetSavedSearchEntries returns the entry IDs for a saved search.
func (c *Client) GetSavedSearchEntries(id int64, opt *GetSavedSearchEntriesOptions) ([]int64, []Entry, *Response, error) {
	path := fmt.Sprintf("saved_searches/%d.json", id)
	params := url.Values{}
	if opt != nil {
		if opt.IncludeEntries {
			params.Add("include_entries", "true")
		}
		if opt.Page > 0 {
			params.Add("page", strconv.Itoa(opt.Page))
		}
	}
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	if opt != nil && opt.IncludeEntries {
		var entries []Entry
		_, resp, err := c.do(req, &entries)
		return nil, entries, resp, err
	}

	var entryIDs []int64
	_, resp, err := c.do(req, &entryIDs)
	return entryIDs, nil, resp, err
}

// CreateSavedSearch creates a new saved search.
func (c *Client) CreateSavedSearch(name, query string) (*SavedSearch, error) {
	body := map[string]string{"name": name, "query": query}
	req, err := c.newRequest("POST", "saved_searches.json", body)
	if err != nil {
		return nil, err
	}
	var search SavedSearch
	_, _, err = c.do(req, &search)
	return &search, err
}

// DeleteSavedSearch deletes a saved search.
func (c *Client) DeleteSavedSearch(id int64) error {
	path := fmt.Sprintf("saved_searches/%d.json", id)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, _, err = c.do(req, nil)
	return err
}

// UpdateSavedSearch updates a saved search.
func (c *Client) UpdateSavedSearch(id int64, name string, usePostAlternative ...bool) (*SavedSearch, error) {
	body := map[string]string{"name": name}
	path := fmt.Sprintf("saved_searches/%d.json", id)
	method := "PATCH"
	if len(usePostAlternative) > 0 && usePostAlternative[0] {
		method = "POST"
		path = fmt.Sprintf("saved_searches/%d/update.json", id)
	}
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	var search SavedSearch
	_, _, err = c.do(req, &search)
	return &search, err
}

// GetRecentlyReadEntries returns recently read entry IDs.
func (c *Client) GetRecentlyReadEntries() ([]int64, *Response, error) {
	req, err := c.newRequest("GET", "recently_read_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}
	var entries []int64
	_, resp, err := c.do(req, &entries)
	return entries, resp, err
}

// CreateRecentlyReadEntries creates recently read entries.
func (c *Client) CreateRecentlyReadEntries(ids []int64) ([]int64, error) {
	body := map[string][]int64{"recently_read_entries": ids}
	req, err := c.newRequest("POST", "recently_read_entries.json", body)
	if err != nil {
		return nil, err
	}
	var entries []int64
	_, _, err = c.do(req, &entries)
	return entries, err
}

type GetUpdatedEntriesOptions struct {
	Since time.Time
}

// GetUpdatedEntries returns updated entry IDs.
func (c *Client) GetUpdatedEntries(opt *GetUpdatedEntriesOptions) ([]int64, *Response, error) {
	path := "updated_entries.json"
	params := url.Values{}
	if opt != nil {
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339Nano))
		}
	}
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	var entries []int64
	_, resp, err := c.do(req, &entries)
	return entries, resp, err
}

// MarkUpdatedEntriesAsRead marks updated entries as read.
func (c *Client) MarkUpdatedEntriesAsRead(ids []int64, usePostAlternative ...bool) ([]int64, error) {
	body := map[string][]int64{"updated_entries": ids}
	path := "updated_entries.json"
	method := "DELETE"
	if len(usePostAlternative) > 0 && usePostAlternative[0] {
		method = "POST"
		path = "updated_entries/delete.json"
	}
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	var entries []int64
	_, _, err = c.do(req, &entries)
	return entries, err
}

// GetIcons returns all feed icons.
func (c *Client) GetIcons() ([]Icon, *Response, error) {
	req, err := c.newRequest("GET", "icons.json", nil)
	if err != nil {
		return nil, nil, err
	}
	var icons []Icon
	_, resp, err := c.do(req, &icons)
	return icons, resp, err
}

// CreateImport creates a new import from an OPML file.
func (c *Client) CreateImport(opmlBody string) (*Import, error) {
	req, err := c.newRequest("POST", "imports.json", strings.NewReader(opmlBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml")
	var imp Import
	_, _, err = c.do(req, &imp)
	return &imp, err
}

// GetImports returns all imports.
func (c *Client) GetImports() ([]Import, *Response, error) {
	req, err := c.newRequest("GET", "imports.json", nil)
	if err != nil {
		return nil, nil, err
	}
	var imports []Import
	_, resp, err := c.do(req, &imports)
	return imports, resp, err
}

// GetImport returns a single import.
func (c *Client) GetImport(id int64) (*Import, error) {
	path := fmt.Sprintf("imports/%d.json", id)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var imp Import
	_, _, err = c.do(req, &imp)
	return &imp, err
}

// CreatePage creates a new page entry.
func (c *Client) CreatePage(pageURL, title string) (*Entry, error) {
	body := Page{URL: pageURL}
	if title != "" {
		body.Title = title
	}
	req, err := c.newRequest("POST", "pages.json", body)
	if err != nil {
		return nil, err
	}
	var entry Entry
	_, _, err = c.do(req, &entry)
	return &entry, err
}

// GetFeed returns a single feed.
func (c *Client) GetFeed(id int64) (*Feed, error) {
	path := fmt.Sprintf("feeds/%d.json", id)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var feed Feed
	_, _, err = c.do(req, &feed)
	return &feed, err
}