package main

import (
	"context"
	"fmt"
	"log"
	"time"

	feedbin "github.com/your-org/feedbin-api/vscode-cline-claude4"
)

func main() {
	// Create a new Feedbin client
	client := feedbin.NewClient(&feedbin.Config{
		Username: "your-email@example.com",
		Password: "your-password",
		// Optional: Enable HTTP caching
		EnableCache: true,
	})

	ctx := context.Background()

	// Example 1: Validate credentials
	fmt.Println("=== Validating Credentials ===")
	if err := client.ValidateCredentials(ctx); err != nil {
		log.Fatal("Invalid credentials:", err)
	}
	fmt.Println("✓ Credentials are valid")

	// Example 2: Get subscriptions
	fmt.Println("\n=== Getting Subscriptions ===")
	subscriptions, pagination, err := client.GetSubscriptions(ctx, nil)
	if err != nil {
		log.Fatal("Failed to get subscriptions:", err)
	}
	fmt.Printf("Found %d subscriptions (total: %d)\n", len(subscriptions), pagination.TotalCount)

	for i, sub := range subscriptions {
		if i >= 3 { // Show only first 3
			break
		}
		fmt.Printf("- %s (%s)\n", sub.Title, sub.FeedURL)
	}

	// Example 3: Create a new subscription
	fmt.Println("\n=== Creating Subscription ===")
	newSub, err := client.CreateSubscription(ctx, "https://feeds.feedburner.com/oreilly/radar")
	if err != nil {
		// Handle multiple choices error
		if multiErr, ok := err.(*feedbin.MultipleChoicesError); ok {
			fmt.Printf("Multiple feeds found:\n")
			for _, choice := range multiErr.GetChoices() {
				fmt.Printf("- %s (%s)\n", choice.Title, choice.FeedURL)
			}
		} else {
			log.Printf("Failed to create subscription: %v", err)
		}
	} else {
		fmt.Printf("✓ Created subscription: %s\n", newSub.Title)
	}

	// Example 4: Get entries with filtering
	fmt.Println("\n=== Getting Recent Entries ===")
	since := time.Now().AddDate(0, 0, -7) // Last 7 days
	entryOpts := &feedbin.EntryOptions{
		Since:   &since,
		PerPage: &[]int{10}[0], // Get 10 entries
	}

	entries, pagination, err := client.GetEntries(ctx, entryOpts)
	if err != nil {
		log.Printf("Failed to get entries: %v", err)
	} else {
		fmt.Printf("Found %d recent entries (total: %d)\n", len(entries), pagination.TotalCount)

		for i, entry := range entries {
			if i >= 3 { // Show only first 3
				break
			}
			title := "No Title"
			if entry.Title != nil {
				title = *entry.Title
			}
			fmt.Printf("- %s\n", title)
		}
	}

	// Example 5: Get unread entries
	fmt.Println("\n=== Getting Unread Entries ===")
	unreadIDs, err := client.GetUnreadEntries(ctx)
	if err != nil {
		log.Printf("Failed to get unread entries: %v", err)
	} else {
		fmt.Printf("Found %d unread entries\n", len(unreadIDs))

		if len(unreadIDs) > 0 {
			// Get details for first few unread entries
			idsToFetch := unreadIDs
			if len(idsToFetch) > 3 {
				idsToFetch = idsToFetch[:3]
			}

			unreadEntries, err := client.GetEntriesByIDs(ctx, idsToFetch)
			if err != nil {
				log.Printf("Failed to get unread entry details: %v", err)
			} else {
				fmt.Println("Recent unread entries:")
				for _, entry := range unreadEntries {
					title := "No Title"
					if entry.Title != nil {
						title = *entry.Title
					}
					fmt.Printf("- %s\n", title)
				}
			}
		}
	}

	// Example 6: Get starred entries
	fmt.Println("\n=== Getting Starred Entries ===")
	starredIDs, err := client.GetStarredEntries(ctx)
	if err != nil {
		log.Printf("Failed to get starred entries: %v", err)
	} else {
		fmt.Printf("Found %d starred entries\n", len(starredIDs))
	}

	// Example 7: Work with tags
	fmt.Println("\n=== Working with Tags ===")
	taggings, err := client.GetTaggings(ctx)
	if err != nil {
		log.Printf("Failed to get taggings: %v", err)
	} else {
		fmt.Printf("Found %d taggings\n", len(taggings))

		// Get unique tag names
		tagNames, err := client.GetAllTagNames(ctx)
		if err != nil {
			log.Printf("Failed to get tag names: %v", err)
		} else {
			fmt.Printf("Unique tags: %v\n", tagNames)
		}
	}

	// Example 8: Pagination with iterator
	fmt.Println("\n=== Using Entry Iterator ===")
	iterator := client.NewEntryIterator(ctx, &feedbin.EntryOptions{
		PerPage: &[]int{5}[0], // Small page size for demo
	})

	pageCount := 0
	totalEntries := 0

	for iterator.HasMore() && pageCount < 3 { // Limit to 3 pages for demo
		entries, err := iterator.Next(ctx)
		if err != nil {
			log.Printf("Iterator error: %v", err)
			break
		}

		pageCount++
		totalEntries += len(entries)
		fmt.Printf("Page %d: %d entries\n", pageCount, len(entries))
	}

	fmt.Printf("Total entries processed: %d\n", totalEntries)

	// Example 9: Error handling
	fmt.Println("\n=== Error Handling Example ===")
	_, err = client.GetSubscription(ctx, 999999) // Non-existent ID
	if err != nil {
		switch e := err.(type) {
		case *feedbin.APIError:
			fmt.Printf("API Error: %d - %s\n", e.StatusCode, e.Message)
			if e.Retryable {
				fmt.Println("This error is retryable")
			}
		default:
			fmt.Printf("Other error: %v\n", err)
		}
	}

	fmt.Println("\n=== Example Complete ===")
}
