package main

import (
	"fmt"
	"log"
	"os"

	feedbin "github.com/feedbin/go-client"
)

func main() {
	// Get credentials from environment variables
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables must be set")
	}

	// Create a new client
	client := feedbin.NewClient(username, password)

	// Check authentication
	if err := client.Authentication.Verify(); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	fmt.Println("Authentication successful!")

	// Get subscriptions
	subscriptions, _, err := client.Subscriptions.List(nil)
	if err != nil {
		log.Fatalf("Failed to get subscriptions: %v", err)
	}

	fmt.Printf("Found %d subscriptions:\n", len(subscriptions))
	for _, sub := range subscriptions {
		fmt.Printf("- %s (%s)\n", sub.Title, sub.FeedURL)
	}

	// Get unread entries
	unreadIDs, _, err := client.UnreadEntries.List()
	if err != nil {
		log.Fatalf("Failed to get unread entries: %v", err)
	}

	fmt.Printf("Found %d unread entries\n", len(unreadIDs))

	// Get starred entries
	starredIDs, _, err := client.StarredEntries.List()
	if err != nil {
		log.Fatalf("Failed to get starred entries: %v", err)
	}

	fmt.Printf("Found %d starred entries\n", len(starredIDs))

	// Get saved searches
	savedSearches, _, err := client.SavedSearches.List()
	if err != nil {
		log.Fatalf("Failed to get saved searches: %v", err)
	}

	fmt.Printf("Found %d saved searches:\n", len(savedSearches))
	for _, search := range savedSearches {
		fmt.Printf("- %s (%s)\n", search.Name, search.Query)
	}

	// Get feed icons
	icons, _, err := client.Icons.List()
	if err != nil {
		log.Fatalf("Failed to get icons: %v", err)
	}

	fmt.Printf("Found %d icons\n", len(icons))

	// Example of using the extract service
	if len(subscriptions) > 0 && subscriptions[0].SiteURL != "" {
		// Set the extract service secret (in a real app, get this from your Feedbin account)
		client.Extract.SetSecret("your-extract-service-secret")

		fmt.Printf("To extract content from %s, you would use the Extract service with your secret key\n", subscriptions[0].SiteURL)
	}
}
