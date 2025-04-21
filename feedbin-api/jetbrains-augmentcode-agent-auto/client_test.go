package feedbin

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("username", "password")

	if c.baseURL.String() != BaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.baseURL.String(), BaseURL)
	}

	if c.username != "username" {
		t.Errorf("NewClient username = %v, want %v", c.username, "username")
	}

	if c.password != "password" {
		t.Errorf("NewClient password = %v, want %v", c.password, "password")
	}

	if c.userAgent != UserAgent {
		t.Errorf("NewClient userAgent = %v, want %v", c.userAgent, UserAgent)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient("username", "password")

	req, err := c.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	// Test that the username and password are set correctly
	username, password, ok := req.BasicAuth()
	if !ok {
		t.Fatal("BasicAuth not set")
	}
	if username != "username" {
		t.Errorf("NewRequest username = %v, want %v", username, "username")
	}
	if password != "password" {
		t.Errorf("NewRequest password = %v, want %v", password, "password")
	}

	// Test that the user agent is set correctly
	if got, want := req.Header.Get("User-Agent"), UserAgent; got != want {
		t.Errorf("NewRequest User-Agent = %v, want %v", got, want)
	}

	// Test that the URL is constructed correctly
	if got, want := req.URL.String(), BaseURL+"/test"; got != want {
		t.Errorf("NewRequest URL = %v, want %v", got, want)
	}
}

func TestParseLinkHeader(t *testing.T) {
	tests := []struct {
		name       string
		linkHeader string
		want       *PaginationLinks
	}{
		{
			name:       "empty",
			linkHeader: "",
			want:       &PaginationLinks{},
		},
		{
			name:       "next and last",
			linkHeader: `<https://api.feedbin.com/v2/feeds/1/entries.json?page=2>; rel="next", <https://api.feedbin.com/v2/feeds/1/entries.json?page=5>; rel="last"`,
			want: &PaginationLinks{
				Next: "https://api.feedbin.com/v2/feeds/1/entries.json?page=2",
				Last: "https://api.feedbin.com/v2/feeds/1/entries.json?page=5",
			},
		},
		{
			name:       "all links",
			linkHeader: `<https://api.feedbin.com/v2/feeds/1/entries.json?page=1>; rel="first", <https://api.feedbin.com/v2/feeds/1/entries.json?page=1>; rel="prev", <https://api.feedbin.com/v2/feeds/1/entries.json?page=3>; rel="next", <https://api.feedbin.com/v2/feeds/1/entries.json?page=5>; rel="last"`,
			want: &PaginationLinks{
				First: "https://api.feedbin.com/v2/feeds/1/entries.json?page=1",
				Prev:  "https://api.feedbin.com/v2/feeds/1/entries.json?page=1",
				Next:  "https://api.feedbin.com/v2/feeds/1/entries.json?page=3",
				Last:  "https://api.feedbin.com/v2/feeds/1/entries.json?page=5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLinkHeader(tt.linkHeader)
			if got.First != tt.want.First {
				t.Errorf("ParseLinkHeader() First = %v, want %v", got.First, tt.want.First)
			}
			if got.Prev != tt.want.Prev {
				t.Errorf("ParseLinkHeader() Prev = %v, want %v", got.Prev, tt.want.Prev)
			}
			if got.Next != tt.want.Next {
				t.Errorf("ParseLinkHeader() Next = %v, want %v", got.Next, tt.want.Next)
			}
			if got.Last != tt.want.Last {
				t.Errorf("ParseLinkHeader() Last = %v, want %v", got.Last, tt.want.Last)
			}
		})
	}
}
