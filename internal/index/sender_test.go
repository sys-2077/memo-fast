package index

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSend_ParsesErrorsFieldFromAPIResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/index" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected authorization header: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"indexed_files": 2,
			"indexed_commits": 1,
			"entities": 4,
			"errors": ["extract error a.py", "embed error b.py"]
		}`))
	}))
	defer server.Close()

	req := IndexRequest{
		Files:  []FilePayload{{Path: "a.py", Content: "print(1)"}},
		Config: ConfigPayload{Collection: "memo_test"},
	}
	resp, err := Send(server.URL, "test-key", req)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if resp.IndexedFiles != 2 || resp.IndexedCommits != 1 || resp.Entities != 4 {
		t.Fatalf("unexpected counters: %+v", resp)
	}
	if len(resp.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(resp.Errors))
	}
}

func TestSend_ReturnsErrorOnNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized"}`))
	}))
	defer server.Close()

	req := IndexRequest{Config: ConfigPayload{Collection: "memo_test"}}
	_, err := Send(server.URL, "bad", req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "status 401") {
		t.Fatalf("unexpected error: %v", err)
	}
}
