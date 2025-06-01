package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/jetbrains/feedbin"
)

func ExampleBasicUsage() {
	// Get credentials from environment variables
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables must be set")
	}

	// Create a new client
	client := feedbin.NewClient(username, password)

	// Validate authentication
	valid, err := client.ValidateAuth()
	if err != nil {
		log.Fatalf("Error validating authentication: %v", err)
	}

	if !valid {
		log.Fatal("Invalid credentials")
	}

	fmt.Println("Authentication successful!")

	// Get subscriptions
	subscriptions, paginationInfo, err := client.GetSubscriptions(nil, false)
	if err != nil {
		log.Fatalf("Error getting subscriptions: %v", err)
	}

	fmt.Printf("Found %d subscriptions\n", len(subscriptions))
	fmt.Printf("Total subscriptions: %d\n", paginationInfo.TotalCount)

	// Print the first 5 subscriptions
	for i, subscription := range subscriptions {
		if i >= 5 {
			break
		}
		fmt.Printf("  %d. %s (%s)\n", i+1, subscription.Title, subscription.FeedURL)
	}

	// Get unread entries
	unreadEntries, _, err := client.GetUnreadEntries(1, 5)
	if err != nil {
		log.Fatalf("Error getting unread entries: %v", err)
	}

	fmt.Printf("\nUnread entries (showing first 5):\n")
	for i, entry := range unreadEntries {
		title := "No title"
		if entry.Title != nil {
			title = *entry.Title
		}
		fmt.Printf("  %d. %s\n", i+1, title)
	}

	// Get tags
	tags, err := client.GetTags()
	if err != nil {
		log.Fatalf("Error getting tags: %v", err)
	}

	fmt.Printf("\nTags:\n")
	for i, tag := range tags {
		fmt.Printf("  %d. %s (ID: %d)\n", i+1, tag.Name, tag.ID)
	}

	fmt.Println("\nExample completed successfully!")
}
