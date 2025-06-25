package feedbin

import (
	"fmt"
	"io"
	"strconv"
)

// Entries

func (c *Client) GetEntries(opts ...RequestOption) ([]*Entry, error) {
	req, err := c.newRequest("GET", "entries.json", nil, opts...)
	if err != nil {
		return nil, err
	}
	var entries []*Entry
	_, err = c.do(req, &entries)
	return entries, err
}

func (c *Client) GetFeedEntries(feedID int64, opts ...RequestOption) ([]*Entry, error) {
	path := fmt.Sprintf("feeds/%d/entries.json", feedID)
	req, err := c.newRequest("GET", path, nil, opts...)
	if err != nil {
		return nil, err
	}
	var entries []*Entry
	_, err = c.do(req, &entries)
	return entries, err
}

func (c *Client) GetEntry(entryID int64, opts ...RequestOption) (*Entry, error) {
	path := fmt.Sprintf("entries/%d.json", entryID)
	req, err := c.newRequest("GET", path, nil, opts...)
	if err != nil {
		return nil, err
	}
	var entry Entry
	_, err = c.do(req, &entry)
	return &entry, err
}

// Feeds

func (c *Client) GetFeed(feedID int64) (*Feed, error) {
	path := fmt.Sprintf("feeds/%d.json", feedID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var feed Feed
	_, err = c.do(req, &feed)
	return &feed, err
}

// Subscriptions

func (c *Client) GetSubscriptions(opts ...RequestOption) ([]*Subscription, error) {
	req, err := c.newRequest("GET", "subscriptions.json", nil, opts...)
	if err != nil {
		return nil, err
	}
	var subs []*Subscription
	_, err = c.do(req, &subs)
	return subs, err
}

func (c *Client) GetSubscription(subID int64) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", subID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	_, err = c.do(req, &sub)
	return &sub, err
}

func (c *Client) CreateSubscription(feedURL string) (*Subscription, error) {
	body := map[string]string{"feed_url": feedURL}
	req, err := c.newRequest("POST", "subscriptions.json", body)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	_, err = c.do(req, &sub)
	return &sub, err
}

func (c *Client) UpdateSubscription(subID int64, title string) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", subID)
	body := map[string]string{"title": title}
	req, err := c.newRequest("PATCH", path, body)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	_, err = c.do(req, &sub)
	return &sub, err
}

