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

	// Create a new Feedbin client
	client := feedbin.NewClient(username, password)

	// Validate authentication
	valid, err := client.Authentication.Validate()
	if err != nil {
		log.Fatalf("Error validating credentials: %v", err)
	}

	if !valid {
		log.Fatal("Invalid credentials")
	}

	fmt.Println("Successfully authenticated with Feedbin API")

	// List subscriptions
	subscriptions, _, err := client.Subscriptions.List(nil)
	if err != nil {
		log.Fatalf("Error listing subscriptions: %v", err)
	}

	fmt.Printf("Found %d subscriptions:\n", len(subscriptions))
	for _, sub := range subscriptions {
		fmt.Printf("- %s (ID: %d, Feed ID: %d)\n", sub.Title, sub.ID, sub.FeedID)
	}

	// Get unread entries count
	unreadIDs, _, err := client.Unread.List()
	if err != nil {
		log.Fatalf("Error listing unread entries: %v", err)
	}

	fmt.Printf("You have %d unread entries\n", len(unreadIDs))

	// Get starred entries count
	starredIDs, _, err := client.Starred.List()
	if err != nil {
		log.Fatalf("Error listing starred entries: %v", err)
	}

	fmt.Printf("You have %d starred entries\n", len(starredIDs))

	// List tags
	tags, _, err := client.Tags.List()
	if err != nil {
		log.Fatalf("Error listing tags: %v", err)
	}

	fmt.Printf("Found %d tags:\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("- %s (ID: %d)\n", tag.Name, tag.ID)
	}

	// List saved searches
	searches, _, err := client.SavedSearches.List()
	if err != nil {
		log.Fatalf("Error listing saved searches: %v", err)
	}

	fmt.Printf("Found %d saved searches:\n", len(searches))
	for _, search := range searches {
		fmt.Printf("- %s: %s (ID: %d)\n", search.Name, search.Query, search.ID)
	}
}
