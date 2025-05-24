package main

import (
	"fmt"
	"log"
	"os"
	"time"

	feedbin "github.com/example/feedbin-api-client"
)

func main() {
	// Get credentials from environment variables
	email := os.Getenv("FEEDBIN_EMAIL")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if email == "" || password == "" {
		log.Fatal("Please set FEEDBIN_EMAIL and FEEDBIN_PASSWORD environment variables")
	}

	// Create client
	client := feedbin.NewClient(email, password)

	// Test authentication
	fmt.Println("Testing authentication...")
	if err := client.CheckAuthentication(); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	fmt.Println("✓ Authentication successful")

	// Get subscriptions
	fmt.Println("\nFetching subscriptions...")
	subscriptions, err := client.GetSubscriptions()
	if err != nil {
		log.Fatalf("Failed to get subscriptions: %v", err)
	}
	fmt.Printf("✓ Found %d subscriptions\n", len(subscriptions))

	// Display first few subscriptions
	for i, sub := range subscriptions {
		if i >= 5 {
			break
		}
		fmt.Printf("  - %s (%s)\n", sub.Title, sub.FeedURL)
	}

	// Get unread entries
	fmt.Println("\nFetching unread entries...")
	unreadIDs, err := client.GetUnreadEntries()
	if err != nil {
		log.Fatalf("Failed to get unread entries: %v", err)
	}
	fmt.Printf("✓ Found %d unread entries\n", len(unreadIDs))

	// Get first 10 unread entries with content
	if len(unreadIDs) > 0 {
		limit := 10
		if len(unreadIDs) < limit {
			limit = len(unreadIDs)
		}

		entries, err := client.GetEntriesByIDs(unreadIDs[:limit])
		if err != nil {
			log.Fatalf("Failed to get entries: %v", err)
		}

		fmt.Println("\nFirst few unread entries:")
		for i, entry := range entries {
			if i >= 3 {
				break
			}
			fmt.Printf("  - %s\n", entry.Title)
			fmt.Printf("    by %s on %s\n", entry.Author, entry.Published.Format("Jan 2, 2006"))
		}
	}

	// Get starred entries
	fmt.Println("\nFetching starred entries...")
	starredIDs, err := client.GetStarredEntries()
	if err != nil {
		log.Fatalf("Failed to get starred entries: %v", err)
	}
	fmt.Printf("✓ Found %d starred entries\n", len(starredIDs))

	// Get tags
	fmt.Println("\nFetching tags...")
	tags, err := client.GetTags()
	if err != nil {
		log.Fatalf("Failed to get tags: %v", err)
	}
	fmt.Printf("✓ Found %d tags\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("  - %s\n", tag.Name)
	}

	// Demonstrate pagination
	fmt.Println("\nDemonstrating pagination...")
	iterator := client.NewPageIterator("/entries.json", 25,
		feedbin.WithRead(false),
		feedbin.WithSince(time.Now().AddDate(0, 0, -7)))

	pageCount := 0
	var allEntries []feedbin.Entry
	for iterator.HasNext() && pageCount < 3 {
		var pageEntries []feedbin.Entry
		err := iterator.NextPage(&pageEntries)
		if err != nil {
			break
		}
		allEntries = append(allEntries, pageEntries...)
		pageCount++
		fmt.Printf("  Page %d: %d entries\n", iterator.CurrentPage(), len(pageEntries))
	}
	fmt.Printf("✓ Retrieved %d entries across %d pages\n", len(allEntries), pageCount)

	// Demonstrate error handling
	fmt.Println("\nDemonstrating error handling...")
	_, err = client.GetEntry(999999999)
	if err != nil {
		if _, ok := err.(*feedbin.NotFoundError); ok {
			fmt.Println("✓ Correctly handled 404 Not Found error")
		} else {
			fmt.Printf("  Unexpected error: %v\n", err)
		}
	}

	// Create a saved search
	fmt.Println("\nCreating a saved search...")
	search, err := client.CreateSavedSearch("Go Articles", "golang OR \"go programming\"")
	if err != nil {
		fmt.Printf("  Failed to create saved search: %v\n", err)
	} else {
		fmt.Printf("✓ Created saved search: %s (ID: %d)\n", search.Name, search.ID)

		// Clean up - delete the saved search
		if err := client.DeleteSavedSearch(search.ID); err != nil {
			fmt.Printf("  Failed to delete saved search: %v\n", err)
		}
	}

	fmt.Println("\nExample completed successfully!")
}