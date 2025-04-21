/*
This is an example of how to use the Feedbin API client.
To run this example, you need to:

1. Set up a proper Go module for the Feedbin client
2. Update the import path to match your module structure
3. Set the FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables

Example of setting up a Go module:
```
mkdir -p ~/go/src/github.com/yourusername/feedbin
cp -r /path/to/feedbin-api/vscode-roocode-claude-3.7-architect/* ~/go/src/github.com/yourusername/feedbin/
cd ~/go/src/github.com/yourusername/feedbin
go mod init github.com/yourusername/feedbin
```

Then you can run this example with:
```
cd ~/go/src/github.com/yourusername/feedbin/example
export FEEDBIN_USERNAME=your-email@example.com
export FEEDBIN_PASSWORD=your-password
go run main.go
```
*/

package main

import (
	"fmt"
	"log"
	"os"

	// Update this import path to match your actual module structure
	// For example: "github.com/yourusername/feedbin"
	"github.com/yourusername/feedbin"
)

func main() {
	// Get credentials from environment variables
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")

	// Check if credentials are provided
	if username == "" || password == "" {
		log.Fatal("FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables must be set")
	}

	// Create a new client
	client := feedbin.NewClient(username, password)

	// Verify authentication
	valid, err := client.Authentication.Verify()
	if err != nil {
		log.Fatalf("Error verifying credentials: %v", err)
	}
	if !valid {
		log.Fatal("Invalid credentials")
	}
	fmt.Println("Authentication successful")

	// List subscriptions
	subscriptions, err := client.Subscriptions.List(nil, false)
	if err != nil {
		log.Fatalf("Error listing subscriptions: %v", err)
	}
	fmt.Printf("Found %d subscriptions\n", len(subscriptions))
	for i, sub := range subscriptions {
		if i < 5 { // Print only the first 5 subscriptions
			fmt.Printf("  - %s (%s)\n", sub.Title, sub.FeedURL)
		}
	}

	// Get entries
	opts := &feedbin.ListEntriesOptions{
		PerPage: feedbin.Int(10),
		Starred: feedbin.Bool(true),
	}
	entries, pagination, err := client.Entries.List(opts)
	if err != nil {
		log.Fatalf("Error listing entries: %v", err)
	}
	fmt.Printf("Found %d starred entries\n", len(entries))
	for i, entry := range entries {
		if i < 3 { // Print only the first 3 entries
			title := "No title"
			if entry.Title != nil {
				title = *entry.Title
			}
			fmt.Printf("  - %s (ID: %d)\n", title, entry.ID)
		}
	}
	if pagination.Next != "" {
		fmt.Printf("More entries available at: %s\n", pagination.Next)
	}

	// Get unread entries
	unreadIDs, err := client.UnreadEntries.List()
	if err != nil {
		log.Fatalf("Error listing unread entries: %v", err)
	}
	fmt.Printf("Found %d unread entries\n", len(unreadIDs))

	// Get starred entries
	starredIDs, err := client.StarredEntries.List()
	if err != nil {
		log.Fatalf("Error listing starred entries: %v", err)
	}
	fmt.Printf("Found %d starred entries\n", len(starredIDs))

	// Get tags
	tags, err := client.Tags.List()
	if err != nil {
		log.Fatalf("Error listing tags: %v", err)
	}
	fmt.Printf("Found %d tags\n", len(tags))
	for i, tag := range tags {
		if i < 5 { // Print only the first 5 tags
			fmt.Printf("  - %s (ID: %d)\n", tag.Name, tag.ID)
		}
	}

	// Get saved searches
	searches, err := client.SavedSearches.List()
	if err != nil {
		log.Fatalf("Error listing saved searches: %v", err)
	}
	fmt.Printf("Found %d saved searches\n", len(searches))
	for i, search := range searches {
		if i < 3 { // Print only the first 3 saved searches
			fmt.Printf("  - %s: %s (ID: %d)\n", search.Name, search.Query, search.ID)
		}
	}

	// Example: Create a subscription (commented out to prevent accidental creation)
	/*
		newSubscription, err := client.Subscriptions.Create("https://example.com/feed.xml")
		if err != nil {
			log.Fatalf("Error creating subscription: %v", err)
		}
		fmt.Printf("Created subscription: %s (ID: %d)\n", newSubscription.Title, newSubscription.ID)
	*/

	// Example: Mark entries as read (commented out to prevent accidental marking)
	/*
		if len(unreadIDs) > 0 {
			// Mark the first unread entry as read
			markedIDs, err := client.UnreadEntries.MarkAsRead(unreadIDs[:1])
			if err != nil {
				log.Fatalf("Error marking entries as read: %v", err)
			}
			fmt.Printf("Marked %d entries as read\n", len(markedIDs))
		}
	*/

	// Example: Star an entry (commented out to prevent accidental starring)
	/*
		if len(entries) > 0 {
			// Star the first entry
			starredIDs, err := client.StarredEntries.Star([]int{entries[0].ID})
			if err != nil {
				log.Fatalf("Error starring entry: %v", err)
			}
			fmt.Printf("Starred %d entries\n", len(starredIDs))
		}
	*/
}
