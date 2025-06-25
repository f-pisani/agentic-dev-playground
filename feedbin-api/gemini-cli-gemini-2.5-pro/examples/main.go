package main

import (
	"fmt"
	"os"

	"feedbin-api/gemini-cli-gemini-2.5-pro"
)

func main() {
	user := os.Getenv("FEEDBIN_USER")
	pass := os.Getenv("FEEDBIN_PASS")

	if user == "" || pass == "" {
		fmt.Println("Please set FEEDBIN_USER and FEEDBIN_PASS environment variables")
		os.Exit(1)
	}

	client := feedbin.New(user, pass)

	authenticated, err := client.Authenticate()
	if err != nil {
		if apiErr, ok := err.(*feedbin.APIError); ok {
			if apiErr.Response != nil {
				fmt.Printf("Failed to authenticate: %s (status code: %d)\n", apiErr.Message, apiErr.Response.StatusCode)
			} else {
				fmt.Printf("Failed to authenticate: %s\n", apiErr.Message)
			}
		} else {
			fmt.Printf("An unexpected error occurred: %s\n", err)
		}
		os.Exit(1)
	}

	if !authenticated {
		fmt.Println("Invalid credentials")
		os.Exit(1)
	}

	fmt.Println("Successfully authenticated!")

	subs, _, err := client.GetSubscriptions(nil)
	if err != nil {
		fmt.Println("Error getting subscriptions:", err)
		os.Exit(1)
	}

	fmt.Println("Subscriptions:")
	for _, sub := range subs {
		fmt.Printf("- %s\n", sub.Title)
	}
}