package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	feedbin "github.com/feedbin/client"
)

func main() {
	// Parse command line flags
	username := flag.String("username", "", "Feedbin API username (email)")
	password := flag.String("password", "", "Feedbin API password")
	flag.Parse()

	// Check if credentials were provided
	if *username == "" || *password == "" {
		fmt.Println("Usage: go run main.go -username your@email.com -password yourpassword")
		fmt.Println("Alternatively, you can set FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables")

		// Try to get credentials from environment variables
		*username = os.Getenv("FEEDBIN_USERNAME")
		*password = os.Getenv("FEEDBIN_PASSWORD")

		if *username == "" || *password == "" {
			fmt.Println("Error: No credentials provided")
			os.Exit(1)
		}
	}

	// Create a new Feedbin client
	client := feedbin.NewClient(*username, *password)

	// Verify authentication
	valid, _, err := client.Authentication.Verify()
	if err != nil {
		log.Fatalf("Error verifying credentials: %v", err)
	}
	if !valid {
		log.Fatalf("Invalid credentials")
	}
	fmt.Println("Authentication successful!")

	// List subscriptions
	subscriptions, _, err := client.Subscriptions.List(nil)
	if err != nil {
		log.Fatalf("Error listing subscriptions: %v", err)
	}
	fmt.Printf("You have %d subscriptions\n", len(subscriptions))

	// Optionally, print subscription details
	if len(subscriptions) > 0 {
		fmt.Println("\nYour subscriptions:")
		for i, sub := range subscriptions {
			fmt.Printf("%d. %s (%s)\n", i+1, sub.Title, sub.SiteURL)
		}
	}

	// Get unread entries count
	unreadIDs, _, err := client.UnreadEntries.List()
	if err != nil {
		log.Fatalf("Error listing unread entries: %v", err)
	}
	fmt.Printf("\nYou have %d unread entries\n", len(unreadIDs))

	// Get starred entries count
	starredIDs, _, err := client.StarredEntries.List()
	if err != nil {
		log.Fatalf("Error listing starred entries: %v", err)
	}
	fmt.Printf("You have %d starred entries\n", len(starredIDs))

	// List saved searches if any
	savedSearches, _, err := client.SavedSearches.List()
	if err != nil {
		log.Fatalf("Error listing saved searches: %v", err)
	}

	if len(savedSearches) > 0 {
		fmt.Println("\nYour saved searches:")
		for i, search := range savedSearches {
			fmt.Printf("%d. %s (query: %s)\n", i+1, search.Name, search.Query)
		}
	} else {
		fmt.Println("\nYou have no saved searches")
	}
}
