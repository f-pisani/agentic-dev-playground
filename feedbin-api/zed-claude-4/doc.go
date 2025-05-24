// Package feedbin provides a Go client library for the Feedbin REST API v2.
//
// The Feedbin API allows you to manage RSS feed subscriptions, entries, and user data.
// This client supports all documented endpoints with proper authentication, pagination,
// and error handling using only the Go standard library.
//
// Quick Start:
//
//	client := feedbin.NewClient("user@example.com", "password")
//
//	// Verify authentication
//	if err := client.Authenticate(context.Background()); err != nil {
//		log.Fatal("Authentication failed:", err)
//	}
//
//	// Get subscriptions
//	subs, _, err := client.GetSubscriptions(context.Background(), nil)
//	if err != nil {
//		log.Fatal("Failed to get subscriptions:", err)
//	}
//
//	// Get recent unread entries
//	opts := &feedbin.EntryOptions{
//		PerPage: feedbin.Int(20),
//		Read:    feedbin.Bool(false),
//	}
//	entries, pagination, err := client.GetEntries(context.Background(), opts)
//
// Authentication:
//
// The client uses HTTP Basic Authentication with email and password.
// All requests require valid credentials.
//
// Error Handling:
//
// The client provides structured error types for different scenarios:
//
//	if err != nil {
//		if apiErr, ok := err.(*feedbin.APIError); ok {
//			switch {
//			case apiErr.IsUnauthorized():
//				// Invalid credentials
//			case apiErr.IsNotFound():
//				// Resource not found
//			case apiErr.IsForbidden():
//				// Access denied
//			}
//		}
//	}
//
// Pagination:
//
// Paginated endpoints return pagination information via Link headers:
//
//	entries, pagination, err := client.GetEntries(ctx, opts)
//	if pagination.HasNext() {
//		nextEntries, _, err := client.GetEntriesFromURL(ctx, pagination.Next)
//	}
//
// For more examples, see the examples/ directory.
package feedbin
