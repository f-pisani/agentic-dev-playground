package main

import (
	"fmt"
	"log"
	"os"
	"time"

	// Adjust the import path based on where `jules-feedbin-client` is located
	// relative to this example, or if it's in GOPATH/GOROOT.
	// If `jules-feedbin-client` is the module name and this `examples` dir
	// is inside it, the import path would be:
	feedbin "jules-feedbin-client/feedbinapi"
	// If the client were hosted on GitHub, it might be:
	// feedbin "github.com/yourusername/jules-feedbin-client"
)

func main() {
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Please set FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables.")
	}

	client := feedbin.NewClient(username, password)
	client.SetUserAgent("Jules Feedbin Go Client Example/1.0")
	client.SetTimeout(30 * time.Second)

	fmt.Println("Attempting to verify authentication...")
	ok, resp, err := client.Authentication.Verify()
	if err != nil {
		if resp != nil {
			log.Printf("Authentication verification HTTP status: %s", resp.Status)
		}
		log.Fatalf("Authentication verification failed: %v", err)
	}
	if !ok {
		log.Fatal("Authentication failed: Invalid credentials or other issue.")
	}
	fmt.Println("Authentication successful!")
	fmt.Println("-------------------------------------")

	// List subscriptions
	fmt.Println("Fetching subscriptions (first 5)...")
	subs, _, err := client.Subscriptions.List(&feedbin.SubscriptionListOptions{
		ListOptions: feedbin.ListOptions{PerPage: 5},
		Mode:        "extended", // Request extended details
	})
	if err != nil {
		log.Fatalf("Error listing subscriptions: %v", err)
	}

	if len(subs) == 0 {
		fmt.Println("No subscriptions found.")
	} else {
		fmt.Printf("Found %d subscriptions:\n", len(subs))
		for _, sub := range subs {
			fmt.Printf("- ID: %d, Title: %q, Feed URL: %s, Site URL: %s\n", sub.ID, sub.Title, sub.FeedURL, sub.SiteURL)
			if sub.IsSpark != nil && *sub.IsSpark {
				fmt.Printf("  (This is a Spark feed!)\n")
			}
			if sub.HasIcon != nil && *sub.HasIcon && sub.Icon != nil {
				fmt.Printf("  Icon Host: %s (Data length: %d, Ext: %s)\n", sub.Icon.Host, len(sub.Icon.Data), sub.Icon.Extension)
			}
		}
	}
	fmt.Println("-------------------------------------")

	// List first 2 entries
	fmt.Println("Fetching entries (first 2)...")
	entries, _, err := client.Entries.List(&feedbin.EntryListOptions{
		ListOptions: feedbin.ListOptions{PerPage: 2},
	})
	if err != nil {
		log.Fatalf("Error listing entries: %v", err)
	}
	if len(entries) == 0 {
		fmt.Println("No entries found.")
	} else {
		fmt.Printf("Found %d entries:\n", len(entries))
		for _, entry := range entries {
			fmt.Printf("- ID: %d, Title: %q, URL: %s, Published: %s\n",
				entry.ID, entry.Title, entry.URL, entry.PublishedAt.Format(time.RFC3339))
		}
	}
	fmt.Println("-------------------------------------")

	// Example: Get unread entry IDs
	fmt.Println("Fetching unread entry IDs...")
	unreadEntryIDs, _, err := client.UnreadEntries.List(nil)
	if err != nil {
		log.Fatalf("Error fetching unread entry IDs: %v", err)
	}
	if len(unreadEntryIDs) > 0 {
		fmt.Printf("Found %d unread entries. First few: %v\n", len(unreadEntryIDs), unreadEntryIDs[:min(len(unreadEntryIDs), 5)])
	} else {
		fmt.Println("No unread entries.")
	}
	fmt.Println("-------------------------------------")

	// Example: Using the Extract service
	// The Extract service token is your Feedbin username, which is handled during client initialization.
	articleURL := "https://www.theverge.com/2023/10/26/23933449/google-ai-search-generative-experience-rollout" // Example URL
	fmt.Printf("Attempting to extract content from: %s\n", articleURL)
	extractedContent, extractResp, err := client.Extract.Extract(articleURL)
	if err != nil {
		if extractResp != nil {
			log.Printf("Extraction HTTP Status: %s\n", extractResp.Status)
		}
		log.Printf("Error extracting content: %v\n", err)
	} else if extractedContent != nil {
		fmt.Printf("Extracted Title: %s\n", extractedContent.Title)
		fmt.Printf("Extracted Author: %s\n", extractedContent.Author)
		// fmt.Printf("Extracted Content (first 100 chars): %.100s...\n", extractedContent.Content)
	}
	fmt.Println("-------------------------------------")
	fmt.Println("Example usage finished.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// To run this example:
// 1. Ensure FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables are set.
// 2. Navigate to the `feedbin-api/jules-feedbin-client/` directory.
// 3. Run `go mod tidy` if you haven't already.
// 4. Run the example: `go run examples/main.go`
//    (Or build: `go build -o example_app examples/main.go` then run `./example_app`)