func (c *Client) DeleteSubscription(subID int64) error {
	path := fmt.Sprintf("subscriptions/%d.json", subID)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// Taggings

func (c *Client) GetTaggings() ([]*Tagging, error) {
	req, err := c.newRequest("GET", "taggings.json", nil)
	if err != nil {
		return nil, err
	}
	var taggings []*Tagging
	_, err = c.do(req, &taggings)
	return taggings, err
}

func (c *Client) GetTagging(taggingID int64) (*Tagging, error) {
	path := fmt.Sprintf("taggings/%d.json", taggingID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var tagging Tagging
	_, err = c.do(req, &tagging)
	return &tagging, err
}

func (c *Client) CreateTagging(feedID int64, name string) (*Tagging, error) {
	body := map[string]string{
		"feed_id": strconv.FormatInt(feedID, 10),
		"name":    name,
	}
	req, err := c.newRequest("POST", "taggings.json", body)
	if err != nil {
		return nil, err
	}
	var tagging Tagging
	_, err = c.do(req, &tagging)
	return &tagging, err
}

func (c *Client) DeleteTagging(taggingID int64) error {
	path := fmt.Sprintf("taggings/%d.json", taggingID)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// Starred Entries

func (c *Client) GetStarredEntryIDs() ([]int64, error) {
	req, err := c.newRequest("GET", "starred_entries.json", nil)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

func (c *Client) StarEntries(entryIDs []int64) ([]int64, error) {
	body := map[string][]int64{"starred_entries": entryIDs}
	req, err := c.newRequest("POST", "starred_entries.json", body)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

func (c *Client) UnstarEntries(entryIDs []int64) ([]int64, error) {
	body := map[string][]int64{"starred_entries": entryIDs}
	req, err := c.newRequest("DELETE", "starred_entries.json", body)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

// Unread Entries

func (c *Client) GetUnreadEntryIDs() ([]int64, error) {
	req, err := c.newRequest("GET", "unread_entries.json", nil)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

func (c *Client) MarkEntriesAsUnread(entryIDs []int64) ([]int64, error) {
	body := map[string][]int64{"unread_entries": entryIDs}
	req, err := c.newRequest("POST", "unread_entries.json", body)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

func (c *Client) MarkEntriesAsRead(entryIDs []int64) ([]int64, error) {
	body := map[string][]int64{"unread_entries": entryIDs}
	req, err := c.newRequest("DELETE", "unread_entries.json", body)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

// Recently Read Entries

func (c *Client) GetRecentlyReadEntryIDs() ([]int64, error) {
	req, err := c.newRequest("GET", "recently_read_entries.json", nil)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

func (c *Client) CreateRecentlyReadEntries(entryIDs []int64) ([]int64, error) {
	body := map[string][]int64{"recently_read_entries": entryIDs}
	req, err := c.newRequest("POST", "recently_read_entries.json", body)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

// Saved Searches

func (c *Client) GetSavedSearches() ([]*SavedSearch, error) {
	req, err := c.newRequest("GET", "saved_searches.json", nil)
	if err != nil {
		return nil, err
	}
	var searches []*SavedSearch
	_, err = c.do(req, &searches)
	return searches, err
}

func (c *Client) GetSavedSearch(searchID int64, opts ...RequestOption) ([]int64, error) {
	path := fmt.Sprintf("saved_searches/%d.json", searchID)
	req, err := c.newRequest("GET", path, nil, opts...)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

func (c *Client) CreateSavedSearch(name, query string) (*SavedSearch, error) {
	body := map[string]string{"name": name, "query": query}
	req, err := c.newRequest("POST", "saved_searches.json", body)
	if err != nil {
		return nil, err
	}
	var search SavedSearch
	_, err = c.do(req, &search)
	return &search, err
}

func (c *Client) UpdateSavedSearch(searchID int64, name string) (*SavedSearch, error) {
	path := fmt.Sprintf("saved_searches/%d.json", searchID)
	body := map[string]string{"name": name}
	req, err := c.newRequest("PATCH", path, body)
	if err != nil {
		return nil, err
	}
	var search SavedSearch
	_, err = c.do(req, &search)
	return &search, err
}

func (c *Client) DeleteSavedSearch(searchID int64) error {
	path := fmt.Sprintf("saved_searches/%d.json", searchID)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}
// Authentication

func (c *Client) Authenticate() (bool, error) {
	req, err := c.newRequest("GET", "authentication.json", nil)
	if err != nil {
		return false, err
	}
	_, err = c.do(req, nil)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == 401 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
// Icons

func (c *Client) GetIcons() ([]*Icon, error) {
	req, err := c.newRequest("GET", "icons.json", nil)
	if err != nil {
		return nil, err
	}
	var icons []*Icon
	_, err = c.do(req, &icons)
	return icons, err
}

// Imports

func (c *Client) CreateImport(opml io.Reader) (*Import, error) {
	req, err := c.newRequest("POST", "imports.json", opml)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml")
	var imp Import
	_, err = c.do(req, &imp)
	return &imp, err
}

func (c *Client) GetImports() ([]*Import, error) {
	req, err := c.newRequest("GET", "imports.json", nil)
	if err != nil {
		return nil, err
	}
	var imports []*Import
	_, err = c.do(req, &imports)
	return imports, err
}

func (c *Client) GetImport(importID int64) (*Import, error) {
	path := fmt.Sprintf("imports/%d.json", importID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var imp Import
	_, err = c.do(req, &imp)
	return &imp, err
}
// Pages

func (c *Client) CreatePage(page *Page) (*Entry, error) {
	req, err := c.newRequest("POST", "pages.json", page)
	if err != nil {
		return nil, err
	}
	var entry Entry
	_, err = c.do(req, &entry)
	return &entry, err
}

func (c *Client) GetPageContent(pageID int64) (string, error) {
	path := fmt.Sprintf("pages/%d/content.json", pageID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	var pageContent PageContent
	_, err = c.do(req, &pageContent)
	return pageContent.Content, err
}
// Tags

func (c *Client) RenameTag(oldName, newName string) ([]*Tagging, error) {
	body := map[string]string{"old_name": oldName, "new_name": newName}
	req, err := c.newRequest("POST", "tags.json", body)
	if err != nil {
		return nil, err
	}
	var taggings []*Tagging
	_, err = c.do(req, &taggings)
	return taggings, err
}

func (c *Client) DeleteTag(name string) ([]*Tagging, error) {
	body := map[string]string{"name": name}
	req, err := c.newRequest("DELETE", "tags.json", body)
	if err != nil {
		return nil, err
	}
	var taggings []*Tagging
	_, err = c.do(req, &taggings)
	return taggings, err
}
// Updated Entries

func (c *Client) GetUpdatedEntryIDs(opts ...RequestOption) ([]int64, error) {
	req, err := c.newRequest("GET", "updated_entries.json", nil, opts...)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

func (c *Client) DeleteUpdatedEntries(entryIDs []int64) ([]int64, error) {
	body := map[string][]int64{"updated_entries": entryIDs}
	req, err := c.newRequest("DELETE", "updated_entries.json", body)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}

// Twitter Accounts

func (c *Client) GetTwitterAccounts() ([]*TwitterAccount, error) {
	req, err := c.newRequest("GET", "twitter_accounts.json", nil)
	if err != nil {
		return nil, err
	}
	var accounts []*TwitterAccount
	_, err = c.do(req, &accounts)
	return accounts, err
}

func (c *Client) GetTwitterAccount(accountID int64) (*TwitterAccount, error) {
	path := fmt.Sprintf("twitter_accounts/%d.json", accountID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var account TwitterAccount
	_, err = c.do(req, &account)
	return &account, err
}

// Twitter Tweets

func (c *Client) GetTwitterAccountTweets(accountID int64, opts ...RequestOption) ([]*Tweet, error) {
	path := fmt.Sprintf("twitter_accounts/%d/tweets.json", accountID)
	req, err := c.newRequest("GET", path, nil, opts...)
	if err != nil {
		return nil, err
	}
	var tweets []*Tweet
	_, err = c.do(req, &tweets)
	return tweets, err
}

func (c *Client) GetTwitterTweet(tweetID int64) (*Tweet, error) {
	path := fmt.Sprintf("twitter_tweets/%d.json", tweetID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var tweet Tweet
	_, err = c.do(req, &tweet)
	return &tweet, err
}
func (c *Client) UpdateSubscriptionAlt(subID int64, title string) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d/update.json", subID)
	body := map[string]string{"title": title}
	req, err := c.newRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	_, err = c.do(req, &sub)
	return &sub, err
}

func (c *Client) DeleteTagAlt(name string) ([]*Tagging, error) {
	body := map[string]string{"name": name}
	req, err := c.newRequest("POST", "tags/delete.json", body)
	if err != nil {
		return nil, err
	}
	var taggings []*Tagging
	_, err = c.do(req, &taggings)
	return taggings, err
}

func (c *Client) DeleteUpdatedEntriesAlt(entryIDs []int64) ([]int64, error) {
	body := map[string][]int64{"updated_entries": entryIDs}
	req, err := c.newRequest("POST", "updated_entries/delete.json", body)
	if err != nil {
		return nil, err
	}
	var ids []int64
	_, err = c.do(req, &ids)
	return ids, err
}