package feedbin

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("username", "password")

	if c.BaseURL.String() != DefaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), DefaultBaseURL)
	}

	if c.Username != "username" {
		t.Errorf("NewClient Username = %v, want %v", c.Username, "username")
	}

	if c.Password != "password" {
		t.Errorf("NewClient Password = %v, want %v", c.Password, "password")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	baseURL := "https://custom.api.example.com/"
	httpClient := &http.Client{}

	c := NewClient(
		"username",
		"password",
		WithBaseURL(baseURL),
		WithHTTPClient(httpClient),
	)

	if c.BaseURL.String() != baseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), baseURL)
	}

	if c.client != httpClient {
		t.Errorf("NewClient httpClient not set correctly")
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient("user", "pass")

	inURL, outURL := "foo", DefaultBaseURL+"foo"
	req, _ := c.NewRequest(http.MethodGet, inURL, nil)

	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL.String(), outURL)
	}

	username, password, ok := req.BasicAuth()
	if !ok {
		t.Error("NewRequest BasicAuth not set")
	}
	if username != "user" {
		t.Errorf("NewRequest BasicAuth username = %v, want %v", username, "user")
	}
	if password != "pass" {
		t.Errorf("NewRequest BasicAuth password = %v, want %v", password, "pass")
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("NewRequest Accept header = %v, want %v", req.Header.Get("Accept"), "application/json")
	}
}

func TestParseLink(t *testing.T) {
	linkHeader := `<https://api.feedbin.com/v2/feeds/1/entries.json?page=2>; rel="next", <https://api.feedbin.com/v2/feeds/1/entries.json?page=5>; rel="last"`

	info := parseLink(linkHeader)

	if info == nil {
		t.Fatal("parseLink returned nil")
	}

	if info.NextPage == nil {
		t.Error("NextPage is nil")
	} else if *info.NextPage != 2 {
		t.Errorf("NextPage = %v, want %v", *info.NextPage, 2)
	}

	if info.LastPage == nil {
		t.Error("LastPage is nil")
	} else if *info.LastPage != 5 {
		t.Errorf("LastPage = %v, want %v", *info.LastPage, 5)
	}
}

func TestClient_Do(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Setup client using test server URL
	baseURL, _ := url.Parse(server.URL + "/")
	client := NewClient("user", "pass")
	client.BaseURL = baseURL

	// Test handler for /test endpoint
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"ok"}`))
	})

	// Make test request
	req, _ := client.NewRequest(http.MethodGet, "test", nil)
	var v map[string]interface{}
	resp, err := client.Do(req, &v)

	// Check response
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Do status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	if v["result"] != "ok" {
		t.Errorf("Do result = %v, want %v", v["result"], "ok")
	}
}
