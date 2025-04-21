package feedbin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("user", "pass")
	
	if client.username != "user" {
		t.Errorf("Expected username to be 'user', got '%s'", client.username)
	}
	
	if client.password != "pass" {
		t.Errorf("Expected password to be 'pass', got '%s'", client.password)
	}
	
	if client.baseURL.String() != BaseURL {
		t.Errorf("Expected baseURL to be '%s', got '%s'", BaseURL, client.baseURL.String())
	}
}

func TestNewRequest(t *testing.T) {
	client := NewClient("user", "pass")
	
	req, err := client.NewRequest("GET", "/test.json", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}
	
	// Test that the base URL was used
	if req.URL.Scheme != "https" {
		t.Errorf("Expected scheme to be 'https', got '%s'", req.URL.Scheme)
	}
	
	if req.URL.Host != "api.feedbin.com" {
		t.Errorf("Expected host to be 'api.feedbin.com', got '%s'", req.URL.Host)
	}
	
	if req.URL.Path != "/v2/test.json" {
		t.Errorf("Expected path to be '/v2/test.json', got '%s'", req.URL.Path)
	}
	
	// Test that the user agent was set
	userAgent := req.Header.Get("User-Agent")
	if userAgent != UserAgent {
		t.Errorf("Expected User-Agent to be '%s', got '%s'", UserAgent, userAgent)
	}
	
	// Test that basic auth was set
	username, password, ok := req.BasicAuth()
	if !ok {
		t.Error("Expected basic auth to be set")
	}
	if username != "user" {
		t.Errorf("Expected username to be 'user', got '%s'", username)
	}
	if password != "pass" {
		t.Errorf("Expected password to be 'pass', got '%s'", password)
	}
}

func TestParseFeedbinTime(t *testing.T) {
	// Test UTC format
	utcStr := "2020-01-02T15:04:05.123456Z"
	utcTime, err := ParseFeedbinTime(utcStr)
	if err != nil {
		t.Fatalf("ParseFeedbinTime returned error: %v", err)
	}
	
	expected := time.Date(2020, 1, 2, 15, 4, 5, 123456000, time.UTC)
	if !utcTime.Equal(expected) {
		t.Errorf("Expected time to be '%v', got '%v'", expected, utcTime)
	}
	
	// Test timezone offset format
	tzStr := "2020-01-02T10:04:05.123456-05:00"
	tzTime, err := ParseFeedbinTime(tzStr)
	if err != nil {
		t.Fatalf("ParseFeedbinTime returned error: %v", err)
	}
	
	// Should be converted to UTC
	if !tzTime.Equal(expected) {
		t.Errorf("Expected time to be '%v', got '%v'", expected, tzTime)
	}
}

func TestFormatFeedbinTime(t *testing.T) {
	timeObj := time.Date(2020, 1, 2, 15, 4, 5, 123456000, time.UTC)
	formatted := FormatFeedbinTime(timeObj)
	expected := "2020-01-02T15:04:05.123456Z"
	
	if formatted != expected {
		t.Errorf("Expected formatted time to be '%s', got '%s'", expected, formatted)
	}
}

func TestGetPaginationLinks(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	
	resp.Header.Set("Links", `<https://api.feedbin.com/v2/entries.json?page=2>; rel="next", <https://api.feedbin.com/v2/entries.json?page=5>; rel="last"`)
	
	links := GetPaginationLinks(resp)
	
	if links["next"] != "https://api.feedbin.com/v2/entries.json?page=2" {
		t.Errorf("Expected next link to be 'https://api.feedbin.com/v2/entries.json?page=2', got '%s'", links["next"])
	}
	
	if links["last"] != "https://api.feedbin.com/v2/entries.json?page=5" {
		t.Errorf("Expected last link to be 'https://api.feedbin.com/v2/entries.json?page=5', got '%s'", links["last"])
	}
}

func TestGetTotalCount(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	
	resp.Header.Set("X-Feedbin-Record-Count", "42")
	
	count := GetTotalCount(resp)
	
	if count != 42 {
		t.Errorf("Expected count to be 42, got %d", count)
	}
}
