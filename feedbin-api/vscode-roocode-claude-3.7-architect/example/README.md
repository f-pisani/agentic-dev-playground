# Feedbin API Client Example

This directory contains an example of how to use the Feedbin API client.

## Setup

To run this example, you need to:

1. Set up a proper Go module for the Feedbin client
2. Update the import path in `main.go` to match your module structure
3. Set the FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables

### Setting up a Go module

```bash
# Create a directory for your project
mkdir -p ~/go/src/github.com/yourusername/feedbin

# Copy the Feedbin client files
cp -r /path/to/feedbin-api/vscode-roocode-claude-3.7-architect/* ~/go/src/github.com/yourusername/feedbin/

# Initialize a Go module
cd ~/go/src/github.com/yourusername/feedbin
go mod init github.com/yourusername/feedbin
```

### Update the import path

Edit `example/main.go` and update the import path to match your module structure:

```go
import (
    "fmt"
    "log"
    "os"

    "github.com/yourusername/feedbin" // Update this to match your module
)
```

### Set environment variables

The example uses environment variables for authentication:

```bash
# Set your Feedbin credentials as environment variables
export FEEDBIN_USERNAME=your-email@example.com
export FEEDBIN_PASSWORD=your-password
```

You can also create a `.env` file and use a tool like [godotenv](https://github.com/joho/godotenv) to load the environment variables:

```
FEEDBIN_USERNAME=your-email@example.com
FEEDBIN_PASSWORD=your-password
```

## Running the example

```bash
cd ~/go/src/github.com/yourusername/feedbin/example
go run main.go
```

Or with environment variables inline:

```bash
FEEDBIN_USERNAME=your-email@example.com FEEDBIN_PASSWORD=your-password go run main.go
```

## Example output

The example will:

1. Verify authentication
2. List subscriptions
3. Get starred entries
4. Get unread entries
5. Get tags
6. Get saved searches

It also includes commented-out examples of:

1. Creating a subscription
2. Marking entries as read
3. Starring entries

These are commented out to prevent accidental modifications to your Feedbin account.