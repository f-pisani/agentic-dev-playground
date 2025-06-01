package examples

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jetbrains/feedbin"
)

func ExampleComprehensive() {
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

	// Example 1: Get subscriptions with extended info
	fmt.Println("\n=== Example 1: Get Subscriptions ===")
	subscriptions, paginationInfo, err := client.GetSubscriptions(nil, true)
	if err != nil {
		log.Fatalf("Error getting subscriptions: %v", err)
	}

	fmt.Printf("Found %d subscriptions\n", len(subscriptions))
	fmt.Printf("Total subscriptions: %d\n", paginationInfo.TotalCount)

	// Print the first 3 subscriptions
	for i, subscription := range subscriptions {
		if i >= 3 {
			break
		}
		fmt.Printf("  %d. %s (%s)\n", i+1, subscription.Title, subscription.FeedURL)

		// If we have JSON feed info, print some of it
		if subscription.JSONFeed != nil {
			fmt.Printf("     - JSON Feed Title: %s\n", subscription.JSONFeed.Title)
			if subscription.JSONFeed.Icon != "" {
				fmt.Printf("     - Icon: %s\n", subscription.JSONFeed.Icon)
			}
		}
	}

	// Example 2: Get entries with various filters
	fmt.Println("\n=== Example 2: Get Entries with Filters ===")

	// Get entries from the last 7 days
	since := time.Now().AddDate(0, 0, -7)
	options := map[string]interface{}{
		"since":    since,
		"per_page": 5,
		"mode":     "extended",
	}

	entries, _, err := client.GetEntries(options)
	if err != nil {
		log.Fatalf("Error getting entries: %v", err)
	}

	fmt.Printf("Found %d entries from the last 7 days (showing first 5)\n", len(entries))
	for i, entry := range entries {
		title := "No title"
		if entry.Title != nil {
			title = *entry.Title
		}

		author := "Unknown author"
		if entry.Author != nil {
			author = *entry.Author
		}

		fmt.Printf("  %d. %s by %s\n", i+1, title, author)
		fmt.Printf("     - Published: %s\n", entry.Published.Format(time.RFC3339))
		fmt.Printf("     - URL: %s\n", entry.URL)

		// If we have enclosure data (podcast), print it
		if entry.Enclosure != nil {
			fmt.Printf("     - Podcast: %s (%s)\n", entry.Enclosure.EnclosureURL, entry.Enclosure.EnclosureType)
			if entry.Enclosure.ItunesDuration != "" {
				fmt.Printf("     - Duration: %s\n", entry.Enclosure.ItunesDuration)
			}
		}

		// If we have Twitter data, print it
		if entry.TwitterID != nil {
			fmt.Printf("     - Tweet ID: %d\n", *entry.TwitterID)
			if len(entry.TwitterThreadIDs) > 0 {
				fmt.Printf("     - Thread with %d tweets\n", len(entry.TwitterThreadIDs))
			}
		}
	}

	// Example 3: Working with tags
	fmt.Println("\n=== Example 3: Working with Tags ===")
	tags, err := client.GetTags()
	if err != nil {
		log.Fatalf("Error getting tags: %v", err)
	}

	fmt.Printf("Found %d tags\n", len(tags))
	for i, tag := range tags {
		if i >= 5 {
			break
		}
		fmt.Printf("  %d. %s (ID: %d)\n", i+1, tag.Name, tag.ID)
	}

	// Get taggings (associations between tags and feeds)
	taggings, err := client.GetTaggings()
	if err != nil {
		log.Fatalf("Error getting taggings: %v", err)
	}

	fmt.Printf("Found %d taggings\n", len(taggings))
	for i, tagging := range taggings {
		if i >= 3 {
			break
		}
		fmt.Printf("  %d. Feed ID %d is tagged with Tag ID %d\n", i+1, tagging.FeedID, tagging.TagID)
	}

	// Example 4: Unread and starred entries
	fmt.Println("\n=== Example 4: Unread and Starred Entries ===")

	// Get unread entry IDs
	unreadIDs, err := client.GetUnreadEntryIDs()
	if err != nil {
		log.Fatalf("Error getting unread entry IDs: %v", err)
	}

	fmt.Printf("You have %d unread entries\n", len(unreadIDs))

	// Get starred entry IDs
	starredIDs, err := client.GetStarredEntryIDs()
	if err != nil {
		log.Fatalf("Error getting starred entry IDs: %v", err)
	}

	fmt.Printf("You have %d starred entries\n", len(starredIDs))

	// Example 5: Saved searches
	fmt.Println("\n=== Example 5: Saved Searches ===")
	searches, err := client.GetSavedSearches()
	if err != nil {
		log.Fatalf("Error getting saved searches: %v", err)
	}

	fmt.Printf("Found %d saved searches\n", len(searches))
	for i, search := range searches {
		if i >= 3 {
			break
		}
		fmt.Printf("  %d. %s (Query: %s)\n", i+1, search.Name, search.Query)
		fmt.Printf("     - Created: %s\n", search.CreatedAt.Format(time.RFC3339))
	}

	fmt.Println("\nComprehensive example completed successfully!")
}
