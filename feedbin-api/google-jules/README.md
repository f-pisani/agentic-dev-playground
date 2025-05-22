# Feedbin API Client for Go (google-jules)

This Go package provides a client for interacting with the [Feedbin API V2](https://github.com/feedbin/feedbin-api/tree/master/specs). It is designed to be idiomatic, use only the Go standard library, and provide access to the various resources offered by the Feedbin API.

## Implementation Plan

The following plan outlines the steps taken to develop this API client:

1.  **Finalize API Endpoint Analysis:**
    *   Thoroughly read all markdown files in `feedbin-api/specs/content/` to fully understand all API endpoints, their parameters, and response structures.
    *   Pay close attention to request types (GET, POST, PUT, DELETE, PATCH), required parameters, and specific response codes for each endpoint.
    *   *Status: Completed.*

2.  **Design Go Package Structure:**
    *   Define the overall package structure for clarity and maintainability. The chosen structure is:
        *   `feedbin-api/google-jules/`
            *   `README.md` (this file)
            *   `feedbinapi/` (The main package directory)
                *   `client.go`: Core `Client` struct, HTTP request handling, authentication, base URL.
                *   `models.go`: Struct definitions for all API resources.
                *   `pagination.go`: Logic for handling paginated responses.
                *   `errors.go`: Custom error types.
                *   `options.go`: Client configuration options.
                *   `resource_*.go`: One file for each API resource type (e.g., `resource_entries.go`, `resource_feeds.go`), containing specific methods for that resource.
    *   *Status: Completed.*

3.  **Write `README.md`:**
    *   Create this `README.md` file in `feedbin-api/google-jules/` detailing the implementation plan and a brief overview of the package.
    *   *Status: In Progress (this step).*

4.  **Implement Core Client (`feedbinapi/client.go`, `feedbinapi/options.go`):**
    *   Define a `Client` struct.
    *   Implement `NewClient(username, password string, options ...ClientOption) (*Client, error)`.
    *   Implement generic request methods (`newRequest`, `do`) for HTTP communication, JSON marshalling/unmarshalling, and basic error handling.
    *   Implement `VerifyCredentials(ctx context.Context) error`.

5.  **Implement Models (`feedbinapi/models.go`):**
    *   Define Go structs for all API resources with correct JSON tags.
    *   Handle nullable fields and timestamps appropriately.

6.  **Implement Pagination (`feedbinapi/pagination.go`):**
    *   Define `PaginationInfo` struct.
    *   Implement `parseLinkHeader` and other helper functions.

7.  **Implement Error Handling (`feedbinapi/errors.go`):**
    *   Define `APIError` and other specific error types.
    *   Ensure non-2xx HTTP responses are converted into Go errors.

8.  **Implement API Endpoints (in `feedbinapi/resource_*.go` files):**
    *   For each resource (Entries, Feeds, Subscriptions, Taggings, Tags, Starred Entries, Unread Entries, Updated Entries, Recently Read Entries, Saved Searches, Imports, Icons, Pages):
        *   Create the corresponding `resource_*.go` file.
        *   Implement methods on the `Client` struct for each endpoint, accepting `context.Context`.
        *   Handle request parameters and pagination.

9.  **Idiomatic Go and Validation:**
    *   Ensure all code is idiomatic Go.
    *   Use `go fmt` to format the code.
    *   Run `go vet` to check for suspicious constructs and ensure no errors.

10. **(Optional) Testing:**
    *   Write basic unit tests for key functionalities (client, representative API calls, pagination, errors).
    *   Consider mocking HTTP requests/responses.

11. **Review and Refine:**
    *   Review the entire codebase for clarity, correctness, and adherence to requirements.

12. **Submit:**
    *   Submit the completed Go package.

## Package Usage (Preliminary)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"path/to/your/project/feedbin-api/google-jules/feedbinapi" // Adjust import path
)

func main() {
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Please set FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables")
	}

	client, err := feedbinapi.NewClient(username, password)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Example: Verify Credentials
	ctx := context.Background()
	if err := client.VerifyCredentials(ctx); err != nil {
		log.Fatalf("Error verifying credentials: %v", err)
	}
	fmt.Println("Credentials verified successfully!")

	// Example: List Subscriptions
	// subscriptions, _, err := client.ListSubscriptions(ctx, nil) // Assuming nil for no params
	// if err != nil {
	//  log.Fatalf("Error listing subscriptions: %v", err)
	// }
	// for _, sub := range subscriptions {
	//  fmt.Printf("Subscription: %s (ID: %d)
", sub.Title, sub.ID)
	// }
}
```

**Note:** This is a preliminary usage example. More detailed documentation will be available via GoDoc once the package is implemented.

## Contributing

This package is being developed by an AI agent. Contributions are not expected at this time.
