package feedbin

import (
	"fmt"
	"strings"
	"time"
)

// ISO8601 is the layout string for ISO 8601 format used by Feedbin API
const ISO8601 = "2006-01-02T15:04:05.000000Z"

// FormatISO8601 formats a time value according to the ISO 8601 format used by Feedbin API
func FormatISO8601(t time.Time) string {
	return t.UTC().Format(ISO8601)
}

// ParseISO8601 parses a string in ISO 8601 format used by Feedbin API
func ParseISO8601(s string) (time.Time, error) {
	return time.Parse(ISO8601, s)
}

// Bool returns a pointer to the given bool value
// Useful when you need to pass a boolean as a pointer
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the given int value
// Useful when you need to pass an int as a pointer
func Int(v int) *int {
	return &v
}

// String returns a pointer to the given string value
// Useful when you need to pass a string as a pointer
func String(v string) *string {
	return &v
}

// SplitIDList converts a comma-separated string of IDs to a slice of ints
func SplitIDList(list string) []int {
	if list == "" {
		return nil
	}

	strIDs := strings.Split(list, ",")
	ids := make([]int, 0, len(strIDs))

	for _, idStr := range strIDs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		var id int
		if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}
