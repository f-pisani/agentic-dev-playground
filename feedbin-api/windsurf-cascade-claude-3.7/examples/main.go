package main

import (
	"fmt"
	"log"
	"os"

	"github.com/feedbin/client"
)

func main() {
	// Get credentials from environment variables
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables must be set")
	}

	// Create a new Feedbin client
	client := feedbin.New(username, password)

	// Verify credentials
	valid, err := client.Authentication.Verify()
	if err != nil {
		log.Fatalf("Error verifying credentials: %v", err)
	}

	if !valid {
		log.Fatal("Invalid credentials")
	}

	fmt.Println("Authentication successful!")

	// Get subscriptions
	subscriptions, err := client.Subscriptions.List(nil)
	if err != nil {
		log.Fatalf("Error getting subscriptions: %v", err)
	}

	fmt.Printf("Found %d subscriptions:\n", len(subscriptions))
	for _, sub := range subscriptions {
		fmt.Printf("- %s (%s)\n", sub.Title, sub.FeedURL)
	}

	// Get unread entries
	unreadIDs, err := client.UnreadEntries.List()
	if err != nil {
		log.Fatalf("Error getting unread entries: %v", err)
	}

	fmt.Printf("You have %d unread entries\n", len(unreadIDs))

	// If there are unread entries, get details for the first 5
	if len(unreadIDs) > 0 {
		// Limit to first 5 entries or less
		count := 5
		if len(unreadIDs) < 5 {
			count = len(unreadIDs)
		}

		// Create options for entry retrieval
		options := &feedbin.EntryListOptions{
			IDs: unreadIDs[:count],
		}

		// Get entry details
		entries, err := client.Entries.List(options)
		if err != nil {
			log.Fatalf("Error getting entries: %v", err)
		}

		fmt.Println("\nRecent unread entries:")
		for _, entry := range entries {
			title := "Untitled"
			if entry.Title != nil {
				title = *entry.Title
			}

			author := "Unknown"
			if entry.Author != nil {
				author = *entry.Author
			}

			fmt.Printf("- %s by %s\n", title, author)
		}
	}

	// List saved searches if any
	savedSearches, err := client.SavedSearches.List()
	if err != nil {
		log.Fatalf("Error getting saved searches: %v", err)
	}

	if len(savedSearches) > 0 {
		fmt.Printf("\nYou have %d saved searches:\n", len(savedSearches))
		for _, search := range savedSearches {
			fmt.Printf("- %s: %s\n", search.Name, search.Query)
		}
	} else {
		fmt.Println("\nYou have no saved searches")
	}
}
