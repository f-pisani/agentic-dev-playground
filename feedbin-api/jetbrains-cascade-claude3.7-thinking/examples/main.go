package main

import (
	"fmt"
	"log"
	"os"
	"time"

	client "github.com/feedbin-api/client"
)

func main() {
	// Get authentication credentials from environment variables
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables must be set")
	}

	// Create a new client with default options
	feedbin, err := client.NewClient(username, password)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Test authentication
	authenticated, err := feedbin.TestAuthentication()
	if err != nil {
		log.Fatalf("Error testing authentication: %v", err)
	}

	if !authenticated {
		log.Fatal("Authentication failed, check your credentials")
	}

	fmt.Println("Authentication successful!")

	// Get all subscriptions
	subscriptions, err := feedbin.Subscriptions.List(nil)
	if err != nil {
		log.Fatalf("Error fetching subscriptions: %v", err)
	}

	fmt.Printf("Found %d subscriptions:\n", len(subscriptions))
	for _, sub := range subscriptions {
		fmt.Printf("- %s (%s)\n", sub.Title, sub.FeedURL)
	}

	// Get unread entries
	unreadIDs, err := feedbin.UnreadEntries.List()
	if err != nil {
		log.Fatalf("Error fetching unread entries: %v", err)
	}

	fmt.Printf("You have %d unread entries\n", len(unreadIDs))

	// Get entries with pagination
	// Limit IDs to max 100 as required by the API
	var entryIDs []int64
	if len(unreadIDs) > 100 {
		entryIDs = unreadIDs[:100]
	} else {
		entryIDs = unreadIDs
	}

	if len(entryIDs) > 0 {
		// Get entries by their IDs
		opts := &client.EntryListOptions{
			Ids:              entryIDs,
			IncludeEnclosure: true,
		}

		entries, pagination, err := feedbin.Entries.List(opts)
		if err != nil {
			log.Fatalf("Error fetching entries: %v", err)
		}

		fmt.Printf("Fetched %d entries (Total: %d)\n", len(entries), pagination.TotalCount)

		for i, entry := range entries {
			if i >= 5 {
				break // Just show the first 5 for the example
			}

			title := "Untitled"
			if entry.Title != nil {
				title = *entry.Title
			}

			author := "Unknown author"
			if entry.Author != nil {
				author = *entry.Author
			}

			fmt.Printf("%d. %s by %s\n", i+1, title, author)
			fmt.Printf("   Published: %s\n", entry.Published.Format(time.RFC1123))
			fmt.Printf("   URL: %s\n\n", entry.URL)
		}

		// Mark the first entry as read if we have any
		if len(entries) > 0 {
			_, err := feedbin.UnreadEntries.MarkAsRead([]int64{entries[0].ID})
			if err != nil {
				log.Fatalf("Error marking entry as read: %v", err)
			}
			fmt.Printf("Marked entry %d as read\n", entries[0].ID)
		}
	}

	// Get starred entries
	starredIDs, err := feedbin.StarredEntries.List()
	if err != nil {
		log.Fatalf("Error fetching starred entries: %v", err)
	}

	fmt.Printf("You have %d starred entries\n", len(starredIDs))

	// Get all tags
	tags, err := feedbin.Tags.List()
	if err != nil {
		log.Fatalf("Error fetching tags: %v", err)
	}

	fmt.Printf("You have %d tags:\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("- %s (ID: %d)\n", tag.Name, tag.ID)
	}
}
