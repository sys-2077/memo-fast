package index

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// IndexRequest is the payload sent to the Cloud Run indexing API.
type IndexRequest struct {
	Files   []FilePayload   `json:"files"`
	Commits []CommitPayload `json:"commits"`
	Config  ConfigPayload   `json:"config"`
}

// FilePayload represents a single file in the index request.
type FilePayload struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// CommitPayload represents a single commit in the index request.
type CommitPayload struct {
	SHA     string   `json:"sha"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Date    string   `json:"date"`
	Files   []string `json:"files"`
}

// ConfigPayload carries the collection name for indexing.
type ConfigPayload struct {
	Collection string `json:"collection"`
}

// IndexResponse is the response from the Cloud Run indexing API.
type IndexResponse struct {
	IndexedFiles   int      `json:"indexed_files"`
	IndexedCommits int      `json:"indexed_commits"`
	Entities       int      `json:"entities"`
	Errors         []string `json:"errors"`
}

// Send posts the index request to the Cloud Run API and returns the response.
func Send(apiURL, apiKey string, req IndexRequest) (*IndexResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	url := apiURL + "/api/index"

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result IndexResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return &result, nil
}
