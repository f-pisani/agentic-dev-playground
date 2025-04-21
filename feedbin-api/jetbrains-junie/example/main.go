// Example usage of the Feedbin API client
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jetbrains-junie/feedbin"
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

	// Verify credentials
	valid, err := client.CheckAuth()
	if err != nil {
		log.Fatalf("Error checking authentication: %v", err)
	}

	if !valid {
		log.Fatal("Invalid credentials")
	}

	fmt.Println("Authentication successful!")

	// Get subscriptions
	subscriptions, _, err := client.Subscriptions.GetSubscriptions(nil)
	if err != nil {
		log.Fatalf("Error getting subscriptions: %v", err)
	}

	fmt.Printf("Found %d subscriptions:\n", len(subscriptions))
	for _, sub := range subscriptions {
		fmt.Printf("- %s (%s)\n", sub.Title, sub.FeedURL)
	}

	// Get unread entries
	unreadIDs, _, err := client.Unread.GetUnreadEntries()
	if err != nil {
		log.Fatalf("Error getting unread entries: %v", err)
	}

	fmt.Printf("Found %d unread entries\n", len(unreadIDs))

	// Get starred entries
	starredIDs, _, err := client.Starred.GetStarredEntries()
	if err != nil {
		log.Fatalf("Error getting starred entries: %v", err)
	}

	fmt.Printf("Found %d starred entries\n", len(starredIDs))

	// Get tags
	tags, err := client.Tags.GetTags()
	if err != nil {
		log.Fatalf("Error getting tags: %v", err)
	}

	fmt.Printf("Found %d tags:\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("- %s (feeds: %d)\n", tag.Name, len(tag.FeedIDs))
	}

	// Get taggings
	taggings, err := client.Taggings.GetTaggings()
	if err != nil {
		log.Fatalf("Error getting taggings: %v", err)
	}

	fmt.Printf("Found %d taggings\n", len(taggings))

	// Get saved searches
	searches, err := client.SavedSearches.GetSavedSearches()
	if err != nil {
		log.Fatalf("Error getting saved searches: %v", err)
	}

	fmt.Printf("Found %d saved searches:\n", len(searches))
	for _, search := range searches {
		fmt.Printf("- %s (query: %s)\n", search.Name, search.Query)
	}

	// Get entries (limited to 10 for this example)
	options := &feedbin.EntryOptions{
		PageOptions: feedbin.PageOptions{
			PerPage: 10,
		},
	}

	entries, resp, err := client.Entries.GetEntries(options)
	if err != nil {
		log.Fatalf("Error getting entries: %v", err)
	}

	// Parse pagination links
	links := feedbin.ParseLinkHeader(resp)
	totalPages, _ := feedbin.GetTotalPages(links)
	totalRecords, _ := feedbin.GetTotalRecords(resp)

	fmt.Printf("Found %d entries (page 1 of %d, %d total records)\n", len(entries), totalPages, totalRecords)

	for _, entry := range entries {
		title := "No Title"
		if entry.Title != nil {
			title = *entry.Title
		}

		author := "Unknown"
		if entry.Author != nil {
			author = *entry.Author
		}

		fmt.Printf("- %s by %s (published: %s)\n", title, author, entry.Published.Format("2006-01-02"))
	}

	// Get recently read entries
	recentlyRead, err := client.RecentlyRead.GetRecentlyRead()
	if err != nil {
		log.Fatalf("Error getting recently read entries: %v", err)
	}

	fmt.Printf("Found %d recently read entries\n", len(recentlyRead))

	// Get updated entries (from 24 hours ago)
	updatedEntries, err := client.Updated.GetUpdatedEntries(time.Now().Add(-24 * time.Hour))
	if err != nil {
		log.Fatalf("Error getting updated entries: %v", err)
	}

	fmt.Printf("Found %d updated entries in the last 24 hours\n", len(updatedEntries))

	// Get icons for the first subscription (if any exist)
	if len(subscriptions) > 0 {
		icons, err := client.Icons.GetIcons([]int{subscriptions[0].FeedID})
		if err != nil {
			log.Fatalf("Error getting icons: %v", err)
		}

		fmt.Printf("Found %d icons\n", len(icons))
	}

	// Get imports
	imports, err := client.Imports.GetImports()
	if err != nil {
		log.Fatalf("Error getting imports: %v", err)
	}

	fmt.Printf("Found %d imports\n", len(imports))

	// Get pages
	pages, err := client.Pages.GetPages()
	if err != nil {
		log.Fatalf("Error getting pages: %v", err)
	}

	fmt.Printf("Found %d pages\n", len(pages))
	for _, page := range pages {
		fmt.Printf("- %s (%s)\n", page.Title, page.URL)
	}
}
