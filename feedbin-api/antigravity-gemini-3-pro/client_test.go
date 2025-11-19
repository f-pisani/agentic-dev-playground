package feedbin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// setup sets up a test HTTP server along with a feedbin.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the Feedbin client being tested.
	client = NewClient("user", "pass")
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func TestNewClient(t *testing.T) {
	c := NewClient("user", "pass")
	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), defaultBaseURL)
	}
	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, want %v", c.UserAgent, userAgent)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient("user", "pass")

	inURL, outURL := "foo", defaultBaseURL+"foo"
	req, err := c.NewRequest("GET", inURL, nil)
	if err != nil {
		t.Fatalf("NewRequest(%q) returned error: %v", inURL, err)
	}

	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%q) URL = %v, want %v", inURL, req.URL.String(), outURL)
	}

	user, pass, ok := req.BasicAuth()
	if !ok || user != "user" || pass != "pass" {
		t.Errorf("NewRequest(%q) BasicAuth = %v, %v, %v; want user, pass, true", inURL, user, pass, ok)
	}

	if got := req.Header.Get("User-Agent"); got != c.UserAgent {
		t.Errorf("NewRequest(%q) User-Agent = %v, want %v", inURL, got, c.UserAgent)
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	type foo struct {
		A string
	}
	var f foo
	resp, err := client.Do(context.Background(), req, &f)
	if err != nil {
		t.Errorf("Do returned error: %v", err)
	}

	if f.A != "a" {
		t.Errorf("want a, got %s", f.A)
	}
	if resp == nil {
		t.Error("Do returned nil response")
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(`{"message":"error"}`)),
		Status:     "400 Bad Request",
	}
	err := CheckResponse(res)

	if err == nil {
		t.Errorf("CheckResponse returned nil error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("expected APIError, got %T", err)
	}
	if apiErr.Response != res {
		t.Errorf("expected response to be equal")
	}
	if apiErr.Message != `{"message":"error"}` {
		t.Errorf("expected error message to be read, got %s", apiErr.Message)
	}
}

func TestPaginationParsing(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/entries.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", `<https://api.feedbin.com/v2/entries.json?page=2>; rel="next", <https://api.feedbin.com/v2/entries.json?page=5>; rel="last"`)
		fmt.Fprint(w, `[]`)
	})

	req, _ := client.NewRequest("GET", "entries.json", nil)
	var v []interface{}
	resp, err := client.Do(context.Background(), req, &v)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	if resp.NextPage != 2 {
		t.Errorf("expected NextPage 2, got %d", resp.NextPage)
	}
	if resp.LastPage != 5 {
		t.Errorf("expected LastPage 5, got %d", resp.LastPage)
	}
}
