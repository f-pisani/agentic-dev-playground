# Feedbin API Go Client Usage Examples

## Setup

```
go get feedbin-api/jetbrains-cascade-gpt-4.1
```

## Basic Client Initialization

```go
import (
	"feedbin-api/jetbrains-cascade-gpt-4.1/client"
)

c := client.NewClient("your@email.com", "yourpassword")
```

## Authentication Check

```go
err := c.CheckAuth()
if err != nil {
	// Invalid credentials
}
```

## Fetch Entries (Paginated)

```go
entries, err := c.GetEntries(client.EntriesOptions{Page: 1, PerPage: 20})
```

## Create a Page (Save URL)

```go
entry, err := c.CreatePage(client.PageRequest{
	URL:   "https://example.com/article",
	Title: "Optional Title",
})
```

## List Subscriptions

```go
subs, err := c.GetSubscriptions()
```

## Get Feed Info

```go
feed, err := c.GetFeed(1) // Feed ID
```

## List Feed Icons

```go
icons, err := c.GetIcons()
```

## Import OPML

```go
importData := "<?xml version=\"1.0\"?><opml>...</opml>"
imp, err := c.CreateImport(importData)
```

## Get Tags, Taggings, Starred, Unread, Recently Read

```go
tags, _ := c.GetTags()
taggs, _ := c.GetTaggings()
starred, _ := c.GetStarredEntries()
unread, _ := c.GetUnreadEntries()
recent, _ := c.GetRecentlyReadEntries()
```

## Generate Extract Full Content URL

```go
import "feedbin-api/jetbrains-cascade-gpt-4.1/utils"

extractURL, err := utils.GenerateExtractContentURL("username", "secret", "https://target.com/article")
```

---

- All methods return Go structs matching the Feedbin API.
- All errors are handled in the standard Go way (check `err`).
- Only the Go standard library is required.
