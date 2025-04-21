package main

import (
	"fmt"
	"log"
	"os"
	
	feedbin "github.com/yourusername/feedbin-api/aider-claude-3.7"
)

func main() {
	// Get credentials from environment variables
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")
	
	if username == "" || password == "" {
		log.Fatal("Please set FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables")
	}
	
	fmt.Println("Using Feedbin API with username:", username)
	
	// Create a new client
	client := feedbin.NewClient(username, password)
	
	// Get all subscriptions
	fmt.Println("Getting subscriptions...")
	subscriptions, err := client.GetSubscriptions()
	if err != nil {
		log.Fatalf("Error getting subscriptions: %v", err)
	}
	
	fmt.Printf("Found %d subscriptions\n", len(subscriptions))
	for i, sub := range subscriptions {
		if i >= 5 {
			fmt.Println("...")
			break
		}
		fmt.Printf("  - %s (%s)\n", sub.Title, sub.FeedURL)
	}
	
	// Get unread entries
	fmt.Println("\nGetting unread entries...")
	unreadIDs, err := client.GetUnreadEntries()
	if err != nil {
		log.Fatalf("Error getting unread entries: %v", err)
	}
	
	fmt.Printf("Found %d unread entries\n", len(unreadIDs))
	
	// Get starred entries
	fmt.Println("\nGetting starred entries...")
	starredIDs, err := client.GetStarredEntries()
	if err != nil {
		log.Fatalf("Error getting starred entries: %v", err)
	}
	
	fmt.Printf("Found %d starred entries\n", len(starredIDs))
	
	// Get tags
	fmt.Println("\nGetting tags...")
	tags, err := client.GetTags()
	if err != nil {
		log.Fatalf("Error getting tags: %v", err)
	}
	
	fmt.Printf("Found %d tags\n", len(tags))
	for i, tag := range tags {
		if i >= 5 {
			fmt.Println("...")
			break
		}
		fmt.Printf("  - %s (ID: %d)\n", tag.Name, tag.ID)
	}
	
	fmt.Println("\nFeedbin API client example completed successfully!")
}
