# Jules Feedbin Go Client

This is a Go client library for the [Feedbin API V2](https://github.com/feedbin/feedbin-api). It provides methods for interacting with all documented API endpoints.

The client is structured with a main `Client` object that provides access to various services, each corresponding to a group of API endpoints (e.g., Subscriptions, Entries, Taggings).

## Installation

To use this client in your Go project:

```bash
go get github.com/your-username/jules-feedbin-client # Replace with actual path if hosted
```
*(Note: The actual import path will depend on where this client is hosted.)*

Initialize the client in your Go code:
```go
import "path/to/jules-feedbin-client" // Replace with actual import path

func main() {
    username := "your-feedbin-username"
    password := "your-feedbin-password"

    client := feedbin.NewClient(username, password)

    // You can now use the client's services
    // For example, to verify authentication:
    ok, _, err := client.Authentication.Verify()
    if err != nil {
        log.Fatalf("Authentication check failed: %v", err)
    }
    if !ok {
        log.Fatal("Authentication failed: Invalid credentials.")
    }
    fmt.Println("Authentication successful!")
}
```

## Features

*   Full coverage of the Feedbin API V2.
*   Helper methods for pagination and date handling.
*   Service-oriented architecture.
*   Uses HTTP Basic Authentication for the main API.
*   Supports Full Content Extraction API (via `extract.feedbin.com`).

## Implemented API Endpoints

The client provides access to the following Feedbin API V2 resources through dedicated services:

*   **Authentication**:
    *   `Verify()`: Verifies client credentials.
*   **Client**:
    *   `NewClient(username, password)`
    *   `SetBaseURL(url)`
    *   `SetUserAgent(ua)`
    *   `SetTimeout(duration)`
*   **Subscriptions**:
    *   `List(opts *SubscriptionListOptions)`
    *   `Get(id int64, opts *SubscriptionGetOptions)`
    *   `Create(opts *CreateSubscriptionOptions)`
    *   `Update(id int64, opts *UpdateSubscriptionOptions)`
    *   `Delete(id int64)`
*   **Entries**:
    *   `List(opts *EntryListOptions)`
    *   `Get(id int64, opts *EntryGetOptions)`
    *   `ListByFeed(feedID int64, opts *EntryListOptions)`
*   **Unread Entries**:
    *   `List(opts *UnreadEntryListOptions)`: Get IDs of unread entries.
    *   `Create(entryIDs []int64)`: Mark entries as unread.
    *   `Delete(entryIDs []int64)`: Mark entries as read.
*   **Starred Entries**:
    *   `List(opts *StarredEntryListOptions)`: Get IDs of starred entries.
    *   `Create(entryIDs []int64)`: Star entries.
    *   `Delete(entryIDs []int64)`: Unstar entries.
*   **Taggings**:
    *   `List(opts *TaggingListOptions)`
    *   `Create(opts *CreateTaggingOptions)`
    *   `Delete(taggingID int64)`
*   **Tags**:
    *   `List(opts *TagListOptions)`
    *   `Delete(tagID int64)`: Deletes a tag and all its taggings.
*   **Saved Searches**:
    *   `List(opts *SavedSearchListOptions)`
    *   `Get(id int64)`
    *   `Create(opts *CreateSavedSearchOptions)`
    *   `Update(id int64, opts *UpdateSavedSearchOptions)`
    *   `Delete(id int64)`
*   **Recently Read Entries**:
    *   `List(opts *RecentlyReadEntryListOptions)`
    *   `Create(entryID int64, interaction *string)`: Record an entry interaction.
*   **Updated Entries**:
    *   `List(opts *UpdatedEntryListOptions)`: Get IDs of entries updated since a timestamp.
*   **Icons**:
    *   `List(opts *IconListOptions)`: Get all favicons.
*   **Imports**:
    *   `List(opts *ImportListOptions)`: List OPML import statuses.
    *   `Create(opmlFilePath string)`: Upload an OPML file for import.
*   **Pages**:
    *   `Get(entryID int64, opts *PageGetOptions)`: Get processed page content for an entry.
*   **Extract (Full Content Extraction)**:
    *   `NewExtractService(httpClient, apiToken, userAgent)`
    *   `Extract(urlToParse string)`: Fetches parsed content from `extract.feedbin.com`.

## Usage Example

```go
package main

import (
	"fmt"
	"log"
	"time"

	feedbin "path/to/jules-feedbin-client" // Replace with actual import path
)

func main() {
	client := feedbin.NewClient("YOUR_USERNAME", "YOUR_PASSWORD")

	// Verify authentication
	ok, _, err := client.Authentication.Verify()
	if err != nil || !ok {
		log.Fatalf("Authentication failed: %v", err)
	}
	fmt.Println("Successfully authenticated!")

	// List first 5 subscriptions
	subs, _, err := client.Subscriptions.List(&feedbin.SubscriptionListOptions{
		ListOptions: feedbin.ListOptions{PerPage: 5},
	})
	if err != nil {
		log.Fatalf("Failed to list subscriptions: %v", err)
	}

	fmt.Println("
First 5 Subscriptions:")
	for _, sub := range subs {
		fmt.Printf("- ID: %d, Title: %s, URL: %s
", sub.ID, sub.Title, sub.FeedURL)
	}

	// List unread entry IDs (if any)
	unreadIDs, _, err := client.UnreadEntries.List(nil)
	if err != nil {
		log.Fatalf("Failed to list unread entry IDs: %v", err)
	}
	if len(unreadIDs) > 0 {
		fmt.Printf("
Found %d unread entries. First few IDs: %v
", len(unreadIDs), unreadIDs[:min(5, len(unreadIDs))])

		// Example: Mark first unread entry as read (removes from unread)
		// For a real scenario, you might mark specific entries read after user interaction.
		// _, _, err = client.UnreadEntries.Delete([]int64{unreadIDs[0]})
		// if err != nil {
		// 	log.Printf("Failed to mark entry %d as read: %v", unreadIDs[0], err)
		// } else {
		// 	fmt.Printf("Successfully marked entry %d as read.
", unreadIDs[0])
		// }
	} else {
		fmt.Println("
No unread entries found.")
	}

	// Example: List entries (first 2, extended mode)
	entries, _, err := client.Entries.List(&feedbin.EntryListOptions{
		ListOptions: feedbin.ListOptions{PerPage: 2},
		Mode:        "extended",
	})
	if err != nil {
		log.Fatalf("Failed to list entries: %v", err)
	}
	fmt.Println("
First 2 Entries (extended mode):")
	for _, entry := range entries {
		fmt.Printf("- ID: %d, Title: %s, Published: %s
", entry.ID, entry.Title, entry.PublishedAt.Format(time.RFC1123))
		if entry.Image != nil {
			fmt.Printf("  Image URL: %s
", entry.Image.URL)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. (A LICENSE file would need to be created separately).
