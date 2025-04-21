package feedbin

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool {
	return &v
}

// Int64 returns a pointer to the given int64 value.
func Int64(v int64) *int64 {
	return &v
}

// Int returns a pointer to the given int value.
func Int(v int) *int {
	return &v
}

// String returns a pointer to the given string value.
func String(v string) *string {
	return &v
}
