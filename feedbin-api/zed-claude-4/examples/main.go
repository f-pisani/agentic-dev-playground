package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	feedbin "github.com/feedbin/feedbin-go"
)

func main() {
	// Get credentials from environment variables
	email := os.Getenv("FEEDBIN_EMAIL")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if email == "" || password == "" {
		log.Fatal("Please set FEEDBIN_EMAIL and FEEDBIN_PASSWORD environment variables")
	}

	// Create a new client
	client := feedbin.NewClient(email, password)
	ctx := context.Background()

	// Test authentication
	fmt.Println("Testing authentication...")
	if err := client.Authenticate(ctx); err != nil {
		log.Fatal("Authentication failed:", err)
	}
	fmt.Println("✓ Authentication successful")

	// Get subscriptions
	fmt.Println("\nFetching subscriptions...")
	subscriptions, _, err := client.GetSubscriptions(ctx, nil)
	if err != nil {
		log.Fatal("Failed to get subscriptions:", err)
	}
	fmt.Printf("✓ Found %d subscriptions\n", len(subscriptions))

	if len(subscriptions) > 0 {
		fmt.Printf("  First subscription: %s (%s)\n", subscriptions[0].Title, subscriptions[0].FeedURL)
	}

	// Get recent unread entries
	fmt.Println("\nFetching unread entries...")
	unreadIDs, err := client.GetUnreadEntries(ctx)
	if err != nil {
		log.Fatal("Failed to get unread entries:", err)
	}
	fmt.Printf("✓ Found %d unread entries\n", len(unreadIDs))

	// Get some recent entries with filtering
	fmt.Println("\nFetching recent entries...")
	entryOpts := &feedbin.EntryOptions{
		PerPage: feedbin.Int(10),
		Read:    feedbin.Bool(false), // Only unread entries
	}

	entries, entryPagination, err := client.GetEntries(ctx, entryOpts)
	if err != nil {
		log.Fatal("Failed to get entries:", err)
	}
	fmt.Printf("✓ Found %d recent unread entries\n", len(entries))

	if len(entries) > 0 {
		entry := entries[0]
		fmt.Printf("  Latest entry: %s\n", getStringValue(entry.Title))
		fmt.Printf("  Published: %s\n", entry.Published.Format(time.RFC3339))
		fmt.Printf("  Author: %s\n", getStringValue(entry.Author))
	}

	// Demonstrate pagination
	if entryPagination.HasNext() {
		fmt.Println("\nTesting pagination...")
		nextEntries, _, err := client.GetEntriesFromURL(ctx, entryPagination.Next)
		if err != nil {
			log.Printf("Failed to get next page: %v", err)
		} else {
			fmt.Printf("✓ Next page has %d entries\n", len(nextEntries))
		}
	}

	// Get starred entries
	fmt.Println("\nFetching starred entries...")
	starredIDs, err := client.GetStarredEntries(ctx)
	if err != nil {
		log.Fatal("Failed to get starred entries:", err)
	}
	fmt.Printf("✓ Found %d starred entries\n", len(starredIDs))

	// Get taggings
	fmt.Println("\nFetching taggings...")
	taggings, err := client.GetTaggings(ctx)
	if err != nil {
		log.Fatal("Failed to get taggings:", err)
	}
	fmt.Printf("✓ Found %d taggings\n", len(taggings))

	if len(taggings) > 0 {
		fmt.Printf("  First tagging: Feed %d tagged as '%s'\n", taggings[0].FeedID, taggings[0].Name)
	}

	// Get unique tag names
	if len(taggings) > 0 {
		fmt.Println("\nFetching unique tag names...")
		tagNames, err := client.GetUniqueTagNames(ctx)
		if err != nil {
			log.Printf("Failed to get tag names: %v", err)
		} else {
			fmt.Printf("✓ Found %d unique tags: %v\n", len(tagNames), tagNames)
		}
	}

	// Example of creating a subscription (commented out to avoid side effects)
	/*
		fmt.Println("\nTesting subscription creation...")
		newSub, choices, err := client.CreateSubscription(ctx, "https://feeds.feedburner.com/TEDTalks_video")
		if err != nil {
			if apiErr, ok := err.(*feedbin.APIError); ok && apiErr.IsMultipleChoices() {
				fmt.Printf("Multiple feeds found: %+v\n", choices)
			} else {
				log.Printf("Failed to create subscription: %v", err)
			}
		} else {
			fmt.Printf("✓ Created subscription: %s\n", newSub.Title)

			// Clean up by deleting the subscription
			if err := client.DeleteSubscription(ctx, newSub.ID); err != nil {
				log.Printf("Failed to delete test subscription: %v", err)
			} else {
				fmt.Println("✓ Cleaned up test subscription")
			}
		}
	*/

	// Example of advanced entry filtering
	fmt.Println("\nTesting advanced entry filtering...")

	// Get entries from the last week
	weekAgo := time.Now().AddDate(0, 0, -7)
	advancedOpts := &feedbin.EntryOptions{
		Since:   feedbin.Time(weekAgo),
		PerPage: feedbin.Int(5),
		Mode:    feedbin.String("extended"), // Get extended metadata
	}

	recentEntries, _, err := client.GetEntries(ctx, advancedOpts)
	if err != nil {
		log.Printf("Failed to get recent entries: %v", err)
	} else {
		fmt.Printf("✓ Found %d entries from the last week\n", len(recentEntries))

		for i, entry := range recentEntries {
			if i >= 3 { // Limit output
				break
			}
			fmt.Printf("  %d. %s (Feed ID: %d)\n", i+1, getStringValue(entry.Title), entry.FeedID)
		}
	}

	// Example of getting entries for a specific feed
	if len(subscriptions) > 0 {
		feedID := subscriptions[0].FeedID
		fmt.Printf("\nGetting entries for feed %d...\n", feedID)

		feedEntries, _, err := client.GetFeedEntries(ctx, feedID, &feedbin.EntryOptions{
			PerPage: feedbin.Int(3),
		})
		if err != nil {
			log.Printf("Failed to get feed entries: %v", err)
		} else {
			fmt.Printf("✓ Found %d entries for feed %d\n", len(feedEntries), feedID)
		}
	}

	fmt.Println("\n🎉 Example completed successfully!")
}

// Helper function to safely get string values from pointers
func getStringValue(s *string) string {
	if s == nil {
		return "[no title]"
	}
	return *s
}
